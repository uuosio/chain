package database

import (
	"errors"

	"github.com/uuosio/chain"
)

type MultiIndex struct {
	code             chain.Name
	scope            chain.Name
	table            chain.Name
	idxDBNameToIndex func(string) int
	IndexTypes       []int
	unpacker         Unpacker

	DB     *DBI64
	IDXDBs []SecondaryDB
}

type MultiIndexInterface interface {
	Store(v MultiIndexValue, payer chain.Name)
	Set(primary uint64, v MultiIndexValue, payer chain.Name)
	Get(id uint64) (Iterator, MultiIndexValue)
	GetByIterator(it Iterator) (MultiIndexValue, error)
	Update(it Iterator, v MultiIndexValue, payer chain.Name)

	Find(primary uint64) Iterator

	Remove(it Iterator)
	Next(it Iterator) (next_iterator Iterator, primary uint64)
	Previous(it Iterator) (previous_iterator Iterator, primary uint64)

	Lowerbound(id uint64) Iterator
	Upperbound(id uint64) Iterator
	End() Iterator

	IdxFind(index int, secondary interface{}) SecondaryIterator
	IdxGet(itSecondary SecondaryIterator) interface{}
	IdxFindByName(idxDBName string, secondary interface{}) SecondaryIterator
	//	UpdateSecondaryValue(idxDB SecondaryDB, primary uint64, secondary interface{}, payer chain.Name)
	IdxUpdate(it SecondaryIterator, secondary interface{}, payer chain.Name)
	GetIdxDBByIndex(index int) SecondaryDB
	GetIdxDB(idxDBName string) SecondaryDB
}

type MultiIndexValue interface {
	DBValue
	GetSecondaryValue(index int) interface{}
	SetSecondaryValue(index int, v interface{})
}

type Unpacker func([]byte) (MultiIndexValue, error)

var (
	ErrNotMultiIndexValue = errors.New("not a MultiIndexValue type")
)

