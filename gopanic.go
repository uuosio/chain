//go:build !eosio

package chain

/*
#include <stdint.h>
*/
import "C"

//export go_panic
func go_panic(funcName *C.char) {
	_funcName := C.GoString(funcName)
	panic(_funcName + " not implemented")
}
