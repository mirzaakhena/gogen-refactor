package gogen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type GogenName string

type PathExpression string

type GogenImport struct {
	Name       string         `json:"name,omitempty"`
	Path       string         `json:"path"`
	Expression PathExpression `json:"expression"`
}

type GogenMethod struct {
	Name    GogenName
	Params  []GogenField
	Results []GogenField
}

func NewGogenMethod() GogenMethod {

	return GogenMethod{}
}

type GogenUsecase struct {
	UsecaseName    string
	InportRequest  GogenStruct
	InportResponse GogenStruct
	Outport        GogenInterface
}

type GogenType struct {
	Type         string `json:"type"`
	DefaultValue string `json:"defaultValue"`
	//Imports      map[PathExpression]GogenImport `json:"imports,omitempty"`
}

type GogenField struct {
	Name     GogenName `json:"name"`
	DataType GogenType `json:"dataType"`
}

type GogenStruct struct {
	Name    GogenName                      `json:"name"`
	Imports map[PathExpression]GogenImport `json:"imports,omitempty"`
	Fields  []GogenField                   `json:"fields"`
}

type GogenInterface struct {
	Name    GogenName
	Imports map[PathExpression]GogenImport `json:"imports,omitempty"`
	Methods []GogenMethod
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
		Expression: PathExpression(expr),
	}
}

func NewGogenType(gomodPath string, astField *ast.Field, importMap map[PathExpression]GogenImport) GogenType {

	// prepare the used import
	usedMap := map[PathExpression]GogenImport{}

	// get type
	myType := GetTypeAsString(astField.Type)

	// get default value
	myDefaultValue := GetDefaultValue(astField.Type)

	switch astFieldType := astField.Type.(type) {

	// it is has selector and must be from external package
	// we want to capture the detail data type
	case *ast.SelectorExpr:

		// the package Expression
		x := astFieldType.X.(*ast.Ident).String()

		// the Selector
		sel := astFieldType.Sel.String()

		// find it from importMap
		gi, exist := importMap[PathExpression(x)]
		if !exist {
			panic(fmt.Sprintf("%v not found in importMap", x))
		}

		// record the import
		usedMap[gi.Expression] = gi

		// only work for path that start with gomod
		if strings.HasPrefix(gi.Path, gomodPath) {

			// take the path part only
			pathWithoutGomod := gi.Path[len(gomodPath)+1:]

			// go to the file
			fset := token.NewFileSet()
			pkgs, err := parser.ParseDir(fset, pathWithoutGomod, nil, parser.ParseComments)
			if err != nil {
				panic(err) // TODO fix later
			}

			for _, pkg := range pkgs {

				for _, file := range pkg.Files {

					ast.Inspect(file, func(node ast.Node) bool {

						// focus only to type
						typeSpec, ok := node.(*ast.TypeSpec)
						if !ok {
							return true
						}

						// with specific name
						if typeSpec.Name.String() != sel {
							return true
						}

						// completing the default value
						myDefaultValue = GetDeepDefaultValue(typeSpec.Type, myDefaultValue)

						return true
					})

				}

			}

		}

	// it does not have selector
	default:

		// take the import path for detail data type
		exprs := getExprForImport(astFieldType)

		for _, s := range exprs {
			importFromMap, exist := importMap[s]
			if exist {
				usedMap[s] = importFromMap
			}
		}
	}

	return GogenType{
		Type:         myType,
		DefaultValue: myDefaultValue,
		//Imports:      usedMap,
	}
}

func NewGogenField(name string, gType GogenType) GogenField {
	return GogenField{Name: GogenName(name), DataType: gType}
}

func NewGogenStruct(gomodPath, path, structName string) (*GogenStruct, error) {

	gs := GogenStruct{
		Name: GogenName(structName),
	}

	// read file
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	found := false

	importMap := map[PathExpression]GogenImport{}

	// typeName the_type (struct|interface)
	//typeMap := map[string]string{}

	for _, pkg := range pkgs {

		for _, file := range pkg.Files {

			ast.Inspect(file, func(node ast.Node) bool {

				// read all import
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

				// focus to type
				typeSpec, ok := node.(*ast.TypeSpec)
				if !ok {
					return true
				}

				// focus to struct only
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					return true
				}

				// focus to specific name only
				if typeSpec.Name.String() != structName {
					return true
				}

				found = true

				// iterate the field
				for _, field := range structType.Fields.List {

					// read the type
					gt := NewGogenType(gomodPath, field, importMap)

					// if names exist, iterate it
					if field.Names != nil {
						for _, name := range field.Names {
							gs.Fields = append(gs.Fields, NewGogenField(name.String(), gt))
						}
					} else {
						nameField := getSel(field.Type)
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

func NewGogenInterface(gomodPath, path, interfaceName string) (*GogenInterface, error) {

	gi := GogenInterface{
		Name:    GogenName(interfaceName),
		Methods: make([]GogenMethod, 0),
	}

	// read file
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	found := false

	for _, pkg := range pkgs {

		for _, file := range pkg.Files {

			importMap := map[PathExpression]GogenImport{}

			ast.Inspect(file, func(node ast.Node) bool {

				// read all import
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

				// focus to type
				typeSpec, ok := node.(*ast.TypeSpec)
				if !ok {
					return true
				}

				// focus to specific name only
				if typeSpec.Name.String() != interfaceName {
					return true
				}

				// focus to interface only
				interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
				if !ok {
					return true
				}

				found = true

				ast.Print(fset, interfaceType)

				for _, method := range interfaceType.Methods.List {
					if len(method.Names) == 0 {
						panic("no method name found")
					}

					gm := GogenMethod{
						Name:    GogenName(method.Names[0].String()),
						Params:  nil,
						Results: nil,
					}

					gi.Methods = append(gi.Methods, gm)
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

	return &gi, nil
}
