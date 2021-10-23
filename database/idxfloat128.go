package database

/*
#include "../structs.h"

int32_t db_idx_long_double_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const float128_t* secondary);
void db_idx_long_double_update(int32_t iterator, uint64_t payer, const float128_t* secondary);
void db_idx_long_double_remove(int32_t iterator);
int32_t db_idx_long_double_next(int32_t iterator, uint64_t* primary);
int32_t db_idx_long_double_previous(int32_t iterator, uint64_t* primary);
int32_t db_idx_long_double_find_primary(uint64_t code, uint64_t scope, uint64_t table, float128_t* secondary, uint64_t primary);
int32_t db_idx_long_double_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const float128_t* secondary, uint64_t* primary);
int32_t db_idx_long_double_lowerbound(uint64_t code, uint64_t scope, uint64_t table, float128_t* secondary, uint64_t* primary);
int32_t db_idx_long_double_upperbound(uint64_t code, uint64_t scope, uint64_t table, float128_t* secondary, uint64_t* primary);
int32_t db_idx_long_double_end(uint64_t code, uint64_t scope, uint64_t table);
*/
import "C"
import (
	"unsafe"

	"github.com/uuosio/chain"
)

type IdxDBFloat128 struct {
	dbIndex int
	code    uint64
	scope   uint64
	table   uint64
}

func NewIdxDBFloat128(index int, code uint64, scope uint64, table uint64) *IdxDBFloat128 {
	v := &IdxDBFloat128{index, code, scope, table}
	return v
}

func (db *IdxDBFloat128) GetIndex() int {
	return db.dbIndex
}

//Store an association of a quadruple-precision floating-point secondary key to a primary key in a secondary quadruple-precision floating-point index table
func (db *IdxDBFloat128) Store(id uint64, secondary interface{}, payer uint64) SecondaryIterator {
	_secondary, ok := secondary.(chain.Float128)
	chain.Check(ok, "bad secondary type")
	ret := C.db_idx_long_double_store(db.scope, db.table, payer, id, (*C.float128_t)(unsafe.Pointer(&_secondary)))
	return SecondaryIterator{ret, id, db.dbIndex}
}

//Update an association for a quadruple-precision floating-point secondary key to a primary key in a secondary quadruple-precision floating-point index table
func (db *IdxDBFloat128) Update(it SecondaryIterator, secondary interface{}, payer uint64) {
	_secondary, ok := secondary.(chain.Float128)
	chain.Check(ok, "bad secondary type")
	C.db_idx_long_double_update(it.I, payer, (*C.float128_t)(unsafe.Pointer(&_secondary)))
}

//Remove a table row from a secondary quadruple-precision floating-point index table
func (db *IdxDBFloat128) Remove(it SecondaryIterator) {
	C.db_idx_long_double_remove(it.I)
}

//Find the table row following the referenced table row in a secondary quadruple-precision floating-point index table
func (db *IdxDBFloat128) Next(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx_long_double_next(it.I, &primary)
	return SecondaryIterator{ret, primary, db.dbIndex}
}

//Find the table row preceding the referenced table row in a secondary quadruple-precision floating-point index table
func (db *IdxDBFloat128) Previous(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx_long_double_previous(it.I, &primary)
	return SecondaryIterator{ret, primary, db.dbIndex}
}

//Find a table row in a secondary quadruple-precision floating-point index table by primary key
func (db *IdxDBFloat128) FindByPrimary(primary uint64) (SecondaryIterator, interface{}) {
	var secondary chain.Float128
	ret := C.db_idx_long_double_find_primary(db.code, db.scope, db.table, (*C.float128_t)(unsafe.Pointer(&secondary)), primary)
	return SecondaryIterator{ret, primary, db.dbIndex}, secondary
}

//Find a table row in a secondary quadruple-precision floating-point index table by secondary key
func (db *IdxDBFloat128) Find(secondary interface{}) SecondaryIterator {
	// var primary uint64 = 0
	_secondary, ok := secondary.(chain.Float128)
	chain.Check(ok, "not a float128_t value")
	// ret := C.db_idx_long_double_find_secondary(db.code, db.scope, db.table, (*C.float128_t)(unsafe.Pointer(&_secondary)), &primary)
	// return SecondaryIterator{ret, primary, db.dbIndex}
	it, value := db.Lowerbound(secondary)
	if it.IsOk() {
		if value.(chain.Float128) == _secondary {
			return it
		}
	}
	return SecondaryIterator{-1, 0, db.dbIndex}
}

//Find the table row in a secondary quadruple-precision floating-point index table that matches the lowerbound condition for a given secondary key
func (db *IdxDBFloat128) Lowerbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(chain.Float128)
	chain.Assert(ok, "not a float128_t value")
	ret := C.db_idx_long_double_lowerbound(db.code, db.scope, db.table, (*C.float128_t)(unsafe.Pointer(&_secondary)), &primary)
	return SecondaryIterator{ret, primary, db.dbIndex}, _secondary
}

//Find the table row in a secondary quadruple-precision floating-point index table that matches the upperbound condition for a given secondary key
func (db *IdxDBFloat128) Upperbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(chain.Float128)
	chain.Assert(ok, "not a float128_t value")
	ret := C.db_idx_long_double_upperbound(db.code, db.scope, db.table, (*C.float128_t)(unsafe.Pointer(&_secondary)), &primary)
	return SecondaryIterator{ret, primary, db.dbIndex}, _secondary
}

//Get an end iterator representing just-past-the-end of the last table row of a secondary quadruple-precision floating-point index table
func (db *IdxDBFloat128) End() SecondaryIterator {
	ret := C.db_idx_long_double_end(db.code, db.scope, db.table)
	return SecondaryIterator{ret, 0, db.dbIndex}
}

type IdxDBFloat128I struct {
	SecondaryDB
}

func (db *IdxDBFloat128I) FindByPrimary(primary uint64) (SecondaryIterator, chain.Float128) {
	it, _secondary := db.SecondaryDB.FindByPrimary(primary)
	return it, _secondary.(chain.Float128)
}

func (db *IdxDBFloat128I) Find(secondary chain.Float128) SecondaryIterator {
	return db.SecondaryDB.Find(secondary)
}

func (db *IdxDBFloat128I) Lowerbound(secondary chain.Float128) (SecondaryIterator, chain.Float128) {
	it, _secondary := db.SecondaryDB.Lowerbound(secondary)
	return it, _secondary.(chain.Float128)
}

func (db *IdxDBFloat128I) Upperbound(secondary chain.Float128) (SecondaryIterator, chain.Float128) {
	it, _secondary := db.SecondaryDB.Upperbound(secondary)
	return it, _secondary.(chain.Float128)
}
