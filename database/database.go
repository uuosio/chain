package database

import "errors"

const (
	IDX64       = iota
	IDX128      = iota
	IDX256      = iota
	IDXFloat64  = iota
	IDXFloat128 = iota
)

var ErrInvalidIterator = errors.New("Invalid iterator")

type Iterator struct {
	I int32
}

func (it Iterator) IsOk() bool {
	return it.I >= 0
}

func (it Iterator) IsEnd() bool {
	return it.I < -1
}

func (it Iterator) IsValid() bool {
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

type SecondaryDB interface {
	GetIndex() int
	GetTable() (uint64, uint64, uint64)
	Store(id uint64, secondary interface{}, payer uint64) SecondaryIterator
	Update(it SecondaryIterator, secondary interface{}, payer uint64)
	Remove(it SecondaryIterator)
	Next(it SecondaryIterator) SecondaryIterator
	Previous(it SecondaryIterator) SecondaryIterator
	FindByPrimary(primary uint64) (SecondaryIterator, interface{})
	Find(secondary interface{}) SecondaryIterator
	Lowerbound(secondary interface{}) (SecondaryIterator, interface{})
	Upperbound(secondary interface{}) (SecondaryIterator, interface{})
	End() SecondaryIterator
}

type DBValue interface {
	GetPrimary() uint64
	Pack() []byte
	Unpack(data []byte) int
}
