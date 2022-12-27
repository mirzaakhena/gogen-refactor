package gogen3

import (
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
	File     *ast.File
	TypeSpec *ast.TypeSpec
	//Imports map[Expression]*GogenImport
}

type GoModProperties struct {
	AbsolutePathProject string
	ModuleName          string
	RequirePath         map[ImportPath]Version
}