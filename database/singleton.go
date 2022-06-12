package database

import "github.com/uuosio/chain"

type SingletonInterface interface {
	Set(data DBValue, payer chain.Name)
	Get() interface{}
	Remove()
}

type SingletonDB struct {
	DBI64
	unpacker Unpacker
}

func NewSingletonDB(code, scope, table chain.Name, unpacker ...Unpacker) *SingletonDB {
	chain.Check(len(unpacker) == 1, "unpacker cannot be nil when save state is enabled")

	if len(unpacker) >= 1 {
		primaryDBUnpacker := func([]byte) DBValue {
			return unpacker[0]([]byte{})
		}
		return &SingletonDB{DBI64: *NewDBI64(code, scope, table, primaryDBUnpacker), unpacker: unpacker[0]}
	} else {
		return &SingletonDB{DBI64: *NewDBI64(code, scope, table, nil), unpacker: nil}
	}
}

func (t *SingletonDB) Set(data DBValue, payer chain.Name) {
	t.DBI64.Set(t.DBI64.GetTableName(), data.Pack(), payer)
}

func (t *SingletonDB) Get() interface{} {
	it, data := t.DBI64.GetByKey(t.DBI64.GetTableName())
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

func (t *SingletonDB) Remove() {
	it := t.DBI64.Find(t.DBI64.GetTableName())
	if it.IsOk() {
		t.DBI64.Remove(it)
	}
}
