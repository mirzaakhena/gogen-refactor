package p1

import (
	"context"
	p3different "mirza/gogen/refactor/p3"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

type Coba struct {
	X gin.HandlerFunc
	Y echo.HandlerFunc
	Z p3different.TheString
	C context.Context
}

type AfterTargetDiffFileSamePackage interface {
	AfterTargetDiffFileSamePackageMethod()
}
