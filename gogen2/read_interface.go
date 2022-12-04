package gogen2

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type GogenInterfaceBuilder struct {
	goModPath   string
	path        string
	importMap   map[string]GogenImport
	typeMap     map[string]*ast.TypeSpec
	foundTarget bool
}

func NewGogenInterfaceBuilder(goModPath, path string) *GogenInterfaceBuilder {

	return &GogenInterfaceBuilder{
		goModPath:   goModPath,
		path:        path,
		importMap:   map[string]GogenImport{},
		typeMap:     map[string]*ast.TypeSpec{},
		foundTarget: false,
	}
}

func (gsi *GogenInterfaceBuilder) Build(interfaceName string) (*GogenInterface, error) {

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, gsi.path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	gf := GogenInterface{
		Name:    interfaceName,
		Imports: make([]GogenImport, 0),
		Methods: make([]GogenMethod, 0),
	}

	for _, pkg := range pkgs {

		// now we try to find the typeSpecName == structName
		for _, file := range pkg.Files {

			gsi.importMap = map[string]GogenImport{}

			ast.Inspect(file, func(node ast.Node) bool {

				genDecl, ok := node.(*ast.GenDecl)
				if ok && genDecl.Tok == token.IMPORT {
					handleImport(genDecl, gsi.importMap)
					return true
				}

				// focus to type
				typeSpec, ok := node.(*ast.TypeSpec)
				if !ok {
					return true
				}

				// get type name
				typeSpecName := typeSpec.Name.String()

				if typeSpecName != interfaceName {
					gsi.typeMap[typeSpecName] = typeSpec
					//gsi.handleUncompleteDefaultValue() // TODO later
					return true
				}

				// -------------- we found the struct target --------------

				gsi.foundTarget = true

				// focus to interface only
				interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
				if !ok {
					err = fmt.Errorf("type %s is not interface", typeSpecName)
					return false
				}

				for _, method := range interfaceType.Methods.List {

					// if names exist, iterate it
					if method.Names != nil {
						for _, name := range method.Names {
							if IsExported(name.String()) {
								fmt.Printf("name : %v\n", name)
							}
						}
					} else {
						nameField := GetSel(method.Type)
						fmt.Printf("name : %v\n", nameField)
					}

					switch theFunc := method.Type.(type) {

					case *ast.FuncType:
						//fmt.Printf("masuk sini as FuncType\n")

						for _, param := range theFunc.Params.List {
							_ = param
						}

						for _, result := range theFunc.Results.List {
							_ = result
						}

					case *ast.SelectorExpr:
						//fmt.Printf("masuk sini as Selector\n")

					case *ast.Ident:
						//fmt.Printf("masuk sini as Ident\n")
					}

				}

				return true
			})

		}

	}

	return &gf, nil
}
