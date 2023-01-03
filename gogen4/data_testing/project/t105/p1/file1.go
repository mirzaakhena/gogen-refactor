package p1

import "mirza/gogen/refactor/t105/p2"

type MyStruct2 struct {
	Name string
}

type MyStruct1 struct {
	Age int
	MyStruct2
	MyStruct3
	MyStruct4
	MyStruct5
	YourInterface
	*p2.MyStruct6
}

type MyStruct3 struct {
}

type YourInterface interface {
	Method1()
}
