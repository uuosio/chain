package database

/*
#include "../structs.h"

int32_t db_idx128_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const uint128* secondary);
void db_idx128_update(int32_t iterator, uint64_t payer, const uint128* secondary);
void db_idx128_remove(int32_t iterator);
int32_t db_idx128_next(int32_t iterator, uint64_t* primary);
int32_t db_idx128_previous(int32_t iterator, uint64_t* primary);
int32_t db_idx128_find_primary(uint64_t code, uint64_t scope, uint64_t table, uint128* secondary, uint64_t primary);
int32_t db_idx128_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const uint128* secondary, uint64_t* primary);
int32_t db_idx128_lowerbound(uint64_t code, uint64_t scope, uint64_t table, uint128* secondary, uint64_t* primary);
int32_t db_idx128_upperbound(uint64_t code, uint64_t scope, uint64_t table, uint128* secondary, uint64_t* primary);
int32_t db_idx128_end(uint64_t code, uint64_t scope, uint64_t table);
*/
import "C"
import (
	"unsafe"

	"github.com/uuosio/chain"
)

type IdxDB128 struct {
	IdxDB
}

func NewIdxDB128(index int, code uint64, scope uint64, table uint64) *IdxDB128 {
	v := &IdxDB128{IdxDB{index, C.uint64_t(code), C.uint64_t(scope), C.uint64_t(table)}}
	return v
}

//Store an association of a 128-bit integer secondary key to a primary key in a secondary 128-bit integer index table
func (db *IdxDB128) Store(id uint64, secondary chain.Uint128, payer uint64) SecondaryIterator {
	ret := C.db_idx128_store(db.scope, db.table, C.uint64_t(payer), C.uint64_t(id), (*C.uint128)(unsafe.Pointer(&secondary)))
	return SecondaryIterator{int32(ret), id, db.dbIndex}
}

func (db *IdxDB128) StoreEx(id uint64, secondary interface{}, payer uint64) SecondaryIterator {
	_secondary, ok := secondary.(chain.Uint128)
	chain.Check(ok, "bad secondary type")
	return db.Store(id, _secondary, payer)
}

//Update an association for a 128-bit integer secondary key to a primary key in a secondary 128-bit integer index table
func (db *IdxDB128) Update(it SecondaryIterator, secondary chain.Uint128, payer uint64) {
	C.db_idx128_update(C.int32_t(it.I), C.uint64_t(payer), (*C.uint128)(unsafe.Pointer(&secondary)))
}

func (db *IdxDB128) UpdateEx(it SecondaryIterator, secondary interface{}, payer uint64) {
	_secondary, ok := secondary.(chain.Uint128)
	chain.Check(ok, "bad secondary type")
	db.Update(it, _secondary, payer)
}

//Remove a table row from a secondary 128-bit integer index table
func (db *IdxDB128) Remove(it SecondaryIterator) {
	var secondary chain.Uint128
	ret := C.db_idx128_find_primary(db.code, db.scope, db.table, (*C.uint128)(unsafe.Pointer(&secondary)), C.uint64_t(it.Primary))
	chain.Check(int32(ret) == it.I, "invalid iterator")
	C.db_idx128_remove(C.int32_t(it.I))
}

//Find the table row following the referenced table row in a secondary 128-bit integer index table
func (db *IdxDB128) Next(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx128_next(C.int32_t(it.I), (*C.uint64_t)(&primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}
}

//Find the table row preceding the referenced table row in a secondary 128-bit integer index table
func (db *IdxDB128) Previous(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx128_previous(C.int32_t(it.I), (*C.uint64_t)(&primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}
}

//Find a table row in a secondary 128-bit integer index table by primary key
func (db *IdxDB128) FindByPrimary(primary uint64) (SecondaryIterator, interface{}) {
	var secondary chain.Uint128
	ret := C.db_idx128_find_primary(db.code, db.scope, db.table, (*C.uint128)(unsafe.Pointer(&secondary)), C.uint64_t(primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}, secondary
}

//Find a table row in a secondary 128-bit integer index table by secondary key
func (db *IdxDB128) Find(secondary chain.Uint128) SecondaryIterator {
	it, value := db.Lowerbound(secondary)
	if it.IsOk() {
		if value == secondary {
			return it
		}
	}
	return SecondaryIterator{-1, 0, db.dbIndex}
}

func (db *IdxDB128) FindEx(secondary interface{}) SecondaryIterator {
	return db.Find(secondary.(chain.Uint128))
}

//Find the table row in a secondary 128-bit integer index table that matches the lowerbound condition for a given secondary key
func (db *IdxDB128) Lowerbound(secondary chain.Uint128) (SecondaryIterator, chain.Uint128) {
	var primary uint64 = 0
	ret := C.db_idx128_lowerbound(db.code, db.scope, db.table, (*C.uint128)(unsafe.Pointer(&secondary)), (*C.uint64_t)(&primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}, secondary
}

//Find the table row in a secondary 128-bit integer index table that matches the upperbound condition for a given secondary key
func (db *IdxDB128) Upperbound(secondary chain.Uint128) (SecondaryIterator, chain.Uint128) {
	var primary uint64 = 0
	ret := C.db_idx128_upperbound(db.code, db.scope, db.table, (*C.uint128)(unsafe.Pointer(&secondary)), (*C.uint64_t)(&primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}, secondary
}

//Get an end iterator representing just-past-the-end of the last table row of a secondary 128-bit integer index table
func (db *IdxDB128) End() SecondaryIterator {
	ret := C.db_idx128_end(db.code, db.scope, db.table)
	return SecondaryIterator{int32(ret), 0, db.dbIndex}
}
