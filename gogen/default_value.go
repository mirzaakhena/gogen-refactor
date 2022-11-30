package gogen

import (
	"fmt"
	"go/ast"
	"strings"
)

func GetDefaultValue(expr ast.Expr) string {

	//defaultValue := "nil"

	switch fieldType := expr.(type) {

	case *ast.SelectorExpr:
		ident, ok := fieldType.X.(*ast.Ident)
		if ok && ident.String() == "context" {
			return "ctx"
		} else {
			return fmt.Sprintf("%s.%s", ident.String(), GetDefaultValue(fieldType.Sel))
		}

	case *ast.StructType:
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

			} else if strings.HasPrefix(fieldType.String(), "int") {
				return "0"

			} else if strings.HasPrefix(fieldType.String(), "float") {
				return "0.0"

			} else if fieldType.String() == "string" {
				return `""`

			} else if fieldType.String() == "bool" {
				return "false"

			} else {
				fmt.Printf("xx>>>>>%s \n", fieldType.String())

				return fieldType.String()
			}

		}

	}

	return "nil"
}
