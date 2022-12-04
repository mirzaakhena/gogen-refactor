package gogen2

import (
	"go/ast"
	"strings"
)

type GogenType struct {
	Type             string `json:"type"`
	JSONDefaultValue string `json:"jsonDefaultValue"`
	TypeDefaultValue string `json:"typeDefaultValue"`
}

func NewGogenType(theType, defaultValue string) *GogenType {
	return &GogenType{
		Type:             theType,
		TypeDefaultValue: defaultValue,
		JSONDefaultValue: defaultValue,
	}
}

type GogenField struct {
	Name     string     `json:"name"`
	DataType *GogenType `json:"dataType"`
}

type GogenImport struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Expression string `json:"expression"`
}

type GogenStruct struct {
	Name    string        `json:"name"`
	Imports []GogenImport `json:"imports"`
	Fields  []GogenField  `json:"fields"`
}

func NewGogenStruct(structName string) *GogenStruct {
	return &GogenStruct{
		Name:    structName,
		Imports: make([]GogenImport, 0),
		Fields:  make([]GogenField, 0),
	}
}

func (g *GogenStruct) AddField(name string, dataType *GogenType) {
	g.Fields = append(g.Fields, GogenField{
		Name:     name,
		DataType: dataType,
	})
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
	Name    string       `json:"name"`
	Params  []GogenField `json:"params"`
	Results []GogenField `json:"results"`
}

type GogenInterface struct {
	Name    string        `json:"name"`
	Imports []GogenImport `json:"imports"`
	Methods []GogenMethod `json:"methods"`
}
