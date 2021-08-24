package database

/*
#include <stdint.h>
typedef uint8_t long_double; //fake uint128_t definition

int32_t db_idx_long_double_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const long_double* secondary);
void db_idx_long_double_update(int32_t iterator, uint64_t payer, const long_double* secondary);
void db_idx_long_double_remove(int32_t iterator);
int32_t db_idx_long_double_next(int32_t iterator, uint64_t* primary);
int32_t db_idx_long_double_previous(int32_t iterator, uint64_t* primary);
int32_t db_idx_long_double_find_primary(uint64_t code, uint64_t scope, uint64_t table, long_double* secondary, uint64_t primary);
int32_t db_idx_long_double_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const long_double* secondary, uint64_t* primary);
int32_t db_idx_long_double_lowerbound(uint64_t code, uint64_t scope, uint64_t table, long_double* secondary, uint64_t* primary);
int32_t db_idx_long_double_upperbound(uint64_t code, uint64_t scope, uint64_t table, long_double* secondary, uint64_t* primary);
int32_t db_idx_long_double_end(uint64_t code, uint64_t scope, uint64_t table);
*/
import "C"
import (
	"chain"
	"unsafe"
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

// int32_t db_idx_long_double_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const long_double* secondary);
func (db *IdxDBFloat128) Store(id uint64, secondary interface{}, payer uint64) SecondaryIterator {
	_secondary, ok := secondary.(chain.Float128)
	chain.Check(ok, "bad secondary type")
	ret := C.db_idx_long_double_store(db.scope, db.table, payer, id, (*uint8)(unsafe.Pointer(&_secondary)))
	return SecondaryIterator{ret, id}
}

// void db_idx_long_double_update(int32_t iterator, uint64_t payer, const long_double* secondary);
func (db *IdxDBFloat128) Update(it SecondaryIterator, secondary interface{}, payer uint64) {
	_secondary, ok := secondary.(chain.Float128)
	chain.Check(ok, "bad secondary type")
	C.db_idx_long_double_update(it.I, payer, (*uint8)(unsafe.Pointer(&_secondary)))
}

// void db_idx_long_double_remove(int32_t iterator);
func (db *IdxDBFloat128) Remove(it SecondaryIterator) {
	C.db_idx_long_double_remove(it.I)
}

// int32_t db_idx_long_double_next(int32_t iterator, uint64_t* primary);
func (db *IdxDBFloat128) Next(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx_long_double_next(it.I, &primary)
	return SecondaryIterator{ret, primary}
}

// int32_t db_idx_long_double_previous(int32_t iterator, uint64_t* primary);
func (db *IdxDBFloat128) Previous(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx_long_double_previous(it.I, &primary)
	return SecondaryIterator{ret, primary}
}

// int32_t db_idx_long_double_find_primary(uint64_t code, uint64_t scope, uint64_t table, long_double* secondary, uint64_t primary);
func (db *IdxDBFloat128) FindByPrimary(primary uint64) (SecondaryIterator, interface{}) {
	var secondary chain.Float128
	ret := C.db_idx_long_double_find_primary(db.code, db.scope, db.table, (*uint8)(unsafe.Pointer(&secondary)), primary)
	return SecondaryIterator{ret, primary}, secondary
}

// int32_t db_idx_long_double_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const long_double* secondary, uint64_t* primary);
func (db *IdxDBFloat128) Find(secondary interface{}) SecondaryIterator {
	var primary uint64 = 0
	_secondary, ok := secondary.(chain.Float128)
	chain.Check(ok, "not a float128 value")
	ret := C.db_idx_long_double_find_secondary(db.code, db.scope, db.table, (*uint8)(unsafe.Pointer(&_secondary)), &primary)
	return SecondaryIterator{ret, primary}
}

// int32_t db_idx_long_double_lowerbound(uint64_t code, uint64_t scope, uint64_t table, long_double* secondary, uint64_t* primary);
func (db *IdxDBFloat128) Lowerbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(chain.Float128)
	chain.Assert(ok, "not a float128 value")
	ret := C.db_idx_long_double_lowerbound(db.code, db.scope, db.table, (*uint8)(unsafe.Pointer(&_secondary)), &primary)
	return SecondaryIterator{ret, primary}, _secondary
}

// int32_t db_idx_long_double_upperbound(uint64_t code, uint64_t scope, uint64_t table, long_double* secondary, uint64_t* primary);
func (db *IdxDBFloat128) Upperbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(chain.Float128)
	chain.Assert(ok, "not a float128 value")
	ret := C.db_idx_long_double_upperbound(db.code, db.scope, db.table, (*uint8)(unsafe.Pointer(&_secondary)), &primary)
	return SecondaryIterator{ret, primary}, _secondary
}

// int32_t db_idx_long_double_end(uint64_t code, uint64_t scope, uint64_t table);
func (db *IdxDBFloat128) End() SecondaryIterator {
	ret := C.db_idx_long_double_end(db.code, db.scope, db.table)
	return SecondaryIterator{ret, 0}
}
