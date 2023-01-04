package p1

import (
	"context"
	"github.com/gin-gonic/gin"
	"mirza/gogen/refactor/interface001/p2"
)

type MyInterface2 interface {
	Method21()
}

type MyInterface1 interface {
	MyInterface2
	MyInterface3
	MyInterface4
	MyInterface5
	p2.MyInterface6
	Method11(x int, y string) (bool, error)
	Method12(int, string) (x bool, y error)
	Method13(ctx3 context.Context, handler gin.HandlerFunc)
}

type MyInterface3 interface {
	Method31()
}
