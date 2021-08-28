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

	DB     *DBI64
	IDXDBs []SecondaryDB
}

type AccessMultiIndexValue interface {
	Store(v MultiIndexValue, payer chain.Name)
	Get(id uint64) (MultiIndexValue, error)
	GetByIterator(it Iterator) (MultiIndexValue, error)
	Update(it Iterator, v MultiIndexValue, payer chain.Name)
}

type MultiIndexValue interface {
	DBValue
	GetSecondaryValue(index int) interface{}
	SetSecondaryValue(index int, v interface{})
}

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
	mi.DB = NewDBI64(code, scope, table, unpacker...)
	mi.idxDBNameToIndex = idxDBNameToIndex
	mi.IndexTypes = indexTypes
	mi.IDXDBs = make([]SecondaryDB, len(indexTypes))
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
	mi.DB.Store(v, payer)
	primary := v.GetPrimary()
	for i, db := range mi.IDXDBs {
		db.Store(primary, v.GetSecondaryValue(i), payer.N)
	}
}

func (mi *MultiIndex) Find(primary uint64) Iterator {
	return mi.DB.Find(primary)
}

func (mi *MultiIndex) Get(id uint64) (MultiIndexValue, error) {
	data, err := mi.DB.Get(id)
	if err != nil {
		return nil, err
	}

	_data, ok := data.(MultiIndexValue)
	if !ok {
		return nil, ErrNotMultiIndexValue
	}
	return _data, nil
}

func (mi *MultiIndex) GetByIterator(it Iterator) (MultiIndexValue, error) {
	v, err := mi.DB.GetByIterator(it)
	if err != nil {
		return nil, err
	}
	_v, ok := v.(MultiIndexValue)
	if !ok {
		return nil, ErrNotMultiIndexValue
	}
	return _v, nil
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
	mi.DB.Update(it, v, payer)
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

func (mi *MultiIndex) IdxFindByName(idxDBName string, secondary interface{}) SecondaryIterator {
	chain.Check(mi.idxDBNameToIndex != nil, "idxDBNameToIndex is nil")
	index := mi.idxDBNameToIndex(idxDBName)
	return mi.IDXDBs[index].Find(secondary)
}

func (mi *MultiIndex) UpdateSecondaryValue(idxDB SecondaryDB, primary uint64, secondary interface{}, payer chain.Name) {
	itPrimary := mi.DB.Find(primary)
	chain.Check(itPrimary.IsOk(), "primary not found!")
	_v, err := mi.DB.GetByIterator(itPrimary)
	chain.Check(err == nil, "get primary value error!")
	v, ok := _v.(MultiIndexValue)
	if !ok {
		panic("Not a MultiIndexValue type")
	}

	newSecondary := v.GetSecondaryValue(idxDB.GetIndex())
	if IsEqual(mi.IndexTypes[idxDB.GetIndex()], newSecondary, secondary) {
		return
	}
	v.SetSecondaryValue(idxDB.GetIndex(), secondary)
	mi.DB.Update(itPrimary, v, payer)
}

func (mi *MultiIndex) IdxUpdate(idxDB SecondaryDB, it SecondaryIterator, secondary interface{}, payer chain.Name) {
	itPrimary := mi.DB.Find(it.Primary)
	chain.Check(itPrimary.IsOk(), "primary not found!")
	v, err := mi.DB.GetByIterator(itPrimary)
	chain.Check(err == nil, "get primary value error!")
	vv, ok := v.(MultiIndexValue)
	if !ok {
		panic("IdxUpdate: Not a MultiIndexValue type")
	}
	newSecondary := vv.GetSecondaryValue(idxDB.GetIndex())
	if IsEqual(mi.IndexTypes[idxDB.GetIndex()], newSecondary, secondary) {
		return
	}
	vv.SetSecondaryValue(idxDB.GetIndex(), secondary)
	mi.DB.Update(itPrimary, v, payer)
	idxDB.Update(it, secondary, payer.N)
}

func (mi *MultiIndex) GetIdxDBByIndex(index int) SecondaryDB {
	return mi.IDXDBs[index]
}

func (mi *MultiIndex) GetIdxDB(idxDBName string) SecondaryDB {
	chain.Check(mi.idxDBNameToIndex != nil, "idxDBNameToIndex is nil")
	index := mi.idxDBNameToIndex(idxDBName)
	return mi.IDXDBs[index]
}
