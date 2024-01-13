package main

import (
	"github.com/thedenisnikulin/sqlchk"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(sqlchk.Analyzer)
}
