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

	gogenInterfaceRoot := new(GogenInterface)

	unknownInterfaces := map[FieldType]*GogenInterface{}

	collectedType := map[FieldType]*TypeProperties{}

	err = traceType(packagePath, interfaceTargetName, collectedType, func(typeSpec *ast.TypeSpec, astFile *ast.File) error {

		err = handleGogenInterface(gogenInterfaceRoot, unknownInterfaces, typeSpec, astFile)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// solve all unknown interface
	for fieldType, unknownInterface := range unknownInterfaces {

		tp, exist := collectedType[fieldType]
		if !exist {

			selectorExpr, ok := unknownInterface.InterfaceType.Expr.(*ast.SelectorExpr)
			if ok {

				importInFile, err := handleImport(unknownInterface.InterfaceType.File.Imports, gomodProperties)
				if err != nil {
					return nil, err
				}

				theX := Expression(selectorExpr.X.(*ast.Ident).String())
				gogenImport := importInFile[theX]

				err = traceType(string(gogenImport.CompletePath), selectorExpr.Sel.String(), collectedType, func(typeSpec *ast.TypeSpec, astFile *ast.File) error {

					err := handleGogenInterface(unknownInterface, unknownInterfaces, typeSpec, astFile)
					if err != nil {
						return err
					}

					return nil
				})

				if err != nil {
					return nil, err
				}

				delete(unknownInterfaces, fieldType)

				continue
			}

			return nil, fmt.Errorf("field type %v is not exist anywhere", fieldType)
		}

		err = handleGogenInterface(unknownInterface, unknownInterfaces, tp.TypeSpec, tp.File)
		if err != nil {
			return nil, err
		}

		delete(unknownInterfaces, fieldType)

	}

	PrintGogenInteface(0, gogenInterfaceRoot)

	return gogenInterfaceRoot, nil
}
