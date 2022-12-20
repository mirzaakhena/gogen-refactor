package gogen2

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sync"
)

type GogenInterfaceBuilder struct {
	GogenBuilder
	unknownInterface map[FieldType]ast.Expr
}

func NewGogenInterfaceBuilder(goModPath, path string) *GogenInterfaceBuilder {
	return &GogenInterfaceBuilder{
		unknownInterface: map[FieldType]ast.Expr{},
		GogenBuilder: GogenBuilder{
			path:          path,
			goModPath:     goModPath,
			importMap:     map[Expression]GogenImport{},
			usedImport:    map[Expression]GogenImport{},
			mapOfRequire:  map[RequirePath]CompletePath{},
			expressionMap: map[Expression][]string{},
			typeMap:       map[FieldType]ast.Expr{},
			unknownTypes:  map[FieldName]*GogenField{},
		},
	}
}

func (gsb *GogenInterfaceBuilder) Build(interfaceName string) (*GogenInterface, error) {

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, gsb.path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	gc := NewGogenInterface(interfaceName)

	hasUnknownInterface := false

	for _, pkg := range pkgs {

		// now we try to find the typeSpecName == structName
		for _, file := range pkg.Files {

			logDebug("file %s", fset.File(file.Package).Name())

			gsb.importMap = map[Expression]GogenImport{}

			hasUnknownIndent := false

			ast.Inspect(file, func(node ast.Node) bool {

				// if there is an error, just ignore everything
				if err != nil {
					logDebug("dipaksa keluar %v", err.Error())
					return false
				}

				genDecl, ok := node.(*ast.GenDecl)
				if ok && genDecl.Tok == token.IMPORT {
					gsb.handleImport(genDecl)
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

					logDebug("menemukan %s", typeSpecName)

					gsb.typeMap[FieldType(typeSpecName)] = typeSpec.Type

					if hasUnknownInterface {

						logDebug("mulai mencari hasUnknownInterface")

						_, ok := typeSpec.Type.(*ast.InterfaceType)
						if !ok {
							return false
						}

						for k, _ := range gsb.unknownInterface {

							ts, exist := gsb.typeMap[k]
							if !exist {
								logDebug("dataType %v belum ditemukan dalam typeMap. mungkin di loop berikutnya", k)
								continue
							}

							err := gsb.handleMethodInterface(k, ts, gc)
							if err != nil {
								return false
							}

							delete(gsb.unknownInterface, k)

						}

					}

					if hasUnknownIndent {
						logDebug(">>>>> masuk hasUnknownIndent")
						// gsb.handleUncompleteDefaultValue() // TODO solve later
					}

					return true
				}

				logDebug("yang dicari ketemu!!! ----> %s", typeSpecName)

				err = gsb.handleMethodInterface(FieldType(typeSpecName), typeSpec.Type, gc)
				if err != nil {
					return false
				}

				for _, v := range gsb.unknownTypes {
					_, ok := v.DataType.Expr.(*ast.Ident)
					if ok {
						hasUnknownIndent = true
						break
					}
				}

				logDebug(">>>>> len unknownInterface %d", len(gsb.unknownInterface))
				if len(gsb.unknownInterface) > 0 {
					hasUnknownInterface = true
				}

				return true
			})

		}

	}

	for fieldType, _ := range gsb.unknownInterface {
		logDebug("unknown %v", fieldType)
	}

	gsb.handleSelector(gc)

	return gc, nil
}

func (gsb *GogenInterfaceBuilder) handleMethodInterface(typeSpecName FieldType, expr ast.Expr, gc *GogenInterface) error {

	// kita hanya fokus ke interface saja
	interfaceType, ok := expr.(*ast.InterfaceType)
	if !ok {
		return fmt.Errorf("type %s is not interface", typeSpecName)
	}

	for _, method := range interfaceType.Methods.List {

		switch methodType := method.Type.(type) {
		case *ast.FuncType:

			// jika masuk kesini berarti ini adalah inline function

			if method.Names == nil && len(method.Names) > 0 {
				return fmt.Errorf("method must have name")
			}
			methodName := method.Names[0].String()
			if !IsExported(methodName) {
				continue
			}
			gm := NewGogenMethod(typeSpecName, MethodName(methodName))
			gc.Methods = append(gc.Methods, gm)
			gsb.handleFuncParamResultType(methodType, gm)

		case *ast.Ident:
			logDebug(">> ident %v", methodType.String())

			// tidak mungkin ada import disini
			// disini kita berharap ident yg ditemukan sudah pernah didaftarkan pada typeMap
			// dan ident disini sudah pasti adalah sebuah interface
			im, exist := gsb.typeMap[FieldType(methodType.String())]
			if !exist {
				// jika masuk disini, maka ident belum pernah ditemukan, tapi mgk nanti akan ketemu
				// mungkin ada di package yg sama dan file yg sama,
				// mungkin ada di package yg sama namun file yg berbeda
				// dan tidak mungkin ada di package yang berbeda
				// kita belum tahu methodType itu apa, masukin aj dulu nanti akan kita cek.
				gsb.unknownInterface[FieldType(methodType.String())] = methodType
				continue
			}

			// jika masuk kesini, maka ident ini sudah pasti sebuah interface yang sudah pernah ditemukan diawal
			// tapi belum pernah ditelusuri lebih lanjut
			// kita akan selesaikan dengan dirinya sendiri sebagai interfaceType (recursive)
			err := gsb.handleMethodInterface(FieldType(methodType.String()), im, gc)
			if err != nil {
				return err
			}

		case *ast.SelectorExpr:

			logDebug("masuk sebagai selector")

			// kalau masuk sini sudah pasti belum pernah ditemukan dalam interfaceType
			// sudah pasti ada import yang akan kita pakai disini
			for _, s := range gsb.handleUsedImport(methodType) {
				importFromMap, exist := gsb.importMap[s]
				if exist {
					gsb.usedImport[s] = importFromMap
				}
			}

			m := fmt.Sprintf("%v.%v", methodType.X.(*ast.Ident).String(), methodType.Sel.String())
			logDebug("selector %v", m)

			// sudah pasti ada di package yg berbeda yang akan kita telusuri nanti
			// method Type disini sudah pasti selector,
			// tapi kita belum tahu Selectornya type apa, masukin aj dulu nanti akan kita cek.
			gsb.unknownInterface[FieldType(m)] = methodType

		default:
			return fmt.Errorf("unsupported type %v", methodType)
		}

	}

	return nil
}

