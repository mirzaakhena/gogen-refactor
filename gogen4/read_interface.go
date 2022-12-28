package gogen4

type GogenInterfaceBuilder struct {
	GoModProperties GoModProperties
	//UnknownInterfaces map[FieldType]*GogenInterface
	//UnknownFields map[FieldMethodSignature]*GogenField
	//CollectedType map[FieldType]*TypeProperties
}

func NewGogenInterfaceBuilder() *GogenInterfaceBuilder {
	return &GogenInterfaceBuilder{
		GoModProperties: NewGoModProperties(),
		//UnknownInterfaces: map[FieldType]*GogenInterface{},
		//UnknownFields: map[FieldMethodSignature]*GogenField{},
		//CollectedType: map[FieldType]*TypeProperties{},
	}
}

func (g *GogenInterfaceBuilder) Build(packagePath, goModFilePath, interfaceTargetName string) (*GogenInterface, error) {

	err := g.handleGoMod(goModFilePath)
	if err != nil {
		return nil, err
	}

	gogenInterfaceRoot, err := g.traceInterfaceType(packagePath, interfaceTargetName)
	if err != nil {
		return nil, err
	}

	PrintGogenInterface(0, gogenInterfaceRoot)

	return gogenInterfaceRoot, nil
}
