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
	"encoding/hex"
	"fmt"
	"github.com/uuosio/chain"
	"github.com/uuosio/chaintester"
	"os"
	"os/exec"
	"testing"
)

var ctx = context.Background()

func OnApply(receiver, firstReceiver, action uint64) {
	println(chain.N2S(receiver), chain.N2S(firstReceiver), chain.N2S(action))
	contract_apply(receiver, firstReceiver, action)
	println("++++++++apply end!")
}

func init() {
	chaintester.SetApplyFunc(OnApply)
}

func initTest(test string, abi string, debug bool) *chaintester.ChainTester {
	tester := chaintester.NewChainTester()

	tester.EnableDebugContract("hello", debug)

	err := tester.DeployContract("hello", "tests.wasm", abi)
	if err != nil {
		panic(err)
	}
	tester.ProduceBlock()

	permissions := `
	{
		"hello": "active"
	}
	`
	_, err = tester.PushAction("hello", "settest", hex.EncodeToString([]byte(test)), permissions)
	if err != nil {
		panic(err)
	}
	return tester
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

	ret, err := tester.PushAction("hello", "settest", hex.EncodeToString([]byte("testhello")), permissions)
	if err != nil {
		panic(err)
	}

	ret, err = tester.PushAction("hello", "sayhello", "", permissions)
	if err != nil {
		panic(err)
	}
	elapsed, _ := ret.GetString("elapsed")
	t.Logf("++++++++%v", elapsed)
	tester.ProduceBlock()
}

func CheckAssertError(err error, msg string) {
	_err := err.(*chaintester.TransactionError)
	__err := _err.Json()
	_msg, _ := __err.GetString("action_traces", 0, "except", "stack", 0, "data", "s")
	if _msg != msg {
		panic(fmt.Errorf("invalid error: %s %s", _msg, msg))
	}
}

func TestAsset(t *testing.T) {
	// t.Errorf("++++++enable_debug: %v", os.Getenv("enable_debug"))
	permissions := `
	{
		"hello": "active"
	}
	`

	tester := initTest("testasset", "tests.abi", true)
	defer tester.FreeChain()
	var err error

	//	var ret *chaintester.JsonValue
	_, err = tester.PushAction("hello", "test1", "", permissions)
	CheckAssertError(err, "addition overflow")

	_, err = tester.PushAction("hello", "test2", "", permissions)
	CheckAssertError(err, "subtraction underflow")

	_, err = tester.PushAction("hello", "test3", "", permissions)
	CheckAssertError(err, "magnitude of asset amount must be less than 2^62")

	_, err = tester.PushAction("hello", "test4", "", permissions)
	CheckAssertError(err, "divide by zero")

	_, err = tester.PushAction("hello", "test5", "", permissions)
	CheckAssertError(err, "divide by negative value")

	_, err = tester.PushAction("hello", "test11", "", permissions)
	CheckAssertError(err, "bad symbol")

	_, err = tester.PushAction("hello", "test12", "", permissions)
	CheckAssertError(err, "multiplication overflow")

	_, err = tester.PushAction("hello", "test13", "", permissions)
	CheckAssertError(err, "multiplication underflow")
}

func TestCrypto(t *testing.T) {
	// t.Errorf("++++++enable_debug: %v", os.Getenv("enable_debug"))
	permissions := `
	{
		"hello": "active"
	}
	`
	tester := initTest("testcrypto", "testcrypto/test.abi", true)
	defer tester.FreeChain()
	var err error

	var ret *chaintester.JsonValue
	ret, err = tester.PushAction("hello", "testhash", "", permissions)
	if err != nil {
		panic(err)
	}
	t.Logf("+++++++%v", ret)

	// data = b'hello,world'
	args := `{
		"data": "68656c6c6f2c776f726c64",
		"sig": "SIG_K1_KiXXExwMGG5NvAngS3X58fXVVcnmPc7fxgwLQAbbkSDj9gwcxWHxHwgpUegSCfgp4nFMMgjLDAKSQWZ2NLEmcJJn1m2UUg",
		"pub": "EOS7wy4M8ZTYqtoghhDRtE37yRoSNGc6zC2zFgdVmaQnKV5ZXe4kV"
	}`
	_, err = tester.PushAction("hello", "testrecover", args, permissions)
	if err != nil {
		panic(err)
	}

	// r = self.chain.push_action('hello', 'testrecover', args)
	// print_console(r)
}

func TestMI(t *testing.T) {
	// t.Errorf("++++++enable_debug: %v", os.Getenv("enable_debug"))
	permissions := `
	{
		"hello": "active"
	}
	`

	tester := initTest("testmi", "tests.abi", true)
	defer tester.FreeChain()

	ret, err := tester.PushAction("hello", "test1", "", permissions)
	if err != nil {
		panic(err)
	}
	t.Logf("+++++:%v", ret)
}
