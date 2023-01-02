package p1

import (
	"context"
	"github.com/labstack/echo/v4"
)

type MyStruct1 struct {
	ctx8       context.Context
	theHandler echo.HandlerFunc
}
