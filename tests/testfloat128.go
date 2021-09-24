package main

import (
	"github.com/uuosio/chain"
)

func IsEqual(a, b *chain.Float128) bool {
	minValue := chain.NewFloat128(0.000000001)
	sub := chain.NewFloat128(0.0)
	sub.Sub(a, b)
	if minValue.Cmp(sub.Abs(sub)) > 0 {
		return true
	}
	return false
}

func main() {
	a := chain.NewFloat128(100.11)
	b := chain.NewFloat128(3.11)
	a.Add(a, b)
	chain.Check(IsEqual(a, chain.NewFloat128(103.22)), "a.Add(a, b)")

	a = chain.NewFloat128(100.11)
	b = chain.NewFloat128(3.11)
	a.Sub(a, b)
	chain.Check(IsEqual(a, chain.NewFloat128(97.0)), "a.Sub(a, b)")

	a = chain.NewFloat128(100.11)
	b = chain.NewFloat128(3.11)
	a.Mul(a, b)
	chain.Check(IsEqual(a, chain.NewFloat128(311.34209999999996)), "a.Mul(a, b)")

	a = chain.NewFloat128(100.11)
	b = chain.NewFloat128(3.11)
	a.Div(a, b)
	chain.Check(IsEqual(a, chain.NewFloat128(32.18971061093248)), "a.Div(a, b)")
}
