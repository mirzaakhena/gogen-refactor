package gogen4

import (
	"fmt"
	"go/ast"
)

type GogenInterfaceBuilder struct {
	GogenBase
	//UnknownInterfaces map[GogenFieldType]*GogenAnyType
	//UnknownFields map[FieldSignature]*GogenField
	//CollectedType map[GogenFieldType]*TypeProperties
}

func NewGogenInterfaceBuilder() *GogenInterfaceBuilder {
	return &GogenInterfaceBuilder{
		GogenBase: GogenBase{GoModProperties: NewGoModProperties()},
		//UnknownInterfaces: map[GogenFieldType]*GogenAnyType{},
		//UnknownFields: map[FieldSignature]*GogenField{},
		//CollectedType: map[GogenFieldType]*TypeProperties{},
	}
}

func (g *GogenInterfaceBuilder) Build(packagePath, goModFilePath, interfaceTargetName string) (*GogenAnyType, error) {

	err := g.handleGoMod(goModFilePath)
	if err != nil {
		return nil, err
	}

	// kita start dari sini..
	gogenInterfaceRoot, err := g.traceInterfaceType(packagePath, interfaceTargetName)
	if err != nil {
		return nil, err
	}

	PrintGogenAnyType(0, gogenInterfaceRoot)

	return gogenInterfaceRoot, nil
}

func (g *GogenInterfaceBuilder) traceInterfaceType(packagePath string, interfaceTargetName string) (*GogenAnyType, error) {

	gogenInterfaceTarget := NewGogenAnyType()

	unknownInterfaces := map[FieldType]*GogenAnyType{}

	unknownFields := map[FieldSignature]*GogenFieldType{}

	collectedType := map[FieldType]*TypeProperties{}

	// Ada 3 hal yg dilakukan
	// pertama kita trace dulu as generalType dengan membaca interface tsb

	err := traceGeneralType(packagePath, interfaceTargetName, collectedType, func(tp TypeProperties) error {
		return g.handleGogenInterface(gogenInterfaceTarget, unknownInterfaces, unknownFields, &tp)
	})
	if err != nil {
		return nil, err
	}

	// kedua, jika ada interface yg tidak dikenal maka kita akan retracing ulang
	for fieldType, theInterface := range unknownInterfaces {

		typeProperties, exist := collectedType[fieldType]
		if !exist {
			return nil, fmt.Errorf("field type %v is not exist anywhere", fieldType)
		}

		err = g.handleGogenInterface(theInterface, unknownInterfaces, unknownFields, typeProperties)
		if err != nil {
			return nil, err
		}

		delete(unknownInterfaces, fieldType)

	}

	// ketiga, kita akan fill semua default value utk semua field yang ditemui
	for _, uf := range unknownFields {
		uf.DefaultValue, err = g.handleDefaultValue(uf, uf.Expr, uf.File, collectedType)
		if err != nil {
			return nil, err
		}
	}

	return gogenInterfaceTarget, nil
}

func (g *GogenInterfaceBuilder) handleGogenInterface(gi *GogenAnyType, unknownInterfaces map[FieldType]*GogenAnyType, unknownFields map[FieldSignature]*GogenFieldType, typeProperties *TypeProperties) error {

	name := FieldType(typeProperties.TypeSpec.Name.String())

	gi.GogenFieldType = &GogenFieldType{
		Name:         name,
		Expr:         typeProperties.TypeSpec.Type,
		DefaultValue: "nil",
		File:         typeProperties.AstFile,
	}

	switch ts := typeProperties.TypeSpec.Type.(type) {

	case *ast.InterfaceType:

		LogDebug(3, "as interface with name %s", gi.GogenFieldType.Name)

		for _, method := range ts.Methods.List {
			switch methodType := method.Type.(type) {

			case *ast.FuncType:
				err := g.handleDirectInlineMethod(gi, unknownFields, method)
				if err != nil {
					return err
				}

			case *ast.Ident:
				err := g.handleIdent(gi, unknownInterfaces, unknownFields, typeProperties, methodType)
				if err != nil {
					return err
				}

			case *ast.SelectorExpr:
				err := g.handleSelector(gi, typeProperties, methodType)
				if err != nil {
					return err
				}

			default:
				// TODO what about type alias?
				err := fmt.Errorf("un-handled method type %T", methodType)
				return err
			}
		}

	case *ast.Ident:

		// type alias will goes here
		err := g.handleIdent(gi, unknownInterfaces, unknownFields, typeProperties, ts)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("this is not an interface but %T", ts)

	}

	return nil
}

