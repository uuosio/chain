package chain

/*
#include <stdint.h>

typedef uint64_t capi_name;
typedef uint8_t capi_checksum256;

void get_resource_limits( capi_name account, int64_t* ram_bytes, int64_t* net_weight, int64_t* cpu_weight );
void set_resource_limits( capi_name account, int64_t ram_bytes, int64_t net_weight, int64_t cpu_weight );
int64_t set_proposed_producers( char *producer_data, uint32_t producer_data_size );
int64_t set_proposed_producers_ex( uint64_t producer_data_format, char *producer_data, uint32_t producer_data_size );
char is_privileged( capi_name account );
void set_privileged( capi_name account, char is_priv );
void set_blockchain_parameters_packed( char* data, uint32_t datalen );
uint32_t get_blockchain_parameters_packed( char* data, uint32_t datalen );
void set_kv_parameters_packed( const char* data, uint32_t datalen );
void preactivate_feature( const capi_checksum256* feature_digest );
*/
import "C"
import "unsafe"

// void get_resource_limits( capi_name account, int64_t* ram_bytes, int64_t* net_weight, int64_t* cpu_weight );
func GetResourceLimits(account Name) (int64, int64, int64) {
	var (
		ram_bytes  int64
		net_weight int64
		cpu_weight int64
	)
	C.get_resource_limits(account.N, (*C.int64_t)(unsafe.Pointer(&ram_bytes)), (*C.int64_t)(unsafe.Pointer(&net_weight)), (*C.int64_t)(unsafe.Pointer(&cpu_weight)))
	return ram_bytes, net_weight, cpu_weight
}

// void set_resource_limits( capi_name account, int64_t ram_bytes, int64_t net_weight, int64_t cpu_weight );
func SetResourceLimits(account Name, ram_bytes, net_weight, cpu_weight int64) {
	C.set_resource_limits(account.N, C.int64_t(ram_bytes), C.int64_t(net_weight), C.int64_t(cpu_weight))
}

// int64_t set_proposed_producers( char *producer_data, uint32_t producer_data_size );
// int64_t set_proposed_producers_ex( uint64_t producer_data_format, char *producer_data, uint32_t producer_data_size );
// bool is_privileged( capi_name account );
func IsPrivileged(account Name) bool {
	return C.is_privileged(account.N) != 0
}

// void set_privileged( capi_name account, bool is_priv );
func SetPrivileged(account Name, is_priv bool) {
	var _is_priv = C.int(0)
	if is_priv {
		_is_priv = 1
	}
	C.set_privileged(account.N, C.char(_is_priv))
}

// void set_blockchain_parameters_packed( char* data, uint32_t datalen );
func SetBlockchainParametersPacked(data []byte) {
	C.set_blockchain_parameters_packed(
		(*C.char)(unsafe.Pointer(&data[0])),
		C.uint32_t(len(data)),
	)
}

// uint32_t get_blockchain_parameters_packed( char* data, uint32_t datalen );
func GetBlockchainParametersPacked(data []byte) uint32 {
	return uint32(C.get_blockchain_parameters_packed(
		(*C.char)(unsafe.Pointer(&data[0])),
		C.uint32_t(len(data)),
	))
}

// void set_kv_parameters_packed( const char* data, uint32_t datalen );
func SetKVParametersPacked(data []byte) {
	C.set_kv_parameters_packed(
		(*C.char)(unsafe.Pointer(&data[0])),
		C.uint32_t(len(data)),
	)
}

// void preactivate_feature( const struct capi_checksum256* feature_digest );
func PreactivateFeature(feature_digest [32]byte) {
	C.preactivate_feature(
		(*C.uint8_t)(unsafe.Pointer(&feature_digest[0])),
	)
}
