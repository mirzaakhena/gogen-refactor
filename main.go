package main

import (
	"encoding/json"
	"fmt"
	"gen/gogen"
)

func main() {

	gs, err := gogen.NewGogenStruct(gogen.GetPackagePath(), "domain_authservice/usecase/runuserlogin", "InportRequest")
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
