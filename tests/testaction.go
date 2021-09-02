package main

import (
	"github.com/uuosio/chain"
)

var gContractName = chain.NewName("hello")
var gActionName = chain.NewName("sayhello2")

/*
   def get_secondary_values(self):
       return (self.a, self.b, self.c, self.d)

   @staticmethod
   def get_secondary_indexes():
       return (db.idx64, db.idx128, db.idx256, db.idx_double)
*/

type uint128 [16]byte
type uint256 [32]byte

func (n *uint128) Pack() []byte {
	return n[:]
}

func IsEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

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
		chain.NewName("hello"), //gContractName,
		gActionName,
		[]chain.PermissionLevel{{gContractName, chain.ActiveName}},
		uint32(1),
		[]byte("hello,world"),
	)
	a.Send()
}
