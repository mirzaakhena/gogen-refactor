package gogen5

import (
	"fmt"
	"go/ast"
)

type gogenAnyTypeBuilder struct {
	goMod         GoModProperties
	defaultValues map[ast.Expr]string
}

type gogenData struct {
	unknownTypes              map[GogenFieldTypeName]*GogenAnyType
	unknownFieldDefaultValues map[*GogenFieldType]ast.Expr
	collectedTypes            map[GogenFieldTypeName]TypeProperties
	//collectedImports          map[util.Expression][]*ast.ImportSpec
}

func (r *gogenData) AddUnknownDefaultValue(gf *GogenFieldType, at gogenAnyTypeBuilder, expr ast.Expr, imports []*ast.ImportSpec) error {
	dv, err := GetDefaultValue(gf, expr, r.collectedTypes, at.handleSelectorDefaultValue(imports))
	if err != nil {
		return err
	}

	at.defaultValues[expr] = dv
	r.unknownFieldDefaultValues[gf] = expr
	return nil
}

func newGogenData() *gogenData {
	return &gogenData{
		unknownTypes:              map[GogenFieldTypeName]*GogenAnyType{},
		unknownFieldDefaultValues: map[*GogenFieldType]ast.Expr{},
		collectedTypes:            map[GogenFieldTypeName]TypeProperties{},
		//collectedImports:          map[util.Expression][]*ast.ImportSpec{},
	}
}

func Build(packagePath, goModFilePath, typeTargetName string) (*GogenAnyType, error) {

	goMod, err := GetGoModProperties(goModFilePath)
	if err != nil {
		return nil, err
	}

	gat := NewGogenAnyType(typeTargetName)

	return gat, gogenAnyTypeBuilder{goMod: goMod, defaultValues: map[ast.Expr]string{}}.traceTypeInPath(packagePath, gat, typeTargetName)
}

func (r gogenAnyTypeBuilder) traceTypeInPath(packagePath string, gat *GogenAnyType, typeTargetName string) error {

	gd := newGogenData()

	err := TraceNode(packagePath, func(astPackage *ast.Package, astFile *ast.File, node ast.Node) (bool, error) {

		switch nodeTypeSpec := node.(type) {
		case *ast.TypeSpec:

			_, exist := gd.collectedTypes[GogenFieldTypeName(nodeTypeSpec.Name.String())]
			if !exist {
				gd.collectedTypes[GogenFieldTypeName(nodeTypeSpec.Name.String())] = TypeProperties{
					TypeSpec: nodeTypeSpec,
					AstFile:  astFile,
				}
			}

			if nodeTypeSpec.Name.String() != typeTargetName {
				return true, nil
			}

			err := r.handleGogenType(gat, gd, nodeTypeSpec.Type, astFile)
			if err != nil {
				return false, err
			}

		case *ast.FuncDecl:
			// TODO found method
		}

		return true, nil
	})
	if err != nil {
		return err
	}

	for k, ct := range gd.collectedTypes {

		_ = k
		_ = ct
	}

	for fieldType, theType := range gd.unknownTypes {

		ct, exist := gd.collectedTypes[fieldType]

		if !exist {
			return fmt.Errorf("tidak ada")
		}

		err = r.handleGogenType(theType, gd, ct.TypeSpec.Type, ct.AstFile)
		if err != nil {
			return err
		}

	}

	for fieldType, gt := range gd.unknownFieldDefaultValues {
		fieldType.DefaultValue = r.defaultValues[gt]
	}

	return nil

}

func (r gogenAnyTypeBuilder) handleGogenType(gat *GogenAnyType, gd *gogenData, expr ast.Expr, astFile *ast.File) error {

	switch exprType := expr.(type) {

	case *ast.StructType:
		// TODO handle type as struct
		err := r.handleStructField(gat, gd, exprType, astFile)
		if err != nil {
			return err
		}

	case *ast.InterfaceType:
		// handle type as interface
		err := r.handleInterfaceField(gat, gd, exprType, astFile)
		if err != nil {
			return err
		}

	case *ast.Ident:
		// TODO handle type as ident alias

		LogDebug(1, ">>>>>>>>>>> %v", exprType.String())

		err := r.handleIdent(gat, gd, exprType, astFile)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsupported type %T", exprType)

	}

	return nil
}

func (r gogenAnyTypeBuilder) handleInterfaceField(gat *GogenAnyType, gd *gogenData, iType *ast.InterfaceType, astFile *ast.File) error {

	for _, method := range iType.Methods.List {
		switch methodType := method.Type.(type) {
		case *ast.FuncType:

			// handle interface method as inline func
			err := r.handleInterfaceMethod(gat, gd, method, astFile)
			if err != nil {
				return err
			}

		case *ast.Ident:

			// handle interface method as ident
			err := r.handleIdent(gat, gd, methodType, astFile)
			if err != nil {
				return err
			}

		case *ast.SelectorExpr:

			// handle interface method as selector
			err := r.handleSelector(gat, methodType, astFile)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unsupported interface field %T", methodType)

		}
	}

	return nil
}

func (r gogenAnyTypeBuilder) handleSelector(gat *GogenAnyType, methodType *ast.SelectorExpr, astFile *ast.File) error {
	gi, err := GetGogenImport(methodType, astFile.Imports, r.goMod)
	if err != nil {
		return err
	}

	expr := Expression(methodType.X.(*ast.Ident).String())

	gat.Imports[expr] = *gi

	newGat := NewGogenAnyType(GetTypeAsString(methodType))

	err = r.traceTypeInPath(gi.CompletePath, newGat, GetBasicType(methodType))
	if err != nil {
		return err
	}

	gat.AddCompositionType(newGat)
	return nil
}

