package main

import (
	"github.com/uuosio/chain"
)

var gContractName = chain.NewName("hello")
var gActionName = chain.NewName("sayhello2")

//contract test
type ActionTest struct {
	self   chain.Name
	code   chain.Name
	action chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *ActionTest {
	return &ActionTest{
		self:   receiver,
		code:   firstReceiver,
		action: action,
	}
}

//action sayhello
func (c *ActionTest) SayHello() {
	a := chain.Action{
		chain.NewName("hello"), //gContractName,
		gActionName,
		[]chain.PermissionLevel{{gContractName, chain.ActiveName}},
		[]byte("hello,world"),
	}
	a.Send()
}

//action sayhello2
func (c *ActionTest) SayHello2() {
	chain.Println(chain.ReadActionData())
}

//action sayhello3
func (c *ActionTest) SayHello3() {
	a := chain.NewAction(
		chain.PermissionLevel{gContractName, chain.ActiveName},
		chain.NewName("eosio.token"),
		chain.NewName("transfer"),
		chain.NewName("hello"),                           //from
		chain.NewName("eosio"),                           //to
		chain.NewAsset(10000, chain.NewSymbol("EOS", 4)), //quantity 1.0000 EOS
		"hello,world",                                    //memo
	)
	a.Send()
}
