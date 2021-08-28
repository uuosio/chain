package main

import (
	"github.com/uuosio/chain"
	"github.com/uuosio/chain/logger"
)

//table mydata
type MyData struct {
	primary uint64         //primary:t.primary
	a1      uint64         //IDX64:bya1:t.a1:t.a1=%v
	a2      chain.Uint128  //IDX128:bya2:t.a2:t.a2=%v
	a3      chain.Uint256  //IDX256:bya3:t.a3:t.a3=%v
	a4      float64        //IDXFloat64:bya4:t.a4:t.a4=%v
	a5      chain.Float128 //IDXFloat128:bya5:t.a5:t.a5=%v
}

func main() {
	// a := big.NewInt(1)
	// b := big.NewInt(2)
	// a = a.Add(a, b)
	// logger.Println(a, b, a.Bytes())

	code := chain.NewName("hello")
	scope := chain.NewName("helloo")
	payer := code
	mi := NewMyDataDB(code, scope)

	primary := uint64(1000)
	secondary := uint64(0)
	{
		it := mi.Find(primary)
		if it.IsOk() {
			logger.Println("+++update")
			data, _ := mi.GetByIterator(it)
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
		idxDB := mi.GetIdxDB("bya1")
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
		idxDB := mi.GetIdxDB("bya2")
		it, _secondary := idxDB.FindByPrimary(primary)
		// it, _secondary := idxDB.Lowerbound(uint64(0))
		secondary := _secondary.(chain.Uint128)
		n := secondary.Uint64()
		n += 2
		secondary.SetUint64(n)
		logger.Println(idxDB.GetIndex(), it.I, it.Primary, n)
		if it.IsOk() {
			mi.IdxUpdate(idxDB, it, secondary, payer)
		}
	}
	{
		idxDB := mi.GetIdxDB("bya3")
		it, _secondary := idxDB.FindByPrimary(primary)
		// it, _secondary := idxDB.Lowerbound(uint64(0))
		secondary := _secondary.(chain.Uint256)
		n := secondary.Uint64()
		n += 3
		secondary.SetUint64(n)
		logger.Println(idxDB.GetIndex(), it.I, it.Primary, n)
		if it.IsOk() {
			mi.IdxUpdate(idxDB, it, secondary, payer)
		}
	}
	{
		idxDB := mi.GetIdxDB("bya4")
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
		idxDB := mi.GetIdxDB("bya5")
		it, _secondary := idxDB.FindByPrimary(primary)
		// it, _secondary := idxDB.Lowerbound(uint64(0))
		secondary := _secondary.(chain.Float128)
		secondary[0] += 5
		logger.Println(idxDB.GetIndex(), it.I, it.Primary, secondary[0])
		if it.IsOk() {
			mi.IdxUpdate(idxDB, it, secondary, payer)
		}
	}
}
