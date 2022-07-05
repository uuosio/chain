package database

/*
#include "../structs.h"

int32_t db_idx256_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const uint128* data, uint32_t data_len );
void db_idx256_update(int32_t iterator, uint64_t payer, const uint128* data, uint32_t data_len);
void db_idx256_remove(int32_t iterator);
int32_t db_idx256_next(int32_t iterator, uint64_t* primary);
int32_t db_idx256_previous(int32_t iterator, uint64_t* primary);
int32_t db_idx256_find_primary(uint64_t code, uint64_t scope, uint64_t table, uint128* data, uint32_t data_len, uint64_t primary);
int32_t db_idx256_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const uint128* data, uint32_t data_len, uint64_t* primary);
int32_t db_idx256_lowerbound(uint64_t code, uint64_t scope, uint64_t table, uint128* data, uint32_t data_len, uint64_t* primary);
int32_t db_idx256_upperbound(uint64_t code, uint64_t scope, uint64_t table, uint128* data, uint32_t data_len, uint64_t* primary);
int32_t db_idx256_end(uint64_t code, uint64_t scope, uint64_t table);
*/
import "C"
import (
	"unsafe"

	"github.com/uuosio/chain"
)

type IdxTable256 struct {
	IdxTable
}

func NewIdxTable256(index int, code uint64, scope uint64, table uint64) *IdxTable256 {
	v := &IdxTable256{IdxTable{index, C.uint64_t(code), C.uint64_t(scope), C.uint64_t(table)}}
	return v
}

//Store an association of a 256-bit secondary key to a primary key in a secondary 256-bit index table
func (db *IdxTable256) Store(id uint64, secondary chain.Uint256, payer uint64) *SecondaryIterator {
	ret := C.db_idx256_store(db.scope, db.table, C.uint64_t(payer), C.uint64_t(id), (*C.uint128)(unsafe.Pointer(&secondary)), 2)
	return &SecondaryIterator{int32(ret), id, db.dbIndex}
}

func (db *IdxTable256) StoreEx(id uint64, secondary interface{}, payer uint64) *SecondaryIterator {
	_secondary, ok := secondary.(chain.Uint256)
	chain.Check(ok, "bad secondary type")
	return db.Store(id, _secondary, payer)
}

//Update an association for a 256-bit secondary key to a primary key in a secondary 256-bit index table
func (db *IdxTable256) Update(it *SecondaryIterator, secondary chain.Uint256, payer uint64) {
	C.db_idx256_update(C.int32_t(it.I), C.uint64_t(payer), (*C.uint128)(unsafe.Pointer(&secondary)), 2)
}

func (db *IdxTable256) UpdateEx(it *SecondaryIterator, secondary interface{}, payer uint64) {
	_secondary, ok := secondary.(chain.Uint256)
	chain.Check(ok, "bad secondary type")
	db.Update(it, _secondary, payer)
}

//Remove a table row from a secondary 256-bit index table
func (db *IdxTable256) Remove(it *SecondaryIterator) {
	C.db_idx256_remove(C.int32_t(it.I))
}

//Find the table row following the referenced table row in a secondary 256-bit index table
func (db *IdxTable256) Next(it *SecondaryIterator) *SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx256_next(C.int32_t(it.I), (*C.uint64_t)(&primary))
	return &SecondaryIterator{int32(ret), primary, db.dbIndex}
}

//Find the table row preceding the referenced table row in a secondary 256-bit index table
func (db *IdxTable256) Previous(it *SecondaryIterator) *SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx256_previous(C.int32_t(it.I), (*C.uint64_t)(&primary))
	return &SecondaryIterator{int32(ret), primary, db.dbIndex}
}

//Find a table row in a secondary 128-bit integer index table by primary key
func (db *IdxTable256) FindByPrimary(primary uint64) (*SecondaryIterator, interface{}) {
	var secondary chain.Uint256
	ret := C.db_idx256_find_primary(db.code, db.scope, db.table, (*C.uint128)(unsafe.Pointer(&secondary)), 2, C.uint64_t(primary))
	return &SecondaryIterator{int32(ret), primary, db.dbIndex}, secondary
}

// Find a table row in a secondary 256-bit index table by secondary key
func (db *IdxTable256) Find(secondary chain.Uint256) *SecondaryIterator {
	it, value := db.Lowerbound(secondary)
	if it.IsOk() {
		if value == secondary {
			return it
		}
	}
	return &SecondaryIterator{-1, 0, db.dbIndex}
}

func (db *IdxTable256) FindEx(secondary interface{}) *SecondaryIterator {
	return db.Find(secondary.(chain.Uint256))
}

//Find the table row in a secondary 256-bit index table that matches the lowerbound condition for a given secondary key
func (db *IdxTable256) Lowerbound(secondary chain.Uint256) (*SecondaryIterator, chain.Uint256) {
	var primary uint64 = 0
	ret := C.db_idx256_lowerbound(db.code, db.scope, db.table, (*C.uint128)(unsafe.Pointer(&secondary)), 2, (*C.uint64_t)(&primary))
	return &SecondaryIterator{int32(ret), primary, db.dbIndex}, secondary
}

//Find the table row in a secondary 256-bit index table that matches the upperbound condition for a given secondary key
func (db *IdxTable256) Upperbound(secondary chain.Uint256) (*SecondaryIterator, chain.Uint256) {
	var primary uint64 = 0
	ret := C.db_idx256_upperbound(db.code, db.scope, db.table, (*C.uint128)(unsafe.Pointer(&secondary)), 2, (*C.uint64_t)(&primary))
	return &SecondaryIterator{int32(ret), primary, db.dbIndex}, secondary
}

//Get an end iterator representing just-past-the-end of the last table row of a secondary 256-bit index table
func (db *IdxTable256) End() *SecondaryIterator {
	ret := C.db_idx256_end(db.code, db.scope, db.table)
	return &SecondaryIterator{int32(ret), 0, db.dbIndex}
}
