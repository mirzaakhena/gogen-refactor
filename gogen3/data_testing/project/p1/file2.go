package p1

type BeforeTargetSameFileSamePackage interface {
	BeforeTargetSameFileSamePackageMethod()
}

type SomeStruct struct{}

type SaveTodoRepo interface {
	SaveTodo(x int)
}

type AnAlias Other

type MyInterfaceInFile2 interface {
	AnAlias
	MethodOne(x int, y string) (bool, error)
	BeforeTargetSameFileSamePackage
	AfterTargetSameFileSamePackage
}

type AfterTargetSameFileSamePackage interface {
	AfterTargetSameFileSamePackageMethod()
}
