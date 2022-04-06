package database

import (
	"errors"

	"github.com/uuosio/chain"
)

type MultiIndex struct {
	code             chain.Name
	scope            chain.Name
	table            chain.Name
	IdxDBNameToIndex func(string) int
	IndexTypes       []int
	Unpack           Unpacker

	DB     *DBI64
	IDXDBs []SecondaryDB
}

type MultiIndexInterface interface {
	Store(v MultiIndexValue, payer chain.Name) *Iterator
	Set(primary uint64, v MultiIndexValue, payer chain.Name)
	Get(id uint64) (*Iterator, MultiIndexValue)
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
	IdxFindByName(idxDBName string, secondary interface{}) *SecondaryIterator
	//	UpdateSecondaryValue(idxDB SecondaryDB, primary uint64, secondary interface{}, payer chain.Name)
	IdxUpdate(it *SecondaryIterator, secondary interface{}, payer chain.Name)
	GetIdxDBByIndex(index int) SecondaryDB
	GetIdxDB(idxDBName string) SecondaryDB
}

type MultiIndexValue interface {
	DBValue
	GetSecondaryValue(index int) interface{}
	SetSecondaryValue(index int, v interface{})
}

type Unpacker func([]byte) MultiIndexValue

var (
	ErrNotMultiIndexValue = errors.New("not a MultiIndexValue type")
)

func NewMultiIndex(code chain.Name, scope chain.Name, table chain.Name, idxDBNameToIndex func(string) int, indexTypes []int, saveState bool, unpacker ...Unpacker) *MultiIndex {
	chain.Check(code != chain.Name{0}, "bad code name")

	if table.N&uint64(0x0f) != 0 {
		// Limit table names to 12 characters so that the last character (4 bits) can be used to distinguish between the secondary indices.
		panic("NewMultiIndex:Invalid multi-index table name ")
	}

	mi := &MultiIndex{}
	mi.code = code
	mi.scope = scope
	mi.table = table
	mi.DB = NewDBI64(code, scope, table, func(data []byte) DBValue {
		return mi.Unpack(data)
	})
	mi.IdxDBNameToIndex = idxDBNameToIndex
	mi.IndexTypes = indexTypes
	mi.IDXDBs = make([]SecondaryDB, len(indexTypes))
	if len(unpacker) == 0 {
		mi.Unpack = nil
	} else {
		mi.Unpack = unpacker[0]
	}
	for i, v := range indexTypes {
		idxTable := (table.N & uint64(0xfffffffffffffff0)) | uint64(i)
		switch v {
		case IDX64:
			mi.IDXDBs[i] = NewIdxDB64(i, code.N, scope.N, idxTable)
			break
		case IDX128:
			mi.IDXDBs[i] = NewIdxDB128(i, code.N, scope.N, idxTable)
			break
		case IDX256:
			mi.IDXDBs[i] = NewIdxDB256(i, code.N, scope.N, idxTable)
			break
		case IDXFloat64:
			mi.IDXDBs[i] = NewIdxDBFloat64(i, code.N, scope.N, idxTable)
			break
		case IDXFloat128:
			mi.IDXDBs[i] = NewIdxDBFloat128(i, code.N, scope.N, idxTable)
		default:
			panic("invalid index")
		}
		// mi.secondaryDBs[v] = NewSecondaryDB(code.N, scope.N, table.N, v)
	}
	return mi
}

func (mi *MultiIndex) SetTable(code chain.Name, scope chain.Name, table chain.Name) {
	mi.code = code
	mi.scope = scope
	mi.table = table
}

func (mi *MultiIndex) Store(v MultiIndexValue, payer chain.Name) *Iterator {
	it := mi.DB.Store(v.GetPrimary(), v.Pack(), payer)
	primary := v.GetPrimary()
	for i, db := range mi.IDXDBs {
		db.StoreEx(primary, v.GetSecondaryValue(i), payer.N)
	}
	return it
}

func (mi *MultiIndex) Set(primary uint64, v MultiIndexValue, payer chain.Name) {
	chain.Check(primary == v.GetPrimary(), "mi.Store: Invalid primary key")
	it := mi.Find(primary)
	if !it.IsOk() {
		mi.DB.Store(primary, v.Pack(), payer)
		for i, db := range mi.IDXDBs {
			db.StoreEx(primary, v.GetSecondaryValue(i), payer.N)
		}
	} else {
		mi.Update(it, v, payer)
	}
}

func (mi *MultiIndex) Find(primary uint64) *Iterator {
	return mi.DB.Find(primary)
}

//Get value by primary index
//Returns iterator and value
//iterator can be used by MultiIndex.Update method to update value
func (mi *MultiIndex) Get(id uint64) (*Iterator, MultiIndexValue) {
	it, data := mi.DB.Get(id)
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
	v := mi.DB.GetByIterator(it)
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
			panic("not a Uint256 value")
		}

		_b, ok := b.(chain.Uint256)
		if !ok {
			panic("not a Uint256 value")
		}
		return _a == _b
	case IDXFloat64:
		return a.(float64) == b.(float64)
	case IDXFloat128:
		_a, ok := a.(chain.Float128)
		if !ok {
			panic("not a Uint256 value")
		}

		_b, ok := b.(chain.Float128)
		if !ok {
			panic("not a Uint256 value")
		}
		return _a == _b
	default:
		panic("invalid index type")
	}
	return false
}

