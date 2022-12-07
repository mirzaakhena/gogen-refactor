package gogen2

import (
	"go/ast"
	"strings"
)

type GogenType struct {
	Type         string `json:"type"`
	DefaultValue string `json:"defaultValue"`
}

func newGogenType(theType, defaultValue string) *GogenType {
	return &GogenType{
		Type:         theType,
		DefaultValue: defaultValue,
	}
}

type GogenField struct {
	Name     string     `json:"name"`
	DataType *GogenType `json:"dataType"`
}

func (g *GogenField) SetNewDefaultValue(newDefaultValue string) {
	g.DataType.DefaultValue = newDefaultValue
}

func NewGogenField(fieldName, theType, defaultValue string) *GogenField {
	return &GogenField{
		Name:     fieldName,
		DataType: newGogenType(theType, defaultValue),
	}
}

type GogenImport struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Expression string `json:"expression"`
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

func NewGogenImport(importSpec *ast.ImportSpec) GogenImport {

	importPath := strings.Trim(importSpec.Path.Value, `"`)

	name := ""
	expr := importPath[strings.LastIndex(importPath, "/")+1:]
	if importSpec.Name != nil {
		name = importSpec.Name.String()
		expr = name
	}

	return GogenImport{
		Name:       name,
		Path:       importPath,
		Expression: expr,
	}
}

type GogenMethod struct {
	Name    string        `json:"name"`
	Params  []*GogenField `json:"params"`
	Results []*GogenField `json:"results"`
	//TypeParams []*GogenField `json:"typeParams"`
}

func NewGogenMethod(methodName string) *GogenMethod {
	return &GogenMethod{
		Name:    methodName,
		Params:  make([]*GogenField, 0),
		Results: make([]*GogenField, 0),
	}
}

//func (g *GogenMethod) AddParam(name string, dataType *GogenType) {
//	g.Params = append(g.Params, &GogenField{
//		Name:     name,
//		DataType: dataType,
//	})
//}
//
//func (g *GogenMethod) AddResult(name string, dataType *GogenType) {
//	g.Results = append(g.Results, &GogenField{
//		Name:     name,
//		DataType: dataType,
//	})
//}

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
