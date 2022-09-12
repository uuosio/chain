package testlargecode

import (
	"strconv"

	"github.com/uuosio/chain"
)

//table mydata2 singleton
type MySingleton struct {
	a1 uint64
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
func (t *MyContract) test() {
	v, err := strconv.ParseFloat("123", 64)
	if err != nil {
		panic(err)
	}
	chain.Println(v)
	return
}
