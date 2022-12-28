package gogen4

import (
	"fmt"
	"os"
	"testing"
)

func TestCase1(t *testing.T) {

	gi := NewGogenInterfaceBuilder()
	actGi, err := gi.Build("./data_testing/project/p1", "./data_testing/project/go.mod", "MyInterfaceInFile2")
	if err != nil {
		fmt.Printf("ERROR : %v", err.Error())
		os.Exit(1)
	}

	_ = actGi

	fmt.Printf("\n\n\n")

	for _, m := range PrintAllMethod(actGi) {
		LogDebug(0, "method : %v", m.Name)
		for _, p := range m.Params {
			LogDebug(1, "field : %v %v = %v", p.Name, p.DataType.Name, p.DataType.DefaultValue)
		}
	}

	//jsonInBytes, err := json.MarshalIndent(actGi, "", " ")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%v\n", string(jsonInBytes))

}
