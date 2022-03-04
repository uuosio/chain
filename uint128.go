package chain

/*
#include <stdint.h>

typedef struct {
	uint64_t lo;
	uint64_t hi;
} uint128;

void uint128_from_uint64(uint128* a, uint64_t* b);
void uint128_to_uint64(uint128* a, uint64_t* b);
void uint128_add(uint128* a, uint128* b, uint128* c);
void uint128_sub(uint128* a, uint128* b, uint128* c);
void uint128_abs(uint128* a, uint128* b);
void uint128_mul(uint128* a, uint128* b, uint128* c);
void uint128_div(uint128* a, uint128* b, uint128* c);
int uint128_cmp(uint128* a, uint128* b);
*/
import "C"
import (
	"encoding/binary"
	"unsafe"
)

type Uint128 [16]byte

func NewUint128(lo uint64, hi uint64) Uint128 {
	var a Uint128
	binary.LittleEndian.PutUint64(a[:], lo)
	binary.LittleEndian.PutUint64(a[8:], hi)
	return a
}

func NewUint128FromBytes(bs []byte) Uint128 {
	Assert(len(bs) <= 16, "bytes too long")
	var a Uint128
	copy(a[:], bs)
	return a
}

func (n *Uint128) Pack() []byte {
	return n[:]
}

func (n *Uint128) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Read(n[:])
	return 16
}

func (t *Uint128) Size() int {
	return 16
}

func (n *Uint128) SetUint64(v uint64) {
	tmp := Uint128{}
	copy(n[:], tmp[:]) //memset
	binary.LittleEndian.PutUint64(n[:], v)
}

func (n *Uint128) Uint64() uint64 {
	return binary.LittleEndian.Uint64(n[:])
}

func (t *Uint128) Add(a, b *Uint128) *Uint128 {
	_a := unsafe.Pointer(a)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.uint128_add((*C.uint128)(_a), (*C.uint128)(_b), (*C.uint128)(_c))
	return t
}

func (t *Uint128) Sub(a, b *Uint128) *Uint128 {
	_a := unsafe.Pointer(a)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.uint128_sub((*C.uint128)(_a), (*C.uint128)(_b), (*C.uint128)(_c))
	return t
}

func (t *Uint128) Set(x *Uint128) *Uint128 {
	if t != x {
		*t = *x
	}
	return t
}

// func (t *Uint128) Abs(a *Uint128) *Uint128 {
// 	t.Set(a)
// 	_a := unsafe.Pointer(t)
// 	C.uint128_abs((*C.uint128)(_a), (*C.uint128)(_a))
// 	return t
// }

func (t *Uint128) Mul(a, b *Uint128) *Uint128 {
	_a := unsafe.Pointer(a)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.uint128_mul((*C.uint128)(_a), (*C.uint128)(_b), (*C.uint128)(_c))
	return t
}

func (t *Uint128) Div(a, b *Uint128) *Uint128 {
	_a := unsafe.Pointer(a)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.uint128_div((*C.uint128)(_a), (*C.uint128)(_b), (*C.uint128)(_c))
	return t
}

func (t *Uint128) Cmp(b *Uint128) int {
	_a := unsafe.Pointer(t)
	_b := unsafe.Pointer(b)
	return int(C.uint128_cmp((*C.uint128)(_a), (*C.uint128)(_b)))
}
