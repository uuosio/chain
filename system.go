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

func Check(b bool, msg string) {
	EosioAssert(b, msg)
}

//Aborts processing of this action and unwinds all pending changes if the test condition is true
func Assert(test bool, msg string) {
	EosioAssert(test, msg)
}

//Aborts processing of this action and unwinds all pending changes if the test condition is true
func EosioAssert(test bool, msg string) {
	_test := uint32(0)
	if test {
		_test = 1
	}
	_msg := (*StringHeader)(unsafe.Pointer(&msg))
	C.eosio_assert_message(C.uint32_t(_test), (*C.char)(_msg.data), C.uint32_t(len(msg)))
}

//Aborts processing of this action and unwinds all pending changes if the test condition is true
func EosioAssertCode(test bool, code uint64) {
	_test := uint32(0)
	if test {
		_test = 1
	}
	C.eosio_assert_code(C.uint32_t(_test), C.uint64_t(code))
}

//Returns the time in microseconds from 1970 of the current block
func CurrentTime() TimePoint {
	return TimePoint{uint64(C.current_time())}
}

//Returns the time in microseconds from 1970 of the current block
func Now() TimePoint {
	return CurrentTime()
}

func NowSeconds() uint32 {
	t := CurrentTime().Elapsed / 1000000
	return uint32(t)
}

func CurrentTimeSeconds() uint32 {
	t := CurrentTime().Elapsed / 1000000
	return uint32(t)
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
