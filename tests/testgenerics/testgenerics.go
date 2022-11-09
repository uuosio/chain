package testgenerics

import (
	"github.com/uuosio/chain"
)

//table test
type MyData struct {
	id    uint64 //primary: t.id
	value uint64 //secondary
}

//table test2
type MyData2 struct {
	id    uint64        //primary: t.id
	value chain.Uint128 //secondary
}

func MapKeys[K comparable, V any](m map[K]V) []K {
    r := make([]K, 0, len(m))
    for k := range m {
        r = append(r, k)
    }
    return r
}

type List[T any] struct {
    head, tail *element[T]
}
type element[T any] struct {
    next *element[T]
    val  T
}

func (lst *List[T]) Push(v T) {
    if lst.tail == nil {
        lst.head = &element[T]{val: v}
        lst.tail = lst.head
    } else {
        lst.tail.next = &element[T]{val: v}
        lst.tail = lst.tail.next
    }
}
func (lst *List[T]) GetAll() []T {
    var elems []T
    for e := lst.head; e != nil; e = e.next {
        elems = append(elems, e.val)
    }
    return elems
}

func main() {
	receiver, firstReceiver, action := chain.GetApplyArgs()
	ContractApply(receiver.N, firstReceiver.N, action.N)
}

func ContractApply(_receiver, firstReceiver, action uint64) {
	receiver := chain.Name{_receiver}
	db := NewMyDataTable(receiver)
	db.Store(&MyData{2, 3}, receiver)
	db.Store(&MyData{4, 5}, receiver)

	idxTable := db.GetIdxTableByvalue()
	a := 0
	it, secondary := idxTable.Lowerbound(uint64(a))
	chain.Println(it.I, secondary, a)

	a = 4
	it, secondary = idxTable.Lowerbound(uint64(a))
	chain.Println(it.I, secondary, a)

    var m = map[int]string{1: "2", 2: "4", 4: "8"}
    _ = MapKeys[int, string](m)
	
    lst := List[int]{}
    lst.Push(10)
    lst.Push(13)
    lst.Push(23)
	chain.Check(lst.head.val == 10, "lst.head.val == 10")
	chain.Check(lst.head.next.val == 13, "lst.head.val == 10")
    chain.Check(lst.tail.val == 23, "lst.head.val == 23")
}
