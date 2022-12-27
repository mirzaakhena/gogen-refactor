package gogen3

import (
	"fmt"
	"go/ast"
)

func NewGogenInterface(packagePath, goModFilePath, interfaceTargetName string) (*GogenInterface, error) {

	gomodProperties, err := handleGoMod(goModFilePath)
	if err != nil {
		return nil, err
	}

	_ = gomodProperties

	gi := new(GogenInterface)

	unknownInterface := map[FieldType]*GogenInterface{}

	collectedType := map[FieldType]*TypeProperties{}

	err = traceType(packagePath, interfaceTargetName, collectedType, func(tp *TypeProperties) error {
		//importInFile, err := handleImport(astFile.Imports, gomodProperties)
		//if err != nil {
		//	return err
		//}

		err = handleGogenInterface(gi, unknownInterface, tp.TypeSpec, tp.File)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	for fieldType, uki := range unknownInterface {

		tp, exist := collectedType[fieldType]
		if !exist {

			selectorExpr, ok := uki.InterfaceType.Expr.(*ast.SelectorExpr)
			if ok {

				importInFile, err := handleImport(uki.InterfaceType.File.Imports, gomodProperties)
				if err != nil {
					return nil, err
				}

				theX := Expression(selectorExpr.X.(*ast.Ident).String())
				gogenImport := importInFile[theX]

				err = traceType(string(gogenImport.CompletePath), selectorExpr.Sel.String(), collectedType, func(tp *TypeProperties) error {

					err = handleGogenInterface(uki, unknownInterface, tp.TypeSpec, tp.File)
					if err != nil {
						return err
					}

					return nil
				})

				if err != nil {
					return nil, err
				}

				delete(unknownInterface, fieldType)

				continue
			}

			return nil, fmt.Errorf("field type %v is not exist anywhere", fieldType)
		}

		err = handleGogenInterface(uki, unknownInterface, tp.TypeSpec, tp.File)
		if err != nil {
			return nil, err
		}

		delete(unknownInterface, fieldType)

	}

	PrintGogenInteface(0, gi)

	return gi, nil
}
