package gogen2

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type GogenInterfaceBuilder struct {
	goModPath    string
	path         string
	importMap    map[string]GogenImport
	usedImport   map[string]GogenImport
	typeMap      map[string]*ast.TypeSpec
	unknownTypes map[string]*GogenType
	foundTarget  bool
}

func NewGogenInterfaceBuilder(goModPath, path string) *GogenInterfaceBuilder {

	return &GogenInterfaceBuilder{
		goModPath:   goModPath,
		path:        path,
		importMap:   map[string]GogenImport{},
		usedImport:  map[string]GogenImport{},
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
		Methods: make([]*GogenMethod, 0),
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

				// -------------- we found the interface target --------------

				gsi.foundTarget = true

				// focus to interface only
				interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
				if !ok {
					err = fmt.Errorf("type %s is not interface", typeSpecName)
					return false
				}

				for _, method := range interfaceType.Methods.List {

					methodName := ""

					// if names exist, iterate it
					if method.Names != nil {
						for _, name := range method.Names {
							if IsExported(name.String()) {
								methodName = name.String()
							}
						}
					} else {
						methodName = GetSel(method.Type)
					}

					gm := GogenMethod{
						Name:    methodName,
						Params:  make([]*GogenField, 0),
						Results: make([]*GogenField, 0),
						//TypeParams: make([]*GogenField, 0),
					}

					gf.Methods = append(gf.Methods, &gm)

					switch theFunc := method.Type.(type) {

					case *ast.FuncType:
						fmt.Printf("masuk sini as FuncType\n")

						if theFunc.Params != nil {
							for i, param := range theFunc.Params.List {

								theType := GetTypeAsString(param.Type)
								defaultValue := GetDefaultValue(param.Type)

								gt := NewGogenType(theType, defaultValue)

								if param.Names != nil {
									for _, name := range param.Names {
										gm.AddParam(GetParamName(-1, name), gt)
									}
								} else {
									gm.AddParam(GetParamName(i, param.Type), gt)

								}

								switch fieldType := param.Type.(type) {

								case *ast.ArrayType:
									HandleArray(gt, fieldType, defaultValue)

								case *ast.Ident:
									HandleIdent(defaultValue, theType, gt, gsi.typeMap, gsi.unknownTypes)

								case *ast.SelectorExpr:

									err = HandleSelector(gsi.goModPath, fieldType, gt, defaultValue, gsi.importMap, gsi.usedImport)
									if err != nil {
										return false
									}

								default:
									// take the import path for detail data type
									for _, s := range GetExprForImport(fieldType) {
										importFromMap, exist := gsi.importMap[s]
										if exist {
											gsi.usedImport[s] = importFromMap
										}
									}

								}

							}

						}

						if theFunc.Results != nil {
							for _, result := range theFunc.Results.List {

								theType := GetTypeAsString(result.Type)
								defaultValue := GetDefaultValue(result.Type)

								gt := NewGogenType(theType, defaultValue)

								if result.Names != nil {
									for _, name := range result.Names {
										gm.AddResult(GetResultName(name), gt)
									}
								} else {

									gm.AddResult("", gt)
								}

							}
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
