package gogen3

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
	unknownTypes map[string]*GogenField
	foundTarget  bool

	selectorMap map[string][]string //X, Sel
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

					if field.Names != nil {
						// fieldNameIdent is exist
						for _, fieldNameIdent := range field.Names {
							if IsExported(fieldNameIdent.String()) {
								logDebug("sudah punya nama    : %v\n", fieldNameIdent.String())
								gsb.addField(fieldNameIdent.String(), field, gs)
							}
						}
					} else {
						// name does not exist, use Selector as Name
						fieldNameStr := GetSel(field.Type)
						logDebug("karena tidak punya nama. maka diberi nama: %v\n", fieldNameStr)
						gsb.addField(fieldNameStr, field, gs)
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

	for x, ui := range gsb.usedImport {
		gs.Imports = append(gs.Imports, ui)

		path := gsb.getPathBasedOnImport(ui, x)

		fmt.Printf("call path %v\n", path)

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

	}

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

func (gsb *GogenStructBuilder) addField(name string, field *ast.Field, gs *GogenStruct) {

	dataTypeStr := getTypeAsString(field.Type)
	logDebug("tipe data           : %v\n", dataTypeStr)

	logDebug("first time handleDefaultValue\n")
	defaultValue := handleDefaultValue(dataTypeStr, field.Type)

	logDebug("default value       : %v\n", defaultValue)

	gf := NewGogenField(name, dataTypeStr, defaultValue)

	gs.Fields = append(gs.Fields, gf)

	for _, s := range gsb.handleUsedImport(field.Type) {
		importFromMap, exist := gsb.importMap[s]
		if exist {
			gsb.usedImport[s] = importFromMap

		}

	}

	//switch fieldType := field.Type.(type) {
	//case *ast.SelectorExpr:
	//	x, sel := GetXAndSelAsString(fieldType)
	//
	//	// find it from importMap
	//	gi, exist := gsb.importMap[x]
	//	if !exist {
	//		panic(fmt.Sprintf("%v not found in importMap", x))
	//	}
	//
	//	gsb.usedImport[gi.Expression] = gi
	//
	//	path := gsb.getPathBasedOnImport(gi, x)
	//
	//}

	if gf.DataType.Type == gf.DataType.DefaultValue {

		logDebug("karena defaultValue utk field %v dengan type %v belum final, kita cek ke map\n", gf.Name, gf.DataType.Type)

		typeSpecFromMap, exist := gsb.typeMap[gf.DataType.Type]
		if !exist {
			logDebug("dataType %v belum ditemukan dalam map. pencarian default value utk var %v ditunda dan sudah didaftarkan dalam unknownTypes\n", gf.DataType.Type, gf.Name)
			gsb.unknownTypes[gf.Name] = gf

			logDebug("status unknownTypes : %+v\n", gsb.unknownTypes)
			return
		}

		oldDefaultValue := gf.DataType.DefaultValue
		newDefaultValue := handleDefaultValue(gf.DataType.Type, typeSpecFromMap.Type)
		gf.DataType.DefaultValue = newDefaultValue

		logDebug("dataType %v ada di map. defaultValue %v sudah di replace dengan %v\n", gf.DataType.Type, oldDefaultValue, newDefaultValue)

		return
	}

	logDebug("default value utk field %v dengan dataType %v sudah final, yaitu : %v\n", gf.Name, gf.DataType.Type, gf.DataType.DefaultValue)

}

func handleDefaultValue(typeAsString string, expr ast.Expr) string {

	logDebug("handleDefaultValue %v %v\n", typeAsString, expr)

	switch exprType := expr.(type) {
	case *ast.Ident:

		logDebug("as ident            : %v\n", exprType.String())

		// found type in the same file
		if exprType.Obj != nil {
			logDebug("dataType %s ada di file yg sama dengan struct target\n", exprType.String())

			typeSpec, ok := exprType.Obj.Decl.(*ast.TypeSpec)
			if !ok {
				panic("cannot assert to TypeSpec")
			}

			logDebug("start recursive handleDefaultValue utk dataType %v\n", exprType.String())
			defaultValue := handleDefaultValue(typeAsString, typeSpec.Type)
			logDebug("end   recursive handleDefaultValue dari type %v hasil recursive adalah %v\n", exprType.String(), defaultValue)

			return defaultValue
		}

		basicType := ""

		for {

			if strings.HasPrefix(exprType.String(), "int") || strings.HasPrefix(exprType.String(), "uint") {
				logDebug("as int / uint\n")
				basicType = "0"
				break
			}

			if strings.HasPrefix(exprType.String(), "float") {
				logDebug("as float\n")
				basicType = "0.0"
				break
			}

			if exprType.String() == "string" {
				logDebug("as string\n")
				basicType = `""`
				break
			}

			if exprType.String() == "bool" {
				logDebug("as bool\n")
				basicType = `false`
				break
			}

			if exprType.String() == "any" {
				logDebug("as any\n")
				basicType = `nil`
				break
			}

			break
		}

		if basicType != "" {
			if typeAsString != exprType.String() {
				return fmt.Sprintf("%s(%s)", typeAsString, basicType)
			}
			return basicType
		}

		logDebug("tipe data dasar masih belum diketahui\n")

		return exprType.String()

	case *ast.StructType:
		v := fmt.Sprintf("%v{}", typeAsString)
		logDebug("as struct %v\n", v)
		return v

	case *ast.ArrayType:
		v := fmt.Sprintf("%s{}", typeAsString)
		logDebug("as array %v\n", v)
		return v

	case *ast.SelectorExpr:

		ident, ok := exprType.X.(*ast.Ident)

		// hardcoded fo context
		if ok && ident.String() == "context" {
			return "ctx"
		}

		v := fmt.Sprintf("%s.%s", ident.String(), exprType.Sel.String())

		logDebug("as selector %v\n", v)

		return v

	case *ast.StarExpr:
		//a := getTypeAsString(exprType.X)
		//if a == "nil" {
		//	return "nil"
		//}
		//v := fmt.Sprintf("&%s{}", a)
		//return v
		return "nil"

	case *ast.InterfaceType:
		return "nil"

	case *ast.MapType:
		return "nil"

	case *ast.ChanType:
		return "nil"

	case *ast.FuncType:
		return "nil"

	}

	return "unknown"

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
		v.DataType.DefaultValue = newDefaultValue

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

	//for _, s := range expressions {
	//	importFromMap, exist := gsb.importMap[s]
	//	if exist {
	//		gsb.usedImport[s] = importFromMap
	//	}
	//}

	return []string{}
}

func logDebug(format string, a ...any) {
	fmt.Printf(format, a...)
}
