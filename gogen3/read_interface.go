package gogen3

func NewGogenInterfaceBuilder(packagePath, goModFilePath, interfaceTargetName string) (*GogenInterface, error) {

	gomodProperties, err := handleGoMod(goModFilePath)
	if err != nil {
		return nil, err
	}

	gogenInterfaceRoot, err := traceType(packagePath, gomodProperties, interfaceTargetName)
	if err != nil {
		LogDebug(1, ">>>>>>>>> masuk sini")
		return nil, err
	}

	PrintGogenInterface(0, gogenInterfaceRoot)

	return gogenInterfaceRoot, nil
}
