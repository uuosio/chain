package main

import (
	"github.com/uuosio/chain"
)

type B struct {
	a uint32
	b uint64
	c uint64
}

func main() {
	a := chain.Float128{}
	b := chain.Float128{}
	a.Set(100.11)
	b.Set(3.11)
	a.Add(&b)
	chain.Println(a)
	a.Mul(&b)
	chain.Println(a)
}
