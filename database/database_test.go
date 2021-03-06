package database_test

import (
	"github.com/uuosio/chain"
)

//table accounts
type MyData struct {
	primary uint64 //primary: t.primary
	n       uint64
}

func ExampleTable() {
	code := chain.NewName("hello")
	payer := code
	mydb := NewMyDataTable(code)
	primary := uint64(1)
	if it, data := mydb.GetByKey(primary); it.IsOk() {
		data.n += 1
		mydb.Update(it, data, payer)
	} else {
		data := &MyData{primary, 111}
		mydb.Store(data, payer)
	}
}
