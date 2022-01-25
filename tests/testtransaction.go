package main

import (
	"bytes"
	"encoding/binary"

	"github.com/uuosio/chain"
	"github.com/uuosio/chain/database"
	"github.com/uuosio/chain/logger"
)

/*
   class transaction_header {
   public:
      time_point_sec  expiration;
      uint16_t        ref_block_num;
      uint32_t        ref_block_prefix;
      unsigned_int    max_net_usage_words = 0UL; /// number of 8 byte words this transaction can serialize into after compressions
      uint8_t         max_cpu_usage_ms = 0UL; /// number of CPU usage units to bill transaction for
      unsigned_int    delay_sec = 0UL; /// number of seconds to delay transaction, default: 0
*/

func testTransaction() {
	rawTx := chain.ReadTransaction()
	chain.Println(rawTx)

	tx := chain.Transaction{}
	tx.Unpack(rawTx)
	tx.Print()

	rawTx2 := tx.Pack()
	chain.Println(rawTx2)

	chain.Check(len(rawTx) == len(rawTx2), "bad")
	for i, v := range rawTx {
		chain.Check(v == rawTx2[i], "bad")
	}
}

func main() {
	_, _, action := chain.GetApplyArgs()

	code := chain.NewName("hello")
	scope := code
	table := chain.NewName("mytable1")
	payer := code

	id := uint64(123)
	if action == chain.NewName("sayhello1") {
		payer := chain.NewName("hello")

		// a := chain.Action{}
		// a.Account = chain.NewName("hello")
		// a.Name = chain.NewName("sayhello2")
		// a.AddPermission(chain.NewName("hello"), chain.ActiveName)
		// a.Data = []byte("hello,world")

		a := chain.NewAction(
			chain.NewPermissionLevel(chain.NewName("hello"), chain.ActiveName),
			chain.NewName("hello"),
			chain.NewName("sayhello2"),
			[]byte("hello,world"),
		)

		t := chain.NewTransaction(1)
		t.Actions = []*chain.Action{a}
		id := chain.NewUint128(1, 0)
		t.Send(id, false, payer)
		logger.Println("send done!")

		rawTx := t.Pack()
		t2 := chain.Transaction{}
		t2.Unpack(rawTx)
		rawTx2 := t2.Pack()
		chain.Check(bytes.Compare(rawTx, rawTx2) == 0, "bad transaction")
	} else if action == chain.NewName("sayhello2") {
		db := database.NewDBI64(code, scope, table, nil)
		it := db.Find(id)
		if it.IsOk() {
			data := db.GetByIterator(it)

			if len(data) != 8 {
				buf := [8]byte{}
				binary.LittleEndian.PutUint64(buf[:], 0)
				db.Update(it, buf[:], payer)
			} else {
				buf := [8]byte{}
				n := binary.LittleEndian.Uint64(data)
				n += 1
				binary.LittleEndian.PutUint64(buf[:], n)
				db.Update(it, buf[:], payer)
			}
			logger.Println("+++++Update")
		} else {
			db.Set(id, []byte("hello,world!!!"), payer)
			logger.Println("+++++Set")
		}
	} else if action == chain.NewName("sayhello3") {
		db := database.NewDBI64(code, scope, table, nil)
		it := db.Find(id)
		if it.IsOk() {
			value := db.GetByIterator(it)
			n := binary.LittleEndian.Uint64(value)
			logger.Println(len(value), n, string(value))
		} else {
			logger.Println("not found")
		}
	} else {
		chain.Println("hello,worldddd")
	}
}
