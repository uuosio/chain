package database

/*
#include <stdint.h>
typedef uint8_t uint128_t; //fake uint128_t definition

int32_t db_idx_double_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const double* secondary);
void db_idx_double_update(int32_t iterator, uint64_t payer, const double* secondary);
void db_idx_double_remove(int32_t iterator);
int32_t db_idx_double_next(int32_t iterator, uint64_t* primary);
int32_t db_idx_double_previous(int32_t iterator, uint64_t* primary);
int32_t db_idx_double_find_primary(uint64_t code, uint64_t scope, uint64_t table, double* secondary, uint64_t primary);
int32_t db_idx_double_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const double* secondary, uint64_t* primary);
int32_t db_idx_double_lowerbound(uint64_t code, uint64_t scope, uint64_t table, double* secondary, uint64_t* primary);
int32_t db_idx_double_upperbound(uint64_t code, uint64_t scope, uint64_t table, double* secondary, uint64_t* primary);
int32_t db_idx_double_end(uint64_t code, uint64_t scope, uint64_t table);
*/
import "C"
import (
	"chain"
	"unsafe"
)

type IdxDBFloat64 struct {
	dbIndex int
	code    uint64
	scope   uint64
	table   uint64
}

func NewIdxDBFloat64(index int, code uint64, scope uint64, table uint64) *IdxDBFloat64 {
	v := &IdxDBFloat64{index, code, scope, table}
	return v
}

func (db *IdxDBFloat64) GetIndex() int {
	return db.dbIndex
}

// int32_t db_idx_double_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const double* secondary);
func (db *IdxDBFloat64) Store(id uint64, secondary interface{}, payer uint64) SecondaryIterator {
	_secondary, ok := secondary.(float64)
	chain.Check(ok, "bad secondary type")
	ret := C.db_idx_double_store(db.scope, db.table, payer, id, &_secondary)
	return SecondaryIterator{ret, id}
}

// void db_idx_double_update(int32_t iterator, uint64_t payer, const double* secondary);
func (db *IdxDBFloat64) Update(it SecondaryIterator, secondary interface{}, payer uint64) {
	_secondary, ok := secondary.(float64)
	chain.Check(ok, "IdxFloat64: bad secondary type")
	C.db_idx_double_update(it.I, payer, &_secondary)
}

// void db_idx_double_remove(int32_t iterator);
func (db *IdxDBFloat64) Remove(it SecondaryIterator) {
	C.db_idx_double_remove(it.I)
}

// int32_t db_idx_double_next(int32_t iterator, uint64_t* primary);
func (db *IdxDBFloat64) Next(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx_double_next(it.I, &primary)
	return SecondaryIterator{ret, primary}
}

// int32_t db_idx_double_previous(int32_t iterator, uint64_t* primary);
func (db *IdxDBFloat64) Previous(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx_double_previous(it.I, &primary)
	return SecondaryIterator{ret, primary}
}

// int32_t db_idx_double_find_primary(uint64_t code, uint64_t scope, uint64_t table, double* secondary, uint64_t primary);
func (db *IdxDBFloat64) FindByPrimary(primary uint64) (SecondaryIterator, interface{}) {
	var secondary float64
	ret := C.db_idx_double_find_primary(db.code, db.scope, db.table, (*float64)(unsafe.Pointer(&secondary)), primary)
	return SecondaryIterator{ret, primary}, secondary
}

// int32_t db_idx_double_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const double* secondary, uint64_t* primary);
func (db *IdxDBFloat64) Find(secondary interface{}) SecondaryIterator {
	var primary uint64 = 0
	_secondary, ok := secondary.(float64)
	chain.Check(ok, "not a float64 value")
	ret := C.db_idx_double_find_secondary(db.code, db.scope, db.table, &_secondary, &primary)
	return SecondaryIterator{ret, primary}
}

// int32_t db_idx_double_lowerbound(uint64_t code, uint64_t scope, uint64_t table, double* secondary, uint64_t* primary);
func (db *IdxDBFloat64) Lowerbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(float64)
	chain.Assert(ok, "not a float64 value")
	ret := C.db_idx_double_lowerbound(db.code, db.scope, db.table, &_secondary, &primary)
	return SecondaryIterator{ret, primary}, _secondary
}

// int32_t db_idx_double_upperbound(uint64_t code, uint64_t scope, uint64_t table, double* secondary, uint64_t* primary);
func (db *IdxDBFloat64) Upperbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(float64)
	chain.Assert(ok, "not a float64 value")
	ret := C.db_idx_double_upperbound(db.code, db.scope, db.table, &_secondary, &primary)
	return SecondaryIterator{ret, primary}, _secondary
}

// int32_t db_idx_double_end(uint64_t code, uint64_t scope, uint64_t table);
func (db *IdxDBFloat64) End() SecondaryIterator {
	ret := C.db_idx_double_end(db.code, db.scope, db.table)
	return SecondaryIterator{ret, 0}
}
