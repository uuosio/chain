package chain

/*
#include <stddef.h>
#include <stdint.h>
#define bool char

uint32_t read_action_data( void* msg, uint32_t len );

uint32_t action_data_size( void );

void require_recipient( uint64_t name );

void require_auth( uint64_t name );

char has_auth( uint64_t name );

void require_auth2( uint64_t name, uint64_t permission );

char is_account( uint64_t name );
void send_inline(char *serialized_action, size_t size);

void send_context_free_inline(char *serialized_action, size_t size);
uint64_t  publication_time( void );

uint64_t current_receiver( void );
void set_action_return_value(char *return_value, size_t size);
*/
import "C"
import (
	"runtime"
	"unsafe"
)

func ReadActionData() []byte {
	n := C.action_data_size()
	buf := runtime.Alloc(uintptr(n))
	C.read_action_data(buf, uint32(n))
	pp := (*[1 << 30]byte)(buf)
	return pp[:n]
}

func ActionDataSize() uint32 {
	return C.action_data_size()
}

func RequireRecipient(name Name) {
	C.require_recipient(name.N)
}

func RequireAuth(name Name) {
	C.require_auth(name.N)
}

func HasAuth(name Name) bool {
	ret := C.has_auth(name.N)
	if ret == 0 {
		return false
	}
	return true
}

func RequireAuth2(name Name, permission Name) {
	C.require_auth2(name.N, permission.N)
}

func IsAccount(name Name) bool {
	ret := C.is_account(name.N)
	if ret == 0 {
		return false
	}
	return true
}

func SendInline(data []byte) {
	//	a := (*sliceHeader)(unsafe.Pointer(&data))
	p := (*C.char)(unsafe.Pointer(&data[0]))
	C.send_inline(p, C.size_t(len(data)))
}

func SendContextFreeInline(data []byte) {
	a := (*SliceHeader)(unsafe.Pointer(&data))
	C.send_context_free_inline((*C.char)(unsafe.Pointer(a.Data)), C.size_t(a.Len))
}

func PublicationTime() uint64 {
	return C.publication_time()
}

func CurrentReceiver() Name {
	n := C.current_receiver()
	return Name{n}
}

func SetActionReturnValue(return_value []byte) {
	a := (*SliceHeader)(unsafe.Pointer(&return_value))
	C.set_action_return_value((*C.char)(unsafe.Pointer(a.Data)), C.size_t(a.Len))
}

// type PermissionLevel struct {
// 	Actor      Name
// 	Permission Name
// }

// func (a *PermissionLevel) Pack() []byte {
// 	enc := NewEncoder(16)
// 	enc.PackName(a.Actor)
// 	enc.PackName(a.Permission)
// 	return enc.GetBytes()
// }

// func (a *PermissionLevel) Unpack(data []byte) (int, error) {
// 	dec := NewDecoder(data)
// 	dec.Unpack(&a.Actor)
// 	dec.Unpack(&a.Permission)
// 	return 16, nil
// }

type Action struct {
	Account       Name
	Name          Name
	Authorization []PermissionLevel
	Data          []byte
}

func NewAction(account Name, name Name) *Action {
	a := &Action{}
	a.Account = account
	a.Name = name
	return a
}

func PackUint64(n uint64) []byte {
	p := [8]byte{}
	pp := (*[8]byte)(unsafe.Pointer(&n))
	copy(p[:], pp[:])
	return p[:]
}

func PackArray(a []Serializer) []byte {
	buf := []byte{byte(len(a))}
	for _, v := range a {
		buf = append(buf, v.Pack()...)
	}
	return buf
}

func (a *Action) EstimatePackedSize() int {
	return 8 + 8 + 5 + len(a.Authorization)*8 + 5 + len(a.Data)
}

func (a *Action) Pack() []byte {
	enc := NewEncoder(8 + 8 + 5 + len(a.Authorization)*8 + 5 + len(a.Data))
	enc.PackName(a.Account)
	enc.PackName(a.Name)
	enc.PackLength(len(a.Authorization))
	for _, v := range a.Authorization {
		enc.Pack(&v)
	}
	enc.Pack(a.Data)
	return enc.GetBytes()
	// buf := []byte{}
	// buf = append(buf, PackUint64(a.Account)...)
	// buf = append(buf, PackUint64(a.Name)...)

	// buf = append(buf, PackUint32(uint32(len(a.Authorization)))...)
	// for _, v := range a.Authorization {
	// 	buf = append(buf, v.Pack()...)
	// }

	// buf = append(buf, a.Data.Pack()...)
	// return buf
}

func (a *Action) Unpack(b []byte) (int, error) {
	dec := NewDecoder(b)
	dec.Unpack(&a.Account)
	dec.Unpack(&a.Name)
	length, err := dec.UnpackLength()
	if err != nil {
		return 0, err
	}
	a.Authorization = make([]PermissionLevel, length)
	for i := 0; i < length; i++ {
		dec.Unpack(&a.Authorization[i])
	}
	dec.Unpack(&a.Data)
	return dec.Pos(), nil
}

func (a *Action) Print() {
	Print("{")
	Print(a.Account, a.Name)
	Print("[")
	for _, v := range a.Authorization {
		Print("[", v.Actor, v.Permission, "]")
	}
	Print("]")
	Print(a.Data)
	Print("}")
}

func (a *Action) AddPermission(actor Name, permission Name) {
	a.Authorization = append(a.Authorization, PermissionLevel{actor, permission})
}

func (a *Action) Send() {
	data := a.Pack()
	SendInline(data)
}
