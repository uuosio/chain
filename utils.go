package chain

import "unsafe"

type PTR unsafe.Pointer

func GetStringPtr(str string) unsafe.Pointer {
	_str := (*stringHeader)(unsafe.Pointer(&str))
	return _str.data
}

func GetBytesPtr(bs []byte) unsafe.Pointer {
	return unsafe.Pointer(&bs[0])
}
