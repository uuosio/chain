package chain

/*
#include <stdint.h>

typedef struct {
	int64_t lo;
	int64_t hi;
} int128;

void int128_from_uint64(int128* a, uint64_t* b);
void int128_to_uint64(int128* a, uint64_t* b);
void int128_add(int128* a, int128* b, int128* c);
void int128_sub(int128* a, int128* b, int128* c);
void int128_abs(int128* a, int128* b);
void int128_mul(int128* a, int128* b, int128* c);
void int128_div(int128* a, int128* b, int128* c);
int int128_cmp(int128* a, int128* b);
*/
import "C"
import (
	"encoding/binary"
	"unsafe"
)

type Int128 [16]byte

func NewInt128(lo int64, hi int64) *Int128 {
	a := &Int128{}
	binary.LittleEndian.PutUint64(a[:], uint64(lo))
	binary.LittleEndian.PutUint64(a[8:], uint64(hi))
	return a
}

func NewInt128FromBytes(bs []byte) *Int128 {
	Assert(len(bs) <= 16, "bytes too long")
	a := &Int128{}
	copy(a[:], bs)
	return a
}

func (n *Int128) Pack() []byte {
	return n[:]
}

func (n *Int128) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Read(n[:])
	return 16
}

func (t *Int128) Size() int {
	return 16
}

func (n *Int128) SetInt64(v int64) {
	tmp := Int128{}
	copy(n[:], tmp[:]) //memset
	binary.LittleEndian.PutUint64(n[:], uint64(v))
}

func (n *Int128) Int64() int64 {
	return int64(binary.LittleEndian.Uint64(n[:]))
}

func (t *Int128) Add(a, b *Int128) *Int128 {
	_a := unsafe.Pointer(a)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.int128_add((*C.int128)(_a), (*C.int128)(_b), (*C.int128)(_c))
	return t
}

func (t *Int128) Sub(a, b *Int128) *Int128 {
	_a := unsafe.Pointer(a)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.int128_sub((*C.int128)(_a), (*C.int128)(_b), (*C.int128)(_c))
	return t
}

func (t *Int128) Set(x *Int128) *Int128 {
	if t != x {
		*t = *x
	}
	return t
}

// func (t *Int128) Abs(a *Int128) *Int128 {
// 	t.Set(a)
// 	_a := unsafe.Pointer(t)
// 	C.int128_abs((*C.int128)(_a), (*C.int128)(_a))
// 	return t
// }

func (t *Int128) Mul(a, b *Int128) *Int128 {
	_a := unsafe.Pointer(a)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.int128_mul((*C.int128)(_a), (*C.int128)(_b), (*C.int128)(_c))
	return t
}

func (t *Int128) Div(a, b *Int128) *Int128 {
	_a := unsafe.Pointer(a)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.int128_div((*C.int128)(_a), (*C.int128)(_b), (*C.int128)(_c))
	return t
}

func (t *Int128) Cmp(b *Int128) int {
	_a := unsafe.Pointer(t)
	_b := unsafe.Pointer(b)
	return int(C.int128_cmp((*C.int128)(_a), (*C.int128)(_b)))
}
