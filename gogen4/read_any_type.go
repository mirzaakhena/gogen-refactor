package gogen4

import (
	"fmt"
	"go/ast"
)

type GogenAnyTypeBuilder struct {
	GogenBase
}

func NewGogenAnyTypeBuilder() *GogenAnyTypeBuilder {
	return &GogenAnyTypeBuilder{
		GogenBase: GogenBase{GoModProperties: NewGoModProperties()},
	}
}

func (g *GogenAnyTypeBuilder) Build(packagePath, goModFilePath, typeTargetName string) (*GogenAnyType, error) {

	err := g.handleGoMod(goModFilePath)
	if err != nil {
		return nil, err
	}

	gogenRootType, err := g.traceTypeInPath(packagePath, typeTargetName)
	if err != nil {
		return nil, err
	}

	PrintGogenAnyType(0, gogenRootType)

	return gogenRootType, nil
}

func (g *GogenAnyTypeBuilder) traceTypeInPath(packagePath string, typeTargetName string) (*GogenAnyType, error) {

	gogenTypeTarget := NewGogenAnyType()

	unknownTypes := map[FieldType]*GogenAnyType{}

	unknownFields := map[FieldSignature]*GogenFieldType{}

	collectedType := map[FieldType]*TypeProperties{}

	// pertama kita trace dulu as generalType dengan membaca interface tsb
	err := findTargetType(packagePath, typeTargetName, collectedType, func(tp TypeProperties) error {
		return g.handleGogenType(gogenTypeTarget, unknownTypes, unknownFields, collectedType, &tp)
	})
	if err != nil {
		return nil, err
	}

	// kedua, jika ada interface yg tidak dikenal maka kita akan retracing ulang
	for fieldType, theType := range unknownTypes {

		typeProperties, exist := collectedType[fieldType]
		if !exist {
			return nil, fmt.Errorf("field type %v is not exist anywhere", fieldType)
		}

		err = g.handleGogenType(theType, unknownTypes, unknownFields, collectedType, typeProperties)
		if err != nil {
			return nil, err
		}

		delete(unknownTypes, fieldType)

	}

	// ketiga, kita akan fill semua default value utk semua field yang ditemui
	for _, uf := range unknownFields {
		uf.DefaultValue, err = g.handleDefaultValue(uf, uf.Expr, uf.File, collectedType)
		if err != nil {
			return nil, err
		}
	}

	return gogenTypeTarget, nil
}

func (g *GogenAnyTypeBuilder) handleGogenType(gi *GogenAnyType, unknownTypes map[FieldType]*GogenAnyType, unknownFields map[FieldSignature]*GogenFieldType, collectedType map[FieldType]*TypeProperties, typeProperties *TypeProperties) error {

	typeSpecName := FieldType(typeProperties.TypeSpec.Name.String())

	gi.GogenFieldType = &GogenFieldType{
		Name:         typeSpecName,
		Expr:         typeProperties.TypeSpec.Type,
		DefaultValue: "Hahahahaha",
		File:         typeProperties.AstFile,
	}

	unknownFields[NewFieldSignature(gi.GogenFieldType.Name.String(), typeSpecName.String())] = gi.GogenFieldType

	switch ts := typeProperties.TypeSpec.Type.(type) {

	case *ast.StructType:

		err := g.handleStructField(gi, unknownTypes, unknownFields, collectedType, typeProperties, ts, typeSpecName)
		if err != nil {
			return err
		}

	case *ast.InterfaceType:

		err, done := g.handleInterfaceField(gi, unknownTypes, unknownFields, collectedType, typeProperties, ts)
		if done {
			return err
		}

	case *ast.Ident:

		// type alias will goes here
		err := g.handleIdent(gi, unknownTypes, unknownFields, typeProperties, ts, collectedType)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("this is not an interface but %T", ts)

	}

	return nil
}

func (g *GogenAnyTypeBuilder) handleInterfaceField(gi *GogenAnyType, unknownTypes map[FieldType]*GogenAnyType, unknownFields map[FieldSignature]*GogenFieldType, collectedType map[FieldType]*TypeProperties, typeProperties *TypeProperties, ts *ast.InterfaceType) (error, bool) {
	name := FieldType(typeProperties.TypeSpec.Name.String())

	gi.GogenFieldType = &GogenFieldType{
		Name:         name,
		Expr:         typeProperties.TypeSpec.Type,
		DefaultValue: "nil",
		File:         typeProperties.AstFile,
	}

	//LogDebug(3, "as interface with name %s", gi.GogenFieldType.Name)

	for _, method := range ts.Methods.List {
		switch methodType := method.Type.(type) {

		case *ast.FuncType:
			err := g.handleInterfaceMethod(gi, unknownFields, method)
			if err != nil {
				return err, true
			}

		case *ast.Ident:
			err := g.handleIdent(gi, unknownTypes, unknownFields, typeProperties, methodType, collectedType)
			if err != nil {
				return err, true
			}

		case *ast.SelectorExpr:
			err := g.handleSelector(gi, typeProperties, methodType)
			if err != nil {
				return err, true
			}

		default:
			// TODO what about type alias?
			err := fmt.Errorf("un-handled method type %T", methodType)
			return err, true
		}
	}
	return nil, false
}

