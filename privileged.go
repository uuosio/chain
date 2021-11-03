package chain

/*
#include "chain.h"

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
import (
	"unsafe"
)

//Get the resource limits of an account
func GetResourceLimits(account Name) (int64, int64, int64) {
	var (
		ram_bytes  int64
		net_weight int64
		cpu_weight int64
	)
	C.get_resource_limits(C.uint64_t(account.N), (*C.int64_t)(unsafe.Pointer(&ram_bytes)), (*C.int64_t)(unsafe.Pointer(&net_weight)), (*C.int64_t)(unsafe.Pointer(&cpu_weight)))
	return ram_bytes, net_weight, cpu_weight
}

//Set the resource limits of an account
func SetResourceLimits(account Name, ram_bytes, net_weight, cpu_weight int64) {
	C.set_resource_limits(C.uint64_t(account.N), C.int64_t(ram_bytes), C.int64_t(net_weight), C.int64_t(cpu_weight))
}

type BlockSigningAuthorityV0 struct {
	Threshold uint32
	Keys      []KeyWeight
}

func (t *BlockSigningAuthorityV0) IsValid() bool {
	return true
}

// using block_signing_authority = std::variant<block_signing_authority_v0>;

type ProducerAuthority struct {
	ProducerName Name
	Authority    BlockSigningAuthorityV0
}

func (t *BlockSigningAuthorityV0) Pack() []byte {
	enc := NewEncoder(t.Size())
	enc.PackUint32(t.Threshold)

	{
		enc.PackLength(len(t.Keys))
		for _, v := range t.Keys {
			enc.Pack(&v)
		}
	}
	return enc.GetBytes()
}

func (t *BlockSigningAuthorityV0) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Unpack(&t.Threshold)
	{
		length := dec.UnpackLength()
		t.Keys = make([]KeyWeight, length)
		for i := 0; i < length; i++ {
			dec.Unpack(&t.Keys[i])
		}
	}
	return dec.Pos()
}

func (t *BlockSigningAuthorityV0) Size() int {
	size := 0
	size += 4
	size += PackedVarUint32Length(uint32(len(t.Keys)))

	for i := range t.Keys {
		size += t.Keys[i].Size()
	}
	return size
}

func (t *ProducerAuthority) Pack() []byte {
	enc := NewEncoder(t.Size())
	enc.PackUint64(t.ProducerName.N)
	enc.PackLength(0) //variant
	enc.Pack(&t.Authority)
	return enc.GetBytes()
}

func (t *ProducerAuthority) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Unpack(&t.ProducerName)
	length := dec.UnpackLength()
	Check(length == 0, "bad variant index")
	dec.Unpack(&t.Authority)
	return dec.Pos()
}

func (t *ProducerAuthority) Size() int {
	size := 0
	size += 8
	size += t.Authority.Size()
	return size
}

//Proposes a schedule change
func SetProposedProducers(producers []Name) int64 {
	_producers := make([]uint64, len(producers))
	for i, p := range producers {
		_producers[i] = p.N
	}

	return int64(C.set_proposed_producers(
		(*C.char)(unsafe.Pointer(&_producers[0])),
		C.uint32_t(len(producers)*8),
	))
}

// Proposes a schedule change with extended features
func SetProposedProducersEx(producers []ProducerAuthority) int64 {
	enc := NewEncoder(len(producers) * int(unsafe.Sizeof(producers[0])))
	for i := range producers {
		enc.Write(producers[i].Pack())
	}
	producer_data := enc.GetBytes()
	ret := C.set_proposed_producers_ex(C.uint64_t(1), (*C.char)(unsafe.Pointer(&producer_data[0])), C.uint32_t(len(producer_data)))
	return int64(ret)
}

//Check if an account is privileged
func IsPrivileged(account Name) bool {
	return C.is_privileged(C.uint64_t(account.N)) != 0
}

//Set the privileged status of an account
func SetPrivileged(account Name, is_priv bool) {
	var _is_priv = C.int(0)
	if is_priv {
		_is_priv = 1
	}
	C.set_privileged(C.uint64_t(account.N), C.char(_is_priv))
}

//Set the blockchain parameters
func SetBlockchainParametersPacked(data []byte) {
	C.set_blockchain_parameters_packed(
		(*C.char)(unsafe.Pointer(&data[0])),
		C.uint32_t(len(data)),
	)
}

//Retrieve the blolckchain parameters
func GetBlockchainParametersPacked(data []byte) uint32 {
	return uint32(C.get_blockchain_parameters_packed(
		(*C.char)(unsafe.Pointer(&data[0])),
		C.uint32_t(len(data)),
	))
}

//Set the KV parameters
func SetKVParametersPacked(data []byte) {
	C.set_kv_parameters_packed(
		(*C.char)(unsafe.Pointer(&data[0])),
		C.uint32_t(len(data)),
	)
}

//Pre-activate protocol feature
func PreactivateFeature(feature_digest [32]byte) {
	C.preactivate_feature(
		(*C.capi_checksum256)(unsafe.Pointer(&feature_digest[0])),
	)
}
