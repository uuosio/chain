package database

import (
	"github.com/uuosio/chain"
	"github.com/uuosio/chain/eosio"
)

type IdxTable128 struct {
	IdxTable
}

func NewIdxTable128(index int, code uint64, scope uint64, table uint64) *IdxTable128 {
	v := &IdxTable128{IdxTable{index, code, scope, table}}
	return v
}

//Store an association of a 128-bit integer secondary key to a primary key in a secondary 128-bit integer index table
func (db *IdxTable128) Store(id uint64, secondary chain.Uint128, payer uint64) *SecondaryIterator {
	ret := eosio.DBIdx128Store(db.scope, db.table, id, secondary, payer)
	return &SecondaryIterator{int32(ret), id, db.dbIndex}
}

func (db *IdxTable128) StoreEx(id uint64, secondary interface{}, payer uint64) *SecondaryIterator {
	_secondary, ok := secondary.(chain.Uint128)
	chain.Check(ok, "bad secondary type")
	return db.Store(id, _secondary, payer)
}

//Update an association for a 128-bit integer secondary key to a primary key in a secondary 128-bit integer index table
func (db *IdxTable128) Update(it *SecondaryIterator, secondary chain.Uint128, payer uint64) {
	eosio.DBIdx128Update(it.I, secondary, payer)
}

func (db *IdxTable128) UpdateEx(it *SecondaryIterator, secondary interface{}, payer uint64) {
	_secondary, ok := secondary.(chain.Uint128)
	chain.Check(ok, "bad secondary type")
	db.Update(it, _secondary, payer)
}

//Remove a table row from a secondary 128-bit integer index table
func (db *IdxTable128) Remove(it *SecondaryIterator) {
	// var secondary chain.Uint128
	// ret := C.db_idx128_find_primary(C.uint64_t(db.code), C.uint64_t(db.scope), C.uint64_t(db.table), (*C.uint128)(unsafe.Pointer(&secondary)), C.uint64_t(it.Primary))
	// chain.Check(int32(ret) == it.I, "invalid iterator")
	eosio.DBIdx128Remove(it.I)
}

//Find the table row following the referenced table row in a secondary 128-bit integer index table
func (db *IdxTable128) Next(it *SecondaryIterator) *SecondaryIterator {
	ret, primary := eosio.DBIdx128Next(it.I)
	return &SecondaryIterator{int32(ret), primary, db.dbIndex}
}

//Find the table row preceding the referenced table row in a secondary 128-bit integer index table
func (db *IdxTable128) Previous(it *SecondaryIterator) *SecondaryIterator {
	ret, primary := eosio.DBIdx128Previous(it.I)
	return &SecondaryIterator{int32(ret), primary, db.dbIndex}
}

//Find a table row in a secondary 128-bit integer index table by primary key
func (db *IdxTable128) FindByPrimary(primary uint64) (*SecondaryIterator, interface{}) {
	ret, secondary := eosio.DBIdx128FindByPrimary(db.code, db.scope, db.table, primary)
	return &SecondaryIterator{int32(ret), primary, db.dbIndex}, chain.Uint128(secondary)
}

//Find a table row in a secondary 128-bit integer index table by secondary key
func (db *IdxTable128) Find(secondary chain.Uint128) *SecondaryIterator {
	it, value := db.Lowerbound(secondary)
	if it.IsOk() {
		if value == secondary {
			return it
		}
	}
	return &SecondaryIterator{-1, 0, db.dbIndex}
}

func (db *IdxTable128) FindEx(secondary interface{}) *SecondaryIterator {
	return db.Find(secondary.(chain.Uint128))
}

//Find the table row in a secondary 128-bit integer index table that matches the lowerbound condition for a given secondary key
func (db *IdxTable128) Lowerbound(secondary chain.Uint128) (*SecondaryIterator, chain.Uint128) {
	it, secondary, primary := eosio.DBIdx128Lowerbound(db.code, db.scope, db.table, secondary)
	return &SecondaryIterator{it, primary, db.dbIndex}, chain.Uint128(secondary)
}

//Find the table row in a secondary 128-bit integer index table that matches the upperbound condition for a given secondary key
func (db *IdxTable128) Upperbound(secondary chain.Uint128) (*SecondaryIterator, chain.Uint128) {
	it, secondary, primary := eosio.DBIdx128Upperbound(db.code, db.scope, db.table, secondary)
	return &SecondaryIterator{it, primary, db.dbIndex}, chain.Uint128(secondary)
}

//Get an end iterator representing just-past-the-end of the last table row of a secondary 128-bit integer index table
func (db *IdxTable128) End() *SecondaryIterator {
	it := eosio.DBIdx128End(db.code, db.scope, db.table)
	return &SecondaryIterator{it, 0, db.dbIndex}
}
