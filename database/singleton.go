package database

import "github.com/uuosio/chain"

type SingletonInterface interface {
	Set(data DBValue, payer chain.Name)
	Get() (interface{}, error)
}

type SingletonDB struct {
	DB *DBI64
}

func NewSingletonDB(code, scope, table chain.Name, unpacker ...Unpacker) *SingletonDB {
	return &SingletonDB{DB: NewDBI64(code, scope, table, unpacker...)}
}

func (t *SingletonDB) Set(data DBValue, payer chain.Name) {
	t.DB.Set(t.DB.table, data.Pack(), payer)
}

func (t *SingletonDB) Get() (interface{}, error) {
	return t.DB.Get(t.DB.table)
}
