package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gen/gogen"
)

func main() {

	flag.Parse()

	values := flag.Args()

	// "domain_authservice/usecase/runuserlogin"
	// "InportRequest"

	//gs, err := gogen.NewGogenStruct(gogen.GetPackagePath(), values[0], values[1])
	gs, err := gogen.NewGogenInterface(gogen.GetPackagePath(), values[0], values[1])
	if err != nil {
		panic(err)
	}

	_ = gs

	jsonInBytes, err := json.MarshalIndent(gs, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", string(jsonInBytes))

}
