package gogen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type GogenName string

type GogenImport struct {
	Name       string `json:"name,omitempty"`
	Path       string `json:"path"`
	Expression string `json:"expression"`
}

func NewGogenImport(importSpec *ast.ImportSpec) GogenImport {

	importPath := importSpec.Path.Value

	name := ""
	expr := strings.Trim(importPath[strings.LastIndex(importPath, "/")+1:], `"`)
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

type GogenType struct {
	Type         string                 `json:"type"`
	DefaultValue string                 `json:"defaultValue"`
	Imports      map[string]GogenImport `json:"imports,omitempty"`
}

func NewGogenType(gomodPath string, astField *ast.Field, importMap map[string]GogenImport) GogenType {

	usedMap := map[string]GogenImport{}

	//var theType = astField.Type

	myType := GetTypeAsString(astField.Type)
	myDefaultValue := GetDefaultValue(astField.Type)

	switch astFieldType := astField.Type.(type) {
	case *ast.SelectorExpr:

		x := astFieldType.X.(*ast.Ident).String()
		sel := astFieldType.Sel.String()
		gi := importMap[x]

		usedMap[gi.Expression] = gi

		path := strings.Trim(gi.Path, "\"")
		if strings.HasPrefix(path, gomodPath) {

			fset := token.NewFileSet()
			pkgs, err := parser.ParseDir(fset, path[len(gomodPath)+1:], nil, parser.ParseComments)
			if err != nil {
				panic(err) // TODO fix later
			}

			for _, pkg := range pkgs {

				for _, file := range pkg.Files {

					ast.Inspect(file, func(node ast.Node) bool {

						typeSpec, ok := node.(*ast.TypeSpec)
						if !ok {
							return true
						}

						if typeSpec.Name.String() != sel {
							return true
						}

						myDefaultValue = getInternalType(typeSpec, myDefaultValue)

						return true
					})

				}

			}

		}
	default:

		selectors := GetExprForImport(astFieldType)

		for _, s := range selectors {
			importFromMap, exist := importMap[s]
			if exist {
				usedMap[s] = importFromMap
			}
		}
	}

	return GogenType{
		Type:         myType,
		DefaultValue: myDefaultValue,
		Imports:      usedMap,
	}
}

func getInternalType(typeSpec *ast.TypeSpec, myDefaultValue string) string {
	switch theType := typeSpec.Type.(type) {
	case *ast.StructType:
		theFields := ""
		for _, field := range theType.Fields.List {
			for _, name := range field.Names {
				theFields += fmt.Sprintf("%s: %s, ", name, GetDefaultValue(field.Type))
			}
		}

		myDefaultValue = fmt.Sprintf("%s{ %s }", myDefaultValue, theFields)
	case *ast.InterfaceType:
		myDefaultValue = "nil"
	case *ast.Ident:
		myDefaultValue = fmt.Sprintf("%s(%s)", myDefaultValue, GetDefaultValue(typeSpec.Type))
	case *ast.FuncType:
		defRetVal := ""

		if theType.Results.NumFields() > 0 {
			for i, retList := range theType.Results.List {
				v := GetDefaultValue(retList.Type)

				if i < len(theType.Results.List)-1 {
					defRetVal += v + ", "
				} else {
					defRetVal += v
				}
			}
		}

		myDefaultValue = fmt.Sprintf("%s{return %s}", GetTypeAsString(typeSpec.Type), defRetVal)
	}
	return myDefaultValue
}

type GogenField struct {
	Name     GogenName `json:"name"`
	DataType GogenType `json:"dataType"`
}

func NewGogenField(name string, gType GogenType) GogenField {
	return GogenField{Name: GogenName(name), DataType: gType}
}

type GogenStruct struct {
	Name   GogenName    `json:"name"`
	Fields []GogenField `json:"fields"`
}

func NewGogenStruct(gomodPath, path, structName string) (*GogenStruct, error) {

	gs := GogenStruct{
		Name: GogenName(structName),
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	found := false

	for _, pkg := range pkgs {

		for _, file := range pkg.Files {

			importMap := map[string]GogenImport{}

			ast.Inspect(file, func(node ast.Node) bool {

				genDecl, ok := node.(*ast.GenDecl)
				if ok && genDecl.Tok == token.IMPORT {

					for _, spec := range genDecl.Specs {

						importSpec, ok := spec.(*ast.ImportSpec)
						if !ok {
							continue
						}

						gi := NewGogenImport(importSpec)
						importMap[gi.Expression] = gi
					}

					return true
				}

				typeSpec, ok := node.(*ast.TypeSpec)
				if !ok {
					return true
				}

				if typeSpec.Name.String() != structName {
					return true
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					return true
				}

				found = true

				for _, field := range structType.Fields.List {

					gt := NewGogenType(gomodPath, field, importMap)

					if field.Names != nil {
						for _, name := range field.Names {
							gs.Fields = append(gs.Fields, NewGogenField(name.String(), gt))
						}
					} else {
						nameField := GetSel(field.Type)
						gs.Fields = append(gs.Fields, NewGogenField(nameField, gt))

					}

				}

				return true
			})

			if found {
				break
			}

		}

		if found {
			break
		}

	}

	return &gs, nil
}

func GetSel(expr ast.Expr) string {
	switch fieldType := expr.(type) {
	case *ast.SelectorExpr:
		return fieldType.Sel.String()
	case *ast.StarExpr:
		return GetSel(fieldType.X)
	}
	return ""
}

func GetExprForImport(expr ast.Expr) []string {

	switch fieldType := expr.(type) {
	case *ast.SelectorExpr:
		return []string{fieldType.X.(*ast.Ident).String()}

	case *ast.StarExpr:
		return GetExprForImport(fieldType.X)

	case *ast.MapType:
		str := make([]string, 0)
		key := GetExprForImport(fieldType.Key)
		if key != nil {
			str = append(str, key...)
		}
		value := GetExprForImport(fieldType.Value)
		if value != nil {
			str = append(str, value...)
		}

		return str

	case *ast.ArrayType:
		return GetExprForImport(fieldType.Elt)

	case *ast.ChanType:
		return GetExprForImport(fieldType.Value)

	case *ast.FuncType:
		str := make([]string, 0)

		if fieldType.Params.NumFields() > 0 {
			for _, x := range fieldType.Params.List {
				str = append(str, GetExprForImport(x.Type)...)
			}
		}

		if fieldType.Results.NumFields() > 0 {
			for _, x := range fieldType.Results.List {
				str = append(str, GetExprForImport(x.Type)...)
			}
		}

		return str

	}

	return nil

}

type GogenMethod struct {
	Name    GogenName
	Params  []GogenField
	Results []GogenField
}

type GogenInterface struct {
	Name    GogenName
	Methods []GogenMethod
}

type GogenUsecase struct {
	UsecaseName    string
	InportRequest  GogenStruct
	InportResponse GogenStruct
	Outport        GogenInterface
}
