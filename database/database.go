package database

import (
	"errors"

	"github.com/uuosio/chain"
)

const (
	IDX64       = iota
	IDX128      = iota
	IDX256      = iota
	IDXFloat64  = iota
	IDXFloat128 = iota
)

var ErrInvalidIterator = errors.New("Invalid iterator")
var UNKNOWN_PRIMARY_KEY = uint64(0xffffffffffffffff)

type Iterator struct {
	I                 int32
	primary           uint64
	isValidPrimaryKey bool
	db                *TableI64
}

func (it *Iterator) GetPrimary() uint64 {
	chain.Check(it.IsOk(), "Iterator.GetPrimary: Invalid iterator")
	if !it.isValidPrimaryKey {
		it.isValidPrimaryKey = true
		data := it.db.GetByIterator(it)
		value := it.db.unpacker(data)
		it.primary = value.GetPrimary()
	}

	return it.primary
}

func (it *Iterator) IsOk() bool {
	return it.I >= 0
}

func (it *Iterator) IsEnd() bool {
	return it.I < -1
}

func (it *Iterator) IsValid() bool {
	return it.I != -1
}

type SecondaryIterator struct {
	I       int32
	Primary uint64
	dbIndex int
}

func (it *SecondaryIterator) IsOk() bool {
	return it.I >= 0
}

func (it *SecondaryIterator) IsEnd() bool {
	return it.I < -1
}

func (it *SecondaryIterator) IsValid() bool {
	return it.I != -1
}

type SecondaryTable interface {
	GetIndex() int
	GetTable() (uint64, uint64, uint64)
	StoreEx(id uint64, secondary interface{}, payer uint64) *SecondaryIterator
	UpdateEx(it *SecondaryIterator, secondary interface{}, payer uint64)
	Remove(it *SecondaryIterator)
	FindByPrimary(primary uint64) (*SecondaryIterator, interface{})
	FindEx(secondary interface{}) *SecondaryIterator
	End() *SecondaryIterator
}

type TableValue interface {
	GetPrimary() uint64
	Pack() []byte
	Unpack(data []byte) int
}
