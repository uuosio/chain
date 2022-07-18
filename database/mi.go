package database

import (
	"errors"

	"github.com/uuosio/chain"
)

type MultiIndex struct {
	code                chain.Name
	scope               chain.Name
	table               chain.Name
	IdxTableNameToIndex func(string) int
	IndexTypes          []int
	Unpack              Unpacker

	Table     *TableI64
	IDXTables []SecondaryTable
}

type MultiIndexInterface interface {
	Store(v MultiIndexValue, payer chain.Name) *Iterator
	Set(primary uint64, v MultiIndexValue, payer chain.Name)
	GetByKey(id uint64) (*Iterator, MultiIndexValue)
	GetByIterator(it *Iterator) MultiIndexValue
	Update(it *Iterator, v MultiIndexValue, payer chain.Name)

	Find(primary uint64) *Iterator

	Remove(it *Iterator)
	Next(it *Iterator) (next_iterator *Iterator, primary uint64)
	Previous(it *Iterator) (previous_iterator *Iterator, primary uint64)

	Lowerbound(id uint64) *Iterator
	Upperbound(id uint64) *Iterator
	End() *Iterator

	IdxFind(index int, secondary interface{}) *SecondaryIterator
	IdxGet(itSecondary *SecondaryIterator) interface{}
	IdxFindByName(idxTableName string, secondary interface{}) *SecondaryIterator
	//	UpdateSecondaryValue(idxTable SecondaryTable, primary uint64, secondary interface{}, payer chain.Name)
	IdxUpdate(it *SecondaryIterator, secondary interface{}, payer chain.Name)
	GetIdxTableByIndex(index int) SecondaryTable
	GetIdxTable(idxTableName string) SecondaryTable
}

type MultiIndexValue interface {
	TableValue
	GetSecondaryValue(index int) interface{}
	SetSecondaryValue(index int, v interface{})
}

type Unpacker func([]byte) MultiIndexValue

var (
	ErrNotMultiIndexValue = errors.New("not a MultiIndexValue type")
)

func NewMultiIndex(code chain.Name, scope chain.Name, table chain.Name, idxTableNameToIndex func(string) int, indexTypes []int, saveState bool, unpacker ...Unpacker) *MultiIndex {
	chain.Check(code != chain.Name{0}, "bad code name")

	if table.N&uint64(0x0f) != 0 {
		// Limit table names to 12 characters so that the last character (4 bits) can be used to distinguish between the secondary indices.
		panic("NewMultiIndex:Invalid multi-index table name ")
	}

	mi := &MultiIndex{}
	mi.code = code
	mi.scope = scope
	mi.table = table
	mi.Table = NewTableI64(code, scope, table, func(data []byte) TableValue {
		return mi.Unpack(data)
	})
	mi.IdxTableNameToIndex = idxTableNameToIndex
	mi.IndexTypes = indexTypes
	mi.IDXTables = make([]SecondaryTable, len(indexTypes))
	if len(unpacker) == 0 {
		mi.Unpack = nil
	} else {
		mi.Unpack = unpacker[0]
	}
	for i, v := range indexTypes {
		idxTable := (table.N & uint64(0xfffffffffffffff0)) | uint64(i)
		switch v {
		case IDX64:
			mi.IDXTables[i] = NewIdxTable64(i, code.N, scope.N, idxTable)
			break
		case IDX128:
			mi.IDXTables[i] = NewIdxTable128(i, code.N, scope.N, idxTable)
			break
		case IDX256:
			mi.IDXTables[i] = NewIdxTable256(i, code.N, scope.N, idxTable)
			break
		case IDXFloat64:
			mi.IDXTables[i] = NewIdxTableFloat64(i, code.N, scope.N, idxTable)
			break
		case IDXFloat128:
			mi.IDXTables[i] = NewIdxTableFloat128(i, code.N, scope.N, idxTable)
		default:
			panic("invalid index")
		}
		// mi.secondaryTables[v] = NewSecondaryTable(code.N, scope.N, table.N, v)
	}
	return mi
}

func (mi *MultiIndex) SetTable(code chain.Name, scope chain.Name, table chain.Name) {
	mi.code = code
	mi.scope = scope
	mi.table = table
}

func (mi *MultiIndex) Store(v MultiIndexValue, payer chain.Name) *Iterator {
	it := mi.Table.Store(v.GetPrimary(), v.Pack(), payer)
	primary := v.GetPrimary()
	for i, db := range mi.IDXTables {
		db.StoreEx(primary, v.GetSecondaryValue(i), payer.N)
	}
	return it
}

