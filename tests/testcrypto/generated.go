
package testcrypto
import (
	"github.com/uuosio/chain"
    "github.com/uuosio/chain/database"
    "unsafe"
)

type testhash struct {
}


func (t *testhash) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
    return enc.GetBytes()
}

func (t *testhash) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
    return dec.Pos()
}

func (t *testhash) Size() int {
    size := 0
    return size
}

type testrecover struct {
    data []byte
    sig chain.Signature
    pub chain.PublicKey
}


func (t *testrecover) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.PackBytes(t.data)
	enc.Pack(&t.sig)
	enc.Pack(&t.pub)
    return enc.GetBytes()
}

func (t *testrecover) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	t.data = dec.UnpackBytes()
	dec.Unpack(&t.sig)
	dec.UnpackI(&t.pub)
    return dec.Pos()
}

func (t *testrecover) Size() int {
    size := 0
	size += chain.PackedVarUint32Length(uint32(len(t.data)))
    size += len(t.data)
	size += t.sig.Size() //sig
	size += t.pub.Size() //pub
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
        case uint64(14605620659221430272): //testhash
            t := testhash{}
            t.Unpack(data)
            contract.TestHash()
        case uint64(14605626217040883200): //testrecover
            t := testrecover{}
            t.Unpack(data)
            contract.TestRecover(t.data, &t.sig, &t.pub)
        }
    }
    if receiver != firstReceiver {
        switch action.N {
        }
    }
}
