package gogen4

import (
	"fmt"
	"go/ast"
)

type GogenStructBuilder struct {
	GogenBase
}

func NewGogenStructBuilder() *GogenStructBuilder {
	return &GogenStructBuilder{
		GogenBase: GogenBase{GoModProperties: NewGoModProperties()},
		//UnknownInterfaces: map[GogenFieldType]*GogenAnyType{},
		//UnknownFields: map[FieldSignature]*GogenField{},
		//CollectedType: map[GogenFieldType]*TypeProperties{},
	}
}

func (g *GogenStructBuilder) Build(packagePath, goModFilePath, structTargetName string) (*GogenAnyType, error) {

	err := g.handleGoMod(goModFilePath)
	if err != nil {
		return nil, err
	}

	gogenStructRoot, err := g.traceStructType(packagePath, structTargetName)
	if err != nil {
		return nil, err
	}

	PrintGogenAnyType(1, gogenStructRoot)

	return gogenStructRoot, nil
}

func (g *GogenStructBuilder) traceStructType(packagePath string, structTargetName string) (*GogenAnyType, error) {

	gogenStructTarget := NewGogenAnyType()

	unknownStructs := map[FieldType]*GogenAnyType{}

	unknownFields := map[FieldSignature]*GogenFieldType{} // we use map because we want to take advantage of passing by value

	collectedType := map[FieldType]*TypeProperties{}

	err := traceGeneralType(packagePath, structTargetName, collectedType, func(tp TypeProperties) error {
		return g.handleGogenStruct(gogenStructTarget, unknownStructs, unknownFields, collectedType, &tp)
	})
	if err != nil {
		return nil, err
	}

	// solve all unknown structs
	for fieldType, theStruct := range unknownStructs {

		typeProperties, exist := collectedType[fieldType]
		if !exist {
			return nil, fmt.Errorf("field type %v is not exist anywhere", fieldType)
		}

		LogDebug(1, ">>>>>>>>> 2 masuk sini unknownStruct Loop %v", theStruct.GogenFieldType.Name.String())

		err = g.handleGogenStruct(theStruct, unknownStructs, unknownFields, collectedType, typeProperties)
		if err != nil {
			return nil, err
		}

		delete(unknownStructs, fieldType)

	}

	// solve all unknown fields
	for _, uf := range unknownFields {
		uf.DefaultValue, err = g.handleDefaultValue(uf, uf.Expr, uf.File, collectedType)
		if err != nil {
			return nil, err
		}
	}

	return gogenStructTarget, nil
}

func (g *GogenStructBuilder) handleGogenStruct(gi *GogenAnyType, unknownStructs map[FieldType]*GogenAnyType, unknownFields map[FieldSignature]*GogenFieldType, collectedType map[FieldType]*TypeProperties, typeProperties *TypeProperties) error {

	switch ts := typeProperties.TypeSpec.Type.(type) {

	case *ast.StructType:

		typeSpecName := FieldType(typeProperties.TypeSpec.Name.String())

		gi.GogenFieldType = &GogenFieldType{
			Name:         typeSpecName,
			Expr:         typeProperties.TypeSpec.Type,
			DefaultValue: fmt.Sprintf("%s{}", typeSpecName),
			File:         typeProperties.AstFile,
		}

		for _, field := range ts.Fields.List {

			if field.Names != nil {

				for _, n := range field.Names {
					gf := NewGogenField(n.String(), field.Type, FieldType(getTypeAsString(field.Type)), typeProperties.AstFile)
					gi.Fields = append(gi.Fields, gf)
					unknownFields[NewFieldSignature(gi.GogenFieldType.Name.String(), gf.Name.String())] = gf.DataType
				}

			} else {

				// has no field name

				switch fieldType := field.Type.(type) {

				case *ast.Ident:
					// TODO handle ident
					LogDebug(1, "===> as ident %v %v", fieldType, field.Type)

					newGi := new(GogenAnyType)

					fieldTypeName := FieldType(fieldType.String())

					newGi.GogenFieldType = &GogenFieldType{
						Name:         fieldTypeName,
						Expr:         field.Type,
						DefaultValue: "",
						File:         nil,
					}

					gi.CompositionTypes = append(gi.CompositionTypes, newGi)

					unknownFields[NewFieldSignature(typeSpecName.String(), fieldTypeName.String())] = newGi.GogenFieldType

					// we still don't know what it is

					if fieldType.Obj == nil {

						x, exist := collectedType[FieldType(getTypeAsString(field.Type))]
						if !exist {
							LogDebug(1, ">>>>>>>>> 2 masuk sini !exist %v", fieldType.String())
							unknownStructs[fieldTypeName] = newGi
							continue
						}
						newGi.GogenFieldType.File = x.AstFile
						newGi.GogenFieldType.Expr = x.TypeSpec.Type

						continue
					}

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

					err := g.handleGogenStruct(newGi, unknownStructs, unknownFields, collectedType, &newTp)
					if err != nil {
						return err
					}

				case *ast.StarExpr:
					// TODO handle star
					LogDebug(1, "===> as star %v", fieldType)

				case *ast.SelectorExpr:
					// TODO handle selector
					LogDebug(1, "===> as selector %v", fieldType)
				}

				//gf := NewGogenField(getSel(field.Type), field.Type, fieldType, typeProperties.AstFile)
				//gi.Fields = append(gi.Fields, gf)
				//unknownFields[NewFieldSignature(string(gi.GogenFieldType.Name), gf.Name)] = gf

			}

		}

	case *ast.Ident:
	// TODO type alias will goes here

	case *ast.InterfaceType:
		// TODO type alias will goes here
		LogDebug(1, ">> this is interface")

	default:
		return fmt.Errorf("this is not an struct but %T", ts)

	}

	return nil
}
