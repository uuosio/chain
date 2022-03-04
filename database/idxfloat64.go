package database

/*
#include "../structs.h"
typedef double Double;

int32_t db_idx_double_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const Double* secondary);
void db_idx_double_update(int32_t iterator, uint64_t payer, const Double* secondary);
void db_idx_double_remove(int32_t iterator);
int32_t db_idx_double_next(int32_t iterator, uint64_t* primary);
int32_t db_idx_double_previous(int32_t iterator, uint64_t* primary);
int32_t db_idx_double_find_primary(uint64_t code, uint64_t scope, uint64_t table, Double* secondary, uint64_t primary);
int32_t db_idx_double_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const Double* secondary, uint64_t* primary);
int32_t db_idx_double_lowerbound(uint64_t code, uint64_t scope, uint64_t table, Double* secondary, uint64_t* primary);
int32_t db_idx_double_upperbound(uint64_t code, uint64_t scope, uint64_t table, Double* secondary, uint64_t* primary);
int32_t db_idx_double_end(uint64_t code, uint64_t scope, uint64_t table);
*/
import "C"
import (
	"unsafe"

	"github.com/uuosio/chain"
)

type IdxDBFloat64 struct {
	IdxDB
}

func NewIdxDBFloat64(index int, code uint64, scope uint64, table uint64) *IdxDBFloat64 {
	v := &IdxDBFloat64{IdxDB{index, C.uint64_t(code), C.uint64_t(scope), C.uint64_t(table)}}
	return v
}

// Store an association of a double-precision floating-point secondary key to a primary key in a secondary double-precision floating-point index table
func (db *IdxDBFloat64) Store(id uint64, secondary float64, payer uint64) SecondaryIterator {
	ret := C.db_idx_double_store(db.scope, db.table, C.uint64_t(payer), C.uint64_t(id), (*C.Double)(&secondary))
	return SecondaryIterator{int32(ret), id, db.dbIndex}
}

func (db *IdxDBFloat64) StoreEx(id uint64, secondary interface{}, payer uint64) SecondaryIterator {
	_secondary, ok := secondary.(float64)
	chain.Check(ok, "bad secondary type")
	return db.Store(id, _secondary, payer)
}

//Update an association for a double-precision floating-point secondary key to a primary key in a secondary double-precision floating-point index table
func (db *IdxDBFloat64) Update(it SecondaryIterator, secondary float64, payer uint64) {
	C.db_idx_double_update(C.int32_t(it.I), C.uint64_t(payer), (*C.Double)(&secondary))
}

func (db *IdxDBFloat64) UpdateEx(it SecondaryIterator, secondary interface{}, payer uint64) {
	_secondary, ok := secondary.(float64)
	chain.Check(ok, "IdxFloat64: bad secondary type")
	db.Update(it, _secondary, payer)
}

//Remove a table row from a secondary double-precision floating-point index table
func (db *IdxDBFloat64) Remove(it SecondaryIterator) {
	C.db_idx_double_remove(C.int32_t(it.I))
}

//Find the table row following the referenced table row in a secondary double-precision floating-point index table
func (db *IdxDBFloat64) Next(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx_double_next(C.int32_t(it.I), (*C.uint64_t)(&primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}
}

//Find the table row preceding the referenced table row in a secondary double-precision floating-point index table
func (db *IdxDBFloat64) Previous(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx_double_previous(C.int32_t(it.I), (*C.uint64_t)(&primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}
}

//Find a table row in a secondary double-precision floating-point index table by primary key
func (db *IdxDBFloat64) FindByPrimary(primary uint64) (SecondaryIterator, interface{}) {
	var secondary float64
	ret := C.db_idx_double_find_primary(db.code, db.scope, db.table, (*C.Double)(unsafe.Pointer(&secondary)), C.uint64_t(primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}, secondary
}

//Find a table row in a secondary double-precision floating-point index table by secondary key
func (db *IdxDBFloat64) Find(secondary float64) SecondaryIterator {
	// var primary uint64 = 0
	it, value := db.Lowerbound(secondary)
	if it.IsOk() {
		if value == secondary {
			return it
		}
	}
	return SecondaryIterator{-1, 0, db.dbIndex}
}

func (db *IdxDBFloat64) FindEx(secondary interface{}) SecondaryIterator {
	_secondary, ok := secondary.(float64)
	chain.Check(ok, "not a float64 value")
	return db.Find(_secondary)
}

//Find the table row in a secondary double-precision floating-point index table that matches the lowerbound condition for a given secondary key
func (db *IdxDBFloat64) Lowerbound(secondary float64) (SecondaryIterator, float64) {
	var primary uint64 = 0
	ret := C.db_idx_double_lowerbound(db.code, db.scope, db.table, (*C.Double)(&secondary), (*C.uint64_t)(&primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}, secondary
}

//Find the table row in a secondary double-precision floating-point index table that matches the upperbound condition for a given secondary key
func (db *IdxDBFloat64) Upperbound(secondary float64) (SecondaryIterator, float64) {
	var primary uint64 = 0
	ret := C.db_idx_double_upperbound(db.code, db.scope, db.table, (*C.Double)(&secondary), (*C.uint64_t)(&primary))
	return SecondaryIterator{int32(ret), primary, db.dbIndex}, secondary
}

//Get an end iterator representing just-past-the-end of the last table row of a secondary double-precision floating-point index table
func (db *IdxDBFloat64) End() SecondaryIterator {
	ret := C.db_idx_double_end(db.code, db.scope, db.table)
	return SecondaryIterator{int32(ret), 0, db.dbIndex}
}
