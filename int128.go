package chain

/*
#include <stdint.h>

typedef struct {
	int64_t lo;
	int64_t hi;
} int128;

//*a = int128(*a)
void int128_from_int64(int128* a, int64_t* b);

//*b = int64(*a)
void int128_to_int64(int128* a, int64_t* b);

//*c = *a + *b
void int128_add(int128* a, int128* b, int128* c);

//*c = *a - *b
void int128_sub(int128* a, int128* b, int128* c);

//*b = abs(*a)
void int128_abs(int128* a, int128* b);

//*c = *a * *b
void int128_mul(int128* a, int128* b, int128* c);

//*c = *a / *b
void int128_div(int128* a, int128* b, int128* c);

//return 1 if *a > *b
//return 0 if *a < *b
//return -1 if *a == *b
int int128_cmp(int128* a, int128* b);
*/
import "C"
import (
	"encoding/binary"
	"unsafe"
)

type Int128 [16]byte

func NewInt128(lo uint64, hi uint64) Int128 {
	a := Int128{}
	binary.LittleEndian.PutUint64(a[:], uint64(lo))
	binary.LittleEndian.PutUint64(a[8:], uint64(hi))
	return a
}

func NewInt128FromInt64(n int64) Int128 {
	a := Int128{}
	C.int128_from_int64((*C.int128)(unsafe.Pointer(&a)), (*C.int64_t)(&n))
	return a
}

func NewInt128FromBytes(bs []byte) Int128 {
	Assert(len(bs) <= 16, "bytes too long")
	a := Int128{}
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
	C.int128_from_int64((*C.int128)(unsafe.Pointer(n)), (*C.int64_t)(&v))
}

func (n *Int128) Int64() int64 {
	ret := int64(0)
	C.int128_to_int64((*C.int128)(unsafe.Pointer(n)), (*C.int64_t)(&ret))
	return ret
}

func (t *Int128) Add(b *Int128) *Int128 {
	return t.AddEx(t, b)
}

func (t *Int128) AddEx(a, b *Int128) *Int128 {
	_a := unsafe.Pointer(a)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.int128_add((*C.int128)(_a), (*C.int128)(_b), (*C.int128)(_c))
	return t
}

func (t *Int128) Sub(a, b *Int128) *Int128 {
	return t.SubEx(t, b)
}

func (t *Int128) SubEx(a, b *Int128) *Int128 {
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

func (t *Int128) Abs() *Int128 {
	return t.AbsEx(t)
}

func (t *Int128) AbsEx(a *Int128) *Int128 {
	t.Set(a)
	_a := unsafe.Pointer(t)
	C.int128_abs((*C.int128)(_a), (*C.int128)(_a))
	return t
}

func (t *Int128) Mul(b *Int128) *Int128 {
	return t.MulEx(t, b)
}

func (t *Int128) MulEx(a, b *Int128) *Int128 {
	_a := unsafe.Pointer(a)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.int128_mul((*C.int128)(_a), (*C.int128)(_b), (*C.int128)(_c))
	return t
}

func (t *Int128) Div(b *Int128) *Int128 {
	return t.DivEx(t, b)
}

func (t *Int128) DivEx(a, b *Int128) *Int128 {
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
