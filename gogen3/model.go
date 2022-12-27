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
	Name         FieldType
	Expr         ast.Expr
	DefaultValue string
	File         *ast.File
	//TypeProperties *TypeProperties
	//Imports map[Expression]*GogenImport
}

type GogenField struct {
	Name     GogenFieldName
	DataType *GogenFieldType
}

type GogenMethod struct {
	Name    GogenMethodName
	Params  []*GogenField
	Results []*GogenField
}

type GogenInterface struct {
	//CurrentPackage *PackageName
	InterfaceType *GogenFieldType
	Interfaces    []*GogenInterface
	Methods       []*GogenMethod
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
