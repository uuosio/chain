package database

/*
#include <stdint.h>

int32_t db_idx64_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const uint64_t* secondary);
void db_idx64_update(int32_t iterator, uint64_t payer, const uint64_t* secondary);
void db_idx64_remove(int32_t iterator);
int32_t db_idx64_next(int32_t iterator, uint64_t* primary);
int32_t db_idx64_previous(int32_t iterator, uint64_t* primary);
int32_t db_idx64_find_primary(uint64_t code, uint64_t scope, uint64_t table, uint64_t* secondary, uint64_t primary);
int32_t db_idx64_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const uint64_t* secondary, uint64_t* primary);
int32_t db_idx64_lowerbound(uint64_t code, uint64_t scope, uint64_t table, uint64_t* secondary, uint64_t* primary);
int32_t db_idx64_upperbound(uint64_t code, uint64_t scope, uint64_t table, uint64_t* secondary, uint64_t* primary);
int32_t db_idx64_end(uint64_t code, uint64_t scope, uint64_t table);
*/
import "C"
import (
	"chain"
)

type IdxDB64 struct {
	dbIndex int
	code    uint64
	scope   uint64
	table   uint64
	mi      *MultiIndex
}

func NewIdxDB64(index int, code uint64, scope uint64, table uint64) *IdxDB64 {
	v := &IdxDB64{index, code, scope, table, nil}
	return v
}

func (db *IdxDB64) GetIndex() int {
	return db.dbIndex
}

func (db *IdxDB64) Store(id uint64, secondary interface{}, payer uint64) SecondaryIterator {
	_secondary, ok := secondary.(uint64)
	chain.Check(ok, "bad secondary type")
	chain.Check(db.code == chain.CurrentReceiver().N, "bad code name")
	ret := C.db_idx64_store(db.scope, db.table, payer, id, &_secondary)
	return SecondaryIterator{ret, id}
}

func (db *IdxDB64) Update(it SecondaryIterator, secondary interface{}, payer uint64) {
	// db.mi.UpdateSecondaryValue(db, it.Primary, secondary, chain.Name{payer})
	_secondary := secondary.(uint64)
	C.db_idx64_update(it.I, payer, &_secondary)
}

// void db_idx64_remove(int32_t iterator);
func (db *IdxDB64) Remove(it SecondaryIterator) {
	C.db_idx64_remove(it.I)
}

// int32_t db_idx64_next(int32_t iterator, uint64_t* primary);
func (db *IdxDB64) Next(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx64_next(it.I, &primary)
	return SecondaryIterator{ret, primary}
}

// int32_t db_idx64_previous(int32_t iterator, uint64_t* primary);
func (db *IdxDB64) Previous(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx64_previous(it.I, (*C.uint64_t)(&primary))
	return SecondaryIterator{ret, primary}
}

// int32_t db_idx64_find_primary(uint64_t code, uint64_t scope, uint64_t table, uint64_t* secondary, uint64_t primary);
func (db *IdxDB64) FindByPrimary(primary uint64) (SecondaryIterator, interface{}) {
	var secondary uint64 = 0
	ret := C.db_idx64_find_primary(db.code, db.scope, db.table, &secondary, primary)
	return SecondaryIterator{ret, primary}, secondary
}

// int32_t db_idx64_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const uint64_t* secondary, uint64_t* primary);
func (db *IdxDB64) Find(secondary interface{}) SecondaryIterator {
	var primary uint64 = 0
	_secondary, ok := secondary.(uint64)
	chain.Check(ok, "bad secondary type")
	ret := C.db_idx64_find_secondary(db.code, db.scope, db.table, &_secondary, &primary)
	return SecondaryIterator{ret, primary}
}

// int32_t db_idx64_lowerbound(uint64_t code, uint64_t scope, uint64_t table, uint64_t* secondary, uint64_t* primary);
func (db *IdxDB64) Lowerbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(uint64)
	chain.Assert(ok, "bad secondary type")
	ret := C.db_idx64_lowerbound(db.code, db.scope, db.table, &_secondary, &primary)
	return SecondaryIterator{ret, primary}, _secondary
}

// int32_t db_idx64_upperbound(uint64_t code, uint64_t scope, uint64_t table, uint64_t* secondary, uint64_t* primary);
func (db *IdxDB64) Upperbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(uint64)
	chain.Assert(ok, "bad secondary type")
	ret := C.db_idx64_upperbound(db.code, db.scope, db.table, &_secondary, &primary)
	return SecondaryIterator{ret, primary}, _secondary
}

// int32_t db_idx64_end(uint64_t code, uint64_t scope, uint64_t table);
func (db *IdxDB64) End() SecondaryIterator {
	ret := C.db_idx64_end(db.code, db.scope, db.table)
	return SecondaryIterator{ret, 0}
}
