package p2

import p3different "mirza/gogen/refactor/p3"

type DiffPackage interface {
	p3different.OtherPackage
	MethodThree(x int, y string) (bool, error)
	//MyMethod()
}
