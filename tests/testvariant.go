package main

import (
	"github.com/uuosio/chain"
)

//variant uint64 chain.Uint128
type MyVariant struct {
	value interface{}
}

//table mytable singleton
type MyTable struct {
	a MyVariant
}

//contract test
type MyContract struct {
	Receiver      chain.Name
	FirstReceiver chain.Name
	Action        chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *MyContract {
	return &MyContract{receiver, firstReceiver, action}
}

//action testvariant
func (c *MyContract) TestVariant(v MyVariant) {
	chain.Check(*v.value.(*uint64) == 123, "bad value")
	chain.Println("+++value:", *v.value.(*uint64))
	payer := c.Receiver
	db := NewMyTableTable(c.Receiver)

	item := MyTable{v}
	db.Set(&item, payer)
}
