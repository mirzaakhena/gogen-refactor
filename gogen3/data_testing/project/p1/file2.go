package p1

import (
	"context"
)

type BeforeTargetSameFileSamePackage interface {
	BeforeTargetSameFileSamePackageMethod(ctx context.Context, aaa SomeStruct)
}

type SomeStruct int

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
