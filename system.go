package chain

/*
#include <stdint.h>

void  eosio_assert( uint32_t test, const char* msg );
void  eosio_assert_message( uint32_t test, const char* msg, uint32_t msg_len );
void  eosio_assert_code( uint32_t test, uint64_t code );

void eosio_exit( int32_t code );
uint64_t  current_time( void );
char is_feature_activated( const uint8_t* feature_digest ); //checksum 32 bytes
uint64_t get_sender( void );
*/
import "C"
import (
	"unsafe"
)

// func EosioAssert(test bool, msg string) {
// 	buf := runtime.Alloc(uintptr(len(msg) + 1))
// 	pp := (*[1 << 30]byte)(buf)
// 	copy(pp[:], msg)
// 	pp[len(msg)] = 0
// 	if !test {
// 		C.eosio_assert(0, (*C.char)(buf))
// 	}
// }

func Assert(test bool, msg string) {
	EosioAssert(test, msg)
}

// void  eosio_assert_message( uint32_t test, const char* msg, uint32_t msg_len );
func EosioAssert(test bool, msg string) {
	_test := uint32(0)
	if test {
		_test = 1
	}
	_msg := (*StringHeader)(unsafe.Pointer(&msg))
	C.eosio_assert_message(_test, (*C.char)(_msg.data), C.uint32_t(len(msg)))
}

// void  eosio_assert_code( uint32_t test, uint64_t code );
func EosioAssertCode(test bool, code uint64) {
	_test := uint32(0)
	if test {
		_test = 1
	}
	C.eosio_assert_code(_test, C.uint64_t(code))
}

// void eosio_exit( int32_t code );
// uint64_t  current_time( void );
func CurrentTime() uint64 {
	return uint64(C.current_time())
}

func CurrentTimeSeconds() uint32 {
	t := CurrentTime() / 1000000
	return uint32(t)
}

// char is_feature_activated( const uint8_t* feature_digest ); //checksum 32 bytes
func IsFeatureActivated(featureDigest [32]byte) bool {
	_featureDigest := (*C.uint8_t)(unsafe.Pointer(&featureDigest[0]))
	return C.is_feature_activated(_featureDigest) != 0
}

// uint64_t get_sender( void );
func GetSender() uint64 {
	return uint64(C.get_sender())
}
