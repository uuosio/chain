package chain

/*
#include <stddef.h>
#include <stdint.h>
void printn(uint64_t a);
void printui( uint64_t value );
void prints( const char* cstr );
uint32_t action_data_size(void);
uint32_t read_action_data( void* msg, uint32_t len );
void send_inline(char *serialized_action, size_t size);
uint32_t get_active_producers( uint8_t* producers, uint32_t datalen );
*/
import "C"
import (
	"runtime"
	"unsafe"
)

var DEBUG = true

// var DEBUG = false

func char_to_symbol(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return (c - 'a') + 6
	}

	if c >= '1' && c <= '5' {
		return (c - '1') + 1
	}
	return 0
}

func string_to_name(str string) uint64 {
	length := len(str)
	value := uint64(0)

	for i := 0; i <= 12; i++ {
		c := uint64(0)
		if i < length && i <= 12 {
			c = uint64(char_to_symbol(str[i]))
		}
		if i < 12 {
			c &= 0x1f
			c <<= 64 - 5*(i+1)
		} else {
			c &= 0x0f
		}

		value |= c
	}

	return value
}

func S2N(s string) uint64 {
	return string_to_name(s)
}

func N2S(value uint64) string {
	charmap := []byte(".12345abcdefghijklmnopqrstuvwxyz")
	// 13 dots
	str := []byte{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.'}

	tmp := value
	for i := 0; i <= 12; i++ {
		var c byte
		if i == 0 {
			c = charmap[tmp&0x0f]
		} else {
			c = charmap[tmp&0x1f]
		}
		str[12-i] = c
		if i == 0 {
			tmp >>= 4
		} else {
			tmp >>= 5
		}
	}

	i := len(str) - 1
	for ; i >= 0; i-- {
		if str[i] != '.' {
			break
		}
	}
	return string(str[:i+1])
}

func N(s string) uint64 {
	return string_to_name(s)
}

func PackUint32(val uint32) []byte {
	result := make([]byte, 0, 5)
	for {
		b := byte(val & 0x7f)
		val >>= 7
		if val > 0 {
			b |= byte(1 << 7)
		}
		result = append(result, b)
		if val <= 0 {
			break
		}
	}
	return result
}

func UnpackUint32(val []byte) (n int, v uint32) {
	var by int = 0
	// if len(val) > 5 {
	// 	val = val[:5]
	// }
	n = 0
	for _, b := range val {
		v |= uint32(b&0x7f) << by
		by += 7
		n += 1
		if b&0x80 == 0 {
			break
		}
	}
	return
}

func PackedSizeLength(val uint32) int {
	n := 0
	for {
		b := byte(val & 0x7f)
		val >>= 7
		if val > 0 {
			b |= byte(1 << 7)
		}
		n += 1
		if val <= 0 {
			break
		}
	}
	return n
}

type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

type StringHeader struct {
	data unsafe.Pointer
	len  uintptr
}

// var gPrintBuffer = unsafe.Pointer(uintptr(0))
// var gPrintBufferSize = 513

// func Prints(s string) {
// 	if len(s)+1 > gPrintBufferSize {
// 		gPrintBufferSize = len(s) + 1
// 		if gPrintBuffer != unsafe.Pointer(uintptr(0)) {
// 			runtime.Free(gPrintBuffer)
// 			gPrintBuffer = unsafe.Pointer(uintptr(0))
// 		}
// 	}

// 	if gPrintBuffer == unsafe.Pointer(uintptr(0)) {
// 		gPrintBuffer = runtime.Alloc(uintptr(gPrintBufferSize))
// 	}

// 	pp := (*[1 << 30]byte)(gPrintBuffer)
// 	copy(pp[:], s)
// 	pp[len(s)] = 0
// 	// C.printui(uint64(uintptr(unsafe.Pointer(&tmp))))
// 	C.prints((*C.char)(gPrintBuffer))
// }

func Printui(n uint64) {
	C.printui(n)
}

type Serializer interface {
	Pack() []byte
	Unpack([]byte) (int, uint64)
}

type T struct {
	a int
	b int
}

func SayHello() {
	Prints("hello,world!")
}

func GetApplyArgs() (Name, Name, Name) {
	receiver, code, action := runtime.GetApplyArgs()
	return Name{receiver}, Name{code}, Name{action}
}

func Check(b bool, msg string) {
	Assert(b, msg)
}
