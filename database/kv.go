package database

/*
#include <stdint.h>
#define bool char

int64_t kv_erase(uint64_t contract, const char* key, uint32_t key_size);
int64_t kv_set(uint64_t contract, const char* key, uint32_t key_size, const char* value, uint32_t value_size, uint64_t payer);
//bool kv_get(uint64_t contract, const char* key, uint32_t key_size, uint32_t& value_size);
bool kv_get(uint64_t contract, const char* key, uint32_t key_size, uint32_t* value_size);
uint32_t kv_get_data(uint32_t offset, char* data, uint32_t data_size);

uint32_t kv_it_create(uint64_t contract, const char* prefix, uint32_t size);
void kv_it_destroy(uint32_t itr);
int32_t kv_it_status(uint32_t itr);
int32_t kv_it_compare(uint32_t itr_a, uint32_t itr_b);
int32_t kv_it_key_compare(uint32_t itr, const char* key, uint32_t size);
int32_t kv_it_move_to_end(uint32_t itr);
//int32_t kv_it_next(uint32_t itr, uint32_t& found_key_size, uint32_t& found_value_size);
//int32_t kv_it_prev(uint32_t itr, uint32_t& found_key_size, uint32_t& found_value_size);
int32_t kv_it_next(uint32_t itr, uint32_t* found_key_size, uint32_t* found_value_size);
int32_t kv_it_prev(uint32_t itr, uint32_t* found_key_size, uint32_t* found_value_size);

//int32_t kv_it_lower_bound(uint32_t itr, const char* key, uint32_t size, uint32_t& found_key_size, uint32_t& found_value_size);
//int32_t kv_it_key(uint32_t itr, uint32_t offset, char* dest, uint32_t size, uint32_t& actual_size);
//int32_t kv_it_value(uint32_t itr, uint32_t offset, char* dest, uint32_t size, uint32_t& actual_size);

int32_t kv_it_lower_bound(uint32_t itr, const char* key, uint32_t size, uint32_t* found_key_size, uint32_t* found_value_size);
int32_t kv_it_key(uint32_t itr, uint32_t offset, char* dest, uint32_t size, uint32_t* actual_size);
int32_t kv_it_value(uint32_t itr, uint32_t offset, char* dest, uint32_t size, uint32_t* actual_size);
*/
import "C"

import (
	"unsafe"

	"github.com/uuosio/chain"
)

type KV struct {
	Contract chain.Name
}

type KVIterator struct {
	Prefix        []byte
	Handle        uint32
	CurrentStatus int32
}

type KVIteratorStatus struct {
	Status int32
}

const (
	ITER_OK     = 0  // Iterator is positioned at a key-value pair
	ITER_ERASED = -1 // The key-value pair that the iterator used to be positioned at was erased
	ITER_END    = -2 // Iterator is out-of-bounds
)

func (t *KVIteratorStatus) IsOk() bool {
	return t.Status == ITER_OK
}

func (t *KVIteratorStatus) IsErased() bool {
	return t.Status == ITER_ERASED
}

func (t *KVIteratorStatus) IsEnd() bool {
	return t.Status == ITER_END
}

func (t *KVIteratorStatus) Print() {
	chain.Print(t.Status)
}

func NewKV(contract chain.Name) *KV {
	return &KV{contract}
}

// int64_t kv_erase(uint64_t contract, const char* key, uint32_t key_size);
func (t *KV) Erase(key []byte) int64 {
	ret := C.kv_erase(C.uint64_t(t.Contract.N), (*C.char)(unsafe.Pointer(&key[0])), C.uint32_t(len(key)))
	return ret
}

// int64_t kv_set(uint64_t contract, const char* key, uint32_t key_size, const char* value, uint32_t value_size, uint64_t payer);
func (t *KV) Set(key []byte, value []byte, payer chain.Name) int64 {
	ret := C.kv_set(C.uint64_t(t.Contract.N), (*C.char)(unsafe.Pointer(&key[0])), C.uint32_t(len(key)), (*C.char)(unsafe.Pointer(&value[0])), C.uint32_t(len(value)), C.uint64_t(payer.N))
	return ret
}

type stringHeader struct {
	data unsafe.Pointer
	len  uintptr
}

func GetStringPointer(s string) *C.char {
	_s := (*stringHeader)(unsafe.Pointer(&s))
	return (*C.char)(unsafe.Pointer(_s.data))
}

func GetBytesPointer(s []byte) *C.char {
	return (*C.char)(unsafe.Pointer(&s[0]))
}

// bool kv_get(uint64_t contract, const char* key, uint32_t key_size, uint32_t* value_size);
func (t *KV) Find(key []byte) (int, bool) {
	var size C.uint32_t
	ret := C.kv_get(C.uint64_t(t.Contract.N), GetBytesPointer(key), C.uint32_t(len(key)), &size)
	if ret == 0 {
		return 0, false
	}
	return int(size), true
}

func (t *KV) Get(key []byte) ([]byte, bool) {
	size, ok := t.Find(key)
	if !ok {
		return nil, false
	}
	value := make([]byte, size)
	C.kv_get_data(0, (*C.char)(unsafe.Pointer(&value[0])), C.uint32_t(len(value)))
	return value, true
}

