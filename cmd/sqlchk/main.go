package main

import (
	sqlchk "sqlchk"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(sqlchk.Analyzer)
}
