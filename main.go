package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gen/gogen"
	"gen/gogen2"
	"golang.org/x/mod/modfile"
	"os"
)

func main() {

	flag.Parse()

	values := flag.Args()

	// "domain_authservice/usecase/runuserlogin"
	// "InportRequest"

	// "/usr/local/go/src/time/time.go"
	// "Time"

	gs, err := gogen2.NewGogenStructBuilder(gogen.GetPackagePath(), values[0]).Build(values[1])
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

func ReadGoMod() string {

	gomodfile := "go.mod"

	fileInBytes, err := os.ReadFile(gomodfile)
	if err != nil {
		panic(err)
	}

	parseGoMod, err := modfile.Parse(gomodfile, fileInBytes, nil)
	if err != nil {
		panic(err)
	}

	return parseGoMod.Module.Mod.String()
}
