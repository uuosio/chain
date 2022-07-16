package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/learnforpractice/chaintester"
	"github.com/uuosio/chain"
)

var ctx = context.Background()

func OnApply(receiver, firstReceiver, action uint64) {
	chaintester.GetVMAPI() // connect to vm api server
	contract_apply(chain.Name{receiver}, chain.Name{firstReceiver}, chain.Name{action})
}

func init() {
	chaintester.SetApplyFunc(OnApply)
}

func TestHello(t *testing.T) {
	// t.Errorf("++++++enable_debug: %v", os.Getenv("enable_debug"))
	permissions := `
	{
		"hello": "active"
	}
	`
	SayHelloFromCpp()

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

	wasm, _ := os.ReadFile("test.wasm")
	abi, _ := os.ReadFile("test.abi")

	hexWasm := make([]byte, len(wasm)*2)
	hex.Encode(hexWasm, wasm)
	setCodeArgs := fmt.Sprintf(
		`
		{
			"account": "%s",
			"vmtype": 0,
			"vmversion": 0,
			"code": "%s"
		 }
		`,
		"hello",
		string(hexWasm),
	)

	ret := tester.PushAction("eosio", "setcode", setCodeArgs, permissions)
	tester.ProduceBlock()

	rawAbi, _ := tester.PackAbi(string(abi))
	hexRawAbi := make([]byte, len(rawAbi)*2)
	hex.Encode(hexRawAbi, rawAbi)
	setAbiArgs := fmt.Sprintf(
		`
		{
			"account": "%s",
			"abi": "%s"
		 }
		`,
		"hello",
		string(hexRawAbi),
	)
	ret = tester.PushAction("eosio", "setabi", setAbiArgs, permissions)

	args := `
	{
		"name": "rust"
	}
	`
	permissions2 := string(permissions)
	for i := 0; i < 1; i++ {
		ret = tester.PushAction("hello", "sayhello", args, permissions)
		tester.PackAbi(string(abi))
		tester.ProduceBlock()
	}

	// defer chaintester.GetApplyRequestServer().Stop()
	// defer chaintester.CloseVMAPI()

	ret = tester.PushAction("hello", "sayhello", args, permissions2)
	tester.ProduceBlock()

	value := &chaintester.JsonValue{}
	err := json.Unmarshal(ret, value)
	if err != nil {
		panic(err)
	}

	id, _ := value.GetString("receipt", "cpu_usage_us")
	// t.Logf("++++++++++ret:%s\n", string(ret))
	t.Logf("++++++++++id:%s\n", id)
}
