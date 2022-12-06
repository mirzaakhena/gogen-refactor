package gogen3

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

//func (g *GogenStruct) tryAssignDefaultValue(name string, dataType *GogenType) {
//
//	if dataType.UncompleteType() {
//
//		typeSpecFromMap, exist := gsb.typeMap[theType]
//		if !exist {
//			// register to unknown type to be find later
//			gsb.unknownTypes[theType] = &gt
//			return &gt, nil
//		}
//
//		gt.DefaultValue, _, err = gsb.handleDefaultValue(defaultValue, typeSpecFromMap.Type)
//		if err != nil {
//			return nil, err
//		}
//
//	}
//
//	g.Fields = append(g.Fields, GogenField{
//		Name:     name,
//		DataType: dataType,
//	})
//}

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

func (g *GogenMethod) AddParam(name string, dataType *GogenType) {
	g.Params = append(g.Params, &GogenField{
		Name:     name,
		DataType: dataType,
	})
}

func (g *GogenMethod) AddResult(name string, dataType *GogenType) {
	g.Results = append(g.Results, &GogenField{
		Name:     name,
		DataType: dataType,
	})
}

type GogenInterface struct {
	Name    string         `json:"name"`
	Imports []GogenImport  `json:"imports"`
	Methods []*GogenMethod `json:"methods"`
}
