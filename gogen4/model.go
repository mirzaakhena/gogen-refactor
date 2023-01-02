package gogen4

import (
	"fmt"
	"go/ast"
)

type (
	GogenMethodName    string
	GogenInterfaceName string
	GogenFieldName     string
	FieldType          string
	PackageName        string
	ImportType         string
	ImportPath         string
	ImportName         string
	AbsolutePath       string
	Expression         string
	Version            string
	FieldSignature     string
)

func (r FieldType) String() string {
	return string(r)
}

func (r GogenMethodName) String() string {
	return string(r)
}

func (r GogenFieldName) String() string {
	return string(r)
}

const (
	ImportTypeGoSDK           ImportType = "GO_SDK"
	ImportTypeExternalModule  ImportType = "EXTERNAL_MODULE"
	ImportTypeInternalProject ImportType = "INTERNAL_PROJECT"
)

type GogenImport struct {
	Name         string       `json:"name"`
	Path         ImportPath   `json:"path"`
	CompletePath AbsolutePath `json:"completePath"`
	Expression   Expression   `json:"expression"`
	ImportType   ImportType   `json:"importType"`
}

type GogenFieldType struct {
	Name         FieldType `json:"name,omitempty"`
	Expr         ast.Expr  `json:"-"`
	DefaultValue string    `json:"defaultValue,omitempty"`
	File         *ast.File `json:"-"`
}

type GogenField struct {
	Name     GogenFieldName  `json:"name,omitempty"`
	DataType *GogenFieldType `json:"dataType,omitempty"`
}

type GogenMethod struct {
	Name    GogenMethodName `json:"name,omitempty"`
	Params  []*GogenField   `json:"params,omitempty"`
	Results []*GogenField   `json:"results,omitempty"`
}

//type GogenStruct struct {
//	StructType *GogenFieldType `json:"structType,omitempty"`
//	Types      []*GogenStruct  `json:"types,omitempty"`
//	Fields     []*GogenField   `json:"fields,omitempty"`
//	Methods    []*GogenMethod  `json:"methods,omitempty"`
//}

type GogenAnyType struct {
	GogenFieldType   *GogenFieldType `json:"fieldType,omitempty"`
	CompositionTypes []*GogenAnyType `json:"compositionTypes,omitempty"`
	Fields           []*GogenField   `json:"fields,omitempty"`
	Methods          []*GogenMethod  `json:"methods,omitempty"`
}

type TypeProperties struct {
	AstFile  *ast.File
	TypeSpec *ast.TypeSpec
	//Imports map[Expression]*GogenImport
}

type GoModProperties struct {
	AbsolutePathProject string
	ModuleName          string
	RequirePath         map[ImportPath]Version
}

func NewGogenField(name string, expr ast.Expr, fieldType FieldType, astFile *ast.File) *GogenField {

	return &GogenField{
		Name: GogenFieldName(name),
		DataType: &GogenFieldType{
			Name:         fieldType,
			Expr:         expr,
			DefaultValue: "",
			File:         astFile,
		},
	}

}

func newGogenMethod(methodName string) *GogenMethod {
	return &GogenMethod{
		Name:    GogenMethodName(methodName),
		Params:  make([]*GogenField, 0),
		Results: make([]*GogenField, 0),
	}
}

func NewGoModProperties() GoModProperties {
	return GoModProperties{
		AbsolutePathProject: "",
		ModuleName:          "",
		RequirePath:         map[ImportPath]Version{},
	}
}

func NewGogenAnyType() *GogenAnyType {
	return &GogenAnyType{
		GogenFieldType:   nil,
		CompositionTypes: make([]*GogenAnyType, 0),
		Methods:          make([]*GogenMethod, 0),
		Fields:           make([]*GogenField, 0),
	}
}

func NewFieldSignature(m string, f string) FieldSignature {
	return FieldSignature(fmt.Sprintf("%v.%v", m, f))
}