func (mi *MultiIndex) Set(primary uint64, v MultiIndexValue, payer chain.Name) {
	chain.Check(primary == v.GetPrimary(), "mi.Store: Invalid primary key")
	it := mi.Find(primary)
	if !it.IsOk() {
		mi.Table.Store(primary, v.Pack(), payer)
		for i, db := range mi.IDXTables {
			db.StoreEx(primary, v.GetSecondaryValue(i), payer.N)
		}
	} else {
		mi.Update(it, v, payer)
	}
}

func (mi *MultiIndex) Find(primary uint64) *Iterator {
	return mi.Table.Find(primary)
}

//Get value by primary index
//Returns iterator and value
//iterator can be used by MultiIndex.Update method to update value
func (mi *MultiIndex) Get(id uint64) (*Iterator, MultiIndexValue) {
	return mi.GetByKey(id)
}

func (mi *MultiIndex) GetByKey(id uint64) (*Iterator, MultiIndexValue) {
	it, data := mi.Table.GetByKey(id)
	if !it.IsOk() {
		return it, nil
	}
	value := mi.Unpack(data)
	_data, ok := value.(MultiIndexValue)
	chain.Check(ok, "mi.Get: Not a MultiIndexValue type")
	chain.Check(id == _data.GetPrimary(), "mi.Get: bad primary value")
	return it, _data
}

//Get value by primary Iterator
func (mi *MultiIndex) GetByIterator(it *Iterator) MultiIndexValue {
	v := mi.Table.GetByIterator(it)
	vv := mi.Unpack(v)
	return vv
}

func IsEqual(indexType int, a, b interface{}) bool {
	switch indexType {
	case IDX64:
		return a.(uint64) == b.(uint64)
	case IDX128:
		_a, ok := a.(chain.Uint128)
		if !ok {
			panic("not a Uint128 value")
		}

		_b, ok := b.(chain.Uint128)
		if !ok {
			panic("not a Uint128 value")
		}
		return _a == _b
	case IDX256:
		_a, ok := a.(chain.Uint256)
		if !ok {
			panic("not a Uint256 value 1")
		}

		_b, ok := b.(chain.Uint256)
		if !ok {
			panic("not a Uint256 value 2")
		}
		return _a == _b
	case IDXFloat64:
		return a.(float64) == b.(float64)
	case IDXFloat128:
		_a, ok := a.(chain.Float128)
		if !ok {
			panic("not a IDXFloat128 value 1")
		}

		_b, ok := b.(chain.Float128)
		if !ok {
			panic("not a IDXFloat128 value 2")
		}
		return _a == _b
	default:
		panic("invalid index type")
	}
	return false
}

func (mi *MultiIndex) Update(it *Iterator, v MultiIndexValue, payer chain.Name) {
	chain.Check(it.IsOk(), "mi.Update: Invalid iterator")

	primary := it.GetPrimary()
	chain.Check(primary == v.GetPrimary(), "mi.Update: Can not change primary key during update")

	chain.Check(mi.code == chain.CurrentReceiver(), "mi.Update: Can not update other contract")

	mi.Table.Update(it, v.Pack(), payer)
	for i, db := range mi.IDXTables {
		it, oldSecondary := db.FindByPrimary(primary)
		// logger.Println(primary, i, oldSecondary, ":")
		chain.Check(it.IsOk(), "secondary value does not exists!")
		secondary := v.GetSecondaryValue(i)
		if IsEqual(mi.IndexTypes[i], oldSecondary, secondary) {
		} else {
			db.UpdateEx(it, secondary, payer.N)
		}
	}
}

func (mi *MultiIndex) Remove(it *Iterator) {
	v := mi.GetByIterator(it)
	mi.Table.Remove(it)
	for _, db := range mi.IDXTables {
		it, secondary := db.FindByPrimary(v.GetPrimary())
		indexType := mi.GetIndexType(db.GetIndex())
		_secondary := v.GetSecondaryValue(db.GetIndex())
		chain.Check(IsEqual(indexType, secondary, _secondary), "mi.Remove: secondary value not the same")
		db.Remove(it)
	}
}

func (mi *MultiIndex) RemoveEx(primary uint64) {
	it := mi.Find(primary)
	chain.Assert(it.IsOk(), "primary value not found1")
	mi.Table.Remove(it)
	for _, db := range mi.IDXTables {
		it, _ := db.FindByPrimary(primary)
		db.Remove(it)
	}
}

