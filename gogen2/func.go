package gogen2

import (
	"fmt"
	"go/ast"
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

func handleDefaultValue(typeAsString string, expr ast.Expr) string {

	logDebug("handleDefaultValue %v %v\n", typeAsString, expr)

	switch exprType := expr.(type) {
	case *ast.Ident:

		logDebug("as ident            : %v\n", exprType.String())

		// found type in the same file
		if exprType.Obj != nil {
			logDebug("dataType %s ada di file yg sama dengan struct target\n", exprType.String())

			typeSpec, ok := exprType.Obj.Decl.(*ast.TypeSpec)
			if !ok {
				panic("cannot assert to TypeSpec")
			}

			logDebug("start recursive handleDefaultValue utk dataType %v\n", exprType.String())
			defaultValue := handleDefaultValue(typeAsString, typeSpec.Type)
			logDebug("end   recursive handleDefaultValue dari type %v hasil recursive adalah %v\n", exprType.String(), defaultValue)

			return defaultValue
		}

		basicType := ""

		for {

			if strings.HasPrefix(exprType.String(), "int") || strings.HasPrefix(exprType.String(), "uint") {
				logDebug("as int / uint\n")
				basicType = "0"
				break
			}

			if strings.HasPrefix(exprType.String(), "float") {
				logDebug("as float\n")
				basicType = "0.0"
				break
			}

			if exprType.String() == "string" {
				logDebug("as string\n")
				basicType = `""`
				break
			}

			if exprType.String() == "bool" {
				logDebug("as bool\n")
				basicType = `false`
				break
			}

			if exprType.String() == "any" {
				logDebug("as any\n")
				basicType = `nil`
				break
			}

			break
		}

		if basicType != "" {
			if typeAsString != exprType.String() {
				return fmt.Sprintf("%s(%s)", typeAsString, basicType)
			}
			return basicType
		}

		logDebug("tipe data dasar masih belum diketahui\n")

		return exprType.String()

	case *ast.StructType:
		v := fmt.Sprintf("%v{}", typeAsString)
		logDebug("as struct %v\n", v)
		return v

	case *ast.ArrayType:
		v := fmt.Sprintf("%s{}", typeAsString)
		logDebug("as array %v\n", v)
		return v

	case *ast.SelectorExpr:

		ident, ok := exprType.X.(*ast.Ident)

		// hardcoded fo context
		if ok && ident.String() == "context" {
			return "ctx"
		}

		v := fmt.Sprintf("%s.%s", ident.String(), exprType.Sel.String())

		logDebug("as selector %v\n", v)

		return v

	case *ast.StarExpr:
		//a := getTypeAsString(exprType.X)
		//if a == "nil" {
		//	return "nil"
		//}
		//v := fmt.Sprintf("&%s{}", a)
		//return v
		return "nil"

	case *ast.InterfaceType:
		return "nil"

	case *ast.MapType:
		return "nil"

	case *ast.ChanType:
		return "nil"

	case *ast.FuncType:
		return "nil"

	}

	return "unknown"

}

func IsExported(varName string) bool {
	if !unicode.IsUpper(rune(varName[0])) && unicode.IsLetter(rune(varName[0])) {
		return false
	}
	return true
}
