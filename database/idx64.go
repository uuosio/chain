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
	"github.com/uuosio/chain"
)

type IdxDB64 struct {
	dbIndex int
	code    uint64
	scope   uint64
	table   uint64
}

func NewIdxDB64(index int, code uint64, scope uint64, table uint64) *IdxDB64 {
	v := &IdxDB64{index, code, scope, table}
	return v
}

func (db *IdxDB64) GetIndex() int {
	return db.dbIndex
}

//Store an association of a 64-bit integer secondary key to a primary key in a secondary 64-bit integer index table
func (db *IdxDB64) Store(id uint64, secondary interface{}, payer uint64) SecondaryIterator {
	_secondary, ok := secondary.(uint64)
	chain.Check(ok, "bad secondary type")
	chain.Check(db.code == chain.CurrentReceiver().N, "bad code name")
	ret := C.db_idx64_store(db.scope, db.table, payer, id, &_secondary)
	return SecondaryIterator{ret, id, db.dbIndex}
}

//Update an association for a 64-bit integer secondary key to a primary key in a secondary 64-bit integer index table
func (db *IdxDB64) Update(it SecondaryIterator, secondary interface{}, payer uint64) {
	// db.mi.UpdateSecondaryValue(db, it.Primary, secondary, chain.Name{payer})
	_secondary := secondary.(uint64)
	C.db_idx64_update(it.I, payer, &_secondary)
}

//Remove a table row from a secondary 64-bit integer index table
func (db *IdxDB64) Remove(it SecondaryIterator) {
	C.db_idx64_remove(it.I)
}

//Find the table row following the referenced table row in a secondary 64-bit integer index table
func (db *IdxDB64) Next(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx64_next(it.I, &primary)
	return SecondaryIterator{ret, primary, db.dbIndex}
}

//Find the table row preceding the referenced table row in a secondary 64-bit integer index table
func (db *IdxDB64) Previous(it SecondaryIterator) SecondaryIterator {
	var primary uint64 = 0
	ret := C.db_idx64_previous(it.I, (*C.uint64_t)(&primary))
	return SecondaryIterator{ret, primary, db.dbIndex}
}

//Find a table row in a secondary 64-bit integer index table by primary key
func (db *IdxDB64) FindByPrimary(primary uint64) (SecondaryIterator, interface{}) {
	var secondary uint64 = 0
	ret := C.db_idx64_find_primary(db.code, db.scope, db.table, &secondary, primary)
	return SecondaryIterator{ret, primary, db.dbIndex}, secondary
}

//Find a table row in a secondary 64-bit integer index table by secondary key
func (db *IdxDB64) Find(secondary interface{}) SecondaryIterator {
	var primary uint64 = 0
	_secondary, ok := secondary.(uint64)
	chain.Check(ok, "bad secondary type")
	ret := C.db_idx64_find_secondary(db.code, db.scope, db.table, &_secondary, &primary)
	return SecondaryIterator{ret, primary, db.dbIndex}
}

//Find the table row in a secondary 64-bit integer index table that matches the lowerbound condition for a given secondary key
func (db *IdxDB64) Lowerbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(uint64)
	chain.Assert(ok, "bad secondary type")
	ret := C.db_idx64_lowerbound(db.code, db.scope, db.table, &_secondary, &primary)
	return SecondaryIterator{ret, primary, db.dbIndex}, _secondary
}

//Find the table row in a secondary 64-bit integer index table that matches the upperbound condition for a given secondary key
func (db *IdxDB64) Upperbound(secondary interface{}) (SecondaryIterator, interface{}) {
	var primary uint64 = 0
	_secondary, ok := secondary.(uint64)
	chain.Assert(ok, "bad secondary type")
	ret := C.db_idx64_upperbound(db.code, db.scope, db.table, &_secondary, &primary)
	return SecondaryIterator{ret, primary, db.dbIndex}, _secondary
}

//Get an end iterator representing just-past-the-end of the last table row of a secondary 64-bit integer index table
func (db *IdxDB64) End() SecondaryIterator {
	ret := C.db_idx64_end(db.code, db.scope, db.table)
	return SecondaryIterator{ret, 0, db.dbIndex}
}

type IdxDB64I struct {
	SecondaryDB
	I *IdxDB64
}

func (db *IdxDB64I) FindByPrimary(primary uint64) (SecondaryIterator, uint64) {
	it, _secondary := db.I.FindByPrimary(primary)
	return it, _secondary.(uint64)
}

func (db *IdxDB64I) Find(secondary uint64) SecondaryIterator {
	return db.I.Find(secondary)
}

func (db *IdxDB64I) Lowerbound(secondary uint64) (SecondaryIterator, uint64) {
	it, _secondary := db.I.Lowerbound(secondary)
	return it, _secondary.(uint64)
}

func (db *IdxDB64I) Upperbound(secondary uint64) (SecondaryIterator, uint64) {
	it, _secondary := db.I.Upperbound(secondary)
	return it, _secondary.(uint64)
}
