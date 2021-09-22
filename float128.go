package chain

/*
#include <stdint.h>

typedef char float128;
typedef uint64_t float64;

void float128_set(float128* a, float64* b);

void float128_add(float128* a, float128* b, float128* c);
void float128_sub(float128* a, float128* b, float128* c);
void float128_abs(float128* a, float128* b);
void float128_mul(float128* a, float128* b, float128* c);
int float128_cmp(float128* a, float128* b);
*/
import "C"
import "unsafe"

type Float128 [16]byte

func (n *Float128) Pack() []byte {
	return n[:]
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

func (t *Float128) Set(value float64) {
	C.float128_set((*C.float128)(unsafe.Pointer(t)), (*C.float64)(unsafe.Pointer(&value)))
}

func (t *Float128) Add(b *Float128) *Float128 {
	_a := unsafe.Pointer(t)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.float128_add((*C.float128)(_a), (*C.float128)(_b), (*C.float128)(_c))
	return t
}

func (t *Float128) Sub(b *Float128) *Float128 {
	_a := unsafe.Pointer(t)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.float128_sub((*C.float128)(_a), (*C.float128)(_b), (*C.float128)(_c))
	return t
}

func (t *Float128) Abs() *Float128 {
	_a := unsafe.Pointer(t)
	b := &Float128{}
	_b := unsafe.Pointer(b)
	C.float128_abs((*C.float128)(_a), (*C.float128)(_b))
	return b
}

func (t *Float128) Mul(b *Float128) *Float128 {
	_a := unsafe.Pointer(t)
	_b := unsafe.Pointer(b)
	_c := unsafe.Pointer(t)
	C.float128_mul((*C.float128)(_a), (*C.float128)(_b), (*C.float128)(_c))
	return t
}

func (t *Float128) Cmp(b *Float128) int {
	_a := unsafe.Pointer(t)
	_b := unsafe.Pointer(b)
	return int(C.float128_cmp((*C.float128)(_a), (*C.float128)(_b)))
}
