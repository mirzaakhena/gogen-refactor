package p1

type MyInterface2 interface {
}

type TypeAlias MyInterface2

type MyInterface1 interface {
	TypeAlias
}
