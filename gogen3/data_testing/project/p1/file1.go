package p1

type BeforeTargetDiffFileSamePackage interface {
	BeforeTargetDiffFileSamePackageMethod()
}

type Other interface {
	BeforeTargetDiffFileSamePackage
	AfterTargetDiffFileSamePackage
}
