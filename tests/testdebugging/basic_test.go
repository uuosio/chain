package main

import (
	"context"
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
	tester.PushAction(ctx, 0, "hello", "sayhello", args, permissions)
}
