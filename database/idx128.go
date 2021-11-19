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
	"math/big"
	"unsafe"

	"github.com/uuosio/chain"
)

type IdxDB128 struct {
	dbIndex int
	code    C.uint64_t
	scope   C.uint64_t
	table   C.uint64_t
}

func (db *IdxDB128) GetIndex() int {
	return db.dbIndex
}

func NewIdxDB128(index int, code uint64, scope uint64, table uint64) *IdxDB128 {
	v := &IdxDB128{index, C.uint64_t(code), C.uint64_t(scope), C.uint64_t(table)}
	return v
}

//Store an association of a 128-bit integer secondary key to a primary key in a secondary 128-bit integer index table
func (db *IdxDB128) Store(id uint64, secondary interface{}, payer uint64) SecondaryIterator {
	_secondary, ok := secondary.(chain.Uint128)
	chain.Check(ok, "bad secondary type")
	ret := C.db_idx128_store(db.scope, db.table, C.uint64_t(payer), C.uint64_t(id), (*C.uint128)(unsafe.Pointer(&_secondary)))
	return SecondaryIterator{int32(ret), id, db.dbIndex}
}

//Update an association for a 128-bit integer secondary key to a primary key in a secondary 128-bit integer index table
func (db *IdxDB128) Update(it SecondaryIterator, secondary interface{}, payer uint64) {
	_secondary, ok := secondary.(chain.Uint128)
	chain.Check(ok, "bad secondary type")
	C.db_idx128_update(C.int32_t(it.I), C.uint64_t(payer), (*C.uint128)(unsafe.Pointer(&_secondary)))
}

//Remove a table row from a secondary 128-bit integer index table
func (db *IdxDB128) Remove(it SecondaryIterator) {
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
func (db *IdxDB128) Find(secondary interface{}) SecondaryIterator {
	// var primary uint64 = 0
	var rawSecondary chain.Uint128
	switch v := secondary.(type) {
	case big.Int:
		_secondary := v.Bytes()
		chain.Check(len(_secondary) <= 16, "Uint128 overflow")
		copy(rawSecondary[:], _secondary)
	case chain.Uint128:
		rawSecondary = v
	default:
		panic("can not conver to uint128 type")
	}
	// _secondary, ok := secondary.(Uint128)
	// chain.Check(ok, "bad secondary type")
	// ret := C.db_idx128_find_secondary(db.code, db.scope, db.table, (*uint8)(unsafe.Pointer(&rawSecondary[0])), &primary)
	// return SecondaryIterator{ret, primary, db.dbIndex}
	it, value := db.Lowerbound(secondary)
	if it.IsOk() {
		if value.(chain.Uint128) == rawSecondary {
			return it
		}
	}
	return SecondaryIterator{-1, 0, db.dbIndex}
}

//Find the table row in a secondary 128-bit integer index table that matches the lowerbound condition for a given secondary key
func (db *IdxDB128) Lowerbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(chain.Uint128)
	chain.Assert(ok, "bad secondary type")
	ret := C.db_idx128_lowerbound(db.code, db.scope, db.table, (*C.uint128)(unsafe.Pointer(&_secondary)), (*C.uint64_t)(&primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}, _secondary
}

//Find the table row in a secondary 128-bit integer index table that matches the upperbound condition for a given secondary key
func (db *IdxDB128) Upperbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(chain.Uint128)
	chain.Assert(ok, "bad secondary type")
	ret := C.db_idx128_upperbound(db.code, db.scope, db.table, (*C.uint128)(unsafe.Pointer(&_secondary)), (*C.uint64_t)(&primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}, _secondary
}

//Get an end iterator representing just-past-the-end of the last table row of a secondary 128-bit integer index table
func (db *IdxDB128) End() SecondaryIterator {
	ret := C.db_idx128_end(db.code, db.scope, db.table)
	return SecondaryIterator{int32(ret), 0, db.dbIndex}
}

type IdxDB128I struct {
	SecondaryDB
}

func (db *IdxDB128I) FindByPrimary(primary uint64) (SecondaryIterator, chain.Uint128) {
	it, _secondary := db.SecondaryDB.FindByPrimary(primary)
	return it, _secondary.(chain.Uint128)
}

func (db *IdxDB128I) Find(secondary chain.Uint128) SecondaryIterator {
	return db.SecondaryDB.Find(secondary)
}

func (db *IdxDB128I) Lowerbound(secondary chain.Uint128) (SecondaryIterator, chain.Uint128) {
	it, _secondary := db.SecondaryDB.Lowerbound(secondary)
	return it, _secondary.(chain.Uint128)
}

func (db *IdxDB128I) Upperbound(secondary chain.Uint128) (SecondaryIterator, chain.Uint128) {
	it, _secondary := db.SecondaryDB.Upperbound(secondary)
	return it, _secondary.(chain.Uint128)
}
