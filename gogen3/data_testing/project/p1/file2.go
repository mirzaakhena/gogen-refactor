package p1

import (
	"mirza/gogen/refactor/p2"
	p3different "mirza/gogen/refactor/p3"
)

type BeforeTargetSameFileSamePackage interface {
	BeforeTargetSameFileSamePackageMethod()
}

type SomeStruct struct{}

type MyInterfaceInFile2 interface {
	MethodOne(x int, y string) (bool, error)
	BeforeTargetSameFileSamePackage
	AfterTargetSameFileSamePackage
	Other
	p2.DiffPackage
	p3different.OtherPackage
	BeforeTargetDiffFileSamePackage
	AfterTargetDiffFileSamePackage
}

type AfterTargetSameFileSamePackage interface {
	AfterTargetSameFileSamePackageMethod()
}
