package main

import "github.com/uuosio/chain"

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
	a := chain.NewAsset(0, chain.NewSymbol("EOS", 4))
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
