package main

import (
	"github.com/uuosio/chain"
	"github.com/uuosio/chain/database"
)

//contract hello
type MyContract struct {
	Receiver      chain.Name
	FirstReceiver chain.Name
	Action        chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *MyContract {
	return &MyContract{receiver, firstReceiver, action}
}

//action testkv
func (c *MyContract) TestKV(name string) {
	kv := database.NewKV(c.Receiver)
	payer := c.Receiver

	key := []byte("hello:11")
	kv.Set(key, []byte("111"), payer)

	key = []byte("hello:22")
	kv.Set(key, []byte("222"), payer)

	key = []byte("hello:33")
	kv.Set(key, []byte("333"), payer)

	v, ok := kv.Get(key)
	chain.Println("Hello", name, string(v), ok)

	prefix := []byte("hello:")
	it := kv.CreateItr(prefix)
	chain.Println("++++it.Handle:", it.Handle, it.Status())

	for {
		keySize, valueSize, status := it.Prev()
		if !status.IsOk() {
			break
		}
		chain.Println("++++Prev:", it.Handle, keySize, valueSize)
		value, ok := it.Value()
		key, ok := it.Key()
		chain.Println("++++Prev:", ok, string(key), string(value))
	}

	{
		_, _, status := it.LowerBound([]byte("hello:22"))
		for status.IsOk() {
			value, ok := it.Value()
			key, ok := it.Key()
			chain.Println("++++Next:", ok, string(key), string(value))
			_, _, status := it.Next()
			if !status.IsOk() {
				break
			}
		}
	}
}
