
package testaction
import (
	"github.com/uuosio/chain"
    "github.com/uuosio/chain/database"
    "unsafe"
)

type sayhello struct {
}


func (t *sayhello) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
    return enc.GetBytes()
}

func (t *sayhello) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
    return dec.Pos()
}

func (t *sayhello) Size() int {
    size := 0
    return size
}

type sayhello2 struct {
}


func (t *sayhello2) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
    return enc.GetBytes()
}

func (t *sayhello2) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
    return dec.Pos()
}

func (t *sayhello2) Size() int {
    size := 0
    return size
}

type sayhello3 struct {
}


func (t *sayhello3) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
    return enc.GetBytes()
}

func (t *sayhello3) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
    return dec.Pos()
}

func (t *sayhello3) Size() int {
    size := 0
    return size
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
            contract.SayHello()
        case uint64(13960267442439585792): //sayhello2
            t := sayhello2{}
            t.Unpack(data)
            contract.SayHello2()
        case uint64(13960267442440110080): //sayhello3
            t := sayhello3{}
            t.Unpack(data)
            contract.SayHello3()
        }
    }
    if receiver != firstReceiver {
        switch action.N {
        }
    }
}
