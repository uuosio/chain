package main

import (
	"github.com/uuosio/chain"
)

//table mydata2 singleton
type MySingleton struct {
	a1 uint64
}

//contract hello
type MyContract struct {
	Receiver      chain.Name
	FirstReceiver chain.Name
	Action        chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *MyContract {
	return &MyContract{receiver, firstReceiver, action}
}

//action sayhello
func (t *MyContract) SayHello() {
	code := t.Receiver
	payer := t.Receiver
	db := NewMySingletonTable(code, code)

	data := db.Get()
	if data != nil {
		println("+++a1:", data.a1)
		data.a1 += 1
		db.Set(data, payer)
	} else {
		s := MySingleton{}
		db.Set(&s, payer)
	}
}
