package gogen2

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"os"
	"strings"
	"sync"
)

type GogenStructBuilder struct {
	goModPath     string
	path          string
	importMap     map[Expression]GogenImport // menampung semua import tapi tidak semuanya dipakai oleh struct yg kita cari
	usedImport    map[Expression]GogenImport // menampung import yg dipakai saja
	typeMap       map[FieldType]ast.Expr     // menampung semua type yang mungkin akan dipakai oleh struct target
	unknownTypes  map[FieldName]*GogenField  // menampung hanya type yang diperlukan oleh field struct target
	expressionMap map[Expression][]string    // grouping FieldType
	mapOfRequire  map[RequirePath]CompletePath

	//foundTarget   bool
}

func NewGogenStructBuilder(goModPath, path string) *GogenStructBuilder {

	return &GogenStructBuilder{
		goModPath:     goModPath,
		path:          path,
		importMap:     map[Expression]GogenImport{},
		usedImport:    map[Expression]GogenImport{},
		typeMap:       map[FieldType]ast.Expr{},
		unknownTypes:  map[FieldName]*GogenField{},
		expressionMap: map[Expression][]string{},
		mapOfRequire:  map[RequirePath]CompletePath{},
		//foundTarget:   false,
	}
}

func (gsb *GogenStructBuilder) Build(structName string) (*GogenStruct, error) {

	err := gsb.handleGoMod()
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, gsb.path, nil, parser.ParseComments)
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

				// if there is an error, just ignore
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

				logDebug("next type : %v ------------------------------------------\n", typeSpecName)

				if typeSpecName != structName {

					logDebug("simpan kedalam typeMap type %v \n", typeSpecName)
					gsb.typeMap[FieldType(typeSpecName)] = typeSpec.Type

					if hasUnknownIndent {

						// disini harusnya kita cuma focus yg ident aj
						// baik yg same file maupun diff file
						// itupun kalo masih ada setelah kita selesai trace seluruh field pada struct target
						// kalo misal tidak ada ident lagi, maka harusnya ini pun tidak perlu di proses lagi
						// func ini dipanggil 2x dan ini adalah pemanggilan pertama
						gsb.handleUncompleteDefaultValue()
					}

					logDebug(".\n")
					return true
				}

				// -------------- we found the struct target --------------

				//gsb.foundTarget = true

				logDebug("target struct %v sudah ditemukan\n", structName)

				// focus to struct only
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					logDebug("bukan struct\n")
					err = fmt.Errorf("type %s is not struct", typeSpecName)
					return false
				}

				for _, field := range structType.Fields.List {

					//dataTypeStr := getTypeAsString(field.Type)
					//logDebug("tipe data           : %v\n", dataTypeStr)

					//logDebug("first time handleDefaultValue\n")
					//defaultValue := handleDefaultValue(dataTypeStr, field.Type)
					//logDebug("default value       : %v\n", defaultValue)

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
								gf := NewGogenField(FieldName(fieldNameIdent.String()), field.Type)
								gs.Fields = append(gs.Fields, gf)
								gsb.checkDefaultValue(gf)
							}
						}
					} else {
						// name does not exist, use Selector as Name
						fieldNameStr := GetSel(field.Type)
						logDebug("karena tidak punya nama. maka diberi nama: %v\n", fieldNameStr)
						gf := NewGogenField(FieldName(fieldNameStr), field.Type)
						gs.Fields = append(gs.Fields, gf)
						gsb.checkDefaultValue(gf)
					}
					logDebug("\n")

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

	logDebug("masuk ke handler selector\n")

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

func (gsb *GogenStructBuilder) handleGoMod() error {

	const goModFileName = "go.mod"

	dataInBytes, err := os.ReadFile(goModFileName)
	if err != nil {
		return err
	}

	parsedGoMod, err := modfile.Parse(goModFileName, dataInBytes, nil)
	if err != nil {
		return err
	}

	for _, r := range parsedGoMod.Require {

		if len(r.Syntax.Token) == 1 {
			gsb.mapOfRequire[RequirePath(r.Syntax.Token[0])] = CompletePath(fmt.Sprintf("%v/pkg/mod/%v", build.Default.GOPATH, r.Syntax.Token[0]))
			continue
		}

		gsb.mapOfRequire[RequirePath(r.Syntax.Token[0])] = CompletePath(fmt.Sprintf("%v/pkg/mod/%v@%v", build.Default.GOPATH, r.Syntax.Token[0], r.Syntax.Token[1]))
	}

	return nil
}

