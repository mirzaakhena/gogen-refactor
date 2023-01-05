package gogen5

import (
	"fmt"
	"gen/gogen5/util"
	"go/ast"
)

type TypeProperties struct {
	TypeSpec *ast.TypeSpec
	AstFile  *ast.File
}

type gogenAnyTypeBuilder struct {
	goMod         util.GoModProperties
	defaultValues map[ast.Expr]string
}

type gogenData struct {
	unknownTypes              map[util.GogenFieldTypeName]*util.GogenAnyType
	unknownFieldDefaultValues map[*util.GogenFieldType]ast.Expr
	collectedTypes            map[util.GogenFieldTypeName]TypeProperties
	//collectedImports          map[util.Expression][]*ast.ImportSpec
}

func (r *gogenData) AddUnknownDefaultValue(gf *util.GogenFieldType, at gogenAnyTypeBuilder, expr ast.Expr, imports []*ast.ImportSpec) error {
	dv, err := util.GetDefaultValue(gf, expr, at.handleSelectorDefaultValue(imports))
	if err != nil {
		return err
	}
	at.defaultValues[expr] = dv
	r.unknownFieldDefaultValues[gf] = expr
	return nil
}

func newGogenData() *gogenData {
	return &gogenData{
		unknownTypes:              map[util.GogenFieldTypeName]*util.GogenAnyType{},
		unknownFieldDefaultValues: map[*util.GogenFieldType]ast.Expr{},
		collectedTypes:            map[util.GogenFieldTypeName]TypeProperties{},
		//collectedImports:          map[util.Expression][]*ast.ImportSpec{},
	}
}

func Build(packagePath, goModFilePath, typeTargetName string) (*util.GogenAnyType, error) {

	goMod, err := util.GetGoModProperties(goModFilePath)
	if err != nil {
		return nil, err
	}

	gat := util.NewGogenAnyType(typeTargetName)

	return gat, gogenAnyTypeBuilder{goMod: goMod, defaultValues: map[ast.Expr]string{}}.traceTypeInPath(packagePath, gat, typeTargetName)
}

func (r gogenAnyTypeBuilder) traceTypeInPath(packagePath string, gat *util.GogenAnyType, typeTargetName string) error {

	gd := newGogenData()

	err := util.TraceNode(packagePath, func(astPackage *ast.Package, astFile *ast.File, node ast.Node) (bool, error) {

		switch nodeTypeSpec := node.(type) {
		case *ast.TypeSpec:

			_, exist := gd.collectedTypes[util.GogenFieldTypeName(nodeTypeSpec.Name.String())]
			if !exist {
				gd.collectedTypes[util.GogenFieldTypeName(nodeTypeSpec.Name.String())] = TypeProperties{
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

func (r gogenAnyTypeBuilder) handleGogenType(gat *util.GogenAnyType, gd *gogenData, expr ast.Expr, astFile *ast.File) error {

	switch exprType := expr.(type) {

	case *ast.StructType:
		// handle type as struct
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
		//err := g.handleIdent(gi, unknownTypes, unknownFields, typeProperties, ts, collectedType)
		//if err != nil {
		//	return err
		//}

	default:
		return fmt.Errorf("unsupported type %T", exprType)

	}

	return nil
}

func (r gogenAnyTypeBuilder) handleInterfaceField(gat *util.GogenAnyType, gd *gogenData, iType *ast.InterfaceType, astFile *ast.File) error {

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

func (r gogenAnyTypeBuilder) handleSelector(gat *util.GogenAnyType, methodType *ast.SelectorExpr, astFile *ast.File) error {
	gi, err := util.GetGogenImport(methodType, astFile.Imports, r.goMod)
	if err != nil {
		return err
	}

	gat.Imports[util.Expression(methodType.X.(*ast.Ident).String())] = gi

	newGat := util.NewGogenAnyType(util.GetTypeAsString(methodType))

	err = r.traceTypeInPath(gi.CompletePath, newGat, util.GetBasicType(methodType))
	if err != nil {
		return err
	}

	gat.AddCompositionType(newGat)
	return nil
}

func (r gogenAnyTypeBuilder) handleStructField(gat *util.GogenAnyType, gd *gogenData, sType *ast.StructType, astFile *ast.File) error {

	for _, field := range sType.Fields.List {

		if field.Names != nil {
			// simple field
			for _, n := range field.Names {
				gf := util.NewGogenField(n.String(), field.Type)
				gat.AddField(gf)
				err := gd.AddUnknownDefaultValue(gf.DataType, r, field.Type, astFile.Imports)
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

func (r gogenAnyTypeBuilder) handleInterfaceMethod(gat *util.GogenAnyType, gd *gogenData, method *ast.Field, astFile *ast.File) error {

	if method.Names == nil && len(method.Names) > 0 {
		return fmt.Errorf("unusual method signature")
	}
	methodName := method.Names[0].String()

	if !ast.IsExported(methodName) {
		return nil
	}

	gm := util.NewGogenMethod(methodName)

	gat.AddMethod(gm)

	methodType, ok := method.Type.(*ast.FuncType)
	if !ok {
		return fmt.Errorf("cannot convert method to FuncType")
	}

	err := r.handleFuncParamResultType(gm, gd, methodType, astFile)
	if err != nil {
		return err
	}

	return nil
}

func (r gogenAnyTypeBuilder) handleFuncParamResultType(gm *util.GogenMethod, gd *gogenData, methodType *ast.FuncType, astFile *ast.File) error {

	if methodType.Params.NumFields() > 0 {
		for _, param := range methodType.Params.List {

			if param.Names == nil {
				gf := util.NewGogenField(util.GetBasicType(param.Type), param.Type)
				gm.AddParam(gf)
				err := gd.AddUnknownDefaultValue(gf.DataType, r, param.Type, astFile.Imports)
				if err != nil {
					return err
				}

				continue
			}

			for _, n := range param.Names {
				gf := util.NewGogenField(n.String(), param.Type)
				gm.AddParam(gf)
				err := gd.AddUnknownDefaultValue(gf.DataType, r, param.Type, astFile.Imports)
				if err != nil {
					return err
				}
			}

		}
	}

	if methodType.Results.NumFields() > 0 {
		for _, result := range methodType.Results.List {

			if result.Names == nil {
				gf := util.NewGogenField(util.GetBasicType(result.Type), result.Type)
				gm.AddResult(gf)
				err := gd.AddUnknownDefaultValue(gf.DataType, r, result.Type, astFile.Imports)
				if err != nil {
					return err
				}
				continue
			}

			for _, n := range result.Names {
				gf := util.NewGogenField(n.String(), result.Type)
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

func (r gogenAnyTypeBuilder) handleSelectorDefaultValue(imports []*ast.ImportSpec) func(selectorExpr *ast.SelectorExpr) (ast.Expr, error) {

	return func(selectorExpr *ast.SelectorExpr) (ast.Expr, error) {

		gi, err := util.GetGogenImport(selectorExpr, imports, r.goMod)
		if err != nil {
			return nil, err
		}

		var expr ast.Expr

		err = util.TraceNode(gi.CompletePath, func(astPackage *ast.Package, astFile *ast.File, node ast.Node) (bool, error) {

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

func (r gogenAnyTypeBuilder) handleIdent(gat *util.GogenAnyType, gd *gogenData, fieldType *ast.Ident, astFile *ast.File) error {

	newGat := util.NewGogenAnyType(fieldType.String())

	gat.CompositionTypes = append(gat.CompositionTypes, newGat)

	if fieldType.Obj == nil {

		gd.unknownTypes[util.NewGogenFieldTypeName(fieldType)] = newGat

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
