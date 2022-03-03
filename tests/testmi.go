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

//contract testmi
type TestMI struct {
}

func NewContract(receiver, firstReceiver, action chain.Name) *TestMI {
	return &TestMI{}
}

//action test1
func (t *TestMI) test1() {
	// a := big.NewInt(1)
	// b := big.NewInt(2)
	// a = a.Add(a, b)
	// logger.Println(a, b, a.Bytes())

	code := chain.NewName("hello")
	scope := chain.NewName("helloo")
	payer := code
	secondary := uint64(0)

	mi := NewMyDataDB(code, scope)
	primary := uint64(1000)
	maxPrimary := uint64(1020)
	for ; primary < maxPrimary; primary += 2 {
		if it, data := mi.Get(primary); it.IsOk() {
			logger.Println(data.a1, data.a4)
			secondary = data.a1
			data.a4 += float64(1.1)
			mi.Update(it, data, payer)
		} else {
			data := &MyData{}
			data.a1 = primary
			data.a4 = float64(primary)
			data.primary = primary
			mi.Store(data, payer)
		}
	}

	it := mi.Lowerbound(uint64(1001))
	for ; it.IsOk(); it, _ = mi.Next(it) {
		data := mi.GetByIterator(it)
		chain.Println("+++++++next value:", data.a1)
	}
	logger.Println()

	primary = uint64(1000)
	for ; primary < maxPrimary; primary += 2 {
		it, data := mi.Get(primary)
		chain.Check(it.IsOk(), "not found")
		chain.Check(data != nil, "data is nil")
	}

	{
		idxDB := mi.GetIdxDBbya1()
		it := idxDB.End()
		chain.Println("+++++++++bya1", it.I)
		if it.IsValid() {
			for {
				it = idxDB.Previous(it)
				chain.Println("+++++++bya1:", it.I)
				if !it.IsOk() {
					break
				}
				secondary := mi.IdxGet(it)
				chain.Println("+++++++bya1 secondary:", secondary.(uint64))
			}
		}
	}

	primary = uint64(1000)
	logger.Println("+++secondary:", secondary)
	{
		idxDB := mi.GetIdxDB("bya1")
		it, _secondary := idxDB.FindByPrimary(primary)
		chain.Check(it.IsOk(), "Invalid secondary iterator")
		// it, _secondary := idxDB.Lowerbound(uint64(0))
		secondary := _secondary.(uint64)
		secondary += 1
		logger.Println(idxDB.GetIndex(), it.I, it.Primary, secondary)
		if it.IsOk() {
			mi.IdxUpdate(it, secondary, payer)
		}
	}

	{
		idxDB := mi.GetIdxDBbya1()
		secondary := uint64(1001)
		it := idxDB.Find(secondary)
		chain.Check(it.IsOk(), "Invalid secondary iterator")
		// it, _secondary := idxDB.Lowerbound(uint64(0))
		secondary += 1
		logger.Println(idxDB.GetIndex(), it.I, it.Primary, secondary)
		if it.IsOk() {
			mi.IdxUpdate(it, secondary, payer)
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
			mi.IdxUpdate(it, secondary, payer)
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
			mi.IdxUpdate(it, secondary, payer)
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
			mi.IdxUpdate(it, secondary, payer)
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
			mi.IdxUpdate(it, secondary, payer)
		}
	}
}

//action test2
func (t *TestMI) test2() {
	code := chain.NewName("hello")
	scope := chain.NewName("helloo")

	mi := NewMyDataDB(code, scope)

	primary := uint64(1000)
	maxPrimary := uint64(1020)

	for {
		it := mi.Lowerbound(uint64(0))
		if !it.IsOk() {
			break
		}
		mi.Remove(it)
	}

	idxDB := mi.GetIdxDB("bya5")
	primary = uint64(1000)
	for ; primary < maxPrimary; primary += 2 {
		it, _ := idxDB.FindByPrimary(primary)
		chain.Check(!it.IsOk(), "Should be empty")
	}
}
