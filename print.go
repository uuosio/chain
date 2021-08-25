package chain

import (
	"unsafe"
)

/*
#include <stddef.h>
#include <stdint.h>
typedef uint8_t uint128_t;

void prints( const char* cstr );
void prints_l( const char* cstr, uint32_t len);
void printi( int64_t value );
void printui( uint64_t value );
//void printi128( const int128_t* value );
void printi128( const uint8_t* value );

//void printui128( const uint128_t* value );
void printui128( const uint8_t* value );
void printsf(float value);
void printdf(double value);
//void printqf(const long double* value);
void printqf(const uint8_t* value);
void printn( uint64_t name );
void printhex( const void* data, uint32_t datalen );
*/
import "C"

type stringHeader struct {
	data unsafe.Pointer
	len  uintptr
}

// void prints_l( const char* cstr, uint32_t len);
func Prints(str string) {
	_str := (*stringHeader)(unsafe.Pointer(&str))
	C.prints_l((*C.char)(_str.data), C.uint32_t(len(str)))
}

// void printi( int64_t value );
func Printi(value int64) {
	C.printi(C.int64_t(value))
}

// void printui( uint64_t value );
func PrintUi(value uint64) {
	C.printui(C.uint64_t(value))
}

// void printi128( const int128_t* value );
func PrintI128(value [16]byte) {
	C.printi128((*C.uint8_t)(&value[0]))
}

// void printui128( const uint128_t* value );
func PrintUi128(value [16]byte) {
	C.printui128((*C.uint8_t)(&value[0]))
}

// void printsf(float value);
func PrintSf(value float32) {
	C.printsf(value)
}

// void printdf(double value);
func PrintDf(value float64) {
	C.printdf(value)
}

// void printqf(const long double* value);
func PrintQf(value [16]byte) {
	C.printqf((*C.uint8_t)(&value[0]))
}

// void printn( uint64_t name );
func PrintN(name uint64) {
	C.printn(C.uint64_t(name))
}

// void printhex( const void* data, uint32_t datalen );
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

func PrintVariant(v interface{}) {
	switch v := v.(type) {
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
