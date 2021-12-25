package revert

import (
	"runtime"

	"github.com/uuosio/chain"
	"github.com/uuosio/chain/database"
)

func Revert(msg string) {
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

func Init() {
	database.SetSaveState(true)
	runtime.SetRevertFunction(Revert)
}
