package sqlchk

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/jmoiron/sqlx"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = `
Validate SQL queries by running them on a database.

Requires setting DATABASE_URL environment variable with the url of the database
on which the queries are going to be run.
The SQL queries in your code are required to start with "--sql" prefix to be 
recognized by sqlchk.
`

var Analyzer = &analysis.Analyzer{
	Name:      "sqlchk",
	Doc:       doc,
	Run:       run,
	FactTypes: []analysis.Fact{},
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	databaseUrl := os.Getenv("DATABASE_URL")
	scheme, _, ok := strings.Cut(databaseUrl, "://")
	if !ok {
		return nil, fmt.Errorf("database url has no scheme")
	}

	db, err := sqlx.Connect(scheme, databaseUrl+"?sslmode=disable")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	filterNodes := []ast.Node{
		(*ast.BasicLit)(nil),
	}
	inspect.Preorder(filterNodes, func(n ast.Node) {
		node := n.(*ast.BasicLit)
		if node.Kind != token.STRING {
			return
		}

		var str string
		if strings.HasPrefix(node.Value, "`") {
			str = strings.Trim(node.Value, "`")
		} else {
			str = strings.TrimRight(node.Value, "\"")
		}

		if !strings.HasPrefix(str, "--sql") {
			return
		}

		sqlStr := str
		stmt, err := db.Prepare(sqlStr)
		if err != nil {
			pass.Reportf(node.Pos(), "%s", err.Error())
			return
		}
		stmt.Close()
	})

	return nil, nil
}
