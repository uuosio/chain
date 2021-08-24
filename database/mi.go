package database

import (
	"chain"
	"errors"
)

type MultiIndex struct {
	code       chain.Name
	scope      chain.Name
	table      chain.Name
	db         *DBI64
	indexTypes []int
	idxDBs     []SecondaryDB
}

type MultiIndexValue interface {
	DBValue
	GetSecondaryValue(index int) interface{}
	SetSecondaryValue(index int, v interface{})
}

var (
	ErrNotMultiIndexValue = errors.New("not a MultiIndexValue type")
)

func NewMultiIndex(code chain.Name, scope chain.Name, table chain.Name, indexTypes []int, unpacker ...Unpacker) *MultiIndex {
	mi := &MultiIndex{}
	mi.code = code
	mi.scope = scope
	mi.table = table
	mi.db = NewDBI64(code, scope, table, unpacker...)
	mi.indexTypes = indexTypes
	mi.idxDBs = make([]SecondaryDB, len(indexTypes))
	for i, v := range indexTypes {
		switch v {
		case IDX64:
			mi.idxDBs[i] = NewIdxDB64(i, code.N, scope.N, table.N)
			break
		case IDX128:
			mi.idxDBs[i] = NewIdxDB128(i, code.N, scope.N, table.N)
			break
		case IDX256:
			mi.idxDBs[i] = NewIdxDB256(i, code.N, scope.N, table.N)
			break
		case IDXFloat64:
			mi.idxDBs[i] = NewIdxDBFloat64(i, code.N, scope.N, table.N)
			break
		case IDXFloat128:
			mi.idxDBs[i] = NewIdxDBFloat128(i, code.N, scope.N, table.N)
		default:
			panic("invalid index")
		}
		// mi.secondaryDBs[v] = NewSecondaryDB(code.N, scope.N, table.N, v)
	}
	return mi
}

func (mi *MultiIndex) Store(v MultiIndexValue, payer chain.Name) {
	mi.db.Store(v, payer)
	primary := v.GetPrimary()
	for i, db := range mi.idxDBs {
		db.Store(primary, v.GetSecondaryValue(i), payer.N)
	}
}

func (mi *MultiIndex) Find(primary uint64) Iterator {
	return mi.db.Find(primary)
}

func (mi *MultiIndex) Get(id uint64) (MultiIndexValue, error) {
	data, err := mi.db.Get(id)
	if err != nil {
		return nil, err
	}

	_data, ok := data.(MultiIndexValue)
	if !ok {
		return nil, ErrNotMultiIndexValue
	}
	return _data, nil
	// v, err := mi.db.GetByIterator(it)
	// if err != nil {
	// 	return nil
	// }
	// return v
	// TODO verify secondary values
}

func (mi *MultiIndex) GetByIterator(it Iterator) (MultiIndexValue, error) {
	v, err := mi.db.GetByIterator(it)
	if err != nil {
		return nil, err
	}
	_v, ok := v.(MultiIndexValue)
	if !ok {
		return nil, ErrNotMultiIndexValue
	}
	return _v, nil

	// v, err := mi.db.GetByIterator(it)
	// if err != nil {
	// 	return nil
	// }
	// return v
	// TODO verify secondary values
}

func isEqual(indexType int, a, b interface{}) bool {
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
	mi.db.Update(it, v, payer)
	primary := v.GetPrimary()
	for i, db := range mi.idxDBs {
		it, oldSecondary := db.FindByPrimary(primary)
		// logger.Println(primary, i, oldSecondary, ":")
		chain.Check(it.IsOk(), "secondary value does not exists!")
		secondary := v.GetSecondaryValue(i)
		if isEqual(mi.indexTypes[i], oldSecondary, secondary) {
		} else {
			db.Update(it, secondary, payer.N)
		}
	}
}

func (mi *MultiIndex) Lowerbound(id uint64) Iterator {
	return mi.db.Lowerbound(id)
}

func (mi *MultiIndex) Upperbound(id uint64) Iterator {
	return mi.db.Upperbound(id)
}

func (mi *MultiIndex) End() Iterator {
	return mi.db.End()
}

func (mi *MultiIndex) IdxFind(index int, secondary interface{}) SecondaryIterator {
	return mi.idxDBs[index].Find(secondary)
}

func (mi *MultiIndex) UpdateSecondaryValue(idxDB SecondaryDB, primary uint64, secondary interface{}, payer chain.Name) {
	itPrimary := mi.db.Find(primary)
	chain.Check(itPrimary.IsOk(), "primary not found!")
	_v, err := mi.db.GetByIterator(itPrimary)
	chain.Check(err == nil, "get primary value error!")
	v, ok := _v.(MultiIndexValue)
	if !ok {
		panic("Not a MultiIndexValue type")
	}

	newSecondary := v.GetSecondaryValue(idxDB.GetIndex())
	if isEqual(mi.indexTypes[idxDB.GetIndex()], newSecondary, secondary) {
		return
	}
	v.SetSecondaryValue(idxDB.GetIndex(), secondary)
	mi.db.Update(itPrimary, v, payer)
}

func (mi *MultiIndex) IdxUpdate(idxDB SecondaryDB, it SecondaryIterator, secondary interface{}, payer chain.Name) {
	itPrimary := mi.db.Find(it.Primary)
	chain.Check(itPrimary.IsOk(), "primary not found!")
	v, err := mi.db.GetByIterator(itPrimary)
	chain.Check(err == nil, "get primary value error!")
	vv, ok := v.(MultiIndexValue)
	if !ok {
		panic("Not a MultiIndexValue type")
	}
	newSecondary := vv.GetSecondaryValue(idxDB.GetIndex())
	if isEqual(mi.indexTypes[idxDB.GetIndex()], newSecondary, secondary) {
		return
	}
	vv.SetSecondaryValue(idxDB.GetIndex(), secondary)
	mi.db.Update(itPrimary, v, payer)
	idxDB.Update(it, secondary, payer.N)
}

func (mi *MultiIndex) GetIdxDB(index int) SecondaryDB {
	return mi.idxDBs[index]
}