func (mi *MultiIndex) GetIndexType(index int) int {
	return mi.IndexTypes[index]
}

//Find the table row following the referenced table row in a primary 64-bit integer index table
func (mi *MultiIndex) Next(it *Iterator) (next_iterator *Iterator, primary uint64) {
	return mi.Table.Next(it)
}

//Find the table row preceding the referenced table row in a primary 64-bit integer index table
func (mi *MultiIndex) Previous(it *Iterator) (previous_iterator *Iterator, primary uint64) {
	return mi.Table.Previous(it)
}

func (mi *MultiIndex) Lowerbound(id uint64) *Iterator {
	return mi.Table.Lowerbound(id)
}

func (mi *MultiIndex) Upperbound(id uint64) *Iterator {
	return mi.Table.Upperbound(id)
}

func (mi *MultiIndex) End() *Iterator {
	return mi.Table.End()
}

func (mi *MultiIndex) IdxFind(index int, secondary interface{}) *SecondaryIterator {
	return mi.IDXTables[index].FindEx(secondary)
}

func (mi *MultiIndex) UpdateSecondaryValue(idxTable SecondaryTable, primary uint64, secondary interface{}, payer chain.Name) {
	itPrimary := mi.Table.Find(primary)
	chain.Check(itPrimary.IsOk(), "primary not found!")
	v := mi.Table.GetByIterator(itPrimary)
	_v := mi.Unpack(v)

	newSecondary := _v.GetSecondaryValue(idxTable.GetIndex())
	if IsEqual(mi.IndexTypes[idxTable.GetIndex()], newSecondary, secondary) {
		return
	}
	_v.SetSecondaryValue(idxTable.GetIndex(), secondary)
	mi.Table.Update(itPrimary, _v.Pack(), payer)
}

func (mi *MultiIndex) IdxFindByName(idxTableName string, secondary interface{}) *SecondaryIterator {
	chain.Check(mi.IdxTableNameToIndex != nil, "idxTableNameToIndex is nil")
	index := mi.IdxTableNameToIndex(idxTableName)
	return mi.IDXTables[index].FindEx(secondary)
}

func (mi *MultiIndex) IdxUpdate(it *SecondaryIterator, secondary interface{}, payer chain.Name) {
	idxTable := mi.IDXTables[it.dbIndex]

	itPrimary := mi.Table.Find(it.Primary)
	chain.Check(itPrimary.IsOk(), "primary not found!")
	v := mi.Table.GetByIterator(itPrimary)
	_v := mi.Unpack(v)
	newSecondary := _v.GetSecondaryValue(idxTable.GetIndex())
	if IsEqual(mi.IndexTypes[idxTable.GetIndex()], newSecondary, secondary) {
		return
	}
	_v.SetSecondaryValue(idxTable.GetIndex(), secondary)
	mi.Table.Update(itPrimary, _v.Pack(), payer)
	idxTable.UpdateEx(it, secondary, payer.N)
}

func (mi *MultiIndex) IdxGet(itSecondary *SecondaryIterator) interface{} {
	idxTable := mi.IDXTables[itSecondary.dbIndex]
	{
		itSecondary2, secondary := idxTable.FindByPrimary(itSecondary.Primary)
		chain.Check(itSecondary2.IsOk(), "secondary not found!")
		return secondary
	}
	{
		it, v := mi.GetByKey(itSecondary.Primary)
		chain.Check(it.IsOk(), "mi.IdxGet: primary not found!")

		itSecondary2, secondary := idxTable.FindByPrimary(itSecondary.Primary)
		chain.Check(itSecondary2.IsOk(), "mi.IdxGet: secondary not found!")
		newSecondary := v.GetSecondaryValue(idxTable.GetIndex())
		equal := IsEqual(mi.IndexTypes[idxTable.GetIndex()], newSecondary, secondary)
		chain.Check(equal, "mi.IdxGet: secondary not the same!")
		return secondary
	}
}

func (mi *MultiIndex) GetIdxTableByIndex(index int) SecondaryTable {
	return mi.IDXTables[index]
}

func (mi *MultiIndex) GetIdxTable(idxTableName string) SecondaryTable {
	chain.Check(mi.IdxTableNameToIndex != nil, "idxTableNameToIndex is nil")
	index := mi.IdxTableNameToIndex(idxTableName)
	return mi.IDXTables[index]
}

// func (mi *MultiIndex) GetIdxTableByA1() IdxTable64I {
// 	return IdxTable64I{mi.IDXTables[0], mi.IDXTables[0].(*IdxTable64)}
// }
