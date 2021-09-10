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
	b := &B{}
	var i interface{} = b
	c := i.(int)
	chain.Printi(int64(c))
	// chain.PrintUi(chain.NewName("hello").N)
}
