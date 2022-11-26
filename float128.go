package chain

/*
#include <stddef.h>
#include <stdint.h>
#include "structs.h"

typedef double float64;
typedef float128_t float128;

void float128_from_double(float128* a, float64* b);
void float128_to_double(float128* a, float64* b);

void float128_add(float128* a, float128* b, float128* c);
void float128_sub(float128* a, float128* b, float128* c);
void float128_abs(float128* a, float128* b);
void float128_mul(float128* a, float128* b, float128* c);
void float128_div(float128* a, float128* b, float128* c);
int float128_cmp(float128* a, float128* b);
*/
import "C"
import "unsafe"

type Float128 [16]byte

func NewFloat128(v float64) Float128 {
	t := Float128{}
	C.float128_from_double((*C.float128)(unsafe.Pointer(&t)), (*C.float64)(unsafe.Pointer(&v)))
	return t
}

func (n *Float128) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.WriteBytes(n[:])
	return enc.GetSize() - oldSize
}

func (n *Float128) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Read(n[:])
	return 16
}

func (t *Float128) Size() int {
	return 16
}

func (t *Float128) SetBytes(data []byte) {
	copy(t[:], data)
}

//void float128_to_double(float128* a, float64* b);
func (t *Float128) ToFloat64() float64 {
	var v float64
	_a := unsafe.Pointer(t)
	_b := unsafe.Pointer(&v)
	C.float128_to_double((*C.float128)(_a), (*C.float64)(_b))
	return v
}

func (t *Float128) Add(b *Float128) *Float128 {
	return t.AddEx(t, b)
}

func (t *Float128) AddEx(a, b *Float128) *Float128 {
	_a := unsafe.Pointer(a)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.float128_add((*C.float128)(_a), (*C.float128)(_b), (*C.float128)(_c))
	return t
}

func (t *Float128) Sub(b *Float128) *Float128 {
	return t.SubEx(t, b)
}

func (t *Float128) SubEx(a, b *Float128) *Float128 {
	_a := unsafe.Pointer(a)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.float128_sub((*C.float128)(_a), (*C.float128)(_b), (*C.float128)(_c))
	return t
}

func (t *Float128) Set(x *Float128) *Float128 {
	if t != x {
		*t = *x
	}
	return t
}

func (t *Float128) Abs(a *Float128) *Float128 {
	t.Set(a)
	_a := unsafe.Pointer(t)
	C.float128_abs((*C.float128)(_a), (*C.float128)(_a))
	return t
}

func (t *Float128) Mul(b *Float128) *Float128 {
	return t.MulEx(t, b)
}

func (t *Float128) MulEx(a, b *Float128) *Float128 {
	_a := unsafe.Pointer(a)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.float128_mul((*C.float128)(_a), (*C.float128)(_b), (*C.float128)(_c))
	return t
}

func (t *Float128) Div(b *Float128) *Float128 {
	return t.DivEx(t, b)
}

func (t *Float128) DivEx(a, b *Float128) *Float128 {
	_a := unsafe.Pointer(a)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.float128_div((*C.float128)(_a), (*C.float128)(_b), (*C.float128)(_c))
	return t
}

func (t *Float128) Cmp(b *Float128) int {
	_a := unsafe.Pointer(t)
	_b := unsafe.Pointer(b)
	return int(C.float128_cmp((*C.float128)(_a), (*C.float128)(_b)))
}
