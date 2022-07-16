//go:build tinygo.wasm
// +build tinygo.wasm

package eosio

/*
#include <stddef.h>
#include <stdint.h>
#define bool char
typedef float Float;
typedef double Double;

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


//
void prints( const char* cstr );
void prints_l( const char* cstr, uint32_t len);
void printi( int64_t value );
void printui( uint64_t value );
//void printi128( const int128_t* value );
void printi128( const uint8_t* value );

//void printui128( const uint128* value );
void printui128( const uint8_t* value );
void printsf(Float value);
void printdf(Double value);
//void printqf(const long double* value);
void printqf(const uint8_t* value);
void printn( uint64_t name );
void printhex( const void* data, uint32_t datalen );
*/
import "C"
import "unsafe"

type PTR unsafe.Pointer

type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

type sliceHeader struct {
	data unsafe.Pointer
	len  uintptr
	cap  uintptr
}

type stringHeader struct {
	data unsafe.Pointer
	len  uintptr
}

func GetStringPtr(str string) unsafe.Pointer {
	if len(str) != 0 {
		_str := (*stringHeader)(unsafe.Pointer(&str))
		return _str.data
	}
	return unsafe.Pointer(uintptr(0))
}

func GetBytesPtr(bs []byte) unsafe.Pointer {
	if len(bs) != 0 {
		return unsafe.Pointer(&bs[0])
	}
	return unsafe.Pointer(uintptr(0))
}

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
func RequireRecipient(name uint64) {
	C.require_recipient(C.uint64_t(name))
}

func RequireRecipientEx(name uint64) {
	C.require_recipient(C.uint64_t(name))
}

//Verifies that name exists in the set of provided auths on a action. Throws if not found.
func RequireAuth(name uint64) {
	C.require_auth(C.uint64_t(name))
}

//Verifies that name has auth.
func HasAuth(name uint64) bool {
	ret := C.has_auth(C.uint64_t(name))
	if ret == 0 {
		return false
	}
	return true
}

//Verifies that name exists in the set of provided auths on a action. Throws if not found.
func RequireAuth2(name uint64, permission uint64) {
	C.require_auth2(C.uint64_t(name), C.uint64_t(permission))
}

//Verifies that name is an existing account.
func IsAccount(name uint64) bool {
	ret := C.is_account(C.uint64_t(name))
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
func CurrentReceiver() uint64 {
	n := C.current_receiver()
	return uint64(n)
}

//Set the action return value which will be included in the action_receipt
func SetActionReturnValue(return_value []byte) {
	a := (*SliceHeader)(unsafe.Pointer(&return_value))
	C.set_action_return_value((*C.char)(unsafe.Pointer(a.Data)), C.size_t(a.Len))
}
