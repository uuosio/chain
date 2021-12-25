package main

import (
	"github.com/uuosio/chain"
	"github.com/uuosio/chain/revert"
)

//table test
type MyData struct {
	id    uint64 //primary: t.id
	value uint64 //IDX64: ByValue : t.value : t.value
}

func main() {
	receiver, _, action := chain.GetApplyArgs()
	if action == chain.NewName("revert") {
		db := NewMyDataDB(receiver, receiver)
		it := db.Find(1)
		chain.Check(!it.IsOk(), "value should not exists")

		idxDB := db.GetIdxDBByValue()
		it2 := idxDB.Find(2)
		chain.Check(!it2.IsOk(), "secondary value should not exists")
		chain.Println("revert ok")
	} else {
		revert.Init()
		db := NewMyDataDB(receiver, receiver)
		db.Store(&MyData{1, 2}, receiver)

		it := db.Find(1)
		chain.Check(it.IsOk(), "value should exists")

		idxDB := db.GetIdxDBByValue()
		it2 := idxDB.Find(2)
		chain.Check(it2.IsOk(), "secondary value should exists")

		chain.Check(false, "revert")
	}
}
