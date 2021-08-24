package main

func call(a [1]uint64) {

}

func test(n uint64) [1]uint64 {
	a := [1]uint64{n}
	// logger.Println(a)
	call(a)
	return a
}

func test2() []uint64 {
	a := make([]uint64, 1)
	if a[0] != 1 {
	}
	return a
}

func main() {
	test(1122)
	test2()
}