func (g *GogenInterfaceBuilder) handleSelector(gi *GogenAnyType, typeProperties *TypeProperties, methodType *ast.SelectorExpr) error {

	//interfacePath, err := getPathFromSelector(typeProperties.AstFile, gomodProperties, methodType)
	//if err != nil {
	//	return err
	//}

	importInFile, err := g.handleImport(typeProperties.AstFile.Imports)
	if err != nil {
		return err
	}

	theX := Expression(methodType.X.(*ast.Ident).String())

	interfacePath := string(importInFile[theX].CompletePath)

	internalGi, err := g.traceInterfaceType(interfacePath, methodType.Sel.String())
	if err != nil {
		return err
	}

	gi.CompositionTypes = append(gi.CompositionTypes, internalGi)

	return nil
}

func (g *GogenInterfaceBuilder) handleIdent(gi *GogenAnyType, unknownInterfaces map[FieldType]*GogenAnyType, unknownFields map[FieldSignature]*GogenFieldType, typeProperties *TypeProperties, methodType *ast.Ident) error {
	internalGi := new(GogenAnyType)
	gi.CompositionTypes = append(gi.CompositionTypes, internalGi)

	if methodType.Obj == nil {

		name := FieldType(methodType.String())

		internalGi.GogenFieldType = &GogenFieldType{
			Name:         name,
			Expr:         methodType,
			DefaultValue: "nil",
			File:         typeProperties.AstFile,
		}

		unknownInterfaces[name] = internalGi

		LogDebug(3, "unknown %s", methodType.String())

		return nil
	}

	newTypeSpec, ok := methodType.Obj.Decl.(*ast.TypeSpec)
	if !ok {
		return fmt.Errorf("%s is not type", methodType.String())
	}

	newTp := TypeProperties{
		AstFile:  typeProperties.AstFile,
		TypeSpec: newTypeSpec,
	}

	err := g.handleGogenInterface(internalGi, unknownInterfaces, unknownFields, &newTp)
	if err != nil {
		return err
	}

	return nil
}

func (g *GogenInterfaceBuilder) handleDirectInlineMethod(gi *GogenAnyType, unknownFields map[FieldSignature]*GogenFieldType, method *ast.Field) error {

	if method.Names == nil && len(method.Names) > 0 {
		err := fmt.Errorf("method must have name")
		return err
	}
	methodName := method.Names[0].String()

	if !ast.IsExported(methodName) {
		return nil
	}

	LogDebug(4, "as function found method %s", methodName)

	gm := newGogenMethod(methodName)

	gi.Methods = append(gi.Methods, gm)

	methodType, ok := method.Type.(*ast.FuncType)
	if !ok {
		return fmt.Errorf("somehow cannot convert to FuncType")
	}

	g.handleFuncParamResultType(methodType, unknownFields, gm, gi.GogenFieldType.File)

	return nil
}

func (g *GogenInterfaceBuilder) handleFuncParamResultType(methodType *ast.FuncType, unknownFields map[FieldSignature]*GogenFieldType, gm *GogenMethod, astFile *ast.File) {

	if methodType.Params.NumFields() > 0 {
		for _, param := range methodType.Params.List {

			fieldType := FieldType(getTypeAsString(param.Type))

			if param.Names != nil {

				for _, n := range param.Names {
					gf := NewGogenField(n.String(), param.Type, fieldType, astFile)
					gm.Params = append(gm.Params, gf)
					unknownFields[NewFieldSignature(gm.Name.String(), gf.Name.String())] = gf.DataType
				}

			} else {

				gf := NewGogenField(getSel(param.Type), param.Type, fieldType, astFile)
				gm.Params = append(gm.Params, gf)
				unknownFields[NewFieldSignature(gm.Name.String(), gf.Name.String())] = gf.DataType

			}

		}
	}

	if methodType.Results.NumFields() > 0 {
		for _, result := range methodType.Results.List {

			fieldType := FieldType(getTypeAsString(result.Type))

			if result.Names != nil {

				for _, n := range result.Names {
					gf := NewGogenField(n.String(), result.Type, fieldType, astFile)
					gm.Results = append(gm.Results, gf)
					unknownFields[NewFieldSignature(gm.Name.String(), gf.Name.String())] = gf.DataType
				}

			} else {

				gf := NewGogenField(getSel(result.Type), result.Type, fieldType, astFile)
				gm.Results = append(gm.Results, gf)
				unknownFields[NewFieldSignature(gm.Name.String(), gf.Name.String())] = gf.DataType

			}

		}
	}
}
