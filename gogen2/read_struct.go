package gogen2

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"strings"
)

type GogenStructBuilder struct {
	goModPath    string
	path         string
	importMap    map[string]GogenImport
	usedImport   map[string]GogenImport
	typeMap      map[string]*ast.TypeSpec
	unknownTypes map[string]*GogenType
	foundTarget  bool
}

func NewGogenStructBuilder(goModPath, path string) *GogenStructBuilder {

	return &GogenStructBuilder{
		goModPath:    goModPath,
		path:         path,
		importMap:    map[string]GogenImport{},
		usedImport:   map[string]GogenImport{},
		typeMap:      map[string]*ast.TypeSpec{},
		unknownTypes: map[string]*GogenType{},
		foundTarget:  false,
	}
}

func (gsb *GogenStructBuilder) Build(structName string) (*GogenStruct, error) {

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, gsb.path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	gs := NewGogenStruct(structName)

	for _, pkg := range pkgs {

		// now we try to find the typeSpecName == structName
		for _, file := range pkg.Files {

			gsb.importMap = map[string]GogenImport{}

			ast.Inspect(file, func(node ast.Node) bool {

				genDecl, ok := node.(*ast.GenDecl)
				if ok && genDecl.Tok == token.IMPORT {
					handleImport(genDecl, gsb.importMap)
					return true
				}

				// focus to type
				typeSpec, ok := node.(*ast.TypeSpec)
				if !ok {
					return true
				}

				// get type name
				typeSpecName := typeSpec.Name.String()

				if typeSpecName != structName {
					gsb.typeMap[typeSpecName] = typeSpec
					gsb.handleUncompleteDefaultValue()
					return true
				}

				// -------------- we found the struct target --------------

				gsb.foundTarget = true

				// focus to struct only
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					err = fmt.Errorf("type %s is not struct", typeSpecName)
					return false
				}

				ast.Print(fset, structType)

				for _, field := range structType.Fields.List {

					theType := GetTypeAsString(field.Type)
					defaultValue := GetDefaultValue(field.Type)

					gt := NewGogenType(theType, defaultValue)

					// if names exist, iterate it
					if field.Names != nil {
						for _, name := range field.Names {
							if IsExported(name.String()) {
								gs.AddField(name.String(), gt)
							}
						}
					} else {
						nameField := GetSel(field.Type)
						gs.AddField(nameField, gt)
					}

					switch fieldType := field.Type.(type) {

					case *ast.ArrayType:
						HandleArray(gt, fieldType, defaultValue)

					case *ast.Ident:
						HandleIdent(defaultValue, theType, gt, gsb.typeMap, gsb.unknownTypes)

					case *ast.SelectorExpr:

						err = HandleSelector(gsb.goModPath, fieldType, gt, defaultValue, gsb.importMap, gsb.usedImport)
						if err != nil {
							return false
						}

					default:
						// take the import path for detail data type
						for _, s := range GetExprForImport(fieldType) {
							importFromMap, exist := gsb.importMap[s]
							if exist {
								gsb.usedImport[s] = importFromMap
							}
						}

					}

				}

				return true
			})

			if err != nil {
				return nil, err
			}

		}

	}

	if len(gsb.unknownTypes) > 0 {

		arrUnknownTypes := make([]string, 0)

		for s := range gsb.unknownTypes {
			arrUnknownTypes = append(arrUnknownTypes, s)
		}

		return nil, fmt.Errorf("unknown struct : %v", arrUnknownTypes)
	}

	for _, v := range gsb.usedImport {
		gs.Imports = append(gs.Imports, v)
	}

	return gs, nil
}

func HandleSelector(goModPath string, fieldType *ast.SelectorExpr, gt *GogenType, defaultValue string, importMap, usedImport map[string]GogenImport) error {

	x, sel := GetXAndSelAsString(fieldType)

	// find it from importMap
	gi, exist := importMap[x]
	if !exist {
		return fmt.Errorf("%v not found in importMap", x)
	}

	// record the import
	usedImport[gi.Expression] = gi

	path := GetPathBasedOnImport(goModPath, gi, x)

	// go to the file
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	found := false

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
				gt.TypeDefaultValue, gt.JSONDefaultValue = GetDeepDefaultValue(typeSpec.Type, defaultValue)

				found = true

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
	return nil
}

func GetPathBasedOnImport(goModPath string, gi GogenImport, x string) string {
	if strings.HasPrefix(gi.Path, goModPath) {
		return gi.Path[len(goModPath)+1:]
	}
	return fmt.Sprintf("%s/src/%s", build.Default.GOROOT, x)
}

func HandleArray(gt *GogenType, fieldType ast.Expr, defaultValue string) {
	gt.TypeDefaultValue, gt.JSONDefaultValue = GetDeepDefaultValue(fieldType, defaultValue)
}

func HandleIdent(defaultValue string, theType string, gt *GogenType, typeMap map[string]*ast.TypeSpec, unknownTypes map[string]*GogenType) {

	if defaultValue == theType {

		typeSpecFromMap, exist := typeMap[theType]
		if !exist {
			unknownTypes[theType] = gt

		} else {

			gt.TypeDefaultValue, gt.JSONDefaultValue = GetDeepDefaultValue(typeSpecFromMap.Type, defaultValue)
		}

	}
}

func (gsb *GogenStructBuilder) handleUncompleteDefaultValue() {

	if gsb.foundTarget {
		for k, v := range gsb.unknownTypes {
			ts, exist := gsb.typeMap[k]
			if exist {

				v.TypeDefaultValue, v.JSONDefaultValue = GetDeepDefaultValue(ts.Type, v.TypeDefaultValue)
				delete(gsb.unknownTypes, k)
				break
			}
		}
	}
}

//func (gsb *GogenStructBuilder) handleImport(genDecl *ast.GenDecl) {
//
//	for _, spec := range genDecl.Specs {
//		importSpec, ok := spec.(*ast.ImportSpec)
//		if !ok {
//			continue
//		}
//		gi := NewGogenImport(importSpec)
//		gsb.importMap[gi.Expression] = gi
//	}
//
//}
