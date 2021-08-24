package chain

/*
#include <stdint.h>

uint32_t get_active_producers( uint8_t* producers, uint32_t datalen );
*/
import "C"
import "unsafe"

// uint32_t get_active_producers( capi_name* producers, uint32_t datalen );
func GetActiveProducers() []Name {
	var datalen uint32 = 0
	datalen = C.get_active_producers((*C.uint8_t)(unsafe.Pointer(uintptr(0))), 0)

	var producers = make([]uint64, int(datalen)/8)
	C.get_active_producers((*C.uint8_t)(unsafe.Pointer(&producers[0])), datalen)

	var _producers = make([]Name, 0, int(datalen/8))
	for _, v := range producers {
		_producers = append(_producers, Name{v})
	}
	return _producers
}
