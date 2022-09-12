
package testmi
import (
	"github.com/uuosio/chain"
    "github.com/uuosio/chain/database"
    "unsafe"
)

type test1 struct {
}


func (t *test1) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
    return enc.GetBytes()
}

func (t *test1) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
    return dec.Pos()
}

func (t *test1) Size() int {
    size := 0
    return size
}


func (t *MyData) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.PackUint64(t.primary)
	enc.PackUint64(t.a1)
	enc.WriteBytes(t.a2[:])
	enc.Pack(&t.a3)
	enc.PackFloat64(t.a4)
	enc.Pack(&t.a5)
    return enc.GetBytes()
}

func (t *MyData) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	t.primary = dec.UnpackUint64()
	t.a1 = dec.UnpackUint64()
	dec.Unpack(&t.a2)
	dec.UnpackI(&t.a3)
	t.a4 = dec.UnpackFloat64()
	dec.UnpackI(&t.a5)
    return dec.Pos()
}

func (t *MyData) Size() int {
    size := 0
	size += 8 //primary
	size += 8 //a1
	size += 16 //a2
	size += 32 //a3
	size += 8 //a4
	size += t.a5.Size() //a5
    return size
}

var (
	MyDataSecondaryTypes = []int{
		database.IDX64,
		database.IDX128,
		database.IDX256,
		database.IDXFloat64,
		database.IDXFloat128,
	}
)

func MyDataTableNameToIndex(indexName string) int {
	switch indexName {
		case "a1":
			return 0
		case "a2":
			return 1
		case "a3":
			return 2
		case "a4":
			return 3
		case "a5":
			return 4
	default:
		panic("unknow indexName")
	}
}

func MyDataUnpacker(buf []byte) database.MultiIndexValue {
	v := &MyData{}
	v.Unpack(buf)
	return v
}

func (t *MyData) GetSecondaryValue(index int) interface{} {
	switch index {
		case 0:
			return t.a1
		case 1:
			return t.a2
		case 2:
			return t.a3
		case 3:
			return t.a4
		case 4:
			return t.a5
		default:
			panic("index out of bound")
	}
}

func (t *MyData) SetSecondaryValue(index int, v interface{}) {
	switch index {
	case 0:
		t.a1 = v.(uint64)
	case 1:
		t.a2 = v.(chain.Uint128)
	case 2:
		t.a3 = v.(chain.Uint256)
	case 3:
		t.a4 = v.(float64)
	case 4:
		t.a5 = v.(chain.Float128)
	default:
		panic("unknown index")
	}
}

func (t *MyData) GetPrimary() uint64 {
    return t.primary
}

type MyDataTable struct {
	database.MultiIndexInterface
}

func (mi *MyDataTable) Store(v *MyData, payer chain.Name) {
	mi.MultiIndexInterface.Store(v, payer)
}

func (mi *MyDataTable) GetByKey(id uint64) (*database.Iterator, *MyData) {
	it, data := mi.MultiIndexInterface.GetByKey(id)
	if !it.IsOk() {
		return it, nil
	}
	return it, data.(*MyData)
}

func (mi *MyDataTable) GetByIterator(it *database.Iterator) *MyData {
	data := mi.MultiIndexInterface.GetByIterator(it)
	return data.(*MyData)
}

func (mi *MyDataTable) Update(it *database.Iterator, v *MyData, payer chain.Name) {
	mi.MultiIndexInterface.Update(it, v, payer)
}

func NewMyDataTable(code chain.Name, optionalScope ...chain.Name) *MyDataTable {
	var scope chain.Name
	if len(optionalScope) > 0 {
		scope = optionalScope[0]
	} else {
		scope = chain.Name{N: 0}
	}
	table := chain.Name{N:uint64(10921911446417702912)} //table name: MyData
	if table.N&uint64(0x0f) != 0 {
		// Limit table names to 12 characters so that the last character (4 bits) can be used to distinguish between the secondary indices.
		panic("NewMultiIndex:Invalid multi-index table name ")
	}

	mi := &database.MultiIndex{}
	mi.SetTable(code, scope, table)
	mi.Table = database.NewTableI64(code, scope, table, func(data []byte) database.TableValue {
		return mi.Unpack(data)
	})
	mi.IdxTableNameToIndex = MyDataTableNameToIndex
	mi.IndexTypes = MyDataSecondaryTypes
	mi.IDXTables = make([]database.SecondaryTable, len(MyDataSecondaryTypes))
	mi.Unpack = MyDataUnpacker
	mi.IDXTables[0] = database.NewIdxTable64(0, code.N, scope.N, uint64(10921911446417702912)+0)
	mi.IDXTables[1] = database.NewIdxTable128(1, code.N, scope.N, uint64(10921911446417702912)+1)
	mi.IDXTables[2] = database.NewIdxTable256(2, code.N, scope.N, uint64(10921911446417702912)+2)
	mi.IDXTables[3] = database.NewIdxTableFloat64(3, code.N, scope.N, uint64(10921911446417702912)+3)
	mi.IDXTables[4] = database.NewIdxTableFloat128(4, code.N, scope.N, uint64(10921911446417702912)+4)
	return &MyDataTable{mi}
}
func (mi *MyDataTable) GetIdxTableBya1() *database.IdxTable64 {
	return mi.GetIdxTableByIndex(0).(*database.IdxTable64)
}
func (mi *MyDataTable) GetIdxTableBya2() *database.IdxTable128 {
	return mi.GetIdxTableByIndex(1).(*database.IdxTable128)
}
func (mi *MyDataTable) GetIdxTableBya3() *database.IdxTable256 {
	return mi.GetIdxTableByIndex(2).(*database.IdxTable256)
}
func (mi *MyDataTable) GetIdxTableBya4() *database.IdxTableFloat64 {
	return mi.GetIdxTableByIndex(3).(*database.IdxTableFloat64)
}
func (mi *MyDataTable) GetIdxTableBya5() *database.IdxTableFloat128 {
	return mi.GetIdxTableByIndex(4).(*database.IdxTableFloat128)
}


//eliminate unused package errors
func dummy() {
	if false {
		v := 0;
		n := unsafe.Sizeof(v);
		chain.Printui(uint64(n));
		chain.Printui(database.IDX64);
	}
}


func main() {
	receiver, firstReceiver, action := chain.GetApplyArgs()
	ContractApply(receiver.N, firstReceiver.N, action.N)
}

func ContractApply(_receiver, _firstReceiver, _action uint64) {
	receiver := chain.Name{_receiver}
	firstReceiver := chain.Name{_firstReceiver}
	action := chain.Name{_action}

	contract := NewContract(receiver, firstReceiver, action)
	if contract == nil {
		return
	}
	data := chain.ReadActionData()
	
	//Fix data declared but not used error
	if false {
		println(len(data))
	}

    if receiver == firstReceiver {
        switch action.N {
        case uint64(14605613945969442816): //test1
            t := test1{}
            t.Unpack(data)
            contract.test1()
        }
    }
    if receiver != firstReceiver {
        switch action.N {
        }
    }
}
