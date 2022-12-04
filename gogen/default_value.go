package gogen

import (
	"fmt"
	"go/ast"
	"strings"
)

func GetDeepDefaultValue(expr ast.Expr, myDefaultValue string) string {
	switch expr.(type) {

	case *ast.StructType:
		return fmt.Sprintf("%s{}", myDefaultValue)

	case *ast.Ident:
		return fmt.Sprintf("%s(%s)", myDefaultValue, GetDefaultValue(expr))

	}

	return "nil"
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

	case *ast.StructType:
		//theFields := ""
		//for _, field := range fieldType.Fields.List {
		//	for _, name := range field.Names {
		//		theFields += fmt.Sprintf("%s: %s, ", name, GetDefaultValue(field.Type))
		//	}
		//}
		//return fmt.Sprintf("%v{ %s }", GetTypeAsString(fieldType), theFields)

		return fmt.Sprintf("%v{}", GetTypeAsString(fieldType))

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
