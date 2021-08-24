package database

/*
#include <stdint.h>

typedef uint8_t uint128_t; //fake uint128_t definition

int32_t db_store_i64(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id,  const char* data, uint32_t len);
void db_update_i64(int32_t iterator, uint64_t payer, const char* data, uint32_t len);
void db_remove_i64(int32_t iterator);
int32_t db_get_i64(int32_t iterator, const char* data, uint32_t len);
int32_t db_next_i64(int32_t iterator, uint64_t* primary);
int32_t db_previous_i64(int32_t iterator, uint64_t* primary);
int32_t db_find_i64(uint64_t code, uint64_t scope, uint64_t table, uint64_t id);
int32_t db_lowerbound_i64(uint64_t code, uint64_t scope, uint64_t table, uint64_t id);
int32_t db_upperbound_i64(uint64_t code, uint64_t scope, uint64_t table, uint64_t id);
int32_t db_end_i64(uint64_t code, uint64_t scope, uint64_t table);
*/
import "C"

import (
	"chain"
	"errors"
	"unsafe"
)

type DBI64 struct {
	code     uint64
	scope    uint64
	table    uint64
	unpacker Unpacker
}

func NewDBI64(code chain.Name, scope chain.Name, table chain.Name, unpacker ...Unpacker) *DBI64 {
	if len(unpacker) == 0 {
		return &DBI64{code.N, scope.N, table.N, nil}
	} else {
		return &DBI64{code.N, scope.N, table.N, unpacker[0]}
	}
}

func (db *DBI64) Init(code chain.Name, scope chain.Name, table chain.Name, unpacker ...Unpacker) {
	db.code = code.N
	db.scope = scope.N
	db.table = table.N
	if len(unpacker) == 0 {
		db.unpacker = nil
	} else {
		db.unpacker = unpacker[0]
	}
}

func (db *DBI64) Set(id uint64, data []byte, payer chain.Name) Iterator {
	it := db.Find(id)
	if !it.IsOk() {
		return db.storeI64(payer.N, id, data)
	} else {
		db.updateI64(it, payer.N, data)
		return it
	}
}

func (db *DBI64) Get(id uint64) (interface{}, error) {
	it := db.Find(id)
	if !it.IsOk() {
		return nil, errors.New("NotFound")
	}

	raw := db.getI64(it)
	if db.unpacker != nil {
		return db.unpacker(raw)
	}
	return raw, nil
}

func (db *DBI64) GetRaw(id uint64) ([]byte, error) {
	it := db.Find(id)
	if !it.IsOk() {
		return nil, errors.New("NotFound")
	}

	raw := db.getI64(it)
	return raw, nil
}

func (db *DBI64) GetByIterator(it Iterator) (interface{}, error) {
	raw := db.getI64(it)
	if db.unpacker == nil {
		return raw, nil
	}
	return db.unpacker(raw)
}

func (db *DBI64) GetRawByIterator(it Iterator) []byte {
	return db.getI64(it)
}

func (db *DBI64) Store(data DBValue, payer chain.Name) Iterator {
	return db.storeI64(payer.N, data.GetPrimary(), data.Pack())
	// primary := data.GetPrimary()
	// itr := db.Find(primary)
	// if !itr.IsOk() {
	// 	return db.storeI64(payer, primary, data.Pack())
	// } else {
	// 	db.Update(itr, payer, data)
	// 	return itr
	// }
}

func (db *DBI64) Update(it Iterator, data interface{}, payer chain.Name) {
	switch v := data.(type) {
	case DBValue:
		db.updateI64(it, payer.N, v.Pack())
		break
	case []byte:
		db.updateI64(it, payer.N, v)
		break
	default:
		panic("unknown data type")
	}
}

func (db *DBI64) Remove(it Iterator) {
	C.db_remove_i64(it.I)
}

func (db *DBI64) Next(it Iterator) (next_iterator Iterator, primary uint64) {
	ret := C.db_next_i64(int32(it.I), (*uint64)(unsafe.Pointer(&primary)))
	return Iterator{ret}, primary
}

func (db *DBI64) Previous(it Iterator) (previous_iterator Iterator, primary uint64) {
	ret := C.db_previous_i64(it.I, (*uint64)(unsafe.Pointer(&primary)))
	return Iterator{ret}, primary
}

func (db *DBI64) Find(id uint64) Iterator {
	ret := C.db_find_i64(db.code, db.scope, db.table, id)
	return Iterator{ret}
}

func (db *DBI64) Lowerbound(id uint64) Iterator {
	ret := C.db_lowerbound_i64(db.code, db.scope, db.table, id)
	return Iterator{ret}
}

func (db *DBI64) Upperbound(id uint64) Iterator {
	ret := C.db_upperbound_i64(db.code, db.scope, db.table, id)
	return Iterator{ret}
}

func (db *DBI64) End() Iterator {
	ret := C.db_end_i64(db.code, db.scope, db.table)
	return Iterator{ret}
}

func (db *DBI64) storeI64(payer uint64, id uint64, data []byte) Iterator {
	p := (*C.char)(unsafe.Pointer(&data[0]))
	ret := C.db_store_i64(db.scope, db.table, payer, id, p, uint32(len(data)))
	return Iterator{ret}
}

func (db *DBI64) updateI64(iterator Iterator, payer uint64, data []byte) {
	p := (*C.char)(unsafe.Pointer(&data[0]))
	C.db_update_i64(iterator.I, payer, p, uint32(len(data)))
}

func (db *DBI64) getI64(it Iterator) (data []byte) {
	p := (*C.char)(unsafe.Pointer(uintptr(0)))
	data_size := C.db_get_i64(it.I, p, uint32(0))
	if data_size <= 0 {
		return []byte{}
	}

	data = make([]byte, data_size)
	p = (*C.char)(unsafe.Pointer(&data[0]))
	C.db_get_i64(it.I, p, uint32(len(data)))
	return data
}