func (mi *MultiIndex) Update(it *Iterator, v MultiIndexValue, payer chain.Name) {
	oldValue := mi.GetByIterator(it)

	primary := v.GetPrimary()
	chain.Check(oldValue.GetPrimary() == primary, "mi.Update: Can not change primary key duration update")

	chain.Check(mi.code == chain.CurrentReceiver(), "mi.Update: Can not update other contract")

	mi.DB.Update(it, v.Pack(), payer)
	for i, db := range mi.IDXDBs {
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
	mi.DB.Remove(it)
	for _, db := range mi.IDXDBs {
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
	mi.DB.Remove(it)
	for _, db := range mi.IDXDBs {
		it, _ := db.FindByPrimary(primary)
		db.Remove(it)
	}
}

func (mi *MultiIndex) GetIndexType(index int) int {
	return mi.IndexTypes[index]
}

//Find the table row following the referenced table row in a primary 64-bit integer index table
func (mi *MultiIndex) Next(it *Iterator) (next_iterator *Iterator, primary uint64) {
	return mi.DB.Next(it)
}

//Find the table row preceding the referenced table row in a primary 64-bit integer index table
func (mi *MultiIndex) Previous(it *Iterator) (previous_iterator *Iterator, primary uint64) {
	return mi.DB.Previous(it)
}

func (mi *MultiIndex) Lowerbound(id uint64) *Iterator {
	return mi.DB.Lowerbound(id)
}

func (mi *MultiIndex) Upperbound(id uint64) *Iterator {
	return mi.DB.Upperbound(id)
}

func (mi *MultiIndex) End() *Iterator {
	return mi.DB.End()
}

func (mi *MultiIndex) IdxFind(index int, secondary interface{}) *SecondaryIterator {
	return mi.IDXDBs[index].FindEx(secondary)
}

func (mi *MultiIndex) UpdateSecondaryValue(idxDB SecondaryDB, primary uint64, secondary interface{}, payer chain.Name) {
	itPrimary := mi.DB.Find(primary)
	chain.Check(itPrimary.IsOk(), "primary not found!")
	v := mi.DB.GetByIterator(itPrimary)
	_v := mi.Unpack(v)

	newSecondary := _v.GetSecondaryValue(idxDB.GetIndex())
	if IsEqual(mi.IndexTypes[idxDB.GetIndex()], newSecondary, secondary) {
		return
	}
	_v.SetSecondaryValue(idxDB.GetIndex(), secondary)
	mi.DB.Update(itPrimary, _v.Pack(), payer)
}

func (mi *MultiIndex) IdxFindByName(idxDBName string, secondary interface{}) *SecondaryIterator {
	chain.Check(mi.IdxDBNameToIndex != nil, "idxDBNameToIndex is nil")
	index := mi.IdxDBNameToIndex(idxDBName)
	return mi.IDXDBs[index].FindEx(secondary)
}

func (mi *MultiIndex) IdxUpdate(it *SecondaryIterator, secondary interface{}, payer chain.Name) {
	idxDB := mi.IDXDBs[it.dbIndex]

	itPrimary := mi.DB.Find(it.Primary)
	chain.Check(itPrimary.IsOk(), "primary not found!")
	v := mi.DB.GetByIterator(itPrimary)
	_v := mi.Unpack(v)
	newSecondary := _v.GetSecondaryValue(idxDB.GetIndex())
	if IsEqual(mi.IndexTypes[idxDB.GetIndex()], newSecondary, secondary) {
		return
	}
	_v.SetSecondaryValue(idxDB.GetIndex(), secondary)
	mi.DB.Update(itPrimary, _v.Pack(), payer)
	idxDB.UpdateEx(it, secondary, payer.N)
}

func (mi *MultiIndex) IdxGet(itSecondary *SecondaryIterator) interface{} {
	idxDB := mi.IDXDBs[itSecondary.dbIndex]
	{
		itSecondary2, secondary := idxDB.FindByPrimary(itSecondary.Primary)
		chain.Check(itSecondary2.IsOk(), "secondary not found!")
		return secondary
	}
	{
		it, v := mi.Get(itSecondary.Primary)
		chain.Check(it.IsOk(), "mi.IdxGet: primary not found!")

		itSecondary2, secondary := idxDB.FindByPrimary(itSecondary.Primary)
		chain.Check(itSecondary2.IsOk(), "mi.IdxGet: secondary not found!")
		newSecondary := v.GetSecondaryValue(idxDB.GetIndex())
		equal := IsEqual(mi.IndexTypes[idxDB.GetIndex()], newSecondary, secondary)
		chain.Check(equal, "mi.IdxGet: secondary not the same!")
		return secondary
	}
}

func (mi *MultiIndex) GetIdxDBByIndex(index int) SecondaryDB {
	return mi.IDXDBs[index]
}

func (mi *MultiIndex) GetIdxDB(idxDBName string) SecondaryDB {
	chain.Check(mi.IdxDBNameToIndex != nil, "idxDBNameToIndex is nil")
	index := mi.IdxDBNameToIndex(idxDBName)
	return mi.IDXDBs[index]
}

// func (mi *MultiIndex) GetIdxDBByA1() IdxDB64I {
// 	return IdxDB64I{mi.IDXDBs[0], mi.IDXDBs[0].(*IdxDB64)}
// }
