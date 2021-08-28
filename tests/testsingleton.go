package main

import (
	"github.com/uuosio/chain"
)

//table mydata singleton
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
	db := NewMySingletonDB(code, code)

	data, err := db.Get()
	if err == nil {
		println("+++a1:", data.a1)
		data.a1 += 1
		db.Set(data, payer)
	} else {
		chain.Println("+++err:", err)
		s := MySingleton{}
		db.Set(&s, payer)
	}
}
