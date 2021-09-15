package database

import "github.com/uuosio/chain"

type SingletonInterface interface {
	Set(data DBValue, payer chain.Name)
	Get() (interface{}, error)
}

type SingletonDB struct {
	DB       *DBI64
	unpacker Unpacker
}

func NewSingletonDB(code, scope, table chain.Name, unpacker ...Unpacker) *SingletonDB {
	if len(unpacker) > 0 {
		return &SingletonDB{DB: NewDBI64(code, scope, table), unpacker: unpacker[0]}
	} else {
		return &SingletonDB{DB: NewDBI64(code, scope, table), unpacker: nil}
	}
}

func (t *SingletonDB) Set(data DBValue, payer chain.Name) {
	t.DB.Set(t.DB.table, data.Pack(), payer)
}

func (t *SingletonDB) Get() interface{} {
	it, data := t.DB.Get(t.DB.table)
	if !it.IsOk() {
		return nil
	}

	if t.unpacker != nil {
		value := t.unpacker(data)
		return value
	} else {
		return data
	}
}
