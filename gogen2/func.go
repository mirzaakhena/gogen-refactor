package gogen2

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"os"
	"strings"
	"unicode"
)

func logDebug(format string, a ...any) {
	//fmt.Printf(format, a...)
}

func GetSel(expr ast.Expr) string {

	switch fieldType := expr.(type) {

	case *ast.SelectorExpr:
		return GetSel(fieldType.Sel)

	case *ast.StarExpr:
		return GetSel(fieldType.X)

	case *ast.Ident:
		return fieldType.String()
	}

	return ""
}

func IsExported(varName string) bool {
	if !unicode.IsUpper(rune(varName[0])) && unicode.IsLetter(rune(varName[0])) {
		return false
	}
	return true
}

func ReadAllAvailablePackageFromImport(parsedGoFile string) ([]string, error) {

	const goModFileName = "go.mod"

	mapOfRequire := map[string]string{}

	dataInBytes, err := os.ReadFile(goModFileName)
	if err != nil {
		return nil, err
	}

	parsedGoMod, err := modfile.Parse(goModFileName, dataInBytes, nil)
	if err != nil {
		return nil, err
	}

	for _, r := range parsedGoMod.Require {

		if len(r.Syntax.Token) == 1 {
			mapOfRequire[r.Syntax.Token[0]] = fmt.Sprintf("%v/pkg/mod/%v", build.Default.GOPATH, r.Syntax.Token[0])
			continue
		}

		mapOfRequire[r.Syntax.Token[0]] = fmt.Sprintf("%v/pkg/mod/%v@%v", build.Default.GOPATH, r.Syntax.Token[0], r.Syntax.Token[1])
	}

	file, err := parser.ParseFile(token.NewFileSet(), parsedGoFile, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	packages := make([]string, 0)

	ast.Inspect(file, func(node ast.Node) bool {

		if err != nil {
			return false
		}

		switch n := node.(type) {

		case *ast.ImportSpec:

			pathValue := strings.Trim(n.Path.Value, `"`)

			for k, m := range mapOfRequire {

				if strings.HasPrefix(pathValue, k) {

					pathToLib := fmt.Sprintf("%v%v", m, pathValue[len(k):])

					pkgs, err := parser.ParseDir(token.NewFileSet(), pathToLib, nil, parser.PackageClauseOnly)
					if err != nil {
						return false
					}

					for _, pkg := range pkgs {
						packages = append(packages, pkg.Name)
					}

				}
			}

		}

		return true
	})

	return packages, nil
}
