
package testasset
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

type test2 struct {
}


func (t *test2) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
    return enc.GetBytes()
}

func (t *test2) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
    return dec.Pos()
}

func (t *test2) Size() int {
    size := 0
    return size
}

type test3 struct {
}


func (t *test3) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
    return enc.GetBytes()
}

func (t *test3) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
    return dec.Pos()
}

func (t *test3) Size() int {
    size := 0
    return size
}

type test4 struct {
}


func (t *test4) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
    return enc.GetBytes()
}

func (t *test4) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
    return dec.Pos()
}

func (t *test4) Size() int {
    size := 0
    return size
}

type test5 struct {
}


func (t *test5) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
    return enc.GetBytes()
}

func (t *test5) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
    return dec.Pos()
}

func (t *test5) Size() int {
    size := 0
    return size
}

type test11 struct {
}


func (t *test11) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
    return enc.GetBytes()
}

func (t *test11) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
    return dec.Pos()
}

func (t *test11) Size() int {
    size := 0
    return size
}

type test12 struct {
}


func (t *test12) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
    return enc.GetBytes()
}

func (t *test12) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
    return dec.Pos()
}

func (t *test12) Size() int {
    size := 0
    return size
}

type test13 struct {
}


func (t *test13) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
    return enc.GetBytes()
}

func (t *test13) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
    return dec.Pos()
}

func (t *test13) Size() int {
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
        case uint64(14605613945969442816): //test1
            t := test1{}
            t.Unpack(data)
            contract.Test1()
        case uint64(14605614495725256704): //test2
            t := test2{}
            t.Unpack(data)
            contract.Test2()
        case uint64(14605615045481070592): //test3
            t := test3{}
            t.Unpack(data)
            contract.Test3()
        case uint64(14605615595236884480): //test4
            t := test4{}
            t.Unpack(data)
            contract.Test4()
        case uint64(14605616144992698368): //test5
            t := test5{}
            t.Unpack(data)
            contract.Test5()
        case uint64(14605613963149312000): //test11
            t := test11{}
            t.Unpack(data)
            contract.Test6()
        case uint64(14605613980329181184): //test12
            t := test12{}
            t.Unpack(data)
            contract.Test12()
        case uint64(14605613997509050368): //test13
            t := test13{}
            t.Unpack(data)
            contract.Test13()
        }
    }
    if receiver != firstReceiver {
        switch action.N {
        }
    }
}
