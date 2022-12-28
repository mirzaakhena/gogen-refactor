package gogen4

import (
	"fmt"
	"go/ast"
)

type (
	GogenMethodName      string
	GogenInterfaceName   string
	GogenFieldName       string
	FieldType            string
	PackageName          string
	ImportType           string
	ImportPath           string
	ImportName           string
	AbsolutePath         string
	Expression           string
	Version              string
	FieldMethodSignature string
)

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

type GogenInterface struct {
	//CurrentPackage *PackageName
	InterfaceType *GogenFieldType   `json:"interfaceType,omitempty"`
	Interfaces    []*GogenInterface `json:"interfaces,omitempty"`
	Methods       []*GogenMethod    `json:"methods,omitempty"`
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

func NewGogenInterface() *GogenInterface {
	return &GogenInterface{
		InterfaceType: nil,
		Interfaces:    make([]*GogenInterface, 0),
		Methods:       make([]*GogenMethod, 0),
	}
}

func NewFieldMethodSignature(m GogenMethodName, f GogenFieldName) FieldMethodSignature {
	return FieldMethodSignature(fmt.Sprintf("%v.%v", m, f))
}
