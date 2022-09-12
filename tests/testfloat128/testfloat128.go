package testfloat128

import (
	"github.com/uuosio/chain"
)

func IsEqual(a, b chain.Float128) bool {
	minValue := chain.NewFloat128(0.000000001)
	sub := chain.NewFloat128(0.0)
	sub.SubEx(&a, &b)
	if minValue.Cmp(sub.Abs(&sub)) > 0 {
		return true
	}
	return false
}

func main() {
	receiver, firstReceiver, action := chain.GetApplyArgs()
	ContractApply(receiver.N, firstReceiver.N, action.N)
}

func ContractApply(receiver, firstReceiver, action uint64) {
	a := chain.NewFloat128(100.11)
	b := chain.NewFloat128(3.11)
	a.Add(&b)
	chain.Check(IsEqual(a, chain.NewFloat128(103.22)), "a.Add(a, b)")

	a = chain.NewFloat128(100.11)
	b = chain.NewFloat128(3.11)
	a.Sub(&b)
	chain.Check(IsEqual(a, chain.NewFloat128(97.0)), "a.Sub(a, b)")

	a = chain.NewFloat128(100.11)
	b = chain.NewFloat128(3.11)
	a.Mul(&b)
	chain.Check(IsEqual(a, chain.NewFloat128(311.34209999999996)), "a.Mul(a, b)")

	a = chain.NewFloat128(100.11)
	b = chain.NewFloat128(3.11)
	a.Div(&b)
	chain.Check(IsEqual(a, chain.NewFloat128(32.18971061093248)), "a.Div(a, b)")
}
