package chain

import "unsafe"

type PTR unsafe.Pointer

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