func (gsb *GogenInterfaceBuilder) handleFuncParamResultType(methodType *ast.FuncType, gm *GogenMethod) {

	if methodType.Params.NumFields() > 0 {
		for _, param := range methodType.Params.List {

			for _, s := range gsb.handleUsedImport(param.Type) {
				importFromMap, exist := gsb.importMap[s]
				if exist {
					gsb.usedImport[s] = importFromMap
				}
			}

			if param.Names != nil {

				for _, n := range param.Names {
					gf := NewGogenField(FieldName(n.String()), param.Type)
					gm.Params = append(gm.Params, gf)
					gsb.checkDefaultValue(gf)
				}
			} else {
				fieldNameStr := GetSel(param.Type)
				gf := NewGogenField(FieldName(fieldNameStr), param.Type)
				gm.Params = append(gm.Params, gf)
				gsb.checkDefaultValue(gf)

			}

		}
	}

	if methodType.Results.NumFields() > 0 {
		for _, result := range methodType.Results.List {

			for _, s := range gsb.handleUsedImport(result.Type) {
				importFromMap, exist := gsb.importMap[s]
				if exist {
					gsb.usedImport[s] = importFromMap
				}
			}

			if result.Names != nil {
				for _, n := range result.Names {
					gf := NewGogenField(FieldName(n.String()), result.Type)
					gm.Results = append(gm.Results, gf)
					gsb.checkDefaultValue(gf)
				}
			} else {
				fieldNameStr := GetSel(result.Type)
				gf := NewGogenField(FieldName(fieldNameStr), result.Type)
				gm.Results = append(gm.Results, gf)
				gsb.checkDefaultValue(gf)

			}

		}
	}
}

func (gsb *GogenInterfaceBuilder) handleSelector(gs *GogenInterface) {

	wg := sync.WaitGroup{}

	logDebug("melihat list usedImport:")
	for k, _ := range gsb.usedImport {
		logDebug("%10s", k)
	}
	logDebug("")

	logDebug("melihat list unknownTypes:")
	for k, _ := range gsb.unknownTypes {
		logDebug("%10s", k)
	}
	logDebug("")

	logDebug("melihat list expressionMap:")
	for k, v := range gsb.expressionMap {
		logDebug("%10s %v", k, v)
	}
	logDebug("")

	// copy expressionMap
	expressionMap := map[Expression]map[Selector]int{}
	for theX, sels := range gsb.expressionMap {
		for _, sel := range sels {
			if expressionMap[theX] == nil {
				expressionMap[theX] = map[Selector]int{}
			}
			expressionMap[theX][Selector(sel)] = 1
		}
	}

	// kenapa gak pakai unknownTypes aj?
	for x, ui := range gsb.usedImport {

		gs.Imports = append(gs.Imports, ui)

		path := ui.CompletePath

		logDebug("call path %v %v", path, gsb.importMap[x].Path)

		wg.Add(1)

		go func(x Expression, path string) {

			// go to the file
			fset := token.NewFileSet()
			pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
			if err != nil {
				panic(err)
			}

			found := false

			for _, pkg := range pkgs {

				for _, file := range pkg.Files {

					logDebug("untuk %v masuk ke file : %v", x, fset.File(file.Package).Name())

					ast.Inspect(file, func(node ast.Node) bool {

						if found {
							return false
						}

						// focus only to type
						typeSpec, ok := node.(*ast.TypeSpec)
						if !ok {
							return true
						}

						for sel, _ := range expressionMap[x] {

							logDebug("   utk %v mengecek %v == %v", x, typeSpec.Name.String(), sel)

							if Selector(fmt.Sprintf("%s", typeSpec.Name.String())) == sel {

								logDebug("   ketemu %v == %v. ===================> status expressionMap[x] = %v", typeSpec.Name.String(), sel, expressionMap[x])

								selector := fmt.Sprintf("%v.%v", x, typeSpec.Name.String())
								gsb.typeMap[FieldType(selector)] = typeSpec.Type

								delete(expressionMap[x], sel)

								logDebug("expressionMap menghapus %v saat ini len dari expressionMap : %d", sel, len(expressionMap[x]))

								if len(expressionMap[x]) == 0 {
									logDebug("found = true")
									found = true
									break
								}

							}
						}

						return true
					})

					if found {
						logDebug("break file")
						break
					}

				}

				if found {
					logDebug("break pkg")
					break
				}

			}

			logDebug("done for %v", x)
			wg.Done()

		}(x, string(path))

	}

	wg.Wait()
}
