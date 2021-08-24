package main

import (
	"chain"
	"chain/logger"
	"unsafe"
)

func main() {
	chain.Prints("hello")
	chain.PrintHex([]byte("hello"))
	chain.Println()

	// logger.Printf("%d %s %v\n", 1, "hello", "world")

	n := int64(-1)
	nn := *(*uint64)(unsafe.Pointer(&n))
	bs := *(*[8]byte)(unsafe.Pointer(&n))
	logger.Println(bs[:], "hello", []byte("hello"), -1, int64(-1), nn, float64(-1.1), float32(-1.1))
}
