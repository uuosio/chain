package main

import (
	"bytes"
	"sort"

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

func less(a, b []byte) bool {
	if len(a) < len(b) {
		return true
	} else if len(a) > len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] < b[i] {
			return true
		}
		return false
	}
	return false
}

//action test
func (t *MyContract) test(pubs []chain.PublicKey) {

	a := []byte("02056f8d91fc8996bd3f2fa4dc2dd60d6cdd05811db0065dd4f191b8bbd8eda5d4")
	b := []byte("02cbb8488ce044116dd4411c7c76ca44bb42b76400a0eb8353bf62a327423da26c")
	r := bytes.Compare(a, b)
	println("++++bytes.Compare:", r)

	// chain.Println("02cbb8488ce044116dd4411c7c76ca44bb42b76400a0eb8353bf62a327423da26c" < "02056f8d91fc8996bd3f2fa4dc2dd60d6cdd05811db0065dd4f191b8bbd8eda5d4")

	// a := []int{1, 3, 2}
	// sort.Slice(a, func(i, j int) bool {
	// 	return a[i] < a[j]
	// })

	// for _, v := range a {
	// 	chain.Println("test", "v", v)
	// }

	// sort.Slice(pubs, func(i, j int) bool {
	// 	return bytes.Compare(pubs[i].Data[:], pubs[j].Data[:]) < 0
	// })

	// sort.Slice(pubs, func(i, j int) bool {
	// 	return string(pubs[i].Data[:]) < string(pubs[j].Data[:])
	// })

	// for _, pub := range pubs {
	// 	chain.Println(pub.Type, pub.Data[:])
	// }
	// return

	sort.Sort(chain.PublicKeyList(pubs))
	for _, pub := range pubs {
		chain.Println(pub.Type, pub.Data[:])
	}
}
