
package testvariant
import (
	"github.com/uuosio/chain"
    "github.com/uuosio/chain/database"
    "unsafe"
)

type testvariant struct {
    v MyVariant
}


func (t *testvariant) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.Pack(&t.v)
    return enc.GetBytes()
}

func (t *testvariant) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	dec.UnpackI(&t.v)
    return dec.Pos()
}

func (t *testvariant) Size() int {
    size := 0
	size += t.v.Size() //v
    return size
}


func (t *MyTable) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.Pack(&t.a)
    return enc.GetBytes()
}

func (t *MyTable) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	dec.UnpackI(&t.a)
    return dec.Pos()
}

func (t *MyTable) Size() int {
    size := 0
	size += t.a.Size() //a
    return size
}


func NewMyVariant(value interface{}) *MyVariant {
	ret := &MyVariant{}
	switch value.(type) {
		case *uint64:
				ret.value = value
		case *chain.Uint128:
				ret.value = value
		default:
			chain.Check(false, "unknown variant type")	
	}
	return ret
}

func (t *MyVariant) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	if _, ok := t.value.(*uint64); ok {
		enc.PackUint8(uint8(0))
		enc.PackUint64((*(t.value.(*uint64))))
		return enc.GetBytes()
	}
	if _, ok := t.value.(*chain.Uint128); ok {
		enc.PackUint8(uint8(1))
		enc.WriteBytes((*(t.value.(*chain.Uint128)))[:])
		return enc.GetBytes()
	}
    return enc.GetBytes()
}

func (t *MyVariant) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	tp := dec.UnpackUint8()
	if tp == uint8(0) {
		type A struct {
			value uint64
		}
		a := &A{}
		dec.Unpack(&a.value)
		t.value = &a.value
	}
	if tp == uint8(1) {
		type A struct {
			value chain.Uint128
		}
		a := &A{}
		dec.Unpack(&a.value)
		t.value = &a.value
	}
    return dec.Pos()
}

func (t *MyVariant) Size() int {
    size := 1
	if _, ok := t.value.(uint64); ok {
		size += 8 //value.(*uint64)
		return size
	}
	if _, ok := t.value.(chain.Uint128); ok {
		size += 16 //value.(*chain.Uint128)
		return size
	}
    return size
}

var (
	MyTableSecondaryTypes = []int{
	}
)

func MyTableTableNameToIndex(indexName string) int {
	switch indexName {
	default:
		panic("unknow indexName")
	}
}

func MyTableUnpacker(buf []byte) database.MultiIndexValue {
	v := &MyTable{}
	v.Unpack(buf)
	return v
}

func (t *MyTable) GetSecondaryValue(index int) interface{} {
	switch index {
		default:
			panic("index out of bound")
	}
}

func (t *MyTable) SetSecondaryValue(index int, v interface{}) {
	switch index {
	default:
		panic("unknown index")
	}
}

func (d *MyTable) GetPrimary() uint64 {
	return uint64(10930908944415064064)
}

type MyTableTable struct {
	db *database.SingletonTable
}

func NewMyTableTable(code chain.Name, optionalScope ...chain.Name) *MyTableTable {
	var scope chain.Name
	if len(optionalScope) > 0 {
		scope = optionalScope[0]
	} else {
		scope = chain.Name{N: 0}
	}
	chain.Check(code != chain.Name{0}, "bad code name")
	table := chain.Name{N:uint64(10930908944415064064)}
	db := database.NewSingletonTable(code, scope, table, MyTableUnpacker)
	return &MyTableTable{db}
}

func (t *MyTableTable) Set(data *MyTable, payer chain.Name) {
	t.db.Set(data, payer)
}

func (t *MyTableTable) Get() (*MyTable) {
	data := t.db.Get()
	if data == nil {
		return nil
	}
	return data.(*MyTable)
}

func (t *MyTableTable) Remove() {
	t.db.Remove()
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
        case uint64(14605628355286200832): //testvariant
            t := testvariant{}
            t.Unpack(data)
            contract.TestVariant(t.v)
        }
    }
    if receiver != firstReceiver {
        switch action.N {
        }
    }
}
