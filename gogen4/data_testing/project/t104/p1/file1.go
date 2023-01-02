package p1

import p1 "mirza/gogen/refactor/t104/p2"

type MyStruct2 struct {
}

type MyStruct1 struct {
	A01 MyStruct2
	A02 MyStruct3
	A03 p1.MyStruct4
	A04 MyStruct5
	A05 MyStruct6
}

type MyStruct3 struct {
}
