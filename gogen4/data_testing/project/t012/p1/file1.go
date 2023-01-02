package p1

import (
	"github.com/gin-gonic/gin"
)

type MyStruct1 struct {
}

type MyInterface1 interface {
	Method1(handler gin.HandlerFunc, a01 MyStruct1)
}
