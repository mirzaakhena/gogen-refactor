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

type GogenInterfaceBuilder struct {
	goModPath    string
	path         string
	importMap    map[string]GogenImport
	usedImport   map[string]GogenImport
	typeMap      map[string]*ast.TypeSpec
	unknownTypes map[string]*GogenField
	foundTarget  bool
	selectorMap  map[string][]string
	interfaceMap map[string]ast.Expr
}

func NewGogenInterfaceBuilder(goModPath, path string) *GogenInterfaceBuilder {

	return &GogenInterfaceBuilder{
		goModPath:    goModPath,
		path:         path,
		importMap:    map[string]GogenImport{},
		usedImport:   map[string]GogenImport{},
		typeMap:      map[string]*ast.TypeSpec{},
		unknownTypes: map[string]*GogenField{},
		selectorMap:  map[string][]string{},
		interfaceMap: map[string]ast.Expr{},
		foundTarget:  false,
	}
}

func (gsb *GogenInterfaceBuilder) Build(interfaceName string) (*GogenInterface, error) {

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, gsb.path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	gc := NewGogenInterface(interfaceName)

	for _, pkg := range pkgs {

		// now we try to find the typeSpecName == interfaceName
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

				if typeSpecName != interfaceName {

					logDebug("simpan kedalam typeMap type %v \n", typeSpecName)
					gsb.typeMap[typeSpecName] = typeSpec

					// kita coba daftarkan semua interface yang kita temukan di package yang sama meskipun file yang berbeda
					// mungkin nanti akan dipakai sebelum atau sesudah interface target ditemukan

					interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
					if ok {
						// disini interfaceType yang dimasukkan sudah pasti sebuah interface
						// tapi dari package yg sama namun file nya bisa sama atau berbeda
						// baik sebelum maupun sesudah interface target ditemukan
						gsb.interfaceMap[typeSpecName] = interfaceType
					}

					if gsb.foundTarget {

						// TODO handle Uncomplete Interface Map

						gsb.handleUncompleteDefaultValue()

					}

					logDebug(".\n")
					return true
				}

				// -------------- we found the interface target --------------

				gsb.foundTarget = true

				logDebug("target interface %v sudah ditemukan\n", interfaceName)

				err = gsb.handleMethodInterface(typeSpecName, typeSpec.Type, gc)
				if err != nil {
					return false
				}

				return true
			})

			if err != nil {
				return nil, err
			}

		}

	}

	// disini kita akan coba ambil semua interfaceMap yang berisi Expr
	// expr ini berasal dari selector saat nemu target

	//for x, im := range gsb.interfaceMap {
	//
	//	ui, exist := gsb.usedImport[x]
	//	if !exist {
	//		return nil, fmt.Errorf("%s not found in importMap", x)
	//	}
	//
	//	path := gsb.getPathBasedOnImport(ui, x)
	//
	//	fmt.Printf("call path %v\n", path)
	//
	//	// go to the file
	//	fset := token.NewFileSet()
	//	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	found := false
	//
	//	for _, pkg := range pkgs {
	//
	//		for _, file := range pkg.Files {
	//
	//			ast.Inspect(file, func(node ast.Node) bool {
	//
	//				// focus only to type
	//				typeSpec, ok := node.(*ast.TypeSpec)
	//				if !ok {
	//					return true
	//				}
	//
	//				// focus to interface only
	//				_, ok = typeSpec.Type.(*ast.InterfaceType)
	//				if !ok {
	//					return true
	//				}
	//
	//				interfaceName := typeSpec.Name.String()
	//
	//				//gsb.interfaceMap
	//
	//				return true
	//			})
	//
	//		}
	//	}
	//
	//}

	// karena bisa jadi ada selector yg harus dihandle juga

	gsb.handleSelector(gc) // TODO ...

	gsb.handleUncompleteDefaultValue()

	if len(gsb.unknownTypes) > 0 {

		arrUnknownTypes := make([]string, 0)

		for k, s := range gsb.unknownTypes {
			arrUnknownTypes = append(arrUnknownTypes, fmt.Sprintf("%s %s,", k, s.DataType.Type))
		}

		return nil, fmt.Errorf("unknown type for field : %v", arrUnknownTypes)
	}

	return gc, nil
}

