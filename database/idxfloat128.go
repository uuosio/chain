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
	IdxDB
}

func NewIdxDBFloat128(index int, code uint64, scope uint64, table uint64) *IdxDBFloat128 {
	v := &IdxDBFloat128{IdxDB{index, C.uint64_t(code), C.uint64_t(scope), C.uint64_t(table)}}
	return v
}

//Store an association of a quadruple-precision floating-point secondary key to a primary key in a secondary quadruple-precision floating-point index table
func (db *IdxDBFloat128) Store(id uint64, secondary chain.Float128, payer uint64) *SecondaryIterator {
	ret := C.db_idx_long_double_store(db.scope, db.table, C.uint64_t(payer), C.uint64_t(id), (*C.float128_t)(unsafe.Pointer(&secondary)))
	return &SecondaryIterator{int32(ret), id, db.dbIndex}
}

func (db *IdxDBFloat128) StoreEx(id uint64, secondary interface{}, payer uint64) *SecondaryIterator {
	_secondary, ok := secondary.(chain.Float128)
	chain.Check(ok, "bad secondary type")
	return db.Store(id, _secondary, payer)
}

//Update an association for a quadruple-precision floating-point secondary key to a primary key in a secondary quadruple-precision floating-point index table
func (db *IdxDBFloat128) Update(it *SecondaryIterator, secondary chain.Float128, payer uint64) {
	C.db_idx_long_double_update(C.int32_t(it.I), C.uint64_t(payer), (*C.float128_t)(unsafe.Pointer(&secondary)))
}

func (db *IdxDBFloat128) UpdateEx(it *SecondaryIterator, secondary interface{}, payer uint64) {
	_secondary, ok := secondary.(chain.Float128)
	chain.Check(ok, "bad secondary type")
	db.Update(it, _secondary, payer)
}

//Remove a table row from a secondary quadruple-precision floating-point index table
func (db *IdxDBFloat128) Remove(it *SecondaryIterator) {
	C.db_idx_long_double_remove(C.int32_t(it.I))
}

//Find the table row following the referenced table row in a secondary quadruple-precision floating-point index table
func (db *IdxDBFloat128) Next(it *SecondaryIterator) *SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx_long_double_next(C.int32_t(it.I), (*C.uint64_t)(&primary))
	return &SecondaryIterator{int32(ret), primary, db.dbIndex}
}

//Find the table row preceding the referenced table row in a secondary quadruple-precision floating-point index table
func (db *IdxDBFloat128) Previous(it *SecondaryIterator) *SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx_long_double_previous(C.int32_t(it.I), (*C.uint64_t)(&primary))
	return &SecondaryIterator{int32(ret), primary, db.dbIndex}
}

//Find a table row in a secondary quadruple-precision floating-point index table by primary key
func (db *IdxDBFloat128) FindByPrimary(primary uint64) (*SecondaryIterator, interface{}) {
	var secondary chain.Float128
	ret := C.db_idx_long_double_find_primary(db.code, db.scope, db.table, (*C.float128_t)(unsafe.Pointer(&secondary)), C.uint64_t(primary))
	return &SecondaryIterator{int32(ret), primary, db.dbIndex}, secondary
}

//Find a table row in a secondary quadruple-precision floating-point index table by secondary key
func (db *IdxDBFloat128) Find(secondary chain.Float128) *SecondaryIterator {
	// var primary uint64 = 0
	// ret := C.db_idx_long_double_find_secondary(db.code, db.scope, db.table, (*C.float128_t)(unsafe.Pointer(&_secondary)), &primary)
	// return &SecondaryIterator{ret, primary, db.dbIndex}
	it, value := db.Lowerbound(secondary)
	if it.IsOk() {
		if value == secondary {
			return it
		}
	}
	return &SecondaryIterator{-1, 0, db.dbIndex}
}

func (db *IdxDBFloat128) FindEx(secondary interface{}) *SecondaryIterator {
	_secondary, ok := secondary.(chain.Float128)
	chain.Check(ok, "not a float128_t value")
	return db.Find(_secondary)
}

//Find the table row in a secondary quadruple-precision floating-point index table that matches the lowerbound condition for a given secondary key
func (db *IdxDBFloat128) Lowerbound(secondary chain.Float128) (*SecondaryIterator, chain.Float128) {
	var primary uint64 = 0
	ret := C.db_idx_long_double_lowerbound(db.code, db.scope, db.table, (*C.float128_t)(unsafe.Pointer(&secondary)), (*C.uint64_t)(&primary))
	return &SecondaryIterator{int32(ret), primary, db.dbIndex}, secondary
}

//Find the table row in a secondary quadruple-precision floating-point index table that matches the upperbound condition for a given secondary key
func (db *IdxDBFloat128) Upperbound(secondary chain.Float128) (*SecondaryIterator, chain.Float128) {
	var primary uint64 = 0
	ret := C.db_idx_long_double_upperbound(db.code, db.scope, db.table, (*C.float128_t)(unsafe.Pointer(&secondary)), (*C.uint64_t)(&primary))
	return &SecondaryIterator{int32(ret), primary, db.dbIndex}, secondary
}

//Get an end iterator representing just-past-the-end of the last table row of a secondary quadruple-precision floating-point index table
func (db *IdxDBFloat128) End() *SecondaryIterator {
	ret := C.db_idx_long_double_end(db.code, db.scope, db.table)
	return &SecondaryIterator{int32(ret), 0, db.dbIndex}
}
