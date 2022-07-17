package database

import (
	"github.com/uuosio/chain"
	"github.com/uuosio/chain/eosio"
)

const (
	MsgErrorUnpacker = "unpacker cannot be nil when save state is enabled"
)

const (
	OperationStore  = 1
	OperationUpdate = 2
	OperationRemove = 3
)

type TableI64Unpacker func([]byte) TableValue

type TableI64 struct {
	code     uint64
	scope    uint64
	table    uint64
	unpacker TableI64Unpacker
}

func NewTableI64(code chain.Name, scope chain.Name, table chain.Name, unpacker TableI64Unpacker) *TableI64 {
	chain.Assert(chain.IsAccount(code), "code account does not exists!")
	return &TableI64{code.N, scope.N, table.N, unpacker}
}

func (db *TableI64) Init(code chain.Name, scope chain.Name, table chain.Name) {
	db.code = code.N
	db.scope = scope.N
	db.table = table.N
}

func (db *TableI64) GetTable() (uint64, uint64, uint64) {
	return db.code, db.scope, db.table
}

func (db *TableI64) GetTableName() uint64 {
	return db.table
}

func (db *TableI64) Set(id uint64, data []byte, payer chain.Name) *Iterator {
	it := db.Find(id)
	if !it.IsOk() {
		return db.storeI64(payer.N, id, data)
	} else {
		db.updateI64(it, payer.N, data)
		return it
	}
}

//Get a record in a primary 64-bit integer index table
func (db *TableI64) Get(id uint64) (*Iterator, []byte) {
	return db.GetByKey(id)
}

func (db *TableI64) GetByKey(id uint64) (*Iterator, []byte) {
	it := db.Find(id)
	if !it.IsOk() {
		return it, nil
	}

	raw := db.getI64(it)
	return it, raw
}

func (db *TableI64) GetByIterator(it *Iterator) []byte {
	chain.Check(it.IsOk(), "invalid iterator")
	raw := db.getI64(it)
	return raw
}

//Store a record in a primary 64-bit integer index table
func (db *TableI64) Store(primary uint64, data []byte, payer chain.Name) *Iterator {
	return db.storeI64(payer.N, primary, data)
}

//Update a record in a primary 64-bit integer index table
func (db *TableI64) Update(it *Iterator, data []byte, payer chain.Name) {
	db.updateI64(it, payer.N, data)
}

//Remove a record from a primary 64-bit integer index table
func (db *TableI64) Remove(it *Iterator) {
	eosio.DBRemoveI64(it.I)
}

//Find the table row following the referenced table row in a primary 64-bit integer index table
func (db *TableI64) Next(it *Iterator) (next_iterator *Iterator, primary uint64) {
	ret, primary := eosio.DBNextI64(it.I)
	return &Iterator{ret, primary, true, db}, primary
}

//Find the table row preceding the referenced table row in a primary 64-bit integer index table
func (db *TableI64) Previous(it *Iterator) (previous_iterator *Iterator, primary uint64) {
	ret, primary := eosio.DBPreviousI64(it.I)
	return &Iterator{ret, primary, true, db}, primary
}

//Find a table row in a primary 64-bit integer index table by primary key
func (db *TableI64) Find(id uint64) *Iterator {
	ret := eosio.DBFindI64(db.code, db.scope, db.table, id)
	if ret >= 0 {
		return &Iterator{ret, id, true, db}
	}
	return &Iterator{ret, 0, false, db}
}

//Find the table row in a primary 64-bit integer index table that matches the lowerbound condition for a given primary key
func (db *TableI64) Lowerbound(id uint64) *Iterator {
	ret := eosio.DBLowerBoundI64(db.code, db.scope, db.table, id)
	return &Iterator{ret, 0, false, db}
}

//Find the table row in a primary 64-bit integer index table that matches the upperbound condition for a given primary key
func (db *TableI64) Upperbound(id uint64) *Iterator {
	ret := eosio.DBUpperBoundI64(db.code, db.scope, db.table, id)
	return &Iterator{ret, 0, false, db}
}

//Get an iterator representing just-past-the-end of the last table row of a primary 64-bit integer index table
func (db *TableI64) End() *Iterator {
	ret := eosio.DBEndI64(db.code, db.scope, db.table)
	return &Iterator{int32(ret), 0, false, db}
}

func (db *TableI64) storeI64(payer uint64, id uint64, data []byte) *Iterator {
	ret := eosio.DBStoreI64(db.scope, db.table, payer, id, data)
	return &Iterator{int32(ret), id, true, db}
}

func (db *TableI64) updateI64(iterator *Iterator, payer uint64, data []byte) {
	eosio.DBUpdateI64(iterator.I, payer, data)
}

func (db *TableI64) getI64(it *Iterator) []byte {
	return eosio.DBGetI64(it.I)
}
