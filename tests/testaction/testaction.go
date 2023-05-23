package testaction

import (
	"github.com/uuosio/chain"
)

var gContractName = chain.NewName("hello")
var gActionName = chain.NewName("sayhello2")

// contract test
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

// action sayhello
func (c *ActionTest) SayHello() {
	a := chain.Action{
		Account: chain.NewName("hello"), //gContractName,
		Name:    chain.NewName("sayhello2"),
		Authorization: []*chain.PermissionLevel{
			{
				Actor:      gContractName,
				Permission: chain.ActiveName,
			},
		},
		Data: []byte("hello,world"),
	}
	a.Send()
}

// action sayhello2
func (c *ActionTest) SayHello2() {
	chain.Println(chain.ReadActionData())
}

// action sayhello3
func (c *ActionTest) SayHello3() {
	a := chain.NewAction(
		&chain.PermissionLevel{
			Actor:      gContractName,
			Permission: chain.ActiveName,
		},
		chain.NewName("eosio.token"),
		chain.NewName("transfer"),
		chain.NewName("hello"),                           //from
		chain.NewName("eosio"),                           //to
		chain.NewAsset(10000, chain.NewSymbol("EOS", 4)), //quantity 1.0000 EOS
		"hello,world",                                    //memo
	)
	a.Send()
}

// packer
type Transfer struct {
	from     chain.Name
	to       chain.Name
	quantity chain.Asset
	memo     string
}

// action sayhello4
func (c *ActionTest) SayHello4() {
	transfer := Transfer{
		from:     chain.NewName("hello"),
		to:       chain.NewName("eosio"),
		quantity: chain.Asset{Amount: 10000, Symbol: chain.NewSymbol("EOS", 4)},
		memo:     "hello, world",
	}

	a := chain.NewAction(
		&chain.PermissionLevel{
			Actor:      gContractName,
			Permission: chain.ActiveName,
		},
		chain.NewName("eosio.token"),
		chain.NewName("transfer"),
		&transfer,
	)
	a.Send()
}

// action getcodehash
func (c *ActionTest) GetCodeHash(hash chain.Checksum256) {
	ret := chain.GetCodeHash(c.code)
	chain.Check(ret == hash, "bad code hash")
}
