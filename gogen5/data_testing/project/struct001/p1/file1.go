package p1

import (
	"mirza/gogen/refactor/struct001/p2"
	"mirza/gogen/refactor/struct001/p3"
	"sync"
)

type MyStruct3 struct {
	Field1 bool
}

type MyStruct1 struct {
	MyStruct3
	p3.MyStruct4
	*p3.MyStruct5
	Field1 int
	Field2 sync.WaitGroup
	Field3 p2.MyStruct2
	Field4 p2.MyAliasBool
	Field5 []string
	Field6 []*p2.MyStruct2
	Field7 p2.MyInterface1
	Field8 struct {
		x int
		y string
	}
	Field9  func(string) int
	Field10 map[string]int
}
