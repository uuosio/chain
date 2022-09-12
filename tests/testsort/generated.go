
package testsort
import (
	"github.com/uuosio/chain"
    "github.com/uuosio/chain/database"
    "unsafe"
)

type test struct {
    pubs []chain.PublicKey
}


func (t *test) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	
	{
		enc.PackLength(len(t.pubs))
		for i := range t.pubs {
			enc.Pack(&t.pubs[i])
		}
	}
    return enc.GetBytes()
}

func (t *test) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	
	{
		length := dec.UnpackLength()
		t.pubs = make([]chain.PublicKey, length)
		for i:=0; i<length; i++ {
		dec.UnpackI(&t.pubs[i])
		}
	}
    return dec.Pos()
}

func (t *test) Size() int {
    size := 0
	size += chain.PackedVarUint32Length(uint32(len(t.pubs)))
    
    for i := range t.pubs {
        size += t.pubs[i].Size()
    }
    return size
}


func (t *MySingleton) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.PackUint64(t.a1)
    return enc.GetBytes()
}

func (t *MySingleton) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	t.a1 = dec.UnpackUint64()
    return dec.Pos()
}

func (t *MySingleton) Size() int {
    size := 0
	size += 8 //a1
    return size
}

var (
	MySingletonSecondaryTypes = []int{
	}
)

func MySingletonTableNameToIndex(indexName string) int {
	switch indexName {
	default:
		panic("unknow indexName")
	}
}

func MySingletonUnpacker(buf []byte) database.MultiIndexValue {
	v := &MySingleton{}
	v.Unpack(buf)
	return v
}

func (t *MySingleton) GetSecondaryValue(index int) interface{} {
	switch index {
		default:
			panic("index out of bound")
	}
}

func (t *MySingleton) SetSecondaryValue(index int, v interface{}) {
	switch index {
	default:
		panic("unknown index")
	}
}

func (d *MySingleton) GetPrimary() uint64 {
	return uint64(10921911447491444736)
}

type MySingletonTable struct {
	db *database.SingletonTable
}

func NewMySingletonTable(code chain.Name, optionalScope ...chain.Name) *MySingletonTable {
	var scope chain.Name
	if len(optionalScope) > 0 {
		scope = optionalScope[0]
	} else {
		scope = chain.Name{N: 0}
	}
	chain.Check(code != chain.Name{0}, "bad code name")
	table := chain.Name{N:uint64(10921911447491444736)}
	db := database.NewSingletonTable(code, scope, table, MySingletonUnpacker)
	return &MySingletonTable{db}
}

func (t *MySingletonTable) Set(data *MySingleton, payer chain.Name) {
	t.db.Set(data, payer)
}

func (t *MySingletonTable) Get() (*MySingleton) {
	data := t.db.Get()
	if data == nil {
		return nil
	}
	return data.(*MySingleton)
}

func (t *MySingletonTable) Remove() {
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
        case uint64(14605613396213628928): //test
            t := test{}
            t.Unpack(data)
            contract.test(t.pubs)
        }
    }
    if receiver != firstReceiver {
        switch action.N {
        }
    }
}
