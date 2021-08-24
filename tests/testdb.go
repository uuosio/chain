package main

import (
	"chain"
	"chain/database"
	"chain/logger"
	"encoding/binary"
)

const DEBUG = false

type MyData struct {
	primary uint64
	n       uint32
}

func uint32ToBytes(d uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, d)
	return b
}

func uint64ToBytes(d uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, d)
	return b
}

func bytesToUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b[:])
}

func bytesToUint64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b[:])
}

func (d *MyData) GetPrimary() uint64 {
	return d.primary
}

func (d *MyData) GetSecondaryIndexes() []int {
	return []int{0, 1, 2}
}

type MyDB struct {
	database.DBI64
}

func NewMyDB(code, scope, table chain.Name) *MyDB {
	v := &MyDB{}
	v.Init(code, scope, table)
	return v
}

func (db *MyDB) Get(it database.Iterator) *MyData {
	data := db.GetRawByIterator(it)
	if len(data) <= 0 {
		return nil
	}

	_data := &MyData{}
	_data.Unpack(data)
	return _data
}

func (db *MyDB) TryGet(primary uint64) (itr database.Iterator, data *MyData) {
	itr = db.Find(primary)
	if !itr.IsOk() {
		return itr, nil
	}

	rawData := db.GetRawByIterator(itr)
	data = &MyData{}
	data.Unpack(rawData)
	return
}

func main() {
	code := chain.NewName("hello")
	scope := code
	table := chain.NewName("mytablee")
	// logger.Printf("%d %d %d\n", code, scope, table)

	db := NewMyDB(code, scope, table)
	//db := database.NewDBI64(code, scope, table)

	itr := db.End()
	logger.Println("+++end: \n", itr)

	if itr.IsValid() {
		itr, primary := db.Previous(itr)
		logger.Println("+++Previous: ", itr, primary)
	}

	itr, data := db.TryGet(55)
	if itr.IsOk() {
		logger.Println("++++TryGet ", data.primary, data.n)
	}

	if false {
		data := MyData{55, 66} //[]byte("hello")
		itr := db.Find(data.primary)

		payer := code
		if itr.IsOk() {
			db.Store(&data, payer)
		} else {
			_data := db.Get(itr)
			logger.Println(_data.primary, _data.n)
			db.Update(itr, data.Pack(), payer)
		}
		logger.Println(itr)
	}
}
