package main

import "fmt"

func main() {

	a := "mirza/gogen/refactor"
	b := "mirza/gogen/refactor/t04/p2"

	fmt.Printf("%v", b[len(a):])

}
