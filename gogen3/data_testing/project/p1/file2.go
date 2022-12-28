package p1

type BeforeTargetSameFileSamePackage interface {
	BeforeTargetSameFileSamePackageMethod()
}

type SomeStruct struct{}

type SaveTodoRepo interface {
	SaveTodo(x int)
}

type MyInterfaceInFile2 interface {
	Other
	MethodOne(x int, y string) (bool, error)
	//BeforeTargetSameFileSamePackage
	//AfterTargetSameFileSamePackage
}

type AfterTargetSameFileSamePackage interface {
	AfterTargetSameFileSamePackageMethod()
}
