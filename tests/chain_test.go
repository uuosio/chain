//go:generate bash gencode.sh testasset
//go:generate bash gencode.sh testaction
//go:generate bash gencode.sh testasset
//go:generate bash gencode.sh testcrypto
//go:generate bash gencode.sh testtoken
//go:generate bash gencode.sh testdb
//go:generate bash gencode.sh testsort
//go:generate bash gencode.sh testfloat128
//go:generate bash gencode.sh testpacksize
//go:generate bash gencode.sh testsingleton
//go:generate bash gencode.sh testmi
//go:generate bash gencode.sh testprimarykey
//go:generate bash gencode.sh testmath
//go:generate bash gencode.sh testvariant
//go:generate bash gencode.sh testlargecode
//go:generate bash gencode.sh testprivileged
//go:generate bash gencode.sh testtransaction

//export GENCODE=TRUE &&
////go:generate go test -run TestGenCode -v

package main

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/uuosio/chain"
	"github.com/uuosio/chain/tests/testaction"
	"github.com/uuosio/chaintester"
)

var ctx = context.Background()

func OnApply(receiver, firstReceiver, action uint64) {
	println(chain.N2S(receiver), chain.N2S(firstReceiver), chain.N2S(action))
	contract_apply(receiver, firstReceiver, action)
	testaction.ContractApply(receiver, firstReceiver, action)
	println("++++++++apply end!")
}

func init() {
	chaintester.SetApplyFunc(OnApply)
}

func TestGenCode(t *testing.T) {
	cmd := exec.Command("go-contract", "gencode", "-p", "testasset")
	cmd.Dir = "./testasset"
	cmd.Run()

	t.Logf("done!\n")
	if os.Getenv("GENCODE") == "" {
		t.Skip("Skipping TestGenCode")
	}
}

func TestHello(t *testing.T) {
	// t.Errorf("++++++enable_debug: %v", os.Getenv("enable_debug"))
	permissions := `
	{
		"hello": "active"
	}
	`

	tester := chaintester.NewChainTester()
	defer tester.FreeChain()

	tester.EnableDebugContract("hello", true)

	updateAuthArgs := `{
		"account": "hello",
		"permission": "active",
		"parent": "owner",
		"auth": {
			"threshold": 1,
			"keys": [
				{
					"key": "EOS6AjF6hvF7GSuSd4sCgfPKq5uWaXvGM2aQtEUCwmEHygQaqxBSV",
					"weight": 1
				}
			],
			"accounts": [{"permission":{"actor": "hello", "permission": "eosio.code"}, "weight":1}],
			"waits": []
		}
	}`
	tester.PushAction("eosio", "updateauth", updateAuthArgs, permissions)

	err := tester.DeployContract("hello", "tests.wasm", "tests.abi")
	if err != nil {
		panic(err)
	}
	tester.ProduceBlock()
	return

	args := `
	{
		"name": "Go"
	}
	`
	ret, err := tester.PushAction("hello", "sayhello", args, permissions)
	if err != nil {
		panic(err)
	}
	elapsed, _ := ret.GetString("elapsed")
	t.Logf("++++++++%v", elapsed)
	tester.ProduceBlock()
}
