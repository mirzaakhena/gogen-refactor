package gogen

import (
	"fmt"
	"go/ast"
	"strings"
)

func GetDefaultValue(expr ast.Expr) string {

	switch fieldType := expr.(type) {

	case *ast.SelectorExpr:
		ident, ok := fieldType.X.(*ast.Ident)
		if ok && ident.String() == "context" {
			return "ctx"
		} else {
			return fmt.Sprintf("%s.%s", ident.String(), GetDefaultValue(fieldType.Sel))
		}

	case *ast.StructType:
		//theFields := ""
		//for _, field := range fieldType.Fields.List {
		//	for _, name := range field.Names {
		//		theFields += fmt.Sprintf("%s: %s, ", name, GetDefaultValue(field.Type))
		//	}
		//}
		//return fmt.Sprintf("%v{ %s }", GetTypeAsString(fieldType), theFields)

		return fmt.Sprintf("%v{ }", GetTypeAsString(fieldType))

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

			} else if strings.HasPrefix(fieldType.String(), "int") {
				return "0"

			} else if strings.HasPrefix(fieldType.String(), "float") {
				return "0.0"

			} else if fieldType.String() == "string" {
				return `""`

			} else if fieldType.String() == "bool" {
				return "false"

			} else {
				return fieldType.String()
			}

		}

	}

	return "nil"
}
