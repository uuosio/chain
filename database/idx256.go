package database

import (
	"github.com/uuosio/chain"
	"github.com/uuosio/chain/eosio"
)

type IdxTable256 struct {
	IdxTable
}

func NewIdxTable256(index int, code uint64, scope uint64, table uint64) *IdxTable256 {
	v := &IdxTable256{IdxTable{index, code, scope, table}}
	return v
}

//Store an association of a 256-bit secondary key to a primary key in a secondary 256-bit index table
func (db *IdxTable256) Store(id uint64, secondary chain.Uint256, payer uint64) *SecondaryIterator {
	ret := eosio.DBIdx256Store(db.scope, db.table, id, secondary, payer)
	return &SecondaryIterator{int32(ret), id, db.dbIndex}
}

func (db *IdxTable256) StoreEx(id uint64, secondary interface{}, payer uint64) *SecondaryIterator {
	_secondary, ok := secondary.(chain.Uint256)
	chain.Check(ok, "bad secondary type")
	return db.Store(id, _secondary, payer)
}

//Update an association for a 256-bit secondary key to a primary key in a secondary 256-bit index table
func (db *IdxTable256) Update(it *SecondaryIterator, secondary chain.Uint256, payer uint64) {
	eosio.DBIdx256Update(it.I, secondary, payer)
}

func (db *IdxTable256) UpdateEx(it *SecondaryIterator, secondary interface{}, payer uint64) {
	_secondary, ok := secondary.(chain.Uint256)
	chain.Check(ok, "bad secondary type")
	db.Update(it, _secondary, payer)
}

//Remove a table row from a secondary 256-bit index table
func (db *IdxTable256) Remove(it *SecondaryIterator) {
	eosio.DBIdx256Remove(it.I)
}

//Find the table row following the referenced table row in a secondary 256-bit index table
func (db *IdxTable256) Next(it *SecondaryIterator) *SecondaryIterator {
	i, primary := eosio.DBIdx256Next(it.I)
	return &SecondaryIterator{i, primary, db.dbIndex}
}

//Find the table row preceding the referenced table row in a secondary 256-bit index table
func (db *IdxTable256) Previous(it *SecondaryIterator) *SecondaryIterator {
	i, primary := eosio.DBIdx256Previous(it.I)
	return &SecondaryIterator{i, primary, db.dbIndex}
}

//Find a table row in a secondary 128-bit integer index table by primary key
func (db *IdxTable256) FindByPrimary(primary uint64) (*SecondaryIterator, interface{}) {
	it, secondary := eosio.DBIdx256FindByPrimary(db.code, db.scope, db.table, primary)
	return &SecondaryIterator{it, primary, db.dbIndex}, chain.Uint256(secondary)
}

// Find a table row in a secondary 256-bit index table by secondary key
func (db *IdxTable256) Find(secondary chain.Uint256) *SecondaryIterator {
	it, value := db.Lowerbound(secondary)
	if it.IsOk() {
		if value == secondary {
			return it
		}
	}
	return &SecondaryIterator{-1, 0, db.dbIndex}
}

func (db *IdxTable256) FindEx(secondary interface{}) *SecondaryIterator {
	return db.Find(secondary.(chain.Uint256))
}

//Find the table row in a secondary 256-bit index table that matches the lowerbound condition for a given secondary key
func (db *IdxTable256) Lowerbound(secondary chain.Uint256) (*SecondaryIterator, chain.Uint256) {
	it, secondary, primary := eosio.DBIdx256Lowerbound(db.code, db.scope, db.table, secondary)
	return &SecondaryIterator{it, primary, db.dbIndex}, secondary
}

//Find the table row in a secondary 256-bit index table that matches the upperbound condition for a given secondary key
func (db *IdxTable256) Upperbound(secondary chain.Uint256) (*SecondaryIterator, chain.Uint256) {
	it, secondary, primary := eosio.DBIdx256Upperbound(db.code, db.scope, db.table, secondary)
	return &SecondaryIterator{it, primary, db.dbIndex}, secondary
}

//Get an end iterator representing just-past-the-end of the last table row of a secondary 256-bit index table
func (db *IdxTable256) End() *SecondaryIterator {
	ret := eosio.DBIdx256End(db.code, db.scope, db.table)
	return &SecondaryIterator{ret, 0, db.dbIndex}
}
