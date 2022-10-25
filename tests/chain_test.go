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
//go:generate bash gencode.sh testserializer

//export GENCODE=TRUE &&
////go:generate go test -run TestGenCode -v

package main

import (
	"context"
	_ "encoding/hex"
	"fmt"
	"github.com/uuosio/chain"
	"github.com/uuosio/chaintester"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var ctx = context.Background()

func OnApply(receiver, firstReceiver, action uint64) {
	println(chain.N2S(receiver), chain.N2S(firstReceiver), chain.N2S(action))
	contract_apply(receiver, firstReceiver, action)
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

	_, err = tester.PushActionEx("hello", "settest", []byte(test), permissions)
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

func CheckAssertError(err error, msg string) {
	_err := err.(*chaintester.TransactionError)
	__err := _err.Json()
	_msg, _ := __err.GetString("action_traces", 0, "except", "stack", 0, "data", "s")
	if strings.Index(_msg, msg) < 0 {
		panic(fmt.Errorf("invalid error: %s %s", _msg, msg))
	}
}

func TestHello(t *testing.T) {
	permissions := `
	{
		"hello": "active"
	}
	`

	tester := initTest("testhello", "tests.abi", true)
	defer tester.FreeChain()

	ret, err := tester.PushAction("hello", "sayhello", "", permissions)
	if err != nil {
		panic(err)
	}
	elapsed, _ := ret.GetString("elapsed")
	t.Logf("++++++++%v", elapsed)
	tester.ProduceBlock()
}

func TestAsset(t *testing.T) {
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
	t.Logf("+++++++%v", ret.ToString())

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
}

func TestMI(t *testing.T) {
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
	t.Logf("+++++:%v", ret.ToString())
}

func TestSingleton(t *testing.T) {
	permissions := `
	{
		"hello": "active"
	}
	`

	tester := initTest("testsingleton", "tests.abi", true)
	defer tester.FreeChain()

	ret, err := tester.PushAction("hello", "sayhello", "", permissions)
	if err != nil {
		panic(err)
	}
	t.Logf("+++++:%v", ret.ToString())
}

func TestDB(t *testing.T) {
	permissions := `
	{
		"hello": "active"
	}
	`

	tester := initTest("testdb", "tests.abi", true)
	defer tester.FreeChain()

	ret, err := tester.PushAction("hello", "sayhello", "", permissions)
	if err != nil {
		panic(err)
	}
	t.Logf("+++++:%v", ret.ToString())
}

func TestSort(t *testing.T) {
	permissions := `
	{
		"hello": "active"
	}
	`

	tester := initTest("testsort", "testsort/test.abi", true)
	defer tester.FreeChain()

	args := `
	{
		"pubs": [
			"EOS6SD6yzqaZhdPHw2LUVmZxWLeWxnp76KLnnBbqP94TsDsjNLosG",
			"EOS4vtCi4jbaVCLVJ9Moenu9j7caHeoNSWgWY65bJgEW8MupWsRMo",
			"EOS82JTja1SbcUjSUCK8SNLLMcMPF8W5fwUYRXmX32obtjsZMW9nx"
		]
	}
	`
	ret, err := tester.PushAction("hello", "test", args, permissions)
	if err != nil {
		panic(err)
	}
	t.Logf("+++++:%v", ret.ToString())
}

func TestUint128(t *testing.T) {
	permissions := `
	{
		"hello": "active"
	}
	`
	tester := initTest("testuint128", "tests.abi", true)
	defer tester.FreeChain()
	ret, err := tester.PushAction("hello", "test", "", permissions)
	if err != nil {
		panic(err)
	}
	t.Logf("+++++:%v", ret.ToString())
}

func TestPackSize(t *testing.T) {
	permissions := `
	{
		"hello": "active"
	}
	`
	tester := initTest("testpacksize", "tests.abi", true)
	defer tester.FreeChain()
	ret, err := tester.PushAction("hello", "test", "", permissions)
	if err != nil {
		panic(err)
	}
	t.Logf("+++++:%v", ret.ToString())
}

func TestPrint(t *testing.T) {
	permissions := `
	{
		"hello": "active"
	}
	`
	tester := initTest("testprint", "tests.abi", true)
	defer tester.FreeChain()
	ret, err := tester.PushAction("hello", "test", "", permissions)
	if err != nil {
		panic(err)
	}
	t.Logf("+++++:%v", ret.ToString())
}

func TestPrimaryKey(t *testing.T) {
	permissions := `
	{
		"hello": "active"
	}
	`
	tester := initTest("testprimarykey", "tests.abi", true)
	defer tester.FreeChain()
	ret, err := tester.PushAction("hello", "test", "", permissions)
	if err != nil {
		panic(err)
	}
	t.Logf("+++++:%v", ret.ToString())
}

func TestMath(t *testing.T) {
	permissions := `
	{
		"hello": "active"
	}
	`
	tester := initTest("testmath", "tests.abi", true)
	defer tester.FreeChain()
	_, err := tester.PushAction("hello", "test", "", permissions)
	CheckAssertError(err, "divide by zero")
	// t.Logf("+++++:%v", ret.ToString())
}

func TestVariant(t *testing.T) {
	permissions := `
	{
		"hello": "active"
	}
	`
	tester := initTest("testvariant", "testvariant/test.abi", true)
	defer tester.FreeChain()

	args := `{
		"v": ["uint64", 123]
	}`
	ret, err := tester.PushAction("hello", "test", args, permissions)
	if err != nil {
		panic(err)
	}
	t.Logf("+++++:%v", ret.ToString())
	//code string, scope string, table string, lower_bound string, upper_bound string, limit int64
	ret, err = tester.GetTableRows(true, "hello", "", "mytable", "", "", 10)
	if err != nil {
		panic(err)
	}
	//['rows'][0]['a'] == ['uint64', 123]
	{
		value, err := ret.GetString("rows", 0, "a")
		if err != nil {
			panic(err)
		}
		if value != `["uint64",123]` {
			panic(fmt.Errorf("invalid value: %v", value))
		}
	}
}

func TestLargeCode(t *testing.T) {
	permissions := `
	{
		"hello": "active"
	}
	`
	tester := initTest("testlargecode", "tests.abi", true)
	defer tester.FreeChain()
	ret, err := tester.PushAction("hello", "test", "", permissions)
	if err != nil {
		panic(err)
	}
	t.Logf("+++++:%v", ret.ToString())
}

func TestPrivileged(t *testing.T) {
	permissions := `
	{
		"hello": "active"
	}
	`
	tester := initTest("testprivileged", "tests.abi", true)
	defer tester.FreeChain()
	return
	ret, err := tester.PushAction("hello", "test", "", permissions)
	if err != nil {
		panic(err)
	}
	t.Logf("+++++:%v", ret.ToString())
}

func TestTransaction(t *testing.T) {
	permissions := `
	{
		"hello": "active"
	}
	`
	tester := initTest("testtransaction", "tests.abi", true)
	defer tester.FreeChain()

	ret, err := tester.PushAction("hello", "sayhello1", "", permissions)
	if err != nil {
		panic(err)
	}
	t.Logf("+++++:%v", ret.ToString())
	// tester.EnableDebugContract("hello", false)
	tester.ProduceBlock()
	tester.ProduceBlock()
	tester.ProduceBlock()

	// tester.EnableDebugContract("hello", true)
	ret, err = tester.PushAction("hello", "sayhello3", "", permissions)
	if err != nil {
		panic(err)
	}
	t.Logf("+++++++++++ret: %v", ret.ToString())
}

func TestAction(t *testing.T) {
	permissions := `
	{
		"hello": "active"
	}
	`
	tester := initTest("testaction", "tests.abi", true)
	defer tester.FreeChain()

	ret, err := tester.PushAction("hello", "sayhello", "", permissions)
	if err != nil {
		panic(err)
	}
	t.Logf("+++++:%v", ret.ToString())
	oldBalance := tester.GetBalance("hello")
	t.Logf("++++++++old Balance: %v", oldBalance)

	// r = self.chain.push_action('hello', 'sayhello3', b'hello,world')
	ret, err = tester.PushAction("hello", "sayhello3", "", permissions)

	newBalance := tester.GetBalance("hello")
	t.Logf("++++++++new balance: %v", newBalance)
	if oldBalance-newBalance != 10000 {
		panic("invalid balance")
	}
}

func TestApplyCtx(t *testing.T) {
	tester := chaintester.NewChainTester()
	tester.GetInfo()
	{
		defer func() {
			err := recover()
			if err == nil {
				panic(err)
			}
			t.Logf("++++%v", err)
		}()
		chain.Check(false, "oops!")
	}
}
