package main

import (
	"github.com/uuosio/chain"
)

//table test
type MyData struct {
	id    uint64 //primary: t.id
	value uint64 //IDX64: ByValue : t.value : t.value
}

//table test2
type MyData2 struct {
	id    uint64        //primary: t.id
	value chain.Uint128 //IDX128: ByValue : t.value : t.value
}

func main() {
	receiver, _, _ := chain.GetApplyArgs()

	db := NewMyDataTable(receiver)
	db.Store(&MyData{2, 3}, receiver)
	db.Store(&MyData{4, 5}, receiver)

	idxTable := db.GetIdxTableByValue()
	a := 0
	it, secondary := idxTable.Lowerbound(uint64(a))
	chain.Println(it.I, secondary, a)

	a = 4
	it, secondary = idxTable.Lowerbound(uint64(a))
	chain.Println(it.I, secondary, a)

}