func (gsb *GogenInterfaceBuilder) handleMethodInterface(typeSpecName string, expr ast.Expr, gc *GogenInterface) error {

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
			gsb.handleFuncType(methodName, gc, methodType)

		case *ast.Ident:
			fmt.Printf("ident %v\n", methodType.String())

			// tidak mungkin ada import disini
			// disini kita berharap ident yg ditemukan sudah pernah didaftarkan pada interfaceType
			// dan ident disini sudah pasti adalah sebuah interface
			im, exist := gsb.interfaceMap[methodType.String()]
			if !exist {

				// jika masuk disini, maka ident belum pernah ditemukan, tapi mgk nanti akan ketemu
				// mungkin ada di package yg sama dan file yg sama,
				// mungkin ada di package yg sama namun file yg berbeda
				// dan tidak mungkin ada di package yang berbeda
				// kita belum tahu methodType itu apa, masukin aj dulu nanti akan kita cek.
				gsb.interfaceMap[methodType.String()] = methodType
				continue
			}

			// jika masuk kesini, maka ident ini sudah pasti sebuah interface yang sudah pernah ditemukan diawal
			// tapi belum pernah ditelusuri lebih lanjut
			// kita akan selesaikan dengan dirinya sendiri sebagai interfaceType (recursive)
			err := gsb.handleMethodInterface(methodType.String(), im, gc)
			if err != nil {
				return err
			}

		case *ast.SelectorExpr:
			m := fmt.Sprintf("%v.%v", methodType.X.(*ast.Ident).String(), methodType.Sel.String())
			fmt.Printf("selector %v\n", m)

			// kalau masuk sini sudah pasti belum pernah ditemukan dalam interfaceType
			// sudah pasti ada import yang akan kita pakai disini
			for _, s := range gsb.handleUsedImport(methodType) {
				importFromMap, exist := gsb.importMap[s]
				if exist {
					gsb.usedImport[s] = importFromMap
				}
			}

			// sudah pasti ada di package yg berbeda yang akan kita telusuri nanti
			// method Type disini sudah pasti selector,
			// tapi kita belum tahu Selectornya type apa, masukin aj dulu nanti akan kita cek.
			gsb.interfaceMap[m] = methodType

		default:
			return fmt.Errorf("unsupported type %v\n", methodType)
		}

	}

	return nil
}

func (gsb *GogenInterfaceBuilder) handleFuncType(methodName string, gc *GogenInterface, methodType *ast.FuncType) {
	gm := NewGogenMethod(methodName)

	mf := GogenMethodField{
		ExtendInterfaces: "",
		Methods:          []*GogenMethod{gm},
	}

	gc.MethodFields = append(gc.MethodFields, &mf)

	if methodType.Params.NumFields() > 0 {
		for _, param := range methodType.Params.List {

			dataTypeStr := getTypeAsString(param.Type)

			defaultValue := handleDefaultValue(dataTypeStr, param.Type)

			for _, s := range gsb.handleUsedImport(param.Type) {
				importFromMap, exist := gsb.importMap[s]
				if exist {
					gsb.usedImport[s] = importFromMap
				}
			}

			if param.Names != nil {
				for _, n := range param.Names {
					if IsExported(n.String()) {
						gf := NewGogenField(n.String(), dataTypeStr, defaultValue)
						gm.Params = append(gm.Params, gf)
						gsb.checkDefaultValue(gf)
					}
				}
			} else {
				fieldNameStr := GetSel(param.Type)
				gf := NewGogenField(fieldNameStr, dataTypeStr, defaultValue)
				gm.Params = append(gm.Params, gf)
				gsb.checkDefaultValue(gf)

			}

		}
	}

	if methodType.Results.NumFields() > 0 {
		for _, result := range methodType.Results.List {
			dataTypeStr := getTypeAsString(result.Type)

			defaultValue := handleDefaultValue(dataTypeStr, result.Type)

			for _, s := range gsb.handleUsedImport(result.Type) {
				importFromMap, exist := gsb.importMap[s]
				if exist {
					gsb.usedImport[s] = importFromMap
				}
			}

			if result.Names != nil {
				for _, n := range result.Names {
					if IsExported(n.String()) {
						gf := NewGogenField(n.String(), dataTypeStr, defaultValue)
						gm.Params = append(gm.Params, gf)
						gsb.checkDefaultValue(gf)
					}
				}
			} else {
				fieldNameStr := GetSel(result.Type)
				gf := NewGogenField(fieldNameStr, dataTypeStr, defaultValue)
				gm.Params = append(gm.Params, gf)
				gsb.checkDefaultValue(gf)

			}
		}
	}
}

func (gsb *GogenInterfaceBuilder) handleSelector(gs *GogenInterface) {

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

func (gsb *GogenInterfaceBuilder) checkDefaultValue(gf *GogenField) {

	if gf.DataType.Type == gf.DataType.DefaultValue {

		logDebug("karena defaultValue utk field %v dengan type %v belum final, kita cek ke map\n", gf.Name, gf.DataType.Type)

		typeSpecFromMap, exist := gsb.typeMap[gf.DataType.Type]
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

func (gsb *GogenInterfaceBuilder) handleUncompleteDefaultValue() {

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

func (gsb *GogenInterfaceBuilder) getPathBasedOnImport(gi GogenImport, x string) string {
	if strings.HasPrefix(gi.Path, gsb.goModPath) {
		return gi.Path[len(gsb.goModPath)+1:]
	}
	return fmt.Sprintf("%s/src/%s", build.Default.GOROOT, x)
}

func (gsb *GogenInterfaceBuilder) handleUsedImport(expr ast.Expr) []string {

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