func NewMultiIndex(code chain.Name, scope chain.Name, table chain.Name, idxDBNameToIndex func(string) int, indexTypes []int, unpacker ...Unpacker) *MultiIndex {

	if table.N&uint64(0x0f) != 0 {
		// Limit table names to 12 characters so that the last character (4 bits) can be used to distinguish between the secondary indices.
		panic("NewMultiIndex:Invalid multi-index table name ")
	}

	mi := &MultiIndex{}
	mi.code = code
	mi.scope = scope
	mi.table = table
	mi.DB = NewDBI64(code, scope, table)
	mi.idxDBNameToIndex = idxDBNameToIndex
	mi.IndexTypes = indexTypes
	mi.IDXDBs = make([]SecondaryDB, len(indexTypes))
	if len(unpacker) == 0 {
		mi.unpacker = nil
	} else {
		mi.unpacker = unpacker[0]
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

func (mi *MultiIndex) Store(v MultiIndexValue, payer chain.Name) {
	mi.DB.Store(v.GetPrimary(), v.Pack(), payer)
	primary := v.GetPrimary()
	for i, db := range mi.IDXDBs {
		db.Store(primary, v.GetSecondaryValue(i), payer.N)
	}
}

func (mi *MultiIndex) Set(primary uint64, v MultiIndexValue, payer chain.Name) {
	chain.Check(primary == v.GetPrimary(), "mi.Store: Invalid primary key")
	it := mi.Find(primary)
	if !it.IsOk() {
		mi.DB.Store(primary, v.Pack(), payer)
		for i, db := range mi.IDXDBs {
			db.Store(primary, v.GetSecondaryValue(i), payer.N)
		}
	} else {
		mi.Update(it, v, payer)
	}
}

func (mi *MultiIndex) Find(primary uint64) Iterator {
	return mi.DB.Find(primary)
}

//Get value by primary index
//Returns iterator and value
//iterator can be used by MultiIndex.Update method to update value
func (mi *MultiIndex) Get(id uint64) (Iterator, MultiIndexValue) {
	it, data := mi.DB.Get(id)
	if !it.IsOk() {
		return it, nil
	}
	value, err := mi.unpacker(data)
	chain.Check(err == nil, "mi.Get: Unpack error")
	_data, ok := value.(MultiIndexValue)
	chain.Check(ok, "mi.Get: Not a MultiIndexValue type")
	return it, _data
}

//Get value by primary Iterator
func (mi *MultiIndex) GetByIterator(it Iterator) (MultiIndexValue, error) {
	v, err := mi.DB.GetByIterator(it)
	if err != nil {
		return nil, err
	}
	vv, err := mi.unpacker(v)
	chain.Check(err == nil, "mi.GetByIterator: Unpack error")
	return vv, nil
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

func (mi *MultiIndex) Update(it Iterator, v MultiIndexValue, payer chain.Name) {
	mi.DB.Update(it, v.Pack(), payer)
	primary := v.GetPrimary()
	for i, db := range mi.IDXDBs {
		it, oldSecondary := db.FindByPrimary(primary)
		// logger.Println(primary, i, oldSecondary, ":")
		chain.Check(it.IsOk(), "secondary value does not exists!")
		secondary := v.GetSecondaryValue(i)
		if IsEqual(mi.IndexTypes[i], oldSecondary, secondary) {
		} else {
			db.Update(it, secondary, payer.N)
		}
	}
}

func (mi *MultiIndex) Remove(it Iterator) {
	v, err := mi.GetByIterator(it)
	chain.Check(err == nil, "mi.Remove: Invalid Iterator")
	mi.DB.Remove(it)
	for _, db := range mi.IDXDBs {
		it, secondary := db.FindByPrimary(v.GetPrimary())
		indexType := mi.GetIndexType(db.GetIndex())
		_secondary := v.GetSecondaryValue(db.GetIndex())
		chain.Check(IsEqual(indexType, secondary, _secondary), "mi.Remove: secondary value not the same")
		db.Remove(it)
	}
}

func (mi *MultiIndex) GetIndexType(index int) int {
	return mi.IndexTypes[index]
}

//Find the table row following the referenced table row in a primary 64-bit integer index table
func (mi *MultiIndex) Next(it Iterator) (next_iterator Iterator, primary uint64) {
	return mi.DB.Next(it)
}

//Find the table row preceding the referenced table row in a primary 64-bit integer index table
func (mi *MultiIndex) Previous(it Iterator) (previous_iterator Iterator, primary uint64) {
	return mi.DB.Previous(it)
}

func (mi *MultiIndex) Lowerbound(id uint64) Iterator {
	return mi.DB.Lowerbound(id)
}

func (mi *MultiIndex) Upperbound(id uint64) Iterator {
	return mi.DB.Upperbound(id)
}

func (mi *MultiIndex) End() Iterator {
	return mi.DB.End()
}

func (mi *MultiIndex) IdxFind(index int, secondary interface{}) SecondaryIterator {
	return mi.IDXDBs[index].Find(secondary)
}

func (mi *MultiIndex) UpdateSecondaryValue(idxDB SecondaryDB, primary uint64, secondary interface{}, payer chain.Name) {
	itPrimary := mi.DB.Find(primary)
	chain.Check(itPrimary.IsOk(), "primary not found!")
	v, err := mi.DB.GetByIterator(itPrimary)
	chain.Check(err == nil, "get primary value error!")
	_v, err := mi.unpacker(v)
	chain.Check(err == nil, "mi.UpdateSecondaryValue: unpack error!")

	newSecondary := _v.GetSecondaryValue(idxDB.GetIndex())
	if IsEqual(mi.IndexTypes[idxDB.GetIndex()], newSecondary, secondary) {
		return
	}
	_v.SetSecondaryValue(idxDB.GetIndex(), secondary)
	mi.DB.Update(itPrimary, _v.Pack(), payer)
}

func (mi *MultiIndex) IdxFindByName(idxDBName string, secondary interface{}) SecondaryIterator {
	chain.Check(mi.idxDBNameToIndex != nil, "idxDBNameToIndex is nil")
	index := mi.idxDBNameToIndex(idxDBName)
	return mi.IDXDBs[index].Find(secondary)
}

func (mi *MultiIndex) IdxUpdate(it SecondaryIterator, secondary interface{}, payer chain.Name) {
	idxDB := mi.IDXDBs[it.dbIndex]

	itPrimary := mi.DB.Find(it.Primary)
	chain.Check(itPrimary.IsOk(), "primary not found!")
	v, err := mi.DB.GetByIterator(itPrimary)
	chain.Check(err == nil, "get primary value error!")
	_v, err := mi.unpacker(v)
	chain.Check(err == nil, "mi.IdxUpdate: unpack error!")
	newSecondary := _v.GetSecondaryValue(idxDB.GetIndex())
	if IsEqual(mi.IndexTypes[idxDB.GetIndex()], newSecondary, secondary) {
		return
	}
	_v.SetSecondaryValue(idxDB.GetIndex(), secondary)
	mi.DB.Update(itPrimary, _v.Pack(), payer)
	idxDB.Update(it, secondary, payer.N)
}

func (mi *MultiIndex) IdxGet(itSecondary SecondaryIterator) interface{} {
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
	chain.Check(mi.idxDBNameToIndex != nil, "idxDBNameToIndex is nil")
	index := mi.idxDBNameToIndex(idxDBName)
	return mi.IDXDBs[index]
}

func (mi *MultiIndex) GetIdxDBByA1() IdxDB64I {
	return IdxDB64I{mi.IDXDBs[0], mi.IDXDBs[0].(*IdxDB64)}
}
