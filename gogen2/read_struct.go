package gogen2

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"strings"
	"sync"
)

type name struct {
}

type GogenStructBuilder struct {
	goModPath    string
	path         string
	importMap    map[string]GogenImport
	usedImport   map[string]GogenImport
	typeMap      map[string]*ast.TypeSpec
	unknownTypes map[string]*GogenField
	selectorMap  map[string][]string
	foundTarget  bool
}

func NewGogenStructBuilder(goModPath, path string) *GogenStructBuilder {

	return &GogenStructBuilder{
		goModPath:    goModPath,
		path:         path,
		importMap:    map[string]GogenImport{},
		usedImport:   map[string]GogenImport{},
		typeMap:      map[string]*ast.TypeSpec{},
		unknownTypes: map[string]*GogenField{},
		selectorMap:  map[string][]string{},
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

				logDebug("next type : %v ------------------------------------------\n", typeSpecName)

				if typeSpecName != structName {

					logDebug("simpan kedalam typeMap type %v \n", typeSpecName)
					gsb.typeMap[typeSpecName] = typeSpec

					if gsb.foundTarget {
						gsb.handleUncompleteDefaultValue()
					}

					logDebug(".\n")
					return true
				}

				// -------------- we found the struct target --------------

				gsb.foundTarget = true

				logDebug("target struct %v sudah ditemukan\n", structName)

				// focus to struct only
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					err = fmt.Errorf("type %s is not struct", typeSpecName)
					return false
				}

				for _, field := range structType.Fields.List {

					dataTypeStr := getTypeAsString(field.Type)
					logDebug("tipe data           : %v\n", dataTypeStr)

					logDebug("first time handleDefaultValue\n")
					defaultValue := handleDefaultValue(dataTypeStr, field.Type)
					logDebug("default value       : %v\n", defaultValue)

					for _, s := range gsb.handleUsedImport(field.Type) {
						importFromMap, exist := gsb.importMap[s]
						if exist {
							gsb.usedImport[s] = importFromMap
						}
					}

					if field.Names != nil {
						// fieldNameIdent is exist
						for _, fieldNameIdent := range field.Names {
							if IsExported(fieldNameIdent.String()) {
								logDebug("sudah punya nama    : %v\n", fieldNameIdent.String())
								gf := NewGogenField(fieldNameIdent.String(), dataTypeStr, defaultValue)
								gs.Fields = append(gs.Fields, gf)
								gsb.checkDefaultValue(gf)
							}
						}
					} else {
						// name does not exist, use Selector as Name
						fieldNameStr := GetSel(field.Type)
						logDebug("karena tidak punya nama. maka diberi nama: %v\n", fieldNameStr)
						gf := NewGogenField(fieldNameStr, dataTypeStr, defaultValue)
						gs.Fields = append(gs.Fields, gf)
						gsb.checkDefaultValue(gf)
					}
					logDebug("\n")

				}

				return true
			})

			if err != nil {
				return nil, err
			}

		}

	}

	gsb.handleSelector(gs)

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

	for x, ui := range gsb.usedImport {

		gs.Imports = append(gs.Imports, ui)

		path := gsb.getPathBasedOnImport(ui, x)

		fmt.Printf("call path %v\n", path)

		wg.Add(1)

		go func(x, path string) {

			// go to the file
			fset := token.NewFileSet()
			pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
			if err != nil {
				panic(err)
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

						for _, sel := range gsb.selectorMap[x] {

							if fmt.Sprintf("%s", typeSpec.Name.String()) == sel {

								logDebug("ketemu %v == %v\n", typeSpec.Name.String(), sel)

								gsb.typeMap[fmt.Sprintf("%v.%v", x, typeSpec.Name.String())] = typeSpec
							}
						}

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

			wg.Done()

		}(x, path)

	}

	wg.Wait()
}

