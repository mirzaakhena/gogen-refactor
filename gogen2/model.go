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
	Selector     string
	FieldType    string
	FieldName    string
	MethodName   string
	RequirePath  string
	CompletePath string
)

const (
	ImportTypeGoSDK     ImportType = "GO_SDK"
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

	logDebug("handleDefaultValue %v %v", g.DataType.Type, expr)

	g.DataType.Expr = expr

	switch exprType := expr.(type) {
	case *ast.Ident:

		logDebug("as ident            : %v", exprType.String())

		// found type in the same file
		if exprType.Obj != nil {
			logDebug("dataType %s ada di file yg sama dengan struct target", exprType.String())

			typeSpec, ok := exprType.Obj.Decl.(*ast.TypeSpec)
			if !ok {
				panic("cannot assert to TypeSpec")
			}

			logDebug("start recursive handleDefaultValue utk dataType %v", exprType.String())

			g.handleDefaultValue(typeSpec.Type)
			logDebug("end   recursive handleDefaultValue dari type %v hasil recursive adalah %v", exprType.String(), g.DataType.DefaultValue)

			return
		}

		basicDefaultValue := ""

		for {

			if strings.HasPrefix(exprType.String(), "int") || strings.HasPrefix(exprType.String(), "uint") {
				logDebug("as int / uint")
				basicDefaultValue = "0"
				break
			}

			if strings.HasPrefix(exprType.String(), "float") {
				logDebug("as float")
				basicDefaultValue = "0.0"
				break
			}

			if exprType.String() == "string" {
				logDebug("as string")
				basicDefaultValue = `""`
				break
			}

			if exprType.String() == "bool" {
				logDebug("as bool")
				basicDefaultValue = `false`
				break
			}

			if exprType.String() == "any" {
				logDebug("as any")
				basicDefaultValue = `nil`
				break
			}

			if exprType.String() == "error" {
				logDebug("as error")
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

		logDebug("tipe data dasar masih belum diketahui")

		g.DataType.DefaultValue = exprType.String()

	case *ast.StructType:
		v := fmt.Sprintf("%v{}", g.DataType.Type)
		logDebug("as struct %v", v)
		g.DataType.DefaultValue = v

	case *ast.ArrayType:
		v := fmt.Sprintf("%s{}", g.DataType.Type)
		logDebug("as array %v", v)
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

		logDebug("as selector %v", v)

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
	logDebug("tipe data           : %v", dataTypeStr)

	gf := GogenField{
		Name:     fieldName,
		DataType: newGogenType(dataTypeStr),
	}

	logDebug("first time handleDefaultValue")
	gf.handleDefaultValue(expr)
	logDebug("default value       : %v", gf.DataType.DefaultValue)

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
	Name           MethodName    `json:"name"`
	OwnerInterface FieldType     `json:"ownerInterface"`
	Params         []*GogenField `json:"params"`
	Results        []*GogenField `json:"results"`
}

func NewGogenMethod(ownerInterface FieldType, methodName MethodName) *GogenMethod {
	return &GogenMethod{
		OwnerInterface: ownerInterface,
		Name:           methodName,
		Params:         make([]*GogenField, 0),
		Results:        make([]*GogenField, 0),
	}
}

//type GogenMethodField struct {
//	ExtendInterfaces string `json:"extendInterfaces"`
//}

type GogenInterface struct {
	Name    FieldType      `json:"name"`
	Imports []GogenImport  `json:"imports"`
	Methods []*GogenMethod `json:"methods"`
}

func NewGogenInterface(interfaceName string) *GogenInterface {
	return &GogenInterface{
		Name:    FieldType(interfaceName),
		Imports: make([]GogenImport, 0),
		Methods: make([]*GogenMethod, 0),
	}
}
