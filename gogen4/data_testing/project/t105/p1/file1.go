package p1

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
}

type MyStruct3 struct {
}

type YourInterface interface {
	Method1()
}
