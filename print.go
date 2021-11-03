package chain

import (
	"unsafe"
)

/*
#include <stddef.h>
#include <stdint.h>
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
func PrintQf(value Float128) {
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

type Printable interface {
	Print()
}

type ExtendedPrintFunc func(value interface{})

var (
	ExtendedPrint ExtendedPrintFunc
)

func PrintVariant(variant interface{}) {
	switch v := variant.(type) {
	case nil:
		Prints("nil")
	case Printable:
		v.Print()
	case error:
		Prints(v.Error())
	case string:
		Prints(v)
	case bool:
		if v {
			Prints("true")
		} else {
			Prints("false")
		}
	case int8:
		Printi(int64(v))
	case uint8:
		PrintUi(uint64(v))
	case int16:
		Printi(int64(v))
	case uint16:
		PrintUi(uint64(v))
	case int32:
		Printi(int64(v))
	case uint32:
		PrintUi(uint64(v))
	case int:
		Printi(int64(v))
	case int64:
		Printi(v)
	case uint64:
		PrintUi(v)
	case [16]byte:
		PrintI128(v)
	// case [16]byte:
	// 	PrintUi128(v)
	case float32:
		PrintSf(v)
	case float64:
		PrintDf(v)
	case Float128:
		PrintQf(v)
	// case [16]byte:
	// 	PrintQf(v)
	case Name:
		PrintN(v.N)
	case []Name:
		for _, n := range v {
			PrintN(n.N)
		}
	case Symbol:
		v.Print()
	case []byte:
		PrintHex(v)
	default:
		// if DEBUG {
		// 	s := fmt.Sprintf("%v", v)
		// 	Prints(s)
		// }
		Prints("<unprintable>")
	}
}

func PrintNoEndSpace(args ...interface{}) {
	for _, arg := range args {
		PrintVariant(arg)
	}
}

func Print(args ...interface{}) {
	for i, v := range args {
		PrintVariant(v)
		if i < len(args)-1 {
			PrintVariant(" ")
		}
	}
}

func Println(args ...interface{}) {
	Print(args...)
	Print("\n")
}
