package testasset

import (
	"bytes"

	"github.com/uuosio/chain"
)

const MAX_AMOUNT = (1 << 62) - 1

//contract hello
type MyContract struct {
}

func NewContract(receiver, firstReceiver, action chain.Name) *MyContract {
	return &MyContract{}
}

//action test1
func (c *MyContract) Test1() {
	a := chain.NewAsset(MAX_AMOUNT, chain.NewSymbol("EOS", 4))
	b := chain.NewAsset(1, chain.NewSymbol("EOS", 4))
	a.Add(b)
	println(a.Amount)
}

//action test2
func (c *MyContract) Test2() {
	a := chain.NewAsset(-MAX_AMOUNT, chain.NewSymbol("EOS", 4))
	b := chain.NewAsset(1, chain.NewSymbol("EOS", 4))
	a.Sub(b)
	println(a.Amount)
}

//action test3
func (c *MyContract) Test3() {
	a := chain.NewAsset(MAX_AMOUNT+3, chain.NewSymbol("EOS", 4))
	println(a.Amount)
}

//action test4
func (c *MyContract) Test4() {
	a := chain.NewAsset(3, chain.NewSymbol("EOS", 4))
	b := chain.NewAsset(0, chain.NewSymbol("EOS", 4))
	a.Div(b)
	println(a.Amount)
}

//action test5
func (c *MyContract) Test5() {
	a := chain.NewAsset(3, chain.NewSymbol("EOS", 4))
	b := chain.NewAsset(-1, chain.NewSymbol("EOS", 4))
	a.Div(b)
	println(a.Amount)
}

//action test11
func (c *MyContract) Test6() {
	a := chain.NewAsset(3, chain.NewSymbol("EO S", 4))
	println(a.Amount)
}

//multiplication overflow test
//action test12
func (c *MyContract) Test12() {
	a := chain.NewAsset(MAX_AMOUNT, chain.NewSymbol("EOS", 4))
	b := chain.NewAsset(MAX_AMOUNT, chain.NewSymbol("EOS", 4))
	a.Mul(b)
	println(a.Amount)
}

//multiplication underflow test
//action test13
func (c *MyContract) Test13() {
	a := chain.NewAsset(-MAX_AMOUNT, chain.NewSymbol("EOS", 4))
	b := chain.NewAsset(MAX_AMOUNT, chain.NewSymbol("EOS", 4))
	a.Mul(b)
	println(a.Amount)
}

//action test14
func (c *MyContract) Test14() {
	a := chain.NewSymbolCode("EOS")
	chain.Check(a.IsValid(), "bad symbol")
	chain.Check(a.Size() == 8, "a.Size() == 8")

	a = chain.NewSymbolCode("EOS EOS")
	chain.Check(!a.IsValid(), "bad symbol")
	a.Pack()
	_a := chain.SymbolCode{}
	_a.Unpack(a.Pack())
	chain.Check(a.Value == _a.Value, "bad value")

	s := chain.NewSymbol("EOS", 4)

	_s := chain.Symbol{}
	_s.Unpack(s.Pack())
	chain.Check(_s.Value == s.Value, "_s.Value == s.Value")
	s.Print()
	chain.Println()
	chain.Println("+++++++=hello, world")
}

//action test15
func (c *MyContract) Test15() {
	{
		a := chain.NewExtendedAsset(chain.NewAsset(100, chain.NewSymbol("EOS", 4)), chain.NewName("alice"))
		_a := chain.ExtendedAsset{}
		_a.Unpack(a.Pack())

		chain.Check(a.Quantity.Amount == _a.Quantity.Amount, "")
		chain.Check(a.Quantity.Symbol.Value == _a.Quantity.Symbol.Value, "")
		chain.Check(a.Size() == _a.Size(), "a.Size() == _a.Size()")
	}
	{
		a := chain.Transfer{
			From:     chain.NewName("alice"),
			To:       chain.NewName("bob"),
			Quantity: *chain.NewAsset(100, chain.NewSymbol("EOS", 4)),
			Memo:     "hello",
		}
		_a := chain.Transfer{}
		_a.Unpack(a.Pack())
		chain.Check(bytes.Compare(a.Pack(), _a.Pack()) == 0, "bytes.Compare(a.Pack(), _a.Pack())")
		chain.Check(a.Size() == _a.Size(), "a.Size() == _a.Size()")
	}
}
