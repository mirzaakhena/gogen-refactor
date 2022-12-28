package gogen3

func NewGogenInterfaceBuilder(packagePath, goModFilePath, interfaceTargetName string) (*GogenInterface, error) {

	gomodProperties, err := handleGoMod(goModFilePath)
	if err != nil {
		return nil, err
	}

	gogenInterfaceRoot, err := traceInterfaceType(packagePath, gomodProperties, interfaceTargetName)
	if err != nil {
		return nil, err
	}

	//PrintGogenInterface(0, gogenInterfaceRoot)

	return gogenInterfaceRoot, nil
}
