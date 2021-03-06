package chain

/*
#include <stdint.h>

uint32_t get_active_producers( uint64_t* producers, uint32_t datalen );
*/
import "C"
import "unsafe"

//Gets the set of active producers.
func GetActiveProducers() []Name {
	datalen := C.get_active_producers((*C.uint64_t)(unsafe.Pointer(uintptr(0))), 0)
	if datalen == 0 {
		return nil
	}

	var producers = make([]uint64, int(datalen)/8)
	C.get_active_producers((*C.uint64_t)(unsafe.Pointer(&producers[0])), datalen)

	var _producers = make([]Name, 0, int(datalen/8))
	for _, v := range producers {
		_producers = append(_producers, Name{v})
	}
	return _producers
}
