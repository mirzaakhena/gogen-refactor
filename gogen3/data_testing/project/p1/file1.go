package p1

import (
	"mirza/gogen/refactor/p2"
)

type BeforeTargetDiffFileSamePackage interface {
	BeforeTargetDiffFileSamePackageMethod()
}

type Other interface {
	p2.DiffPackage
	MethodTwo(x int, y string) (bool, error)
	//BeforeTargetDiffFileSamePackage
	//AfterTargetDiffFileSamePackage
}
