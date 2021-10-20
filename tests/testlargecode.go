package main

import (
	"strconv"

	"github.com/uuosio/chain"
)

/*
#include <stdlib.h>
long      strtol( const char          *str, char          **str_end, int base );
double atof(const char *nptr);
*/
import "C"

//serializer
// type stringHeader struct {
// 	data unsafe.Pointer
// 	len  uintptr
// }

func Str2Long(str string) int64 {
	var str_end *byte
	return int64(C.strtol((*C.char)(chain.GetStringPtr(str)), (**C.char)(chain.PTR(&str_end)), 10))
}

func Str2Float(str string) float64 {
	return C.atof((*C.char)(chain.GetStringPtr(str)))
}

//table mydata2 singleton
type MySingleton struct {
	a1 uint64
}

//contract test
type MyContract struct {
	Receiver      chain.Name
	FirstReceiver chain.Name
	Action        chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *MyContract {
	return &MyContract{receiver, firstReceiver, action}
}

func less(a, b []byte) bool {
	if len(a) < len(b) {
		return true
	} else if len(a) > len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] < b[i] {
			return true
		}
		return false
	}
	return false
}

//action test
func (t *MyContract) test() {
	chain.Println(Str2Long("123"))
	chain.Println(Str2Float("123.456"))

	v, err := strconv.ParseFloat("123", 64)
	if err != nil {
		panic(err)
	}
	chain.Println(v)
	return
}
