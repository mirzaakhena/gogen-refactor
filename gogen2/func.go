package gogen2

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"
)

func GetXAndSelAsString(fieldType *ast.SelectorExpr) (string, string) {
	return fieldType.X.(*ast.Ident).String(), fieldType.Sel.String()
}

func GetDefaultValue(expr ast.Expr) string {

	switch fieldType := expr.(type) {

	case *ast.SelectorExpr:

		ident, ok := fieldType.X.(*ast.Ident)

		// hardcoded fo context
		if ok && ident.String() == "context" {
			return "ctx"
		}

		return fmt.Sprintf("%s.%s", ident.String(), GetDefaultValue(fieldType.Sel))

	case *ast.ArrayType:
		return fmt.Sprintf("%s{}", GetTypeAsString(fieldType))

	case *ast.StructType:
		//theFields := ""
		//for _, field := range fieldType.Fields.List {
		//	for _, name := range field.Names {
		//		theFields += fmt.Sprintf("%s: %s, ", name, GetDefaultValue(field.Type))
		//	}
		//}
		//return fmt.Sprintf("%v{ %s }", GetTypeAsString(fieldType), theFields)

		return fmt.Sprintf("%s{}", GetTypeAsString(fieldType))

	case *ast.Ident:
		if fieldType.Obj != nil {

			typeSpec, ok := fieldType.Obj.Decl.(*ast.TypeSpec)
			if !ok {
				return ""
			}

			return GetDefaultValue(typeSpec.Type)

		} else {

			if fieldType.String() == "error" {
				return "nil"

			}

			if strings.HasPrefix(fieldType.String(), "int") {
				return "0"

			}

			if strings.HasPrefix(fieldType.String(), "uint") {
				return "0"

			}

			if strings.HasPrefix(fieldType.String(), "float") {
				return "0.0"

			}

			if fieldType.String() == "string" {
				return `""`

			}

			if fieldType.String() == "bool" {
				return "false"

			}

			if fieldType.String() == "any" {
				return "nil"
			}

			return fieldType.String()

		}

	}

	return "nil"
}

func GetDeepDefaultValue(expr ast.Expr, myDefaultValue string) (string, string) {
	switch expr.(type) {

	case *ast.StructType:
		return fmt.Sprintf("%s{}", myDefaultValue), "{}"

	case *ast.ArrayType:
		return myDefaultValue, "[]"

	case *ast.Ident:
		realDefaultValue := GetDefaultValue(expr)
		return fmt.Sprintf("%s(%s)", myDefaultValue, realDefaultValue), realDefaultValue

	}

	return "nil", "null"
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

func getExprForImport(expr ast.Expr) []string {

	switch fieldType := expr.(type) {
	case *ast.SelectorExpr:
		return []string{fieldType.X.(*ast.Ident).String()}

	case *ast.StarExpr:
		return getExprForImport(fieldType.X)

	case *ast.MapType:
		str := make([]string, 0)
		key := getExprForImport(fieldType.Key)
		if key != nil {
			str = append(str, key...)
		}
		value := getExprForImport(fieldType.Value)
		if value != nil {
			str = append(str, value...)
		}

		return str

	case *ast.ArrayType:
		return getExprForImport(fieldType.Elt)

	case *ast.ChanType:
		return getExprForImport(fieldType.Value)

	case *ast.FuncType:
		expressions := make([]string, 0)

		if fieldType.Params.NumFields() > 0 {
			for _, x := range fieldType.Params.List {
				expressions = append(expressions, getExprForImport(x.Type)...)
			}
		}

		if fieldType.Results.NumFields() > 0 {
			for _, x := range fieldType.Results.List {
				expressions = append(expressions, getExprForImport(x.Type)...)
			}
		}

		return expressions

	}

	return nil

}

func handleImport(genDecl *ast.GenDecl, importMap map[string]GogenImport) {

	for _, spec := range genDecl.Specs {
		importSpec, ok := spec.(*ast.ImportSpec)
		if !ok {
			continue
		}
		gi := NewGogenImport(importSpec)
		importMap[gi.Expression] = gi
	}

}

func IsExported(varName string) bool {
	if !unicode.IsUpper(rune(varName[0])) && unicode.IsLetter(rune(varName[0])) {
		return false
	}
	return true
}
