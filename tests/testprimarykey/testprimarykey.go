package testprimarykey

import (
	"github.com/uuosio/chain"
)

//table mytable
type MyData struct {
	primary uint64 //primary: t.primary
	n       uint64
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

//action test
func (c *MyContract) Test(name string) {
	code := chain.NewName("hello")
	payer := code
	mydb := NewMyDataTable(code)
	primary := uint64(1)
	it, data := mydb.GetByKey(primary)
	if !it.IsOk() {
		data := &MyData{primary, 111}
		mydb.Store(data, payer)
	} else {
		data.n += 1
		mydb.Update(it, data, payer)
	}
}
