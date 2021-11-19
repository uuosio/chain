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
	"math/big"
	"unsafe"

	"github.com/uuosio/chain"
)

type IdxDB256 struct {
	dbIndex int
	code    C.uint64_t
	scope   C.uint64_t
	table   C.uint64_t
}

func NewIdxDB256(index int, code uint64, scope uint64, table uint64) *IdxDB256 {
	v := &IdxDB256{index, C.uint64_t(code), C.uint64_t(scope), C.uint64_t(table)}
	return v
}

func (db *IdxDB256) GetIndex() int {
	return db.dbIndex
}

//Store an association of a 256-bit secondary key to a primary key in a secondary 256-bit index table
func (db *IdxDB256) Store(id uint64, secondary interface{}, payer uint64) SecondaryIterator {
	_secondary, ok := secondary.(chain.Uint256)
	chain.Check(ok, "bad secondary type")
	ret := C.db_idx256_store(db.scope, db.table, C.uint64_t(payer), C.uint64_t(id), (*C.uint128)(unsafe.Pointer(&_secondary)), 2)
	return SecondaryIterator{int32(ret), id, db.dbIndex}
}

//Update an association for a 256-bit secondary key to a primary key in a secondary 256-bit index table
func (db *IdxDB256) Update(it SecondaryIterator, secondary interface{}, payer uint64) {
	_secondary, ok := secondary.(chain.Uint256)
	chain.Check(ok, "bad secondary type")
	C.db_idx256_update(C.int32_t(it.I), C.uint64_t(payer), (*C.uint128)(unsafe.Pointer(&_secondary)), 2)
}

//Remove a table row from a secondary 256-bit index table
func (db *IdxDB256) Remove(it SecondaryIterator) {
	C.db_idx256_remove(C.int32_t(it.I))
}

//Find the table row following the referenced table row in a secondary 256-bit index table
func (db *IdxDB256) Next(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx256_next(C.int32_t(it.I), (*C.uint64_t)(&primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}
}

//Find the table row preceding the referenced table row in a secondary 256-bit index table
func (db *IdxDB256) Previous(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx256_previous(C.int32_t(it.I), (*C.uint64_t)(&primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}
}

//Find a table row in a secondary 128-bit integer index table by primary key
func (db *IdxDB256) FindByPrimary(primary uint64) (SecondaryIterator, interface{}) {
	var secondary chain.Uint256
	ret := C.db_idx256_find_primary(db.code, db.scope, db.table, (*C.uint128)(unsafe.Pointer(&secondary)), 2, C.uint64_t(primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}, secondary
}

// Find a table row in a secondary 256-bit index table by secondary key
func (db *IdxDB256) Find(secondary interface{}) SecondaryIterator {
	// var primary uint64 = 0
	var rawSecondary chain.Uint256
	switch v := secondary.(type) {
	case big.Int:
		_secondary := v.Bytes()
		chain.Check(len(_secondary) <= 32, "Uint256 overflow")
		copy(rawSecondary[:], _secondary)
	case chain.Uint256:
		rawSecondary = v
	default:
		panic("can not conver to uint128 type")
	}
	// ret := C.db_idx256_find_secondary(db.code, db.scope, db.table, (*uint8)(unsafe.Pointer(&rawSecondary[0])), 2, &primary)
	// return SecondaryIterator{ret, primary, db.dbIndex}
	it, value := db.Lowerbound(secondary)
	if it.IsOk() {
		if value.(chain.Uint256) == rawSecondary {
			return it
		}
	}
	return SecondaryIterator{-1, 0, db.dbIndex}
}

//Find the table row in a secondary 256-bit index table that matches the lowerbound condition for a given secondary key
func (db *IdxDB256) Lowerbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(chain.Uint256)
	chain.Assert(ok, "not a Uint256 value")
	ret := C.db_idx256_lowerbound(db.code, db.scope, db.table, (*C.uint128)(unsafe.Pointer(&_secondary)), 2, (*C.uint64_t)(&primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}, _secondary
}

//Find the table row in a secondary 256-bit index table that matches the upperbound condition for a given secondary key
func (db *IdxDB256) Upperbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(chain.Uint256)
	chain.Assert(ok, "not a Uint256 value")
	ret := C.db_idx256_upperbound(db.code, db.scope, db.table, (*C.uint128)(unsafe.Pointer(&_secondary)), 2, (*C.uint64_t)(&primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}, _secondary
}

//Get an end iterator representing just-past-the-end of the last table row of a secondary 256-bit index table
func (db *IdxDB256) End() SecondaryIterator {
	ret := C.db_idx256_end(db.code, db.scope, db.table)
	return SecondaryIterator{int32(ret), 0, db.dbIndex}
}

type IdxDB256I struct {
	SecondaryDB
}

func (db *IdxDB256I) FindByPrimary(primary uint64) (SecondaryIterator, chain.Uint256) {
	it, _secondary := db.SecondaryDB.FindByPrimary(primary)
	return it, _secondary.(chain.Uint256)
}

func (db *IdxDB256I) Find(secondary chain.Uint256) SecondaryIterator {
	return db.SecondaryDB.Find(secondary)
}

func (db *IdxDB256I) Lowerbound(secondary chain.Uint256) (SecondaryIterator, chain.Uint256) {
	it, _secondary := db.SecondaryDB.Lowerbound(secondary)
	return it, _secondary.(chain.Uint256)
}

func (db *IdxDB256I) Upperbound(secondary chain.Uint256) (SecondaryIterator, chain.Uint256) {
	it, _secondary := db.SecondaryDB.Upperbound(secondary)
	return it, _secondary.(chain.Uint256)
}
