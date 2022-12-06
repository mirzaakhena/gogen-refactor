package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gen/gogen"
	"gen/gogen3"
)

func main() {

	flag.Parse()

	values := flag.Args()

	// "domain_authservice/usecase/runuserlogin"
	// "InportRequest"

	// "/usr/local/go/src/time/time.go"
	// "Time"

	gs, err := gogen3.NewGogenStructBuilder(gogen.GetPackagePath(), values[0]).Build(values[1])
	if err != nil {
		panic(err)
	}

	//gs, err := gogen2.NewGogenInterfaceBuilder(gogen.GetPackagePath(), values[0]).Build(values[1])
	//if err != nil {
	//	panic(err)
	//}

	_ = gs

	jsonInBytes, err := json.MarshalIndent(gs, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", string(jsonInBytes))

}
