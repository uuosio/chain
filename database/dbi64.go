package database

/*
#include <stdint.h>

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
	"unsafe"

	"github.com/uuosio/chain"
)

const (
	MsgErrorUnpacker = "unpacker cannot be nil when save state is enabled"
)

const (
	OperationStore  = 1
	OperationUpdate = 2
	OperationRemove = 3
)

type DBI64Unpacker func([]byte) DBValue

type DBI64 struct {
	code      C.uint64_t
	scope     C.uint64_t
	table     C.uint64_t
	unpacker  DBI64Unpacker
	saveState bool
}

func NewDBI64(code chain.Name, scope chain.Name, table chain.Name, unpacker DBI64Unpacker) *DBI64 {
	return &DBI64{C.uint64_t(code.N), C.uint64_t(scope.N), C.uint64_t(table.N), unpacker, saveState}
}

func (db *DBI64) Init(code chain.Name, scope chain.Name, table chain.Name) {
	db.code = C.uint64_t(code.N)
	db.scope = C.uint64_t(scope.N)
	db.table = C.uint64_t(table.N)
}

func (db *DBI64) GetTable() (uint64, uint64, uint64) {
	return uint64(db.code), uint64(db.scope), uint64(db.table)
}

func (db *DBI64) GetTableName() uint64 {
	return uint64(db.table)
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

//Get a record in a primary 64-bit integer index table
func (db *DBI64) Get(id uint64) (Iterator, []byte) {
	it := db.Find(id)
	if !it.IsOk() {
		return it, nil
	}

	raw := db.getI64(it)
	return it, raw
}

func (db *DBI64) GetByIterator(it Iterator) ([]byte, error) {
	if !it.IsOk() {
		return nil, ErrInvalidIterator
	}
	raw := db.getI64(it)
	return raw, nil
}

//Store a record in a primary 64-bit integer index table
func (db *DBI64) Store(primary uint64, data []byte, payer chain.Name) Iterator {
	return db.storeI64(payer.N, primary, data)
}

//Update a record in a primary 64-bit integer index table
func (db *DBI64) Update(it Iterator, data []byte, payer chain.Name) {
	db.updateI64(it, payer.N, data)
}

//Remove a record from a primary 64-bit integer index table
func (db *DBI64) Remove(it Iterator) {
	C.db_remove_i64(C.int32_t(it.I))
}

//Find the table row following the referenced table row in a primary 64-bit integer index table
func (db *DBI64) Next(it Iterator) (next_iterator Iterator, primary uint64) {
	ret := C.db_next_i64(C.int32_t(it.I), (*C.uint64_t)(unsafe.Pointer(&primary)))
	return Iterator{int32(ret)}, primary
}

//Find the table row preceding the referenced table row in a primary 64-bit integer index table
func (db *DBI64) Previous(it Iterator) (previous_iterator Iterator, primary uint64) {
	ret := C.db_previous_i64(C.int32_t(it.I), (*C.uint64_t)(unsafe.Pointer(&primary)))
	return Iterator{int32(ret)}, primary
}

//Find a table row in a primary 64-bit integer index table by primary key
func (db *DBI64) Find(id uint64) Iterator {
	ret := C.db_find_i64(db.code, db.scope, db.table, C.uint64_t(id))
	return Iterator{int32(ret)}
}

//Find the table row in a primary 64-bit integer index table that matches the lowerbound condition for a given primary key
func (db *DBI64) Lowerbound(id uint64) Iterator {
	ret := C.db_lowerbound_i64(db.code, db.scope, db.table, C.uint64_t(id))
	return Iterator{int32(ret)}
}

//Find the table row in a primary 64-bit integer index table that matches the upperbound condition for a given primary key
func (db *DBI64) Upperbound(id uint64) Iterator {
	ret := C.db_upperbound_i64(db.code, db.scope, db.table, C.uint64_t(id))
	return Iterator{int32(ret)}
}

//Get an iterator representing just-past-the-end of the last table row of a primary 64-bit integer index table
func (db *DBI64) End() Iterator {
	ret := C.db_end_i64(db.code, db.scope, db.table)
	return Iterator{int32(ret)}
}

func (db *DBI64) storeI64(payer uint64, id uint64, data []byte) Iterator {
	GetStateManager().OnStore(db, id)
	p := (*C.char)(unsafe.Pointer(&data[0]))
	ret := C.db_store_i64(db.scope, db.table, C.uint64_t(payer), C.uint64_t(id), p, C.uint32_t(len(data)))
	return Iterator{int32(ret)}
}

func (db *DBI64) updateI64(iterator Iterator, payer uint64, data []byte) {
	GetStateManager().OnUpdate(db, iterator, payer)
	p := (*C.char)(unsafe.Pointer(&data[0]))
	C.db_update_i64(C.int32_t(iterator.I), C.uint64_t(payer), p, C.uint32_t(len(data)))
}

func (db *DBI64) getI64(it Iterator) (data []byte) {
	p := (*C.char)(unsafe.Pointer(uintptr(0)))
	data_size := C.db_get_i64(C.int32_t(it.I), p, C.uint32_t(0))
	if data_size <= 0 {
		return []byte{}
	}

	data = make([]byte, data_size)
	p = (*C.char)(unsafe.Pointer(&data[0]))
	C.db_get_i64(C.int32_t(it.I), p, C.uint32_t(len(data)))
	return data
}
