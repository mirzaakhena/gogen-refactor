package gogen3

import (
	"fmt"
	"os"
	"testing"
)

func TestCase1(t *testing.T) {

	actGi, err := NewGogenInterfaceBuilder("./data_testing/project/p1", "./data_testing/project/go.mod", "MyInterfaceInFile2")
	if err != nil {
		fmt.Printf("ERROR : %v", err.Error())
		os.Exit(1)
	}

	_ = actGi

	for _, m := range PrintAllMethod(actGi) {
		fmt.Printf(">>> %v\n", m.Name)
	}

	//jsonInBytes, err := json.MarshalIndent(actGi, "", " ")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%v\n", string(jsonInBytes))

}
