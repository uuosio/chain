package testmath

import (
	"github.com/uuosio/chain"
)

func main() {
	receiver, firstReceiver, action := chain.GetApplyArgs()
	ContractApply(receiver.N, firstReceiver.N, action.N)
}

func ContractApply(receiver, firstReceiver, action uint64) {
	a := chain.Name{receiver}
	b := 0
	if a.N == 0 {
		b = 1
	}
	chain.Println(1 / b)
}
