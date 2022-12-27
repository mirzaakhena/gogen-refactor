package p1

import "mirza/gogen/refactor/p2"

type BeforeTargetSameFileSamePackage interface {
	BeforeTargetSameFileSamePackageMethod()
}

type SomeStruct struct{}

type MyInterfaceInFile2 interface {
	BeforeTargetSameFileSamePackage
	AfterTargetSameFileSamePackage
	Other
	p2.DiffPackage
	//BeforeTargetDiffFileSamePackage
	//AfterTargetDiffFileSamePackage
}

type AfterTargetSameFileSamePackage interface {
	AfterTargetSameFileSamePackageMethod()
}