func (r gogenAnyTypeBuilder) handleStructField(gat *GogenAnyType, gd *gogenData, sType *ast.StructType, astFile *ast.File) error {

	for _, field := range sType.Fields.List {

		err := r.handleImport(gat, field.Type, astFile)
		if err != nil {
			return err
		}

		if field.Names != nil {
			// simple field
			for _, n := range field.Names {
				gf := NewGogenField(n.String(), field.Type)
				gat.AddField(gf)
				err = gd.AddUnknownDefaultValue(gf.DataType, r, field.Type, astFile.Imports)
				if err != nil {
					return err
				}

			}
			continue
		}

		switch fieldType := field.Type.(type) {
		case *ast.SelectorExpr:
			// TODO handle struct field as Selector
		case *ast.Ident:
			// TODO handle struct field as Ident
		case *ast.StarExpr:
			// TODO handle struct field as Star
		default:
			return fmt.Errorf("unsupported struct field %T", fieldType)

		}

	}

	return nil
}

func (r gogenAnyTypeBuilder) handleInterfaceMethod(gat *GogenAnyType, gd *gogenData, method *ast.Field, astFile *ast.File) error {

	if method.Names == nil && len(method.Names) > 0 {
		return fmt.Errorf("unusual method signature")
	}
	methodName := method.Names[0].String()

	if !ast.IsExported(methodName) {
		return nil
	}

	gm := NewGogenMethod(methodName)

	gat.AddMethod(gm)

	methodType, ok := method.Type.(*ast.FuncType)
	if !ok {
		return fmt.Errorf("cannot convert method to FuncType")
	}

	err := r.handleFuncParamResultType(gat, gm, gd, methodType, astFile)
	if err != nil {
		return err
	}

	return nil
}

func (r gogenAnyTypeBuilder) handleFuncParamResultType(gat *GogenAnyType, gm *GogenMethod, gd *gogenData, methodType *ast.FuncType, astFile *ast.File) error {

	if methodType.Params.NumFields() > 0 {
		for _, param := range methodType.Params.List {

			err := r.handleImport(gat, param.Type, astFile)
			if err != nil {
				return err
			}

			if param.Names == nil {
				gf := NewGogenField(GetBasicType(param.Type), param.Type)
				gm.AddParam(gf)
				err := gd.AddUnknownDefaultValue(gf.DataType, r, param.Type, astFile.Imports)
				if err != nil {
					return err
				}

				continue
			}

			for _, n := range param.Names {

				gf := NewGogenField(n.String(), param.Type)
				gm.AddParam(gf)
				err = gd.AddUnknownDefaultValue(gf.DataType, r, param.Type, astFile.Imports)
				if err != nil {
					return err
				}
			}

		}
	}

	if methodType.Results.NumFields() > 0 {
		for _, result := range methodType.Results.List {

			err := r.handleImport(gat, result.Type, astFile)
			if err != nil {
				return err
			}

			if result.Names == nil {
				gf := NewGogenField(GetBasicType(result.Type), result.Type)
				gm.AddResult(gf)
				err := gd.AddUnknownDefaultValue(gf.DataType, r, result.Type, astFile.Imports)
				if err != nil {
					return err
				}
				continue
			}

			for _, n := range result.Names {
				gf := NewGogenField(n.String(), result.Type)
				gm.AddResult(gf)
				err := gd.AddUnknownDefaultValue(gf.DataType, r, result.Type, astFile.Imports)
				if err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func (r gogenAnyTypeBuilder) handleImport(gat *GogenAnyType, expr ast.Expr, astFile *ast.File) error {

	selectorExpr := GetSelectorExpr(expr)
	if selectorExpr == nil {
		return nil
	}

	gi, err := GetGogenImport(selectorExpr, astFile.Imports, r.goMod)
	if err != nil {
		return err
	}

	exprFromSelector := Expression(selectorExpr.X.(*ast.Ident).String())

	if _, exist := gat.Imports[exprFromSelector]; !exist {
		gat.Imports[exprFromSelector] = *gi
	}

	return nil
}

func (r gogenAnyTypeBuilder) handleSelectorDefaultValue(imports []*ast.ImportSpec) func(selectorExpr *ast.SelectorExpr) (ast.Expr, error) {

	return func(selectorExpr *ast.SelectorExpr) (ast.Expr, error) {

		gi, err := GetGogenImport(selectorExpr, imports, r.goMod)
		if err != nil {
			return nil, err
		}

		var expr ast.Expr

		err = TraceNode(gi.CompletePath, func(astPackage *ast.Package, astFile *ast.File, node ast.Node) (bool, error) {

			switch nodeTypeSpec := node.(type) {
			case *ast.TypeSpec:
				if nodeTypeSpec.Name.String() != selectorExpr.Sel.String() {
					return true, nil
				}
				expr = nodeTypeSpec.Type
			}

			return true, nil
		})
		if err != nil {
			return nil, err
		}

		return expr, nil
	}
}

func (r gogenAnyTypeBuilder) handleIdent(gat *GogenAnyType, gd *gogenData, fieldType *ast.Ident, astFile *ast.File) error {

	newGat := NewGogenAnyType(fieldType.String())

	gat.CompositionTypes = append(gat.CompositionTypes, newGat)

	if fieldType.Obj == nil {

		gd.unknownTypes[NewGogenFieldTypeName(fieldType)] = newGat

		return nil
	}

	newTypeSpec, ok := fieldType.Obj.Decl.(*ast.TypeSpec)
	if !ok {
		return fmt.Errorf("%s is not type", fieldType.String())
	}

	err := r.handleGogenType(newGat, gd, newTypeSpec.Type, astFile)
	if err != nil {
		return err
	}

	return nil
}
