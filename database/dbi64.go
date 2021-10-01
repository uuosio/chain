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

type DBI64 struct {
	code  uint64
	scope uint64
	table uint64
}

func NewDBI64(code chain.Name, scope chain.Name, table chain.Name) *DBI64 {
	return &DBI64{code.N, scope.N, table.N}
}

func (db *DBI64) Init(code chain.Name, scope chain.Name, table chain.Name) {
	db.code = code.N
	db.scope = scope.N
	db.table = table.N
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
	C.db_remove_i64(it.I)
}

//Find the table row following the referenced table row in a primary 64-bit integer index table
func (db *DBI64) Next(it Iterator) (next_iterator Iterator, primary uint64) {
	ret := C.db_next_i64(int32(it.I), (*uint64)(unsafe.Pointer(&primary)))
	return Iterator{ret}, primary
}

//Find the table row preceding the referenced table row in a primary 64-bit integer index table
func (db *DBI64) Previous(it Iterator) (previous_iterator Iterator, primary uint64) {
	ret := C.db_previous_i64(it.I, (*uint64)(unsafe.Pointer(&primary)))
	return Iterator{ret}, primary
}

//Find a table row in a primary 64-bit integer index table by primary key
func (db *DBI64) Find(id uint64) Iterator {
	ret := C.db_find_i64(db.code, db.scope, db.table, id)
	return Iterator{ret}
}

//Find the table row in a primary 64-bit integer index table that matches the lowerbound condition for a given primary key
func (db *DBI64) Lowerbound(id uint64) Iterator {
	ret := C.db_lowerbound_i64(db.code, db.scope, db.table, id)
	return Iterator{ret}
}

//Find the table row in a primary 64-bit integer index table that matches the upperbound condition for a given primary key
func (db *DBI64) Upperbound(id uint64) Iterator {
	ret := C.db_upperbound_i64(db.code, db.scope, db.table, id)
	return Iterator{ret}
}

//Get an iterator representing just-past-the-end of the last table row of a primary 64-bit integer index table
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