func (gsb *GogenStructBuilder) handleSelector(gs *GogenStruct) {

	wg := sync.WaitGroup{}

	logDebug("melihat list usedImport:\n")
	for k, _ := range gsb.usedImport {
		logDebug("%10s\n", k)
	}
	logDebug("\n")

	logDebug("melihat list unknownTypes:\n")
	for k, _ := range gsb.unknownTypes {
		logDebug("%10s\n", k)
	}
	logDebug("\n")

	logDebug("melihat list expressionMap:\n")
	for k, v := range gsb.expressionMap {
		logDebug("%10s %v\n", k, v)
	}
	logDebug("\n")

	// copy expressionMap
	expressionMap := map[Expression]map[string]int{}
	for theX, sels := range gsb.expressionMap {
		for _, sel := range sels {
			if expressionMap[theX] == nil {
				expressionMap[theX] = map[string]int{}
			}
			expressionMap[theX][sel] = 1
		}
	}

	// kenapa gak pakai unknownTypes aj?
	for x, ui := range gsb.usedImport {

		gs.Imports = append(gs.Imports, ui)

		path := ui.CompletePath

		logDebug("call path %v %v\n", path, gsb.importMap[x].Path)

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

					logDebug("untuk %v masuk ke file : %v\n", x, fset.File(file.Package).Name())

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

							logDebug("   utk %v mengecek %v == %v\n", x, typeSpec.Name.String(), sel)

							if fmt.Sprintf("%s", typeSpec.Name.String()) == sel {

								logDebug("   ketemu %v == %v. ===================> status expressionMap[x] = %v\n", typeSpec.Name.String(), sel, expressionMap[x])

								selector := fmt.Sprintf("%v.%v", x, typeSpec.Name.String())
								gsb.typeMap[FieldType(selector)] = typeSpec.Type

								delete(expressionMap[x], sel)

								logDebug("expressionMap menghapus %v saat ini len dari expressionMap : %d\n", sel, len(expressionMap[x]))

								if len(expressionMap[x]) == 0 {
									logDebug("found = true\n")
									found = true
									break
								}

							}
						}

						return true
					})

					if found {
						logDebug("break file\n")
						break
					}

				}

				if found {
					logDebug("break pkg\n")
					break
				}

			}

			logDebug("done for %v\n", x)
			wg.Done()

		}(x, string(path))

	}

	wg.Wait()
}

func (gsb *GogenStructBuilder) checkDefaultValue(gf *GogenField) {

	if string(gf.DataType.Type) == gf.DataType.DefaultValue {

		logDebug("karena defaultValue utk field %v dengan type %v belum final, kita cek ke map\n", gf.Name, gf.DataType.Type)

		typeSpecFromMap, exist := gsb.typeMap[gf.DataType.Type]
		if !exist {
			logDebug("dataType %v belum ditemukan dalam map. pencarian default value utk var %v ditunda dan sudah didaftarkan dalam unknownTypes\n", gf.DataType.DefaultValue, gf.Name)
			gsb.unknownTypes[gf.Name] = gf

			logDebug("status unknownTypes : %+v\n", gsb.unknownTypes)
			return
		}

		oldDefaultValue := gf.DataType.DefaultValue
		gf.handleDefaultValue(typeSpecFromMap)

		logDebug("dataType %v ada di map. defaultValue %v sudah di replace dengan %v\n", gf.DataType.DefaultValue, oldDefaultValue, gf.DataType.DefaultValue)

		return
	}

	logDebug("default value utk field %v dengan dataType %v sudah final, yaitu : %v\n", gf.Name, gf.DataType.DefaultValue, gf.DataType.DefaultValue)

}

func (gsb *GogenStructBuilder) handleUncompleteDefaultValue() {

	//removeUnknownTypes := make([]string, 0)

	for k, v := range gsb.unknownTypes {
		ts, exist := gsb.typeMap[v.DataType.Type]
		if !exist {
			logDebug("dataType %v belum ditemukan dalam typeMap. mungkin di loop berikutnya\n", v.DataType.Type)
			continue
		}

		logDebug("tipe data %v untuk field %v sudah ready di map\n", v.DataType.Type, k)

		oldDefaultValue := v.DataType.DefaultValue
		v.handleDefaultValue(ts)
		logDebug("skg defaultValue yang tadinya %v, sudah direplace dengan %v\n", oldDefaultValue, v.DataType.DefaultValue)

		//removeUnknownTypes = append(removeUnknownTypes, k)
		logDebug("\n")

		delete(gsb.unknownTypes, k)
	}

	logDebug("melihat list typeMap:\n")
	for k, _ := range gsb.typeMap {
		logDebug("%30s\n", k)
	}

	logDebug("\n")

	//for _, ut := range removeUnknownTypes {
	//	logDebug("menghapus %v dari unknown type map\n", ut)
	//	delete(gsb.unknownTypes, ut)
	//}

	logDebug("status unknownTypes : %+v\n", gsb.unknownTypes)
}

