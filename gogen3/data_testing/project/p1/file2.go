package p1

import (
	"context"
	"github.com/gin-gonic/gin"
)

type BeforeTargetSameFileSamePackage interface {
	BeforeTargetSameFileSamePackageMethod(ctx context.Context, aaa gin.RouterGroup)
}

type SomeStruct struct{}

type SaveTodoRepo interface {
	SaveTodo(x int)
}

type AnAlias Other

type MyInterfaceInFile2 interface {
	//TheOnlyOne(ctx context.Context)
	AnAlias
	MethodOne(x int, y string) (bool, error)
	BeforeTargetSameFileSamePackage
	AfterTargetSameFileSamePackage
}

type AfterTargetSameFileSamePackage interface {
	AfterTargetSameFileSamePackageMethod()
}
