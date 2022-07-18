package database

import (
	"github.com/uuosio/chain"
	"github.com/uuosio/chain/eosio"
)

type IdxTableFloat64 struct {
	IdxTable
}

func NewIdxTableFloat64(index int, code uint64, scope uint64, table uint64) *IdxTableFloat64 {
	v := &IdxTableFloat64{IdxTable{index, code, scope, table}}
	return v
}

// Store an association of a double-precision floating-point secondary key to a primary key in a secondary double-precision floating-point index table
func (db *IdxTableFloat64) Store(id uint64, secondary float64, payer uint64) *SecondaryIterator {
	ret := eosio.DBIdxFloat64Store(db.scope, db.table, id, secondary, payer)
	return &SecondaryIterator{int32(ret), id, db.dbIndex}
}

func (db *IdxTableFloat64) StoreEx(id uint64, secondary interface{}, payer uint64) *SecondaryIterator {
	_secondary, ok := secondary.(float64)
	chain.Check(ok, "bad secondary type")
	return db.Store(id, _secondary, payer)
}

//Update an association for a double-precision floating-point secondary key to a primary key in a secondary double-precision floating-point index table
func (db *IdxTableFloat64) Update(it *SecondaryIterator, secondary float64, payer uint64) {
	eosio.DBIdxFloat64Update(it.I, secondary, payer)
}

func (db *IdxTableFloat64) UpdateEx(it *SecondaryIterator, secondary interface{}, payer uint64) {
	_secondary, ok := secondary.(float64)
	chain.Check(ok, "IdxFloat64: bad secondary type")
	db.Update(it, _secondary, payer)
}

//Remove a table row from a secondary double-precision floating-point index table
func (db *IdxTableFloat64) Remove(it *SecondaryIterator) {
	eosio.DBIdxFloat64Remove(it.I)
}

//Find the table row following the referenced table row in a secondary double-precision floating-point index table
func (db *IdxTableFloat64) Next(it *SecondaryIterator) *SecondaryIterator {
	i, primary := eosio.DBIdxFloat64Next(it.I)
	return &SecondaryIterator{i, primary, db.dbIndex}
}

//Find the table row preceding the referenced table row in a secondary double-precision floating-point index table
func (db *IdxTableFloat64) Previous(it *SecondaryIterator) *SecondaryIterator {
	i, primary := eosio.DBIdxFloat64Previous(it.I)
	return &SecondaryIterator{i, primary, db.dbIndex}
}

//Find a table row in a secondary double-precision floating-point index table by primary key
func (db *IdxTableFloat64) FindByPrimary(primary uint64) (*SecondaryIterator, interface{}) {
	i, secondary := eosio.DBIdxFloat64FindByPrimary(db.code, db.scope, db.table, primary)
	return &SecondaryIterator{i, primary, db.dbIndex}, secondary
}

//Find a table row in a secondary double-precision floating-point index table by secondary key
func (db *IdxTableFloat64) Find(secondary float64) *SecondaryIterator {
	// var primary uint64 = 0
	it, value := db.Lowerbound(secondary)
	if it.IsOk() {
		if value == secondary {
			return it
		}
	}
	return &SecondaryIterator{-1, 0, db.dbIndex}
}

func (db *IdxTableFloat64) FindEx(secondary interface{}) *SecondaryIterator {
	_secondary, ok := secondary.(float64)
	chain.Check(ok, "not a float64 value")
	return db.Find(_secondary)
}

//Find the table row in a secondary double-precision floating-point index table that matches the lowerbound condition for a given secondary key
func (db *IdxTableFloat64) Lowerbound(secondary float64) (*SecondaryIterator, float64) {
	i, secondary, primary := eosio.DBIdxFloat64Lowerbound(db.code, db.scope, db.table, secondary)
	return &SecondaryIterator{i, primary, db.dbIndex}, secondary
}

//Find the table row in a secondary double-precision floating-point index table that matches the upperbound condition for a given secondary key
func (db *IdxTableFloat64) Upperbound(secondary float64) (*SecondaryIterator, float64) {
	i, secondary, primary := eosio.DBIdxFloat64Upperbound(db.code, db.scope, db.table, secondary)
	return &SecondaryIterator{i, primary, db.dbIndex}, secondary
}

//Get an end iterator representing just-past-the-end of the last table row of a secondary double-precision floating-point index table
func (db *IdxTableFloat64) End() *SecondaryIterator {
	i := eosio.DBIdxFloat64End(db.code, db.scope, db.table)
	return &SecondaryIterator{i, 0, db.dbIndex}
}