func (gsb *GogenStructBuilder) handleUsedImport(expr ast.Expr) []Expression {

	switch fieldType := expr.(type) {
	case *ast.StructType:

		str := make([]Expression, 0)
		for _, f := range fieldType.Fields.List {
			str = append(str, gsb.handleUsedImport(f.Type)...)
		}
		return str

	case *ast.SelectorExpr:
		x := Expression(fieldType.X.(*ast.Ident).String())
		sel := fieldType.Sel.String()

		gsb.expressionMap[x] = append(gsb.expressionMap[x], sel)

		return []Expression{Expression(fieldType.X.(*ast.Ident).String())}

	case *ast.StarExpr:
		return gsb.handleUsedImport(fieldType.X)

	case *ast.MapType:
		str := make([]Expression, 0)
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
		str := make([]Expression, 0)

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

	return []Expression{}
}

func (gsb *GogenStructBuilder) handleImport(genDecl *ast.GenDecl) {

	for _, spec := range genDecl.Specs {

		importSpec, ok := spec.(*ast.ImportSpec)
		if !ok {
			continue
		}

		// kita ambil import dari file yg sedang dibaca
		importPath := strings.Trim(importSpec.Path.Value, `"`)

		lenImportPath := len(importPath)

		// kita cek apakah dia ada ada di gomod, 3rdlib, atau internal apps
		cp, exist := gsb.mapOfRequire[RequirePath(importPath)]
		if exist {

			pathToLib := fmt.Sprintf("%v%v", cp, importPath[lenImportPath:])

			fmt.Printf("####### %v\n", pathToLib)

			pkgs, err := parser.ParseDir(token.NewFileSet(), pathToLib, nil, parser.PackageClauseOnly)
			if err != nil {
				panic(err)
			}

			for _, pkg := range pkgs {

				name := ""
				expr := pkg.Name
				if importSpec.Name != nil {
					name = importSpec.Name.String()
					expr = name
				}

				gi := GogenImport{
					Name:         name,
					Path:         ImportPath(importPath),
					Expression:   Expression(expr),
					ImportType:   ImportTypeExtModule,
					CompletePath: CompletePath(pathToLib),
				}

				gsb.importMap[gi.Expression] = gi
			}

			continue
		}

		name := ""
		expr := importPath[strings.LastIndex(importPath, "/")+1:]
		if importSpec.Name != nil {
			name = importSpec.Name.String()
			expr = name
		}

		importType := ImportTypeSDK
		completePath := CompletePath(fmt.Sprintf("%s/src/%s", build.Default.GOROOT, expr))
		if strings.HasPrefix(importPath, gsb.goModPath) {
			importType = ImportTypeProject
			completePath = CompletePath(importPath[len(gsb.goModPath)+1:])
		}

		gi := GogenImport{
			Name:         name,
			Path:         ImportPath(importPath),
			CompletePath: completePath,
			Expression:   Expression(expr),
			ImportType:   importType,
		}

		gsb.importMap[gi.Expression] = gi

	}

}

func (gsb *GogenStructBuilder) getPathBasedOnImport(gi GogenImport, x Expression) string {
	if strings.HasPrefix(string(gi.Path), gsb.goModPath) {
		return string(gi.Path[len(gsb.goModPath)+1:])
	}

	// TODO change between build.Default.GOROOT or build.Default.GOPATH

	// sample build.Default.GOROOT : /usr/local/go
	// sample GOROOT full path     : /usr/local/go/src/context/context.go

	// sample build.Default.GOPATH : /Users/mirza/go
	// sample GOPATH full path     : /Users/mirza/go/pkg/mod/github.com/gin-gonic/gin@v1.8.1/auth.go

	return fmt.Sprintf("%s/src/%s", build.Default.GOROOT, x)
}
