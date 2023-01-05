package gogen5

import (
	"go/ast"
)

type (
	GogenAnyTypeName   string
	GogenMethodName    string
	GogenFieldName     string
	GogenFieldTypeName string
	ImportType         string
	ImportPath         string
	Expression         string
	Version            string
)

type TypeProperties struct {
	TypeSpec *ast.TypeSpec
	AstFile  *ast.File
}

func (r GogenFieldTypeName) String() string {
	return string(r)
}

func (r GogenMethodName) String() string {
	return string(r)
}

func (r GogenFieldName) String() string {
	return string(r)
}

func (r GogenAnyTypeName) String() string {
	return string(r)
}

func (r Expression) String() string {
	return string(r)
}

func (r ImportPath) String() string {
	return string(r)
}

func (r ImportType) String() string {
	return string(r)
}

const (
	ImportTypeGoSDK           ImportType = "GO_SDK"
	ImportTypeExternalModule  ImportType = "EXTERNAL_MODULE"
	ImportTypeInternalProject ImportType = "INTERNAL_PROJECT"
)

type GogenImport struct {
	Name         string     `json:"name"`
	Path         ImportPath `json:"path"`
	CompletePath string     `json:"-"`
	Expression   Expression `json:"expression"`
	ImportType   ImportType `json:"importType"`
}

type GogenFieldType struct {
	Name         GogenFieldTypeName `json:"name,omitempty"`         // real type like []*repo.SaveOrderRepo
	DefaultValue string             `json:"defaultValue,omitempty"` //
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

func (r *GogenMethod) AddParam(gf *GogenField) {
	r.Params = append(r.Params, gf)
}

func (r *GogenMethod) AddResult(gf *GogenField) {
	r.Results = append(r.Results, gf)
}

type GogenAnyType struct {
	Name             GogenAnyTypeName           `json:"name"`             // short type []*repo.SaveOrderRepo --> SaveOrderRepo
	CompositionTypes []*GogenAnyType            `json:"compositionTypes"` //
	Fields           []*GogenField              `json:"fields"`           //
	Methods          []*GogenMethod             `json:"methods"`          //
	Imports          map[Expression]GogenImport `json:"imports"`          //
}

type GoModProperties struct {
	AbsolutePathProject string
	ModuleName          string
	RequirePath         map[ImportPath]Version
}

func NewGogenField(name string, expr ast.Expr) *GogenField {

	return &GogenField{
		Name: GogenFieldName(name),
		DataType: &GogenFieldType{
			Name:         NewGogenFieldTypeName(expr),
			DefaultValue: "",
		},
	}

}

func NewGogenMethod(methodName string) *GogenMethod {
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

func NewGogenAnyType(name string) *GogenAnyType {
	return &GogenAnyType{
		Name:             GogenAnyTypeName(name),
		CompositionTypes: make([]*GogenAnyType, 0),
		Methods:          make([]*GogenMethod, 0),
		Fields:           make([]*GogenField, 0),
		Imports:          map[Expression]GogenImport{},
	}
}

func NewGogenFieldTypeName(expr ast.Expr) GogenFieldTypeName {
	return GogenFieldTypeName(GetTypeAsString(expr))
}

func (m *GogenAnyType) AddField(gf *GogenField) {
	m.Fields = append(m.Fields, gf)
}

func (m *GogenAnyType) AddMethod(gm *GogenMethod) {
	m.Methods = append(m.Methods, gm)
}

func (m *GogenAnyType) AddCompositionType(gat *GogenAnyType) {
	m.CompositionTypes = append(m.CompositionTypes, gat)
}
