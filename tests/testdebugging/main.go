package main

import (
	"github.com/uuosio/chain"
)

/*
void say_hello();

// #cgo CXXFLAGS: -Iincs
// #cgo CFLAGS: -Iincs

*/
import "C"

//contract test
type Contract struct {
	receiver      chain.Name
	firstReceiver chain.Name
	action        chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *Contract {
	return &Contract{
		receiver,
		firstReceiver,
		action,
	}
}

//action sayhello
func (c *Contract) SayHello(name string) {
	chain.Prints("++++++hellow,world")
	chain.NewAction(
		chain.NewPermissionLevel(chain.NewName("hello"), chain.NewName("active")),
		chain.NewName("hello"),
		chain.NewName("saygoodbye"),
		"+++++++++++goodbye, world",
	).Send()
}

//action saygoodbye
func (c *Contract) SayGoodbye(name string) {
	chain.Prints("+++++++goodbyte word")
}

//action inc
func (c *Contract) Inc(name string) {
	db := NewCounterTable(c.receiver)
	it := db.Find(1)
	payer := c.receiver
	if it.IsOk() {
		value := db.GetByIterator(it)
		value.count += 1
		db.Update(it, value, payer)
		chain.Println("count: ", value.count)
	} else {
		value := &Counter{
			key:   1,
			count: 1,
		}
		db.Store(value, payer)
		chain.Println("count: ", value.count)
	}
}

func SayHelloFromCpp() {
	C.say_hello()
}
