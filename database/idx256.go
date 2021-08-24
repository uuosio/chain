package database

/*
#include <stdint.h>

typedef uint8_t uint128_t; //fake uint128_t definition

int32_t db_idx256_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const uint128_t* data, uint32_t data_len );
void db_idx256_update(int32_t iterator, uint64_t payer, const uint128_t* data, uint32_t data_len);
void db_idx256_remove(int32_t iterator);
int32_t db_idx256_next(int32_t iterator, uint64_t* primary);
int32_t db_idx256_previous(int32_t iterator, uint64_t* primary);
int32_t db_idx256_find_primary(uint64_t code, uint64_t scope, uint64_t table, uint128_t* data, uint32_t data_len, uint64_t primary);
int32_t db_idx256_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const uint128_t* data, uint32_t data_len, uint64_t* primary);
int32_t db_idx256_lowerbound(uint64_t code, uint64_t scope, uint64_t table, uint128_t* data, uint32_t data_len, uint64_t* primary);
int32_t db_idx256_upperbound(uint64_t code, uint64_t scope, uint64_t table, uint128_t* data, uint32_t data_len, uint64_t* primary);
int32_t db_idx256_end(uint64_t code, uint64_t scope, uint64_t table);
*/
import "C"
import (
	"chain"
	"math/big"
	"unsafe"
)

type IdxDB256 struct {
	dbIndex int
	code    uint64
	scope   uint64
	table   uint64
}

func NewIdxDB256(index int, code uint64, scope uint64, table uint64) *IdxDB256 {
	v := &IdxDB256{index, code, scope, table}
	return v
}

func (db *IdxDB256) GetIndex() int {
	return db.dbIndex
}

// int32_t db_idx256_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const uint128_t* data, uint32_t data_len );
func (db *IdxDB256) Store(id uint64, secondary interface{}, payer uint64) SecondaryIterator {
	_secondary, ok := secondary.(chain.Uint256)
	chain.Check(ok, "bad secondary type")
	ret := C.db_idx256_store(db.scope, db.table, payer, id, (*uint8)(unsafe.Pointer(&_secondary)), 2)
	return SecondaryIterator{ret, id}
}

// void db_idx256_update(int32_t iterator, uint64_t payer, const uint128_t* data, uint32_t data_len);
func (db *IdxDB256) Update(it SecondaryIterator, secondary interface{}, payer uint64) {
	_secondary, ok := secondary.(chain.Uint256)
	chain.Check(ok, "bad secondary type")
	C.db_idx256_update(it.I, payer, (*uint8)(unsafe.Pointer(&_secondary)), 2)
}

// void db_idx256_remove(int32_t iterator);
func (db *IdxDB256) Remove(it SecondaryIterator) {
	C.db_idx256_remove(it.I)
}

// int32_t db_idx256_next(int32_t iterator, uint64_t* primary);
func (db *IdxDB256) Next(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx256_next(it.I, &primary)
	return SecondaryIterator{ret, primary}
}

// int32_t db_idx256_previous(int32_t iterator, uint64_t* primary);
func (db *IdxDB256) Previous(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx256_previous(it.I, &primary)
	return SecondaryIterator{ret, primary}
}

// int32_t db_idx256_find_primary(uint64_t code, uint64_t scope, uint64_t table, uint128_t* data, uint32_t data_len, uint64_t primary);
func (db *IdxDB256) FindByPrimary(primary uint64) (SecondaryIterator, interface{}) {
	var secondary chain.Uint256
	ret := C.db_idx256_find_primary(db.code, db.scope, db.table, (*uint8)(unsafe.Pointer(&secondary)), 2, primary)
	return SecondaryIterator{ret, primary}, secondary
}

// int32_t db_idx256_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const uint128_t* data, uint32_t data_len, uint64_t* primary);
func (db *IdxDB256) Find(secondary interface{}) SecondaryIterator {
	var primary uint64 = 0
	var rawSecondary chain.Uint256
	switch v := secondary.(type) {
	case big.Int:
		_secondary := v.Bytes()
		chain.Check(len(_secondary) <= 32, "Uint256 overflow")
		copy(rawSecondary[:], _secondary)
	case chain.Uint256:
		copy(rawSecondary[:], v[:])
	default:
		panic("can not conver to uint128 type")
	}
	ret := C.db_idx256_find_secondary(db.code, db.scope, db.table, (*uint8)(unsafe.Pointer(&rawSecondary[0])), 2, &primary)
	return SecondaryIterator{ret, primary}
}

// int32_t db_idx256_lowerbound(uint64_t code, uint64_t scope, uint64_t table, uint128_t* data, uint32_t data_len, uint64_t* primary);
func (db *IdxDB256) Lowerbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(chain.Uint256)
	chain.Assert(ok, "not a Uint256 value")
	ret := C.db_idx256_lowerbound(db.code, db.scope, db.table, (*uint8)(unsafe.Pointer(&_secondary)), 2, &primary)
	return SecondaryIterator{ret, primary}, _secondary
}

// int32_t db_idx256_upperbound(uint64_t code, uint64_t scope, uint64_t table, uint128_t* data, uint32_t data_len, uint64_t* primary);
func (db *IdxDB256) Upperbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(chain.Uint256)
	chain.Assert(ok, "not a Uint256 value")
	ret := C.db_idx256_upperbound(db.code, db.scope, db.table, (*uint8)(unsafe.Pointer(&_secondary)), 2, &primary)
	return SecondaryIterator{ret, primary}, _secondary
}

// int32_t db_idx256_end(uint64_t code, uint64_t scope, uint64_t table);
func (db *IdxDB256) End() SecondaryIterator {
	ret := C.db_idx256_end(db.code, db.scope, db.table)
	return SecondaryIterator{ret, 0}
}