// uint32_t kv_it_create(uint64_t contract, const char* prefix, uint32_t size);
func (t *KV) CreateItr(prefix ...[]byte) KVIterator {
	var _prefix []byte
	if len(prefix) == 0 {
		_prefix = nil
	} else {
		_prefix = prefix[0]
	}

	itr := C.kv_it_create(C.uint64_t(t.Contract.N), GetBytesPointer(_prefix), C.uint32_t(len(prefix)))
	return KVIterator{Handle: itr, Prefix: _prefix}
}

// void kv_it_destroy(uint32_t itr);
func (t *KVIterator) Destroy() {
	C.kv_it_destroy(t.Handle)
}

// int32_t kv_it_status(uint32_t itr);
func (t *KVIterator) Status() KVIteratorStatus {
	ret := C.kv_it_status(t.Handle)
	return KVIteratorStatus{ret}
}

// int32_t kv_it_compare(uint32_t itr_a, uint32_t itr_b);
func (t *KVIterator) Compare(itr_b KVIterator) KVIteratorStatus {
	ret := C.kv_it_compare(t.Handle, itr_b.Handle)
	return KVIteratorStatus{ret}
}

// int32_t kv_it_key_compare(uint32_t itr, const char* key, uint32_t size);
func (t *KVIterator) KeyCompare(key []byte) KVIteratorStatus {
	ret := C.kv_it_key_compare(t.Handle, GetBytesPointer(key), C.uint32_t(len(key)))
	return KVIteratorStatus{ret}
}

// int32_t kv_it_move_to_end(uint32_t itr);
func (t *KVIterator) MoveToEnd() KVIteratorStatus {
	ret := C.kv_it_move_to_end(t.Handle)
	return KVIteratorStatus{ret}
}

// int32_t kv_it_next(uint32_t itr, uint32_t* found_key_size, uint32_t* found_value_size);
func (t *KVIterator) Next() (int32, int32, KVIteratorStatus) {
	var key_size C.uint32_t
	var value_size C.uint32_t
	ret := C.kv_it_next(t.Handle, &key_size, &value_size)
	return int32(key_size), int32(value_size), KVIteratorStatus{ret}
}

// int32_t kv_it_prev(uint32_t itr, uint32_t* found_key_size, uint32_t* found_value_size);
func (t *KVIterator) Prev() (int32, int32, KVIteratorStatus) {
	var key_size C.uint32_t
	var value_size C.uint32_t
	ret := C.kv_it_prev(t.Handle, &key_size, &value_size)
	return int32(key_size), int32(value_size), KVIteratorStatus{ret}
}

// int32_t kv_it_lower_bound(uint32_t itr, const char* key, uint32_t size, uint32_t* found_key_size, uint32_t* found_value_size);
func (t *KVIterator) LowerBound(key []byte) (int32, int32, KVIteratorStatus) {
	var key_size C.uint32_t
	var value_size C.uint32_t
	ret := C.kv_it_lower_bound(t.Handle, GetBytesPointer(key), C.uint32_t(len(key)), &key_size, &value_size)
	return int32(key_size), int32(value_size), KVIteratorStatus{ret}
}

// int32_t kv_it_key(uint32_t itr, uint32_t offset, char* dest, uint32_t size, uint32_t* actual_size);
func (t *KVIterator) Key() ([]byte, bool) {
	var actual_size C.uint32_t
	ret := C.kv_it_key(t.Handle, C.uint32_t(0), (*C.char)(unsafe.Pointer(uintptr(0))), C.uint32_t(0), &actual_size)
	if ret != ITER_OK {
		return nil, false
	}

	dest := make([]byte, actual_size)
	C.kv_it_key(t.Handle, C.uint32_t(0), (*C.char)(unsafe.Pointer(&dest[0])), C.uint32_t(len(dest)), &actual_size)
	return dest, true
}

// int32_t kv_it_value(uint32_t itr, uint32_t offset, char* dest, uint32_t size, uint32_t* actual_size);
func (t *KVIterator) Value() ([]byte, bool) {
	var actual_size C.uint32_t
	ret := C.kv_it_value(t.Handle, C.uint32_t(0), (*C.char)(unsafe.Pointer(uintptr(0))), C.uint32_t(0), &actual_size)
	if ret != ITER_OK {
		return nil, false
	}

	dest := make([]byte, actual_size)
	C.kv_it_value(t.Handle, C.uint32_t(0), (*C.char)(unsafe.Pointer(&dest[0])), C.uint32_t(len(dest)), &actual_size)
	return dest, true
}

// lower_bound(k);
// if (itr_key_compare(handle, {k.data(), k.size()}) != 0)
//    seek_to_end();

func (t *KVIterator) Find(key []byte) bool {
	t.LowerBound(key)
	ret := t.KeyCompare(key)
	if !ret.IsOk() {
		t.MoveToEnd()
	}
	if !ret.IsOk() {
		return true
	} else {
		return false
	}
}
