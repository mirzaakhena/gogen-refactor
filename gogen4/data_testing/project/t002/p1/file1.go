package p1

import (
	"context"
	"github.com/gin-gonic/gin"
)

type MyInterface1 interface {
	Method1(x int, y string) (bool, error)
	Method2(ctx3 context.Context, handler gin.HandlerFunc)
	Method3(int, string) (x bool, y error)
}
