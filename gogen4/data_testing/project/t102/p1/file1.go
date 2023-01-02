package p1

type MyStruct1 struct {
	Name       string
	Age        int
	Hobbies    []string
	Wife       *MyStruct1
	ThisStruct struct {
		x int
		y float64
	}
}
