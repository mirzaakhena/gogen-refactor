package gogen2

import (
	"fmt"
	"go/ast"
	"strings"
)

type (
	ImportPath   string
	ImportType   string
	ModulePath   string
	Expression   string
	FieldType    string
	FieldName    string
	RequirePath  string
	CompletePath string
)

const (
	ImportTypeSDK       ImportType = "ImportTypeSDK"
	ImportTypeExtModule ImportType = "EXT_MODULE"
	ImportTypeProject   ImportType = "PROJECT"
)

type GogenType struct {
	Type         FieldType `json:"type"`
	Expr         ast.Expr  `json:"-"`
	DefaultValue string    `json:"defaultValue"`
}

func newGogenType(theType FieldType) *GogenType {
	return &GogenType{
		Type:         theType,
		Expr:         nil,
		DefaultValue: "",
	}
}

type GogenField struct {
	Name     FieldName  `json:"name"`
	DataType *GogenType `json:"dataType"`
}

func (g *GogenField) handleDefaultValue(expr ast.Expr) {

	logDebug("handleDefaultValue %v %v\n", g.DataType.Type, expr)

	g.DataType.Expr = expr

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

			g.handleDefaultValue(typeSpec.Type)
			logDebug("end   recursive handleDefaultValue dari type %v hasil recursive adalah %v\n", exprType.String(), g.DataType.DefaultValue)

			return
		}

		basicDefaultValue := ""

		for {

			if strings.HasPrefix(exprType.String(), "int") || strings.HasPrefix(exprType.String(), "uint") {
				logDebug("as int / uint\n")
				basicDefaultValue = "0"
				break
			}

			if strings.HasPrefix(exprType.String(), "float") {
				logDebug("as float\n")
				basicDefaultValue = "0.0"
				break
			}

			if exprType.String() == "string" {
				logDebug("as string\n")
				basicDefaultValue = `""`
				break
			}

			if exprType.String() == "bool" {
				logDebug("as bool\n")
				basicDefaultValue = `false`
				break
			}

			if exprType.String() == "any" {
				logDebug("as any\n")
				basicDefaultValue = `nil`
				break
			}

			break
		}

		if basicDefaultValue != "" {
			if string(g.DataType.Type) != exprType.String() {
				g.DataType.DefaultValue = fmt.Sprintf("%s(%s)", g.DataType.Type, basicDefaultValue)
				return
			}
			g.DataType.DefaultValue = basicDefaultValue
			return
		}

		logDebug("tipe data dasar masih belum diketahui\n")

		g.DataType.DefaultValue = exprType.String()

	case *ast.StructType:
		v := fmt.Sprintf("%v{}", g.DataType.Type)
		logDebug("as struct %v\n", v)
		g.DataType.DefaultValue = v

	case *ast.ArrayType:
		v := fmt.Sprintf("%s{}", g.DataType.Type)
		logDebug("as array %v\n", v)
		g.DataType.DefaultValue = v

	case *ast.SelectorExpr:

		ident, ok := exprType.X.(*ast.Ident)
		if !ok {
			panic("cannot assert to Ident")
		}

		// hardcoded fo context
		//if ok && ident.String() == "context" {
		//	return "ctx"
		//}

		v := fmt.Sprintf("%s.%s", ident.String(), exprType.Sel.String())

		logDebug("as selector %v\n", v)

		g.DataType.DefaultValue = v

	case *ast.StarExpr:
		//a := getTypeAsString(exprType.X)
		//if a == "nil" {
		//	return "nil"
		//}
		//v := fmt.Sprintf("&%s{}", a)
		//return v
		g.DataType.DefaultValue = "nil"

	case *ast.InterfaceType:
		g.DataType.DefaultValue = "nil"

	case *ast.MapType:
		g.DataType.DefaultValue = "nil"

	case *ast.ChanType:
		g.DataType.DefaultValue = "nil"

	case *ast.FuncType:
		g.DataType.DefaultValue = "nil"

	default:
		g.DataType.DefaultValue = "unknown"
	}

}

func NewGogenField(fieldName FieldName, expr ast.Expr) *GogenField {

	dataTypeStr := FieldType(getTypeAsString(expr))
	logDebug("tipe data           : %v\n", dataTypeStr)

	gf := GogenField{
		Name:     fieldName,
		DataType: newGogenType(dataTypeStr),
	}

	logDebug("first time handleDefaultValue\n")
	gf.handleDefaultValue(expr)
	logDebug("default value       : %v\n", gf.DataType.DefaultValue)

	return &gf
}

type GogenImport struct {
	Name         string       `json:"name"`
	Path         ImportPath   `json:"path"`
	CompletePath CompletePath `json:"completePath"`
	Expression   Expression   `json:"expression"`
	ImportType   ImportType   `json:"importType"`
}

type GogenStruct struct {
	Name    string        `json:"name"`
	Imports []GogenImport `json:"imports"`
	Fields  []*GogenField `json:"fields"`
}

func NewGogenStruct(structName string) *GogenStruct {
	return &GogenStruct{
		Name:    structName,
		Imports: make([]GogenImport, 0),
		Fields:  make([]*GogenField, 0),
	}
}

//func NewGogenImport(importSpec *ast.ImportSpec) GogenImport {
//
//	importPath := strings.Trim(importSpec.Path.Value, `"`)
//
//	name := ""
//	expr := importPath[strings.LastIndex(importPath, "/")+1:]
//	if importSpec.Name != nil {
//		name = importSpec.Name.String()
//		expr = name
//	}
//
//	return GogenImport{
//		Name:       name,
//		Path:       importPath,
//		Expression: Expression(expr),
//	}
//}

type GogenMethod struct {
	Name    string        `json:"name"`
	Params  []*GogenField `json:"params"`
	Results []*GogenField `json:"results"`
}

func NewGogenMethod(methodName string) *GogenMethod {
	return &GogenMethod{
		Name:    methodName,
		Params:  make([]*GogenField, 0),
		Results: make([]*GogenField, 0),
	}
}

type GogenMethodField struct {
	ExtendInterfaces string         `json:"extendInterfaces"`
	Methods          []*GogenMethod `json:"methods"`
}

type GogenInterface struct {
	Name         string              `json:"name"`
	Imports      []GogenImport       `json:"imports"`
	MethodFields []*GogenMethodField `json:"methodFields"`
}

func NewGogenInterface(interfaceName string) *GogenInterface {
	return &GogenInterface{
		Name:         interfaceName,
		Imports:      make([]GogenImport, 0),
		MethodFields: make([]*GogenMethodField, 0),
	}
}
