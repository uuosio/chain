package database

import "github.com/uuosio/chain"

type SingletonInterface interface {
	Set(data TableValue, payer chain.Name)
	Get() interface{}
	Remove()
}

type SingletonTable struct {
	TableI64
	unpacker Unpacker
}

func NewSingletonTable(code, scope, table chain.Name, unpacker ...Unpacker) *SingletonTable {
	chain.Check(len(unpacker) == 1, "unpacker cannot be nil when save state is enabled")

	if len(unpacker) >= 1 {
		primaryTableUnpacker := func([]byte) TableValue {
			return unpacker[0]([]byte{})
		}
		return &SingletonTable{TableI64: *NewTableI64(code, scope, table, primaryTableUnpacker), unpacker: unpacker[0]}
	} else {
		return &SingletonTable{TableI64: *NewTableI64(code, scope, table, nil), unpacker: nil}
	}
}

func (t *SingletonTable) Set(data TableValue, payer chain.Name) {
	t.TableI64.Set(t.TableI64.GetTableName(), chain.EncoderPack(data), payer)
}

func (t *SingletonTable) Get() interface{} {
	it, data := t.TableI64.GetByKey(t.TableI64.GetTableName())
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

func (t *SingletonTable) Remove() {
	it := t.TableI64.Find(t.TableI64.GetTableName())
	if it.IsOk() {
		t.TableI64.Remove(it)
	}
}