func (gsb *GogenStructBuilder) checkDefaultValue(gf *GogenField) {

	if gf.DataType.Type == gf.DataType.DefaultValue {

		logDebug("karena defaultValue utk field %v dengan type %v belum final, kita cek ke map\n", gf.Name, gf.DataType.Type)

		typeSpecFromMap, exist := gsb.typeMap[gf.DataType.DefaultValue]
		if !exist {
			logDebug("dataType %v belum ditemukan dalam map. pencarian default value utk var %v ditunda dan sudah didaftarkan dalam unknownTypes\n", gf.DataType.DefaultValue, gf.Name)
			gsb.unknownTypes[gf.Name] = gf

			logDebug("status unknownTypes : %+v\n", gsb.unknownTypes)
			return
		}

		oldDefaultValue := gf.DataType.DefaultValue
		newDefaultValue := handleDefaultValue(gf.DataType.DefaultValue, typeSpecFromMap.Type)
		gf.SetNewDefaultValue(newDefaultValue)

		logDebug("dataType %v ada di map. defaultValue %v sudah di replace dengan %v\n", gf.DataType.DefaultValue, oldDefaultValue, newDefaultValue)

		return
	}

	logDebug("default value utk field %v dengan dataType %v sudah final, yaitu : %v\n", gf.Name, gf.DataType.DefaultValue, gf.DataType.DefaultValue)

}

func (gsb *GogenStructBuilder) handleUncompleteDefaultValue() {

	removeUnknownTypes := make([]string, 0)

	for k, v := range gsb.unknownTypes {
		ts, exist := gsb.typeMap[v.DataType.Type]
		if !exist {
			logDebug("dataType %v belum ditemukan dalam typeMap. mungkin di loop berikutnya\n", v.DataType.Type)
			continue
		}

		logDebug("tipe data %v untuk field %v sudah ready di map\n", v.DataType.Type, k)
		newDefaultValue := handleDefaultValue(v.DataType.Type, ts.Type)

		logDebug("skg defaultValue yang tadinya %v, sudah direplace dengan %v\n", v.DataType.DefaultValue, newDefaultValue)
		v.SetNewDefaultValue(newDefaultValue)

		removeUnknownTypes = append(removeUnknownTypes, k)
		logDebug("\n")
	}

	for _, ut := range removeUnknownTypes {
		logDebug("menghapus %v dari unknown type map\n", ut)
		delete(gsb.unknownTypes, ut)
	}

	logDebug("status unknownTypes : %+v\n", gsb.unknownTypes)
}

func (gsb *GogenStructBuilder) getPathBasedOnImport(gi GogenImport, x string) string {
	if strings.HasPrefix(gi.Path, gsb.goModPath) {
		return gi.Path[len(gsb.goModPath)+1:]
	}
	return fmt.Sprintf("%s/src/%s", build.Default.GOROOT, x)
}

func (gsb *GogenStructBuilder) handleUsedImport(expr ast.Expr) []string {

	switch fieldType := expr.(type) {
	case *ast.SelectorExpr:
		x := fieldType.X.(*ast.Ident).String()
		sel := fieldType.Sel.String()

		gsb.selectorMap[x] = append(gsb.selectorMap[x], sel)

		return []string{fieldType.X.(*ast.Ident).String()}

	case *ast.StarExpr:
		return gsb.handleUsedImport(fieldType.X)

	case *ast.MapType:
		str := make([]string, 0)
		key := gsb.handleUsedImport(fieldType.Key)
		if key != nil {
			str = append(str, key...)
		}
		value := gsb.handleUsedImport(fieldType.Value)
		if value != nil {
			str = append(str, value...)
		}
		return str

	case *ast.ArrayType:
		return gsb.handleUsedImport(fieldType.Elt)

	case *ast.ChanType:
		return gsb.handleUsedImport(fieldType.Value)

	case *ast.FuncType:
		str := make([]string, 0)

		if fieldType.Params.NumFields() > 0 {
			for _, x := range fieldType.Params.List {
				str = append(str, gsb.handleUsedImport(x.Type)...)
			}
		}

		if fieldType.Results.NumFields() > 0 {
			for _, x := range fieldType.Results.List {
				str = append(str, gsb.handleUsedImport(x.Type)...)
			}
		}

		return str

	}

	return []string{}
}
