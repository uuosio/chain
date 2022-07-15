package main

import (
	"context"
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
	t.Errorf("++++++enable_debug: %v", os.Getenv("enable_debug"))
	SayHelloFromCpp()
	tester := chaintester.NewChainTester()
	args := `
	{
		"name": "rust"
	}
	`
	permissions := `
	{
		"hello": "active"
	}
	`

	fmt.Println("+++++++push sayhello action")
	ret := tester.PushAction("hello", "sayhello", args, permissions)
	t.Errorf("++++++++++ret:%s\n", string(ret))
}
