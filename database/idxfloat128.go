package database

import (
	"github.com/uuosio/chain"
	"github.com/uuosio/chain/eosio"
)

type IdxTableFloat128 struct {
	IdxTable
}

func NewIdxTableFloat128(index int, code uint64, scope uint64, table uint64) *IdxTableFloat128 {
	v := &IdxTableFloat128{IdxTable{index, code, scope, table}}
	return v
}

//Store an association of a quadruple-precision floating-point secondary key to a primary key in a secondary quadruple-precision floating-point index table
func (db *IdxTableFloat128) Store(id uint64, secondary chain.Float128, payer uint64) *SecondaryIterator {
	i := eosio.DBIdxFloat128Store(db.scope, db.table, id, secondary, payer)
	return &SecondaryIterator{i, id, db.dbIndex}
}

func (db *IdxTableFloat128) StoreEx(id uint64, secondary interface{}, payer uint64) *SecondaryIterator {
	_secondary, ok := secondary.(chain.Float128)
	chain.Check(ok, "bad secondary type")
	return db.Store(id, _secondary, payer)
}

//Update an association for a quadruple-precision floating-point secondary key to a primary key in a secondary quadruple-precision floating-point index table
func (db *IdxTableFloat128) Update(it *SecondaryIterator, secondary chain.Float128, payer uint64) {
	eosio.DBIdxFloat128Update(it.I, secondary, payer)
}

func (db *IdxTableFloat128) UpdateEx(it *SecondaryIterator, secondary interface{}, payer uint64) {
	_secondary, ok := secondary.(chain.Float128)
	chain.Check(ok, "bad secondary type")
	db.Update(it, _secondary, payer)
}

//Remove a table row from a secondary quadruple-precision floating-point index table
func (db *IdxTableFloat128) Remove(it *SecondaryIterator) {
	eosio.DBIdxFloat128Remove(it.I)
}

//Find the table row following the referenced table row in a secondary quadruple-precision floating-point index table
func (db *IdxTableFloat128) Next(it *SecondaryIterator) *SecondaryIterator {
	i, primary := eosio.DBIdxFloat128Next(it.I)
	return &SecondaryIterator{i, primary, db.dbIndex}
}

//Find the table row preceding the referenced table row in a secondary quadruple-precision floating-point index table
func (db *IdxTableFloat128) Previous(it *SecondaryIterator) *SecondaryIterator {
	i, primary := eosio.DBIdxFloat128Previous(it.I)
	return &SecondaryIterator{i, primary, db.dbIndex}
}

//Find a table row in a secondary quadruple-precision floating-point index table by primary key
func (db *IdxTableFloat128) FindByPrimary(primary uint64) (*SecondaryIterator, interface{}) {
	var secondary chain.Float128
	i, secondary := eosio.DBIdxFloat128FindByPrimary(db.code, db.scope, db.table, primary)
	return &SecondaryIterator{i, primary, db.dbIndex}, secondary
}

//Find a table row in a secondary quadruple-precision floating-point index table by secondary key
func (db *IdxTableFloat128) Find(secondary chain.Float128) *SecondaryIterator {
	// var primary uint64 = 0
	// ret := C.db_idx_long_double_find_secondary(C.uint64_t(db.code), C.uint64_t(db.scope), C.uint64_t(db.table), (*C.float128_t)(unsafe.Pointer(&_secondary)), &primary)
	// return &SecondaryIterator{ret, primary, db.dbIndex}
	it, value := db.Lowerbound(secondary)
	if it.IsOk() {
		if value == secondary {
			return it
		}
	}
	return &SecondaryIterator{-1, 0, db.dbIndex}
}

func (db *IdxTableFloat128) FindEx(secondary interface{}) *SecondaryIterator {
	_secondary, ok := secondary.(chain.Float128)
	chain.Check(ok, "not a float128_t value")
	return db.Find(_secondary)
}

//Find the table row in a secondary quadruple-precision floating-point index table that matches the lowerbound condition for a given secondary key
func (db *IdxTableFloat128) Lowerbound(secondary chain.Float128) (*SecondaryIterator, chain.Float128) {
	i, secondary, primary := eosio.DBIdxFloat128Lowerbound(db.code, db.scope, db.table, secondary)
	return &SecondaryIterator{i, primary, db.dbIndex}, secondary
}

//Find the table row in a secondary quadruple-precision floating-point index table that matches the upperbound condition for a given secondary key
func (db *IdxTableFloat128) Upperbound(secondary chain.Float128) (*SecondaryIterator, chain.Float128) {
	i, secondary, primary := eosio.DBIdxFloat128Upperbound(db.code, db.scope, db.table, secondary)
	return &SecondaryIterator{i, primary, db.dbIndex}, secondary
}

//Get an end iterator representing just-past-the-end of the last table row of a secondary quadruple-precision floating-point index table
func (db *IdxTableFloat128) End() *SecondaryIterator {
	i := eosio.DBIdxFloat128End(db.code, db.scope, db.table)
	return &SecondaryIterator{i, 0, db.dbIndex}
}
