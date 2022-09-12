
package testprimarykey
import (
	"github.com/uuosio/chain"
    "github.com/uuosio/chain/database"
    "unsafe"
)

type sayhello struct {
    name string
}


func (t *sayhello) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.PackString(t.name)
    return enc.GetBytes()
}

func (t *sayhello) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	t.name = dec.UnpackString()
    return dec.Pos()
}

func (t *sayhello) Size() int {
    size := 0
	size += chain.PackedVarUint32Length(uint32(len(t.name))) + len(t.name) //name
    return size
}


func (t *MyData) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.PackUint64(t.primary)
	enc.PackUint64(t.n)
    return enc.GetBytes()
}

func (t *MyData) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	t.primary = dec.UnpackUint64()
	t.n = dec.UnpackUint64()
    return dec.Pos()
}

func (t *MyData) Size() int {
    size := 0
	size += 8 //primary
	size += 8 //n
    return size
}

var (
	MyDataSecondaryTypes = []int{
	}
)

func MyDataTableNameToIndex(indexName string) int {
	switch indexName {
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
		default:
			panic("index out of bound")
	}
}

func (t *MyData) SetSecondaryValue(index int, v interface{}) {
	switch index {
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
	table := chain.Name{N:uint64(10930908944415064064)} //table name: MyData
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
	return &MyDataTable{mi}
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
        case uint64(13960267442438537216): //sayhello
            t := sayhello{}
            t.Unpack(data)
            contract.SayHello(t.name)
        }
    }
    if receiver != firstReceiver {
        switch action.N {
        }
    }
}
