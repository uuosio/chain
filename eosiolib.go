package chain

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
