package p1

type MyInterface2 interface {
	Method2()
}

type MyInterface1 interface {
	MyInterface2
	MyInterface3
	MyInterface4
	MyInterface5
	Method1()
}

type MyInterface3 interface {
	Method3()
}
