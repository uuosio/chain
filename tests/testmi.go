package main

import (
	"github.com/uuosio/chain"
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

func check(b bool, msg string) {
	chain.Check(b, msg)
}

//contract testmi
type TestMI struct {
	receiver, firstReceiver, action chain.Name
	mi                              *MyDataDB
}

func NewContract(receiver, firstReceiver, action chain.Name) *TestMI {
	code := receiver
	scope := chain.NewName("helloo")
	mi := NewMyDataDB(code, scope)

	return &TestMI{receiver, firstReceiver, action, mi}
}

func (t *TestMI) initTest() *MyDataDB {
	payer := t.receiver

	mi := t.mi

	data := &MyData{}
	data.primary = 1
	data.a1 = 1
	data.a2 = chain.NewUint128(2, 0)
	data.a3 = chain.NewUint256(3, 0, 0, 0)
	data.a4 = 4.0
	data.a5 = chain.NewFloat128(5.0)
	mi.Store(data, payer)

	data = &MyData{}
	data.primary = 11
	data.a1 = 11
	data.a2 = chain.NewUint128(22, 0)
	data.a3 = chain.NewUint256(33, 0, 0, 0)
	data.a4 = 44.0
	data.a5 = chain.NewFloat128(55.0)
	mi.Store(data, payer)

	data = &MyData{}
	data.primary = 111
	data.a1 = 111
	data.a2 = chain.NewUint128(222, 0)
	data.a3 = chain.NewUint256(333, 0, 0, 0)
	data.a4 = 444.0
	data.a5 = chain.NewFloat128(555.0)
	mi.Store(data, payer)
	return mi
}

func (t *TestMI) clearTest() {
	for {
		it := t.mi.Lowerbound(0)
		if !it.IsOk() {
			break
		}
		t.mi.Remove(it)
	}
}

//action test1
func (t *TestMI) test1() {
	payer := t.receiver
	mi := t.initTest()
	{
		check(mi.Find(1).IsOk(), "mi.Find(1).IsOk()")
		it := mi.Lowerbound(1)
		check(it.IsOk(), "mi.Lowerbound(1)")

		it = mi.Lowerbound(2)
		check(it.IsOk(), "mi.Lowerbound(2)")
		check(mi.GetByIterator(it).primary == 11, "mi.Lowerbound(2)")

		it = mi.Upperbound(1)
		check(it.IsOk(), "mi.Upperbound(1)")
		check(mi.GetByIterator(it).primary == 11, "mi.Upperbound(1)")

		it = mi.Find(1)
		check(it.IsOk(), "mi.Find(1).IsOk()")
		it, primary := mi.Next(it)
		check(it.IsOk(), "mi.Next(it)")
		check(primary == 11, "mi.Next(it)")

		data := mi.GetByIterator(it)
		check(data.primary == 11, "data.primary == 11")
		check(data.a1 == 11, "data.a1 == 11")
		check(data.a2 == chain.NewUint128(22, 0), "data.a2 == *chain.NewUint128(22, 0)")
		check(data.a3 == chain.NewUint256(33, 0, 0, 0), "data.a3 == *chain.NewUint256(33, 0, 0, 0)")
		check(data.a4 == 44.0, "data.a4 == 44.0")
		check(data.a5 == chain.NewFloat128(55.0), "data.a5 == *chain.NewFloat128(55.0)")

		it, primary = mi.Previous(it)
		data = mi.GetByIterator(it)
		check(data.primary == 1, "data.primary == 1")
		check(data.a1 == 1, "data.a1 == 1")
		check(data.a2 == chain.NewUint128(2, 0), "data.a2 == *chain.NewUint128(2, 0)")
		check(data.a3 == chain.NewUint256(3, 0, 0, 0), "data.a3 == *chain.NewUint256(3, 0, 0, 0)")
		check(data.a4 == 4.0, "data.a4 == 4.0")
		check(data.a5 == chain.NewFloat128(5.0), "data.a5 == *chain.NewFloat128(5.0)")

		it = mi.End()
		it, primary = mi.Previous(it)
		data = mi.GetByIterator(it)
		check(data.primary == 111, "data.primary == 1")
		check(data.a1 == 111, "data.a1 == 1")
		check(data.a2 == chain.NewUint128(222, 0), "data.a2 == *chain.NewUint128(222, 0)")
		check(data.a3 == chain.NewUint256(333, 0, 0, 0), "data.a3 == *chain.NewUint256(333, 0, 0, 0)")
		check(data.a4 == 444.0, "data.a4 == 4.0")
		check(data.a5 == chain.NewFloat128(555.0), "data.a5 == *chain.NewFloat128(555.0)")
	}

	{
		idx64 := mi.GetIdxDBbya1()
		it, secondary := idx64.Lowerbound(1)
		check(it.IsOk(), "idx64.Lowerbound(1)")
		check(secondary == 1, "secondary == 1")

		it, secondary = idx64.Upperbound(1)
		check(it.IsOk(), "idx64.Upperbound(1)")
		check(secondary == 11, "secondary == 11")

		it = idx64.Find(1)
		check(it.IsOk(), "idx64.Find(1)")
		check(it.Primary == 1, "it.Primary == 1")

		it = idx64.Next(it)
		check(it.IsOk(), "idx64.Next(it)")
		check(it.Primary == 11, "idx64.Next(it)")

		it = idx64.Previous(it)
		check(it.IsOk(), "idx64.Previous(it)")
		check(it.Primary == 1, "idx64.Previous(it)")

		it = idx64.End()
		it = idx64.Previous(it)
		check(it.IsOk(), "idx64.End()")
		check(it.Primary == 111, "idx64.End()")
	}

	{
		idx128 := mi.GetIdxDBbya2()
		it, secondary := idx128.Lowerbound(chain.NewUint128(2, 0))
		check(it.IsOk(), "idx128.Lowerbound(1)")
		check(secondary == chain.NewUint128(2, 0), "secondary == chain.NewUint128(2, 0)")

		it, secondary = idx128.Upperbound(chain.NewUint128(2, 0))
		check(it.IsOk(), "idx128.Upperbound(chain.NewUint128(1, 0))")
		check(secondary == chain.NewUint128(22, 0), "*secondary == *chain.NewUint128(11, 0)")

		it = idx128.Find(chain.NewUint128(2, 0))
		check(it.IsOk(), "idx128.Find(chain.NewUint128(1, 0))")
		check(it.Primary == 1, "it.Primary == 11")

		it = idx128.Next(it)
		check(it.IsOk(), "idx128.Next(it)")
		check(it.Primary == 11, "idx128.Next(it)")

		it = idx128.Previous(it)
		check(it.IsOk(), "idx128.Previous(it)")
		check(it.Primary == 1, "idx128.Previous(it)")

		it = idx128.End()
		it = idx128.Previous(it)
		check(it.IsOk(), "idx128.End()")
		check(it.Primary == 111, "idx128.End()")
	}

	{
		idx256 := mi.GetIdxDBbya3()
		it, secondary := idx256.Lowerbound(chain.NewUint256(3, 0, 0, 0))
		check(it.IsOk(), "idx256.Lowerbound(1)")
		check(secondary == chain.NewUint256(3, 0, 0, 0), "secondary == chain.NewUint256(3, 0, 0, 0)")

		it, secondary = idx256.Upperbound(chain.NewUint256(3, 0, 0, 0))
		check(it.IsOk(), "idx256.Upperbound(chain.NewUint128(1, 0))")
		check(secondary == chain.NewUint256(33, 0, 0, 0), "secondary == *chain.NewUint256(33, 0, 0, 0)")

		it = idx256.Find(chain.NewUint256(3, 0, 0, 0))
		check(it.IsOk(), "idx256.Find(chain.NewUint256(3, 0, 0, 0))")
		check(it.Primary == 1, "it.Primary == 11")

		it = idx256.Next(it)
		check(it.IsOk(), "idx256.Next(it)")
		check(it.Primary == 11, "idx256.Next(it)")

		it = idx256.Previous(it)
		check(it.IsOk(), "idx256.Previous(it)")
		check(it.Primary == 1, "idx256.Previous(it)")

		it = idx256.End()
		it = idx256.Previous(it)
		check(it.IsOk(), "idx256.End()")
		check(it.Primary == 111, "idx256.End()")
	}

	{
		idxf64 := mi.GetIdxDBbya4()
		it, secondary := idxf64.Lowerbound(4.0)
		check(it.IsOk(), "idxf64.Lowerbound(4.0)")
		check(secondary == 4.0, "secondary == 4.0")

		it, secondary = idxf64.Upperbound(4.0)
		check(it.IsOk(), "idxf64.Upperbound(4.0)")
		check(secondary == 44.0, "secondary == 44.0")

		it = idxf64.Find(4.0)
		check(it.IsOk(), "idxf64.Find(4.0)")
		check(it.Primary == 1, "it.Primary == 1")

		it = idxf64.Next(it)
		check(it.IsOk(), "idxf64.Next(it)")
		check(it.Primary == 11, "idxf64.Next(it)")

		it = idxf64.Previous(it)
		check(it.IsOk(), "idxf64.Previous(it)")
		check(it.Primary == 1, "idxf64.Previous(it)")

		it = idxf64.End()
		it = idxf64.Previous(it)
		check(it.IsOk(), "idxf64.End()")
		check(it.Primary == 111, "idxf64.End()")
	}

	{
		idxf128 := mi.GetIdxDBbya5()
		it, secondary := idxf128.Lowerbound(chain.NewFloat128(5.0))
		check(it.IsOk(), "idxf128.Lowerbound(1)")
		check(secondary == chain.NewFloat128(5.0), "secondary == chain.NewFloat128(5.0)")

		it, secondary = idxf128.Upperbound(chain.NewFloat128(5.0))
		check(it.IsOk(), "idxf128.Upperbound(chain.NewFloat128(5.0))")
		check(secondary == chain.NewFloat128(55.0), "secondary == chain.NewFloat128(55.0)")

		it = idxf128.Find(chain.NewFloat128(5.0))
		check(it.IsOk(), "idxf128.Find(chain.NewFloat128(5.0))")
		check(it.Primary == 1, "it.Primary == 1")

		it = idxf128.Next(it)
		check(it.IsOk(), "idxf128.Next(it)")
		check(it.Primary == 11, "idxf128.Next(it)")

		it = idxf128.Previous(it)
		check(it.IsOk(), "idxf128.Previous(it)")
		check(it.Primary == 1, "idxf128.Previous(it)")

		it = idxf128.End()
		it = idxf128.Previous(it)
		check(it.IsOk(), "idxf128.End()")
		check(it.Primary == 111, "idxf128.End()")
	}
	// {
	// 	it, data := mi.GetByKey(1)
	// 	check(it.IsOk(), "bad value")
	// 	check(data.primary == 1, "bad value")
	// 	mi.Update(it, data, payer)
	// }
	{
		idx64 := mi.GetIdxDBbya1()
		it := idx64.Find(1)
		check(it.IsOk(), "idx64.Find(1)")
		mi.IdxUpdate(it, uint64(2), payer)
		itPrimary, data := mi.GetByKey(it.Primary)
		check(itPrimary.IsOk(), "mi.GetByKey(it.Primary)")
		check(data.primary == 1, "data.primary == 1")
		check(data.a1 == 2, "data.a1 == 2")
	}

	{
		idx128 := mi.GetIdxDBbya2()
		it := idx128.Find(chain.NewUint128(2, 0))
		check(it.IsOk(), "idx128.Find(chain.NewUint128(2, 0))")
		mi.IdxUpdate(it, chain.NewUint128(3, 0), payer)
		itPrimary, data := mi.GetByKey(it.Primary)
		check(itPrimary.IsOk(), "mi.GetByKey(it.Primary)")
		check(data.primary == 1, "data.primary == 1")
		check(data.a2 == chain.NewUint128(3, 0), "data.a1 == 2")
	}

	{
		idx256 := mi.GetIdxDBbya3()
		it := idx256.Find(chain.NewUint256(3, 0, 0, 0))
		check(it.IsOk(), "idx256.Find(chain.NewUint256(3, 0, 0, 0))")
		mi.IdxUpdate(it, chain.NewUint256(4, 0, 0, 0), payer)
		itPrimary, data := mi.GetByKey(it.Primary)
		check(itPrimary.IsOk(), "mi.GetByKey(it.Primary)")
		check(data.primary == 1, "data.primary == 1")
		check(data.a3 == chain.NewUint256(4, 0, 0, 0), "data.a3 == chain.NewUint256(4, 0, 0, 0)")
	}

	{
		idxf64 := mi.GetIdxDBbya4()
		it := idxf64.Find(4.0)
		check(it.IsOk(), "idx64.Find(4.0)")
		mi.IdxUpdate(it, float64(5.0), payer)
		itPrimary, data := mi.GetByKey(it.Primary)
		check(itPrimary.IsOk(), "mi.GetByKey(it.Primary)")
		check(data.primary == 1, "data.primary == 1")
		check(data.a4 == 5.0, "data.a4 == 5.0")
	}
	{
		idxf128 := mi.GetIdxDBbya5()
		it := idxf128.Find(chain.NewFloat128(5.0))
		check(it.IsOk(), "idx64.Find(chain.NewFloat128(5.0))")
		mi.IdxUpdate(it, chain.NewFloat128(6.0), payer)
		itPrimary, data := mi.GetByKey(it.Primary)
		check(itPrimary.IsOk(), "mi.GetByKey(it.Primary)")
		check(data.primary == 1, "data.primary == 1")
		check(data.a5 == chain.NewFloat128(6.0), "data.a5 == 6.0")
	}

	{
		it := mi.Find(1)
		mi.Remove(it)

		it = mi.Find(11)
		mi.Remove(it)

		it = mi.Find(111)
		mi.Remove(it)

		it = mi.Find(1)
		check(!it.IsOk(), "bad")

		it = mi.Find(11)
		check(!it.IsOk(), "bad")

		it = mi.Find(111)
		check(!it.IsOk(), "bad")
	}

	t.initTest()
	t.clearTest()
	{
		idx := mi.GetIdxDBbya1()
		it := idx.Find(1)
		check(!it.IsOk(), "bad return")
	}

	{
		idx := mi.GetIdxDBbya2()
		it := idx.Find(chain.NewUint128(2, 0))
		check(!it.IsOk(), "bad return")
	}

	{
		idx := mi.GetIdxDBbya3()
		it := idx.Find(chain.NewUint256(3, 0, 0, 0))
		check(!it.IsOk(), "bad return")
	}
	{
		idx := mi.GetIdxDBbya4()
		it := idx.Find(4.0)
		check(!it.IsOk(), "bad return")
	}
	{
		idx := mi.GetIdxDBbya5()
		it := idx.Find(chain.NewFloat128(5.0))
		check(!it.IsOk(), "bad return")
	}
}
