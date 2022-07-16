//go:build tinygo.wasm
// +build tinygo.wasm

package eosio

/*
#include <stddef.h>
#include <stdint.h>

#define bool char
typedef float Float;
typedef double Double;

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

//Prints string
func Prints(str string) {
	C.prints_l((*C.char)(GetStringPtr(str)), C.uint32_t(len(str)))
}

//Prints value as a 64 bit signed integer
func Printi(value int64) {
	C.printi(C.int64_t(value))
}

//Prints value as a 64 bit unsigned integer
func PrintUi(value uint64) {
	C.printui(C.uint64_t(value))
}

func Printui(value uint64) {
	C.printui(C.uint64_t(value))
}

//Prints value as a 128 bit signed integer
func PrintI128(value [16]byte) {
	C.printi128((*C.uint8_t)(&value[0]))
}

//Prints value as a 128 bit unsigned integer
func PrintUi128(value [16]byte) {
	C.printui128((*C.uint8_t)(&value[0]))
}

//Prints value as single-precision floating point number
func PrintSf(value float32) {
	C.printsf(C.Float(value))
}

//Prints value as double-precision floating point number
func PrintDf(value float64) {
	C.printdf(C.Double(value))
}

//Prints value as quadruple-precision floating point number
func PrintQf(value [16]byte) {
	C.printqf((*C.uint8_t)(&value[0]))
}

//Prints a 64 bit names as base32 encoded string
func PrintN(name uint64) {
	C.printn(C.uint64_t(name))
}

//Prints hexidecimal data of length datalen
func PrintHex(data []byte) {
	C.printhex(unsafe.Pointer(&data[0]), C.uint32_t(len(data)))
}
