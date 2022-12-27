package gogen3

import (
	"fmt"
	"os"
	"testing"
)

func TestCase1(t *testing.T) {

	actGi, err := NewGogenInterface("./data_testing/project/p1", "./data_testing/project/go.mod", "MyInterfaceInFile2")
	if err != nil {
		fmt.Printf("ERROR : %v", err.Error())
		os.Exit(1)
	}

	//packageName := PackageName("p1")

	//expGi := GogenInterface{
	//	//CurrentPackage: &packageName,
	//	InterfaceType: &GogenFieldType{
	//		Name: "MyInterfaceInFile2",
	//		TypeSpec: &ast.InterfaceType{
	//			Methods: nil,
	//		},
	//		DefaultValue: "nil",
	//		Imports:      nil,
	//	},
	//	Interfaces: []*GogenInterface{},
	//	Methods:    []*GogenMethod{},
	//}

	_ = actGi

	//_ = expGi

	//assert.Equal(t, expGi.InterfaceType, actGi.InterfaceType)

}
