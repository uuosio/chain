package database

/*
#include <stdint.h>
typedef uint8_t Uint128_t; //fake Uint128_t definition

int32_t db_idx128_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const Uint128_t* secondary);
void db_idx128_update(int32_t iterator, uint64_t payer, const Uint128_t* secondary);
void db_idx128_remove(int32_t iterator);
int32_t db_idx128_next(int32_t iterator, uint64_t* primary);
int32_t db_idx128_previous(int32_t iterator, uint64_t* primary);
int32_t db_idx128_find_primary(uint64_t code, uint64_t scope, uint64_t table, Uint128_t* secondary, uint64_t primary);
int32_t db_idx128_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const Uint128_t* secondary, uint64_t* primary);
int32_t db_idx128_lowerbound(uint64_t code, uint64_t scope, uint64_t table, Uint128_t* secondary, uint64_t* primary);
int32_t db_idx128_upperbound(uint64_t code, uint64_t scope, uint64_t table, Uint128_t* secondary, uint64_t* primary);
int32_t db_idx128_end(uint64_t code, uint64_t scope, uint64_t table);
*/
import "C"
import (
	"chain"
	"math/big"
	"unsafe"
)

type IdxDB128 struct {
	dbIndex int
	code    uint64
	scope   uint64
	table   uint64
}

func (db *IdxDB128) GetIndex() int {
	return db.dbIndex
}

func NewIdxDB128(index int, code uint64, scope uint64, table uint64) *IdxDB128 {
	v := &IdxDB128{index, code, scope, table}
	return v
}

// int32_t db_idx128_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const Uint128_t* secondary);
func (db *IdxDB128) Store(id uint64, secondary interface{}, payer uint64) SecondaryIterator {
	_secondary, ok := secondary.(chain.Uint128)
	chain.Check(ok, "bad secondary type")
	ret := C.db_idx128_store(db.scope, db.table, payer, id, (*uint8)(unsafe.Pointer(&_secondary)))
	return SecondaryIterator{ret, id}
}

// void db_idx128_update(int32_t iterator, uint64_t payer, const Uint128_t* secondary);
func (db *IdxDB128) Update(it SecondaryIterator, secondary interface{}, payer uint64) {
	_secondary, ok := secondary.(chain.Uint128)
	chain.Check(ok, "bad secondary type")
	C.db_idx128_update(it.I, payer, (*uint8)(unsafe.Pointer(&_secondary)))
}

// void db_idx128_remove(int32_t iterator);
func (db *IdxDB128) Remove(it SecondaryIterator) {
	C.db_idx128_remove(it.I)
}

// int32_t db_idx128_next(int32_t iterator, uint64_t* primary);
func (db *IdxDB128) Next(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx128_next(it.I, &primary)
	return SecondaryIterator{ret, primary}
}

// int32_t db_idx128_previous(int32_t iterator, uint64_t* primary);
func (db *IdxDB128) Previous(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx128_previous(it.I, &primary)
	return SecondaryIterator{ret, primary}
}

// int32_t db_idx128_find_primary(uint64_t code, uint64_t scope, uint64_t table, Uint128_t* secondary, uint64_t primary);
func (db *IdxDB128) FindByPrimary(primary uint64) (SecondaryIterator, interface{}) {
	var secondary chain.Uint128
	ret := C.db_idx128_find_primary(db.code, db.scope, db.table, (*uint8)(unsafe.Pointer(&secondary)), primary)
	return SecondaryIterator{ret, primary}, secondary
}

// int32_t db_idx128_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const Uint128_t* secondary, uint64_t* primary);
func (db *IdxDB128) Find(secondary interface{}) SecondaryIterator {
	var primary uint64 = 0
	var rawSecondary chain.Uint128
	switch v := secondary.(type) {
	case big.Int:
		_secondary := v.Bytes()
		chain.Check(len(_secondary) <= 16, "Uint128 overflow")
		copy(rawSecondary[:], _secondary)
	case chain.Uint128:
		copy(rawSecondary[:], v[:])
	default:
		panic("can not conver to uint128 type")
	}
	// _secondary, ok := secondary.(Uint128)
	// chain.Check(ok, "bad secondary type")
	ret := C.db_idx128_find_secondary(db.code, db.scope, db.table, (*uint8)(unsafe.Pointer(&rawSecondary[0])), &primary)
	return SecondaryIterator{ret, primary}
}

// int32_t db_idx128_lowerbound(uint64_t code, uint64_t scope, uint64_t table, Uint128_t* secondary, uint64_t* primary);
func (db *IdxDB128) Lowerbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(chain.Uint128)
	chain.Assert(ok, "bad secondary type")
	ret := C.db_idx128_lowerbound(db.code, db.scope, db.table, (*uint8)(unsafe.Pointer(&_secondary)), &primary)
	return SecondaryIterator{ret, primary}, _secondary
}

// int32_t db_idx128_upperbound(uint64_t code, uint64_t scope, uint64_t table, Uint128_t* secondary, uint64_t* primary);
func (db *IdxDB128) Upperbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(chain.Uint128)
	chain.Assert(ok, "bad secondary type")
	ret := C.db_idx128_upperbound(db.code, db.scope, db.table, (*uint8)(unsafe.Pointer(&_secondary)), &primary)
	return SecondaryIterator{ret, primary}, _secondary
}

// int32_t db_idx128_end(uint64_t code, uint64_t scope, uint64_t table);
func (db *IdxDB128) End() SecondaryIterator {
	ret := C.db_idx128_end(db.code, db.scope, db.table)
	return SecondaryIterator{ret, 0}
}
