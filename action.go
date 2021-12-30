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
	"unsafe"
)

var gActionCache [][]byte
var gNotifyCache []Name

//Read current action data
func ReadActionData() []byte {
	n := C.action_data_size()
	if n <= 0 {
		return nil
	}
	buf := make([]byte, int(n))
	ptr := GetBytesPtr(buf)
	C.read_action_data(ptr, n)
	return buf
}

//Get the length of the current action's data field
func ActionDataSize() uint32 {
	return uint32(C.action_data_size())
}

//Add the specified account to set of accounts to be notified
func RequireRecipient(name Name) {
	if IsRevertEnabled() {
		AddNotifyToCache(name)
	} else {
		C.require_recipient(C.uint64_t(name.N))
	}
}

func RequireRecipientEx(name Name) {
	C.require_recipient(C.uint64_t(name.N))
}

//Verifies that name exists in the set of provided auths on a action. Throws if not found.
func RequireAuth(name Name) {
	C.require_auth(C.uint64_t(name.N))
}

//Verifies that name has auth.
func HasAuth(name Name) bool {
	ret := C.has_auth(C.uint64_t(name.N))
	if ret == 0 {
		return false
	}
	return true
}

//Verifies that name exists in the set of provided auths on a action. Throws if not found.
func RequireAuth2(name Name, permission Name) {
	C.require_auth2(C.uint64_t(name.N), C.uint64_t(permission.N))
}

//Verifies that name is an existing account.
func IsAccount(name Name) bool {
	ret := C.is_account(C.uint64_t(name.N))
	if ret == 0 {
		return false
	}
	return true
}

//Send an inline action in the context of this action's parent transaction
func SendInline(data []byte) {
	//	a := (*sliceHeader)(unsafe.Pointer(&data))
	p := (*C.char)(unsafe.Pointer(&data[0]))
	C.send_inline(p, C.size_t(len(data)))
}

//Send an inline context free action in the context of this action's parent transaction
func SendContextFreeInline(data []byte) {
	a := (*SliceHeader)(unsafe.Pointer(&data))
	C.send_context_free_inline((*C.char)(unsafe.Pointer(a.Data)), C.size_t(a.Len))
}

//Returns the time in microseconds from 1970 of the publication_time
func PublicationTime() uint64 {
	return uint64(C.publication_time())
}

//Get the current receiver of the action
func CurrentReceiver() Name {
	n := C.current_receiver()
	return Name{uint64(n)}
}

//Set the action return value which will be included in the action_receipt
func SetActionReturnValue(return_value []byte) {
	a := (*SliceHeader)(unsafe.Pointer(&return_value))
	C.set_action_return_value((*C.char)(unsafe.Pointer(a.Data)), C.size_t(a.Len))
}

type Action struct {
	Account       Name
	Name          Name
	Authorization []PermissionLevel
	Data          []byte
}

func NewAction(perm PermissionLevel, account Name, name Name, args ...interface{}) *Action {
	a := &Action{}
	a.Account = account
	a.Name = name

	a.Authorization = append(a.Authorization, perm)
	if len(args) == 0 {
		a.Data = []byte{}
		return a
	}

	size := 0
	for _, v := range args {
		n := CalcPackedSize(v)
		size += n
	}
	enc := NewEncoder(size)
	for _, arg := range args {
		enc.Pack(arg)
	}
	a.Data = enc.GetBytes()

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

func (a *Action) Size() int {
	return 8 + 8 + 5 + len(a.Authorization)*8 + 5 + len(a.Data)
}

func (a *Action) Pack() []byte {
	enc := NewEncoder(a.Size())
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

func (a *Action) Unpack(b []byte) int {
	dec := NewDecoder(b)
	a.Account = dec.UnpackName()
	a.Name = dec.UnpackName()
	length := dec.UnpackLength()
	a.Authorization = make([]PermissionLevel, length)
	for i := 0; i < length; i++ {
		dec.Unpack(&a.Authorization[i])
	}
	dec.Unpack(&a.Data)
	return dec.Pos()
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

func AddNotifyToCache(a Name) {
	if !IsRevertEnabled() {
		return
	}
	if gNotifyCache == nil {
		gNotifyCache = make([]Name, 0, 5)
	}
	gNotifyCache = append(gNotifyCache, a)
}

func AddActionToCache(a []byte) {
	if !IsRevertEnabled() {
		return
	}
	if gActionCache == nil {
		gActionCache = make([][]byte, 0, 5)
	}
	gActionCache = append(gActionCache, a)
}

func (a *Action) Send() {
	data := a.Pack()
	if IsRevertEnabled() {
		AddActionToCache(data)
	} else {
		SendInline(data)
	}
}

//send action directly, no cache
func (a *Action) SendEx() {
	data := a.Pack()
	SendInline(data)
}

func SendAllActions() {
	for _, a := range gActionCache {
		SendInline(a)
	}
}

func Finish() {
	if !IsRevertEnabled() {
		return
	}

	SendCachedTransactions()

	for _, notify := range gNotifyCache {
		C.require_recipient(C.uint64_t(notify.N))
	}
	SendAllActions()
}
