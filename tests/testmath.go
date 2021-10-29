package main

import (
	"github.com/uuosio/chain"
)

func main() {
	a, _, _ := chain.GetApplyArgs()
	b := 0
	if a.N == 0 {
		b = 1
	}
	chain.Println(1 / b)
}
