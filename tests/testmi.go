package main

import (
	"chain"
	"chain/database"
	"chain/logger"
)

var (
	MyDataSecondaryTypes = [5]int{
		database.IDX64, database.IDX128, database.IDX256, database.IDXFloat64, database.IDXFloat128,
	}
)

type MyData struct {
	primary uint64
	a1      uint64
	a2      chain.Uint128
	a3      chain.Uint256
	a4      float64
	a5      chain.Float128
}

func MyDataUnpacker(buf []byte) (database.DBValue, error) {
	v := &MyData{}
	_, err := v.Unpack(buf)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (d *MyData) GetPrimary() uint64 {
	return d.primary
}

func (d *MyData) GetSecondaryValue(index int) interface{} {
	if index >= len(MyDataSecondaryTypes) {
		panic("index overflow")
	}
	switch index {
	case 0:
		return d.a1
	case 1:
		return d.a2
	case 2:
		return d.a3
	case 3:
		return d.a4
	case 4:
		return d.a5
	default:
		panic("unknown index")
	}
}

func (d *MyData) SetSecondaryValue(index int, v interface{}) {
	if index >= len(MyDataSecondaryTypes) {
		panic("index overflow")
	}
	switch index {
	case 0:
		d.a1 = v.(uint64)
	case 1:
		d.a2 = v.(chain.Uint128)
	case 2:
		d.a3 = v.(chain.Uint256)
	case 3:
		d.a4 = v.(float64)
	case 4:
		d.a5 = v.(chain.Float128)
	default:
		panic("unknown index")
	}
}

func main() {
	// a := big.NewInt(1)
	// b := big.NewInt(2)
	// a = a.Add(a, b)
	// logger.Println(a, b, a.Bytes())

	code := chain.NewName("hello")
	scope := chain.NewName("helloo")
	table := chain.NewName("idxtablee5")
	payer := code
	mi := database.NewMultiIndex(code, scope, table, MyDataSecondaryTypes[:], MyDataUnpacker)

	primary := uint64(1000)
	secondary := uint64(0)
	{
		it := mi.Find(primary)
		if it.IsOk() {
			logger.Println("+++update")
			_data, _ := mi.GetByIterator(it)
			data := _data.(*MyData)
			logger.Println(data.a1, data.a4)
			secondary = data.a1
			data.a4 += float64(1.1)
			mi.Update(it, data, payer)
		} else {
			data := &MyData{}
			data.primary = primary
			mi.Store(data, payer)
		}
	}
	logger.Println("+++secondary:", secondary)
	{
		idxDB := mi.GetIdxDB(0)
		it, _secondary := idxDB.FindByPrimary(primary)
		// it, _secondary := idxDB.Lowerbound(uint64(0))
		secondary := _secondary.(uint64)
		secondary += 1
		logger.Println(idxDB.GetIndex(), it.I, it.Primary, secondary)
		if it.IsOk() {
			mi.IdxUpdate(idxDB, it, secondary, payer)
		}
	}
	{
		idxDB := mi.GetIdxDB(1)
		it, _secondary := idxDB.FindByPrimary(primary)
		// it, _secondary := idxDB.Lowerbound(uint64(0))
		secondary := _secondary.(chain.Uint128)
		n := secondary.Uint64()
		n += 1
		secondary.SetUint64(n)
		logger.Println(idxDB.GetIndex(), it.I, it.Primary, n)
		if it.IsOk() {
			mi.IdxUpdate(idxDB, it, secondary, payer)
		}
	}
	{
		idxDB := mi.GetIdxDB(2)
		it, _secondary := idxDB.FindByPrimary(primary)
		// it, _secondary := idxDB.Lowerbound(uint64(0))
		secondary := _secondary.(chain.Uint256)
		n := secondary.Uint64()
		n += 1
		secondary.SetUint64(n)
		logger.Println(idxDB.GetIndex(), it.I, it.Primary, n)
		if it.IsOk() {
			mi.IdxUpdate(idxDB, it, secondary, payer)
		}
	}
	{
		idxDB := mi.GetIdxDB(3)
		it, _secondary := idxDB.FindByPrimary(primary)
		// it, _secondary := idxDB.Lowerbound(uint64(0))
		secondary := _secondary.(float64)
		secondary += 0.1
		logger.Println(idxDB.GetIndex(), it.I, it.Primary, secondary)
		if it.IsOk() {
			mi.IdxUpdate(idxDB, it, secondary, payer)
		}
	}
	{
		idxDB := mi.GetIdxDB(4)
		it, _secondary := idxDB.FindByPrimary(primary)
		// it, _secondary := idxDB.Lowerbound(uint64(0))
		secondary := _secondary.(chain.Float128)
		secondary[0] += 1
		logger.Println(idxDB.GetIndex(), it.I, it.Primary, secondary[0])
		if it.IsOk() {
			mi.IdxUpdate(idxDB, it, secondary, payer)
		}
	}
}
