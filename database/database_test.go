package database_test

import (
	"github.com/uuosio/chain"
)

//table accounts
type MyData struct {
	primary uint64 //primary: t.primary
	n       uint64
}

func ExampleDB() {
	code := chain.NewName("hello")
	scope := code
	payer := code
	mydb := NewMyDataDB(code, scope)
	primary := uint64(1)
	if it, data := mydb.GetByKey(primary); it.IsOk() {
		data.n += 1
		mydb.Update(it, data, payer)
	} else {
		data := &MyData{primary, 111}
		mydb.Store(data, payer)
	}
}
