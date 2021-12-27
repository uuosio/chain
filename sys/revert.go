package sys

import (
	"github.com/uuosio/chain"
	"github.com/uuosio/chain/database"
)

type RevertInterface interface {
	OnRevert(msg string)
}

var gRevert RevertInterface

func Revert(msg string) {
	if gRevert != nil {
		gRevert.OnRevert(msg)
	}
	database.GetStateManager().Revert()
	receiver := chain.CurrentReceiver()
	chain.NewAction(
		chain.NewPermission(receiver, chain.ActiveName),
		receiver,
		chain.NewName("revert"),
		msg,
	).Send()
	chain.Exit()
}

func Init(onRevert RevertInterface) {
	chain.EnableRevert(true)
	gRevert = onRevert
	chain.SetRevertFn(Revert)
	//		runtime.SetRevertOnPanicFn(Revert)
}
