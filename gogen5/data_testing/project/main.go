package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "simple/file1.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Print(fset, file)

}
