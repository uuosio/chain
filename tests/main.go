package main

import (
	"github.com/uuosio/chain"
	"github.com/uuosio/chain/database"

	"github.com/uuosio/chain/tests/testaction"
	"github.com/uuosio/chain/tests/testasset"
	"github.com/uuosio/chain/tests/testcrypto"
	"github.com/uuosio/chain/tests/testdb"
	"github.com/uuosio/chain/tests/testfloat128"
	"github.com/uuosio/chain/tests/testlargecode"
	"github.com/uuosio/chain/tests/testmath"
	"github.com/uuosio/chain/tests/testmi"
	"github.com/uuosio/chain/tests/testpacksize"
	"github.com/uuosio/chain/tests/testprimarykey"
	"github.com/uuosio/chain/tests/testprint"
	"github.com/uuosio/chain/tests/testprivileged"
	"github.com/uuosio/chain/tests/testsingleton"
	"github.com/uuosio/chain/tests/testsort"
	"github.com/uuosio/chain/tests/testtoken"
	"github.com/uuosio/chain/tests/testuint128"
	"github.com/uuosio/chain/tests/testvariant"
)

//contract tests
type MyContract struct {
	Receiver      chain.Name
	FirstReceiver chain.Name
	Action        chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *MyContract {
	testaction.NewContract(receiver, firstReceiver, action)
	return &MyContract{receiver, firstReceiver, action}
}

func DJBH(cp string) uint64 {
	hash := uint64(5381)
	for _, c := range []byte(cp) {
		hash = 33*hash ^ uint64(c)
	}
	return hash
}

func TEST_ACTION(CLASS string, METHOD string) uint64 {
	return (DJBH(CLASS) << 32) | DJBH(METHOD)
}

func main() {
	receiver, firstReceiver, action := chain.GetApplyArgs()
	contract_apply(receiver.N, firstReceiver.N, action.N)
}

func contract_apply(_receiver, _firstReceiver, _action uint64) {
	receiver := chain.Name{_receiver}
	// firstReceiver := chain.Name{_firstReceiver}
	action := chain.Name{_action}

	table := database.NewTableI64(receiver, receiver, chain.NewName("curtest"), nil)
	keyCurTest := chain.NewName("curtest").N
	if action == chain.NewName("settest") {
		data := chain.ReadActionData()
		table.Set(keyCurTest, data, receiver)
		return
	}

	curTest := "none"
	it := table.Find(keyCurTest)
	chain.Check(it.IsOk(), "bad value")
	raw := table.GetByIterator(it)
	curTest = string(raw)

	if curTest == "testaction" {
		testaction.ContractApply(_receiver, _firstReceiver, _action)
	} else if curTest == "testasset" {
		testasset.ContractApply(_receiver, _firstReceiver, _action)
	} else if curTest == "testcrypto" {
		testcrypto.ContractApply(_receiver, _firstReceiver, _action)
	} else if curTest == "testtoken" {
		testtoken.ContractApply(_receiver, _firstReceiver, _action)
	} else if curTest == "testfloat128" {
		testfloat128.ContractApply(_receiver, _firstReceiver, _action)
	} else if curTest == "testdb" {
		testdb.ContractApply(_receiver, _firstReceiver, _action)
	} else if curTest == "testsort" {
		testsort.ContractApply(_receiver, _firstReceiver, _action)
	} else if curTest == "testuint128" {
		testuint128.ContractApply(_receiver, _firstReceiver, _action)
	} else if curTest == "testpacksize" {
		testpacksize.ContractApply(_receiver, _firstReceiver, _action)
	} else if curTest == "testsingleton" {
		testsingleton.ContractApply(_receiver, _firstReceiver, _action)
	} else if curTest == "testprint" {
		testprint.ContractApply(_receiver, _firstReceiver, _action)
	} else if curTest == "testmi" {
		testmi.ContractApply(_receiver, _firstReceiver, _action)
	} else if curTest == "testprimarykey" {
		testprimarykey.ContractApply(_receiver, _firstReceiver, _action)
	} else if curTest == "testmath" {
		testmath.ContractApply(_receiver, _firstReceiver, _action)
	} else if curTest == "testvariant" {
		testvariant.ContractApply(_receiver, _firstReceiver, _action)
	} else if curTest == "testlargecode" {
		testlargecode.ContractApply(_receiver, _firstReceiver, _action)
	} else if curTest == "testprivileged" {
		testprivileged.ContractApply(_receiver, _firstReceiver, _action)
	}
}
