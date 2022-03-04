package main

import (
	"github.com/uuosio/chain"
)

func main() {
	a := chain.NewUint128(3, 4)
	b := chain.NewUint128(1, 2)
	a.Add(&a, &b)
	chain.Check(a == chain.NewUint128(4, 6), "a.Add(a, b)")

	a = chain.NewUint128(3, 4)
	b = chain.NewUint128(1, 2)
	a.Sub(&a, &b)
	chain.Check(a == chain.NewUint128(2, 2), "a.Sub(a, b)")

	a = chain.NewUint128(0xffffffffffffffff, 2)
	b = chain.NewUint128(3, 0)
	a.Mul(&a, &b)
	chain.Check(a == chain.NewUint128(uint64(0xfffffffffffffffd), 0x08), "a.Mul(a, b)")

	a = chain.NewUint128(1, 2)
	b = chain.NewUint128(3, 0)
	a.Div(&a, &b)
	chain.Check(a == chain.NewUint128(0xaaaaaaaaaaaaaaab, 0), "a.Div(a, b)")
}
