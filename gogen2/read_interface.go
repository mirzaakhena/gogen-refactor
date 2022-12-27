package gogen2

import (
	"fmt"
	"go/ast"
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

	gc := NewGogenInterface(interfaceName)

	// trace the primary interface
	err := gsb.traceType(gsb.path, gc.Name, gsb.handleInterfaceTarget(gc))
	if err != nil {
		return nil, err
	}

	return gc, nil
}

func (gsb *GogenInterfaceBuilder) handleInterfaceTarget(gc *GogenInterface) func(fieldType FieldType, expr ast.Expr) error {

	return func(fieldType FieldType, expr ast.Expr) error {

		switch ts := expr.(type) {
		case *ast.InterfaceType:
			logDebug("bertemu interface target %s !! langsung kita olah MethodInterfacenya", fieldType)
			logDebug("start  handleMethodInterface")
			err := gsb.handleMethodInterface(fieldType, ts, gc)
			if err != nil {
				return err
			}
			logDebug("finish handleMethodInterface")

		default:
			return fmt.Errorf("unsupported type %T", ts)
		}

		return nil
	}
}

func (gsb *GogenInterfaceBuilder) handleMethodInterface(typeSpecName FieldType, interfaceType *ast.InterfaceType, gc *GogenInterface) error {

	for _, method := range interfaceType.Methods.List {

		switch methodType := method.Type.(type) {

		case *ast.FuncType:
			logDebug("masuk sebagai inline func")

			// harus ada nama kalau tidak ini tidak normal
			if method.Names == nil && len(method.Names) > 0 {
				return fmt.Errorf("method must have name")
			}
			methodName := method.Names[0].String()

			// nama methodnya harus public
			if !IsExported(methodName) {
				continue
			}

			// ini adalah tujuan akhir dari method ini
			gm := NewGogenMethod(typeSpecName, MethodName(methodName))
			gc.Methods = append(gc.Methods, gm)

			// karena sudah ketemu, jadi kita remove dari unknown interface
			//delete(gsb.unknownInterface, typeSpecName)

			//logDebug("start  handleFuncParamResultType")
			//gsb.handleFuncParamResultType(methodType, gm)
			//logDebug("finish handleFuncParamResultType")

		case *ast.Ident:
			logDebug("masuk sebagai ident %v", methodType.String())

			// tidak mungkin ada import disini
			// disini kita berharap ident yg ditemukan sudah pernah didaftarkan pada typeMap
			im, exist := gsb.typeMap[FieldType(methodType.String())]
			if !exist {
				// jika masuk disini, maka ident belum pernah ditemukan, tapi mgk nanti akan ketemu
				// mungkin ada di package yg sama dan file yg sama (setelah interface target ditemukan),
				// mungkin ada di package yg sama namun file yg berbeda
				// dan tidak mungkin ada di package yang berbeda
				// kita belum tahu methodType itu apa, masukin aj dulu nanti akan kita cek.
				gsb.unknownInterface[FieldType(methodType.String())] = methodType
				continue
			}

			// jika masuk kesini, maka ident ini sudah pasti sebuah interface yang sudah pernah ditemukan diawal
			// tapi belum pernah ditelusuri lebih lanjut
			// kita akan selesaikan dengan dirinya sendiri sebagai interfaceType (recursive)
			// masalahnya, jika misal ident ini ditemukan di file yang berbeda yg dibaca lebih dahulu sebelum interface target maka kita akan miss importnya
			err := gsb.handleMethodInterface(FieldType(methodType.String()), im.(*ast.InterfaceType), gc)
			if err != nil {
				return err
			}

		case *ast.SelectorExpr:
			logDebug("masuk sebagai selector")

			// kalau masuk sini sudah pasti belum pernah ditemukan dalam interfaceType
			// sudah pasti ada import yang akan kita pakai disini
			for _, expr := range gsb.extractAllExpression(methodType) {
				importFromMap, exist := gsb.importMap[expr]
				if !exist {
					return fmt.Errorf("aneh ketemu selector tapi tidak ada di import")
				}
				gsb.usedImport[expr] = importFromMap
			}

			// dapatkan namanya
			//m := fmt.Sprintf("%v.%v", methodType.X.(*ast.Ident).String(), methodType.Sel.String())
			m := methodType.Sel.String()
			logDebug("nama selectornya %v", m)

			theEx := Expression(methodType.X.(*ast.Ident).String())
			gi, exist := gsb.importMap[theEx]
			if !exist {
				return fmt.Errorf("%s is not exist in importmap", theEx)
			}

			err := gsb.traceType(string(gi.CompletePath), FieldType(methodType.Sel.String()), gsb.handleInterfaceTarget(gc))
			if err != nil {
				return err
			}

			// sudah pasti ada di package yg berbeda yang akan kita telusuri nanti
			// method Type disini sudah pasti selector,
			// tapi kita belum tahu Selectornya type apa, masukin aj dulu nanti akan kita cek.
			//gsb.unknownInterface[FieldType(m)] = methodType

		default:
			return fmt.Errorf("unsupported type %v", methodType)

		}

	}

	return nil
}
