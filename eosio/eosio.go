//go:build tinygo.wasm
// +build tinygo.wasm

package eosio

/*
#include <stddef.h>
#include <stdint.h>
#define bool char
typedef float Float;
typedef double Double;

#include "../chain.h"

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


//system.h
void  eosio_assert( uint32_t test, const char* msg );
void  eosio_assert_message( uint32_t test, const char* msg, uint32_t msg_len );
void  eosio_assert_code( uint32_t test, uint64_t code );

void eosio_exit( int32_t code );
uint64_t  current_time( void );
char is_feature_activated( const uint8_t* feature_digest ); //checksum 32 bytes
uint64_t get_sender( void );


//crypto.h
void assert_sha256( const char* data, uint32_t length, const capi_checksum256* hash );
void assert_sha1( const char* data, uint32_t length, const capi_checksum160* hash );
void assert_sha512( const char* data, uint32_t length, const capi_checksum512* hash );
void assert_ripemd160( const char* data, uint32_t length, const capi_checksum160* hash );
void sha256( const char* data, uint32_t length, capi_checksum256* hash );
void sha1( const char* data, uint32_t length, capi_checksum160* hash );
void sha512( const char* data, uint32_t length, capi_checksum512* hash );
void ripemd160( const char* data, uint32_t length, capi_checksum160* hash );
int recover_key( const capi_checksum256* digest, const char* sig, size_t siglen, char* pub, size_t publen );
void assert_recover_key( const capi_checksum256* digest, const char* sig, size_t siglen, const char* pub, size_t publen );


int32_t db_store_i64(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id,  const char* data, uint32_t len);
void db_update_i64(int32_t iterator, uint64_t payer, const char* data, uint32_t len);
void db_remove_i64(int32_t iterator);
int32_t db_get_i64(int32_t iterator, const char* data, uint32_t len);
int32_t db_next_i64(int32_t iterator, uint64_t* primary);
int32_t db_previous_i64(int32_t iterator, uint64_t* primary);
int32_t db_find_i64(uint64_t code, uint64_t scope, uint64_t table, uint64_t id);
int32_t db_lowerbound_i64(uint64_t code, uint64_t scope, uint64_t table, uint64_t id);
int32_t db_upperbound_i64(uint64_t code, uint64_t scope, uint64_t table, uint64_t id);
int32_t db_end_i64(uint64_t code, uint64_t scope, uint64_t table);
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

//system.h
func Check(b bool, msg string) {
	if !b {
		EosioAssert(false, msg)
	}
}

//Aborts processing of this action and unwinds all pending changes if the test condition is true
func Assert(test bool, msg string) {
	if !test {
		EosioAssert(false, msg)
	}
}

type StringHeader struct {
	data unsafe.Pointer
	len  uintptr
}

//Aborts processing of this action and unwinds all pending changes if the test condition is true
func EosioAssert(test bool, msg string) {
	if !test {
		_msg := (*StringHeader)(unsafe.Pointer(&msg))
		C.eosio_assert_message(C.uint32_t(0), (*C.char)(_msg.data), C.uint32_t(len(msg)))
	}
}

//Aborts processing of this action and unwinds all pending changes if the test condition is true
func EosioAssertCode(test bool, code uint64) {
	if !test {
		C.eosio_assert_code(C.uint32_t(0), C.uint64_t(code))
	}
}

//Returns the time in microseconds from 1970 of the current block
func CurrentTime() uint64 {
	return uint64(C.current_time())
}

//Check if specified protocol feature has been activated
func IsFeatureActivated(featureDigest [32]byte) bool {
	_featureDigest := (*C.uint8_t)(unsafe.Pointer(&featureDigest[0]))
	return C.is_feature_activated(_featureDigest) != 0
}

//Return name of account that sent current inline action
func GetSender() uint64 {
	return uint64(C.get_sender())
}

func Exit() {
	C.eosio_exit(0)
}

//crypto.h
//Tests if the sha256 hash generated from data matches the provided checksum.
func AssertSha256(data []byte, hash [32]byte) {
	C.assert_sha256((*C.char)(unsafe.Pointer(&data[0])), C.uint32_t(len(data)), (*C.capi_checksum256)(unsafe.Pointer(&hash)))
}

//Tests if the sha1 hash generated from data matches the provided checksum.
func AssertSha1(data []byte, hash [20]byte) {
	C.assert_sha1((*C.char)(unsafe.Pointer(&data[0])), C.uint32_t(len(data)), (*C.capi_checksum160)(unsafe.Pointer(&hash)))
}

//Tests if the sha512 hash generated from data matches the provided checksum.
func AssertSha512(data []byte, hash [64]byte) {
	C.assert_sha512((*C.char)(unsafe.Pointer(&data[0])), C.uint32_t(len(data)), (*C.capi_checksum512)(unsafe.Pointer(&hash)))
}

//Tests if the ripemod160 hash generated from data matches the provided checksum.
func AssertRipemd160(data []byte, hash [20]byte) {
	C.assert_ripemd160((*C.char)(unsafe.Pointer(&data[0])), C.uint32_t(len(data)), (*C.capi_checksum160)(unsafe.Pointer(&hash)))
}

//Hashes data using sha256 and return hash value.
func Sha256(data []byte) [32]byte {
	var hash [32]byte
	C.sha256((*C.char)(unsafe.Pointer(&data[0])), C.uint32_t(len(data)), (*C.capi_checksum256)(unsafe.Pointer(&hash)))
	return hash
}

//Hashes data using sha1 and return hash value.
func Sha1(data []byte) [20]byte {
	var hash [20]byte
	C.sha1((*C.char)(unsafe.Pointer(&data[0])), C.uint32_t(len(data)), (*C.capi_checksum160)(unsafe.Pointer(&hash)))
	return hash
}

//Hashes data using sha512 and return hash value.
func Sha512(data []byte) [64]byte {
	var hash [64]byte
	C.sha512((*C.char)(unsafe.Pointer(&data[0])), C.uint32_t(len(data)), (*C.capi_checksum512)(unsafe.Pointer(&hash)))
	return hash
}

//Hashes data using ripemd160 and return hash value.
func Ripemd160(data []byte) [20]byte {
	var hash [20]byte
	C.ripemd160((*C.char)(unsafe.Pointer(&data[0])), C.uint32_t(len(data)), (*C.capi_checksum160)(unsafe.Pointer(&hash)))
	return hash
}

//Recover the public key from digest and signature
func RecoverKey(digest [32]byte, sig []byte) []byte {
	//TODO: handle webauth signature
	var pub [34]byte //34
	C.recover_key((*C.capi_checksum256)(unsafe.Pointer(&digest)), (*C.char)(unsafe.Pointer(&sig[0])), C.size_t(len(sig)), (*C.char)(unsafe.Pointer(&pub[0])), C.size_t(len(pub)))
	return pub[:]
}

//Tests a given public key with the generated key from digest and the signature
func AssertRecoverKey(digest [32]byte, sig []byte, pub []byte) {
	C.assert_recover_key((*C.capi_checksum256)(unsafe.Pointer(&digest)), (*C.char)(unsafe.Pointer(&sig[0])), C.size_t(len(sig)), (*C.char)(unsafe.Pointer(&pub[0])), C.size_t(len(pub)))
}

// int32_t db_store_i64(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id,  const char* data, uint32_t len);
func DBStoreI64(scope uint64, table uint64, payer uint64, id uint64, data []byte) int32 {
	return C.db_store_i64(scope, table, payer, id, (*C.char)(unsafe.Pointer(&data[0])), C.uint32_t(len(data)))
}

// void db_update_i64(int32_t iterator, uint64_t payer, const char* data, uint32_t len);
func DBUpdateI64(iterator int32, payer uint64, data []byte) {
	p := (*C.char)(unsafe.Pointer(&data[0]))
	C.db_update_i64(iterator, payer, p, C.uint32_t(len(data)))
}

// void db_remove_i64(int32_t iterator);
func DBRemoveI64(iterator int32) {
	C.db_remove_i64(iterator)
}

// int32_t db_get_i64(int32_t iterator, const char* data, uint32_t len);
func DBGetI64(iterator int32) []byte {
	size := C.db_get_i64(iterator, (*C.char)(unsafe.Pointer(uintptr(0))), 0)
	data := make([]byte, size)
	C.db_get_i64(iterator, (*C.char)(unsafe.Pointer(&data[0])), uint32(size))
	return data
}

// int32_t db_next_i64(int32_t iterator, uint64_t* primary);
func DBNextI64(iterator int32) (int32, uint64) {
	var primary uint64
	next := C.db_next_i64(iterator, &primary)
	return next, primary
}

// int32_t db_previous_i64(int32_t iterator, uint64_t* primary);
func DBPreviousI64(iterator int32) (int32, uint64) {
	var primary uint64
	next := C.db_previous_i64(iterator, &primary)
	return next, primary
}

// int32_t db_find_i64(uint64_t code, uint64_t scope, uint64_t table, uint64_t id);
func DBFindI64(code uint64, scope uint64, table uint64, id uint64) int32 {
	return C.db_find_i64(code, scope, table, id)
}

// int32_t db_lowerbound_i64(uint64_t code, uint64_t scope, uint64_t table, uint64_t id);
func DBLowerBoundI64(code uint64, scope uint64, table uint64, id uint64) int32 {
	return C.db_lowerbound_i64(code, scope, table, id)
}

// int32_t db_upperbound_i64(uint64_t code, uint64_t scope, uint64_t table, uint64_t id);
func DBUpperBoundI64(code uint64, scope uint64, table uint64, id uint64) int32 {
	return C.db_upperbound_i64(code, scope, table, id)
}

// int32_t db_end_i64(uint64_t code, uint64_t scope, uint64_t table);
func DBEndI64(code uint64, scope uint64, table uint64) int32 {
	return C.db_end_i64(code, scope, table)
}
