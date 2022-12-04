package gogen

import "go/ast"

func getSel(expr ast.Expr) string {
	switch fieldType := expr.(type) {
	case *ast.SelectorExpr:
		return fieldType.Sel.String()
	case *ast.StarExpr:
		return getSel(fieldType.X)
	}
	return ""
}

func getExprForImport(expr ast.Expr) []PathExpression {

	switch fieldType := expr.(type) {
	case *ast.SelectorExpr:
		return []PathExpression{PathExpression(fieldType.X.(*ast.Ident).String())}

	case *ast.StarExpr:
		return getExprForImport(fieldType.X)

	case *ast.MapType:
		str := make([]PathExpression, 0)
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
		expressions := make([]PathExpression, 0)

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
