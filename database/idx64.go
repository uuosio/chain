package database

import (
	"github.com/uuosio/chain"
	"github.com/uuosio/chain/eosio"
)

type IdxTable64 struct {
	IdxTable
}

func NewIdxTable64(index int, code uint64, scope uint64, table uint64) *IdxTable64 {
	v := &IdxTable64{IdxTable{index, code, scope, table}}
	return v
}

//Store an association of a 64-bit integer secondary key to a primary key in a secondary 64-bit integer index table
func (db *IdxTable64) Store(id uint64, secondary uint64, payer uint64) *SecondaryIterator {
	chain.Check(uint64(db.code) == chain.CurrentReceiver().N, "bad code name")
	ret := eosio.DBIdx64Store(db.scope, db.table, id, secondary, payer)
	return &SecondaryIterator{int32(ret), id, db.dbIndex}
}

func (db *IdxTable64) StoreEx(id uint64, secondary interface{}, payer uint64) *SecondaryIterator {
	_secondary, ok := secondary.(uint64)
	chain.Check(ok, "IdxTable64.Store: bad secondary type")
	return db.Store(id, _secondary, payer)
}

//Update an association for a 64-bit integer secondary key to a primary key in a secondary 64-bit integer index table
func (db *IdxTable64) Update(it *SecondaryIterator, secondary uint64, payer uint64) {
	eosio.DBIdx64Update(it.I, secondary, payer)
}

func (db *IdxTable64) UpdateEx(it *SecondaryIterator, secondary interface{}, payer uint64) {
	_secondary := secondary.(uint64)
	db.Update(it, _secondary, payer)
}

//Remove a table row from a secondary 64-bit integer index table
func (db *IdxTable64) Remove(it *SecondaryIterator) {
	eosio.DBIdx64Remove(it.I)
}

//Find the table row following the referenced table row in a secondary 64-bit integer index table
func (db *IdxTable64) Next(it *SecondaryIterator) *SecondaryIterator {
	i, primary := eosio.DBIdx64Next(it.I)
	return &SecondaryIterator{i, primary, db.dbIndex}
}

//Find the table row preceding the referenced table row in a secondary 64-bit integer index table
func (db *IdxTable64) Previous(it *SecondaryIterator) *SecondaryIterator {
	i, primary := eosio.DBIdx64Previous(it.I)
	return &SecondaryIterator{i, primary, db.dbIndex}
}

//Find a table row in a secondary 64-bit integer index table by primary key
func (db *IdxTable64) FindByPrimary(primary uint64) (*SecondaryIterator, interface{}) {
	var secondary uint64 = 0
	it, secondary := eosio.DBIdx64FindByPrimary(db.code, db.scope, db.table, primary)
	if it >= 0 {
		return &SecondaryIterator{int32(it), primary, db.dbIndex}, secondary
	} else {
		return &SecondaryIterator{it, 0, db.dbIndex}, 0
	}
}

//Find a table row in a secondary 64-bit integer index table by secondary key
func (db *IdxTable64) Find(secondary uint64) *SecondaryIterator {
	//ret := C.db_idx64_find_secondary(db.code, db.scope, db.table, &_secondary, &primary)
	// return SecondaryIterator{ret, primary, db.dbIndex}
	it, value := db.Lowerbound(secondary)
	if it.IsOk() {
		if value == secondary {
			return it
		}
	}
	return &SecondaryIterator{-1, 0, db.dbIndex}
}

func (db *IdxTable64) FindEx(secondary interface{}) *SecondaryIterator {
	_secondary, ok := secondary.(uint64)
	chain.Check(ok, "IdxTable64.Find: bad secondary type")
	return db.FindEx(_secondary)
}

//Find the table row in a secondary 64-bit integer index table that matches the lowerbound condition for a given secondary key
func (db *IdxTable64) Lowerbound(secondary uint64) (*SecondaryIterator, uint64) {
	it, secondary, primary := eosio.DBIdx64Lowerbound(db.code, db.scope, db.table, secondary)
	return &SecondaryIterator{int32(it), primary, db.dbIndex}, secondary
}

//Find the table row in a secondary 64-bit integer index table that matches the upperbound condition for a given secondary key
func (db *IdxTable64) Upperbound(secondary uint64) (*SecondaryIterator, uint64) {
	it, secondary, primary := eosio.DBIdx64Upperbound(db.code, db.scope, db.table, secondary)
	return &SecondaryIterator{int32(it), primary, db.dbIndex}, secondary
}

//Get an end iterator representing just-past-the-end of the last table row of a secondary 64-bit integer index table
func (db *IdxTable64) End() *SecondaryIterator {
	ret := eosio.DBIdx64End(db.code, db.scope, db.table)
	return &SecondaryIterator{int32(ret), 0, db.dbIndex}
}
