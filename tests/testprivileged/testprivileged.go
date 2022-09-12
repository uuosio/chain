package testprivileged

import (
	"github.com/uuosio/chain"
	"github.com/uuosio/chain/logger"
)

func main() {
	receiver, firstReceiver, action := chain.GetApplyArgs()
	ContractApply(receiver.N, firstReceiver.N, action.N)
}

func ContractApply(receiver, firstReceiver, action uint64) {
	a, b, c := chain.GetResourceLimits(chain.NewName("hello"))
	logger.Println(a, b, c)
}
