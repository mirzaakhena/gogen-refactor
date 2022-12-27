package gogen2

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sync"
)

type GogenStructBuilder struct {
	GogenBuilder
}

func NewGogenStructBuilder(goModPath, path string) *GogenStructBuilder {

	return &GogenStructBuilder{
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

func (gsb *GogenStructBuilder) Build(structName string) (*GogenStruct, error) {

	err := gsb.handleGoMod()
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, string(gsb.path), nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	gs := NewGogenStruct(structName)

	for _, pkg := range pkgs {

		// now we try to find the typeSpecName == structName
		for _, file := range pkg.Files {

			// we prepare the import
			gsb.importMap = map[Expression]GogenImport{}

			hasUnknownIndent := false

			ast.Inspect(file, func(node ast.Node) bool {

				// if there is an error, just ignore everything
				if err != nil {
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

				logDebug("next type : %v ------------------------------------------", typeSpecName)

				if typeSpecName != structName {

					logDebug("simpan kedalam typeMap type %v ", typeSpecName)
					gsb.typeMap[FieldType(typeSpecName)] = typeSpec.Type

					if hasUnknownIndent {

						// disini harusnya kita cuma focus yg ident aj
						// baik yg same file maupun diff file
						// itupun kalo masih ada setelah kita selesai trace seluruh field pada struct target
						// kalo misal tidak ada ident lagi, maka harusnya ini pun tidak perlu di proses lagi
						// func ini dipanggil 2x dan ini adalah pemanggilan pertama
						gsb.handleUncompleteDefaultValue()
					}

					logDebug(".")
					return true
				}

				// -------------- we found the struct target --------------

				//gsb.foundTarget = true

				logDebug("target struct %v sudah ditemukan", structName)

				// focus to struct only
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					logDebug("bukan struct")
					err = fmt.Errorf("type %s is not struct", typeSpecName)
					return false
				}

				for _, field := range structType.Fields.List {

					for _, s := range gsb.extractAllExpression(field.Type) {
						importFromMap, exist := gsb.importMap[s]
						if exist {
							gsb.usedImport[s] = importFromMap
						}
					}

					if field.Names != nil {
						// fieldNameIdent is exist
						for _, fieldNameIdent := range field.Names {
							if IsExported(fieldNameIdent.String()) {
								logDebug("sudah punya nama    : %v", fieldNameIdent.String())
								gf := NewGogenField(FieldName(fieldNameIdent.String()), field.Type)
								gs.Fields = append(gs.Fields, gf)
								gsb.checkDefaultValue(gf)
							}
						}
					} else {
						// name does not exist, use Selector as Name
						fieldNameStr := GetSel(field.Type)
						logDebug("karena tidak punya nama. maka diberi nama: %v", fieldNameStr)
						gf := NewGogenField(FieldName(fieldNameStr), field.Type)
						gs.Fields = append(gs.Fields, gf)
						gsb.checkDefaultValue(gf)
					}
					logDebug("")

				}

				for _, v := range gsb.unknownTypes {
					_, ok := v.DataType.Expr.(*ast.Ident)
					if ok {
						hasUnknownIndent = true
						break
					}
				}

				return true
			})

			if err != nil {
				return nil, err
			}

		}

	}

	logDebug("masuk ke handler selector")

	// kita akan coba pergi ke file yang lain untuk mencari tahu Selector tertentu bertipe apa
	gsb.handleSelector(gs)

	// kita kembali memanggil func ini utk kedua kalinya
	// harusnya disini kita hanya menghandle field yg ada selector-nya saja
	gsb.handleUncompleteDefaultValue()

	if len(gsb.unknownTypes) > 0 {

		arrUnknownTypes := make([]string, 0)

		for k, s := range gsb.unknownTypes {
			arrUnknownTypes = append(arrUnknownTypes, fmt.Sprintf("%s %s,", k, s.DataType.Type))
		}

		return nil, fmt.Errorf("unknown type for field : %v", arrUnknownTypes)
	}

	return gs, nil
}

func (gsb *GogenStructBuilder) handleSelector(gs *GogenStruct) {

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

						sel := Selector(typeSpec.Name.String())
						_, exist := expressionMap[x][sel]
						if !exist {
							return true
						}

						logDebug("   ketemu %v == %v. ===================> status expressionMap[x] = %v", typeSpec.Name.String(), sel, expressionMap[x])

						selector := fmt.Sprintf("%v.%v", x, typeSpec.Name.String())
						gsb.typeMap[FieldType(selector)] = typeSpec.Type

						delete(expressionMap[x], sel)

						logDebug("expressionMap menghapus %v saat ini len dari expressionMap : %d", sel, len(expressionMap[x]))

						if len(expressionMap[x]) == 0 {
							logDebug("found = true")
							found = true
							return false
						}

						//for sel, _ := range expressionMap[x] {
						//
						//	logDebug("   utk %v mengecek %v == %v", x, typeSpec.Name.String(), sel)
						//
						//	if Selector(fmt.Sprintf("%s", typeSpec.Name.String())) == sel {
						//
						//	}
						//}

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

func (gsb *GogenStructBuilder) handleUncompleteDefaultValue() {

	for k, v := range gsb.unknownTypes {
		ts, exist := gsb.typeMap[v.DataType.Type]
		if !exist {
			logDebug("dataType %v belum ditemukan dalam typeMap. mungkin di loop berikutnya", v.DataType.Type)
			continue
		}

		logDebug("tipe data %v untuk field %v sudah ready di map", v.DataType.Type, k)

		oldDefaultValue := v.DataType.DefaultValue
		v.handleDefaultValue(ts)
		logDebug("skg defaultValue yang tadinya %v, sudah direplace dengan %v", oldDefaultValue, v.DataType.DefaultValue)

		logDebug("")

		delete(gsb.unknownTypes, k)
	}

	logDebug("melihat list typeMap:")
	for k, _ := range gsb.typeMap {
		logDebug("%30s", k)
	}

	logDebug("")

	logDebug("status unknownTypes : %+v", gsb.unknownTypes)
}
