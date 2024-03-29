package testdb

import (
	"encoding/binary"

	"github.com/uuosio/chain"
	"github.com/uuosio/chain/database"
	"github.com/uuosio/chain/logger"
)

const DEBUG = false

//packer
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

type MyTable struct {
	database.TableI64
}

func NewMyTable(code, scope, table chain.Name) *MyTable {
	v := &MyTable{}
	v.Init(code, scope, table)
	return v
}

func (db *MyTable) Get(it *database.Iterator) *MyData {
	data := db.GetByIterator(it)
	if len(data) <= 0 {
		return nil
	}

	_data := &MyData{}
	_data.Unpack(data)
	return _data
}

func (db *MyTable) TryGet(primary uint64) (itr *database.Iterator, data *MyData) {
	itr = db.Find(primary)
	if !itr.IsOk() {
		return itr, nil
	}

	rawData := db.GetByIterator(itr)
	data = &MyData{}
	data.Unpack(rawData)
	return
}

func main() {
	receiver, firstReceiver, action := chain.GetApplyArgs()
	ContractApply(receiver.N, firstReceiver.N, action.N)
}

func ContractApply(receiver, firstReceiver, action uint64) {
	code := chain.NewName("hello")
	scope := code
	table := chain.NewName("mytablee")
	// logger.Printf("%d %d %d\n", code, scope, table)

	db := NewMyTable(code, scope, table)
	//db := database.NewTableI64(code, scope, table)

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
			db.Store(data.primary, chain.EncoderPack(&data), payer)
		} else {
			_data := db.Get(itr)
			logger.Println(_data.primary, _data.n)
			db.Update(itr, chain.EncoderPack(&data), payer)
		}
		logger.Println(itr)
	}
}
