package testsort

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
	a := []int{1, 3, 2}
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})

	for _, v := range a {
		chain.Println("test", "v", v)
	}

	sort.Slice(pubs, func(i, j int) bool {
		return bytes.Compare(pubs[i].Data[:], pubs[j].Data[:]) < 0
	})

	// sort.Slice(pubs, func(i, j int) bool {
	// 	return string(pubs[i].Data[:]) < string(pubs[j].Data[:])
	// })

	for _, pub := range pubs {
		chain.Println(pub.Type, pub.Data[:])
	}

	sort.Sort(chain.PublicKeyList(pubs))
	for _, pub := range pubs {
		chain.Println(pub.Type, pub.Data[:])
	}
}
