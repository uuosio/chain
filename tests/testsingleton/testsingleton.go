package testsingleton

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

//action test1
func (t *MyContract) Test1() {
	code := t.Receiver
	payer := t.Receiver
	db := NewMySingletonTable(code)
	chain.Println("++++test1")
	data := db.Get()
	if data != nil {
		chain.Println("+++a1:", data.a1)
		data.a1 += 1
		db.Set(data, payer)
	} else {
		s := MySingleton{}
		db.Set(&s, payer)
	}
}

//action test2
func (t *MyContract) Test2() {
	code := t.Receiver
	// payer := t.Receiver
	chain.Println("++++test2")
	db := NewMySingletonTable(code)
	db.Remove()
	chain.Check(db.Get() == nil, "db.Get() == nil")
}