func (g *GogenAnyTypeBuilder) handleStructField(gi *GogenAnyType, unknownTypes map[FieldType]*GogenAnyType, unknownFields map[FieldSignature]*GogenFieldType, collectedType map[FieldType]*TypeProperties, typeProperties *TypeProperties, ts *ast.StructType, typeSpecName FieldType) error {
	for _, field := range ts.Fields.List {

		if field.Names != nil {

			for _, n := range field.Names {
				gf := NewGogenField(n.String(), field.Type, FieldType(getTypeAsString(field.Type)), typeProperties.AstFile)
				gi.Fields = append(gi.Fields, gf)
				unknownFields[NewFieldSignature(gi.GogenFieldType.Name.String(), gf.Name.String())] = gf.DataType
			}

		} else {

			err := g.handleField(gi, unknownTypes, unknownFields, collectedType, typeProperties, field.Type)
			if err != nil {
				return err
			}

		}

	}
	return nil
}

func (g *GogenAnyTypeBuilder) handleField(gi *GogenAnyType, unknownTypes map[FieldType]*GogenAnyType, unknownFields map[FieldSignature]*GogenFieldType, collectedType map[FieldType]*TypeProperties, typeProperties *TypeProperties, field ast.Expr) error {
	switch fieldType := field.(type) {

	case *ast.Ident:
		//LogDebug(1, "===> as ident %v", fieldType)
		err := g.handleIdent(gi, unknownTypes, unknownFields, typeProperties, fieldType, collectedType)
		if err != nil {
			return err
		}

	case *ast.StarExpr:
		//LogDebug(1, "===> as star %v", fieldType)
		err := g.handleField(gi, unknownTypes, unknownFields, collectedType, typeProperties, fieldType.X)
		if err != nil {
			return err
		}

	case *ast.SelectorExpr:
		//LogDebug(1, "===> as selector %v", fieldType)
		err := g.handleSelector(gi, typeProperties, fieldType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *GogenAnyTypeBuilder) handleSelector(gi *GogenAnyType, typeProperties *TypeProperties, methodType *ast.SelectorExpr) error {

	importInFile, err := g.handleImport(typeProperties.AstFile.Imports)
	if err != nil {
		return err
	}

	theX := Expression(methodType.X.(*ast.Ident).String())

	interfacePath := string(importInFile[theX].CompletePath)

	internalGi, err := g.traceTypeInPath(interfacePath, methodType.Sel.String())
	if err != nil {
		return err
	}

	gi.CompositionTypes = append(gi.CompositionTypes, internalGi)

	return nil
}

func (g *GogenAnyTypeBuilder) handleIdent(gi *GogenAnyType, unknownTypes map[FieldType]*GogenAnyType, unknownFields map[FieldSignature]*GogenFieldType, typeProperties *TypeProperties, fieldType *ast.Ident, collectedType map[FieldType]*TypeProperties) error {

	newGi := new(GogenAnyType)

	fieldTypeName := FieldType(fieldType.String())

	newGi.GogenFieldType = &GogenFieldType{
		Name:         fieldTypeName,
		Expr:         fieldType,
		DefaultValue: getTypeAsString(fieldType),
		File:         nil,
	}

	gi.CompositionTypes = append(gi.CompositionTypes, newGi)

	unknownFields[NewFieldSignature(gi.GogenFieldType.Name.String(), fieldTypeName.String())] = newGi.GogenFieldType

	if fieldType.Obj == nil {

		//x, exist := collectedType[FieldType(getTypeAsString(fieldType))]
		//if !exist {
		//	LogDebug(1, ">>>>>>>>> 2 masuk sini !exist %v", fieldType.String())
		//	unknownTypes[fieldTypeName] = newGi
		//	return nil
		//}
		//newGi.GogenFieldType.File = x.AstFile
		//newGi.GogenFieldType.Expr = x.TypeSpec.Type

		unknownTypes[fieldTypeName] = newGi

	} else {

		newTypeSpec, ok := fieldType.Obj.Decl.(*ast.TypeSpec)
		if !ok {
			return fmt.Errorf("%s is not type", fieldType.String())
		}

		newGi.GogenFieldType.File = typeProperties.AstFile
		newGi.GogenFieldType.Expr = newTypeSpec.Type

		newTp := TypeProperties{
			AstFile:  typeProperties.AstFile,
			TypeSpec: newTypeSpec,
		}

		err := g.handleGogenType(newGi, unknownTypes, unknownFields, collectedType, &newTp)
		if err != nil {
			return err
		}

	}

	return nil
}

func (g *GogenAnyTypeBuilder) handleInterfaceMethod(gi *GogenAnyType, unknownFields map[FieldSignature]*GogenFieldType, method *ast.Field) error {

	if method.Names == nil && len(method.Names) > 0 {
		err := fmt.Errorf("method must have name")
		return err
	}
	methodName := method.Names[0].String()

	if !ast.IsExported(methodName) {
		return nil
	}

	//LogDebug(4, "as function found method %s", methodName)

	gm := newGogenMethod(methodName)

	gi.Methods = append(gi.Methods, gm)

	methodType, ok := method.Type.(*ast.FuncType)
	if !ok {
		return fmt.Errorf("somehow cannot convert to FuncType")
	}

	g.handleFuncParamResultType(methodType, unknownFields, gm, gi.GogenFieldType.File)

	return nil
}

func (g *GogenAnyTypeBuilder) handleFuncParamResultType(methodType *ast.FuncType, unknownFields map[FieldSignature]*GogenFieldType, gm *GogenMethod, astFile *ast.File) {

	if methodType.Params.NumFields() > 0 {
		for _, param := range methodType.Params.List {

			fieldType := FieldType(getTypeAsString(param.Type))

			if param.Names != nil {
				LogDebug(1, ">>>>>>>>>>>>>> masuk sini fieldType %v", param.Type)
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
