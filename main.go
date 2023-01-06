package main

import (
	"bytes"
	"fmt"
	"gen/gogen5"
	"time"
)

func main() {

	actGi, err := gogen5.Build("./gogen5/data_testing/project/struct001/p1", "./gogen5/data_testing/project/go.mod", "MyStruct1")
	if err != nil {
		panic(err)
	}

	//gogen5.PrintGogenAnyType(0, actGi)
	//
	//jsonInBytes, err := json.Marshal(actGi)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("%v\n", string(jsonInBytes))

	var param bytes.Buffer
	var level int
	gogen5.WriteTest(actGi, level, "", &param)

	time.Sleep(1 * time.Second)

	fmt.Printf("%v\n", param.String())

}
