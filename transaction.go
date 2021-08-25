package chain

/*
#include <stddef.h>
#include <stdint.h>
typedef uint8_t uint128_t;
typedef uint64_t capi_name;

void send_deferred(const uint128_t* sender_id, capi_name payer, const char *serialized_transaction, size_t size, uint32_t replace_existing);
int cancel_deferred(const uint128_t* sender_id);
size_t read_transaction(char *buffer, size_t size);
size_t transaction_size( void );
int tapos_block_num( void );
int tapos_block_prefix( void );
uint32_t expiration( void );
int get_action( uint32_t type, uint32_t index, char* buff, size_t size );
int get_context_free_data( uint32_t index, char* buff, size_t size );
*/
import "C"
import "unsafe"

// void send_deferred(const uint128_t* sender_id, capi_name payer, const char *serialized_transaction, size_t size, uint32_t replace_existing);
func SendDeferred(senderID [2]uint64, payer Name, transaction []byte, replaceExisting bool) {
	cReplaceExisting := C.uint32_t(0)
	if replaceExisting {
		cReplaceExisting = C.uint32_t(1)
	}

	C.send_deferred((*C.uint8_t)(unsafe.Pointer(&senderID[0])), payer.N, (*C.char)(unsafe.Pointer(&transaction[0])), C.size_t(len(transaction)), cReplaceExisting)
}

// int cancel_deferred(const uint128_t* sender_id);
func CancelDeferred(senderID [16]byte) int {
	ret := C.cancel_deferred((*C.uint8_t)(unsafe.Pointer(&senderID[0])))
	return int(ret)
}

// size_t read_transaction(char *buffer, size_t size);
func ReadTransaction() []byte {
	ret := C.read_transaction((*C.char)(unsafe.Pointer(uintptr(0))), 0)
	buffer := make([]byte, ret)
	C.read_transaction((*C.char)(unsafe.Pointer(&buffer[0])), C.size_t(len(buffer)))
	return buffer
}

// __attribute__((eosio_wasm_import))
// size_t transaction_size( void );
func TransactionSize() int {
	ret := C.transaction_size()
	return int(ret)
}

// int tapos_block_num( void );
func TaposBlockNum() int {
	return int(C.tapos_block_num())
}

// int tapos_block_prefix( void );
func TaposBlockPrefix() int {
	return int(C.tapos_block_prefix())
}

// uint32_t expiration( void );
func Expiration() uint32 {
	return C.expiration()
}

// int get_action( uint32_t type, uint32_t index, char* buff, size_t size );
func GetAction(_type uint32, index uint32) []byte {
	var buff []byte
	ret := C.get_action(C.uint32_t(_type), C.uint32_t(index), (*C.char)(unsafe.Pointer(uintptr(0))), 0)

	buf := make([]byte, ret)
	C.get_action(C.uint32_t(_type), C.uint32_t(index), (*C.char)(unsafe.Pointer(&buff[0])), C.size_t(len(buff)))
	return buf
}

// int get_context_free_data( uint32_t index, char* buff, size_t size );
func GetContextFreeData(index uint32) []byte {
	var buff []byte
	ret := C.get_context_free_data(C.uint32_t(index), (*C.char)(unsafe.Pointer(uintptr(0))), 0)

	buf := make([]byte, ret)
	C.get_context_free_data(C.uint32_t(index), (*C.char)(unsafe.Pointer(&buff[0])), C.size_t(len(buff)))
	return buf
}

type TransactionExtension struct {
	Type uint16
	Data []byte
}

func (a *TransactionExtension) EstimatePackedSize() int {
	return 2 + 5 + len(a.Data)
}

func (t *TransactionExtension) Pack() []byte {
	enc := NewEncoder(2 + 5 + len(t.Data))
	enc.Pack(t.Type)
	enc.Pack(t.Data)
	return enc.GetBytes()
}

func (t *TransactionExtension) Unpack(data []byte) (int, error) {
	var err error
	dec := NewDecoder(data)
	t.Type, err = dec.UnpackUint16()
	if err != nil {
		return 0, err
	}

	t.Data, err = dec.UnpackBytes()
	if err != nil {
		return 0, err
	}
	return dec.Pos(), nil
}

func (t *TransactionExtension) Print() {
	Print("{", t.Type, ",", string(t.Data), "}")
}

type Transaction struct {
	// time_point_sec  expiration;
	// uint16_t        ref_block_num;
	// uint32_t        ref_block_prefix;
	// unsigned_int    max_net_usage_words = 0UL; /// number of 8 byte words this transaction can serialize into after compressions
	// uint8_t         max_cpu_usage_ms = 0UL; /// number of CPU usage units to bill transaction for
	// unsigned_int    delay_sec = 0UL; /// number of seconds to delay transaction, default: 0
	Expiration     uint32
	RefBlockNum    uint16
	RefBlockPrefix uint32
	//[VLQ or Base-128 encoding](https://en.wikipedia.org/wiki/Variable-length_quantity)
	//unsigned_int vaint (eosio.cdt/libraries/eosiolib/core/eosio/varint.hpp)
	MaxNetUsageWords   uint32
	MaxCpuUsageMs      uint8
	DelaySec           uint32 //unsigned_int
	ContextFreeActions []Action
	Actions            []Action
	Extention          []TransactionExtension
}

func NewTransaction(delaySec int) *Transaction {
	t := &Transaction{}
	t.Expiration = CurrentTimeSeconds() + uint32(60*60)
	t.RefBlockNum = uint16(TaposBlockNum())
	t.RefBlockPrefix = uint32(TaposBlockPrefix())
	t.MaxNetUsageWords = uint32(0)
	t.MaxCpuUsageMs = uint8(0)
	t.DelaySec = uint32(delaySec)
	return t
}

func (t *Transaction) Send(senderId uint64, replaceExisting bool, payer Name) {
	SendDeferred([2]uint64{senderId, 0}, payer, t.Pack(), replaceExisting)
}

func (t *Transaction) SendEx(senderId [2]uint64, payer Name, replaceExisting bool) {
	SendDeferred(senderId, payer, t.Pack(), replaceExisting)
}

func (t *Transaction) Pack() []byte {
	initSize := 4 + 2 + 4 + 5 + 1 + 5

	initSize += 5 // Max varint size
	for _, action := range t.ContextFreeActions {
		initSize += action.EstimatePackedSize()
	}

	initSize += 5 // Max varint size
	for _, action := range t.Actions {
		initSize += action.EstimatePackedSize()
	}

	initSize += 5 // Max varint size
	for _, extention := range t.Extention {
		initSize += extention.EstimatePackedSize()
	}
	enc := NewEncoder(initSize)
	enc.Pack(t.Expiration)
	enc.Pack(t.RefBlockNum)
	enc.Pack(t.RefBlockPrefix)
	enc.PackVarInt(t.MaxNetUsageWords)
	enc.PackUint8(t.MaxCpuUsageMs)
	enc.PackVarInt(t.DelaySec)

	enc.PackLength(len(t.ContextFreeActions))
	for _, action := range t.ContextFreeActions {
		enc.Pack(&action)
	}

	enc.PackLength(len(t.Actions))
	for _, action := range t.Actions {
		enc.Pack(&action)
	}

	enc.PackLength(len(t.Extention))
	for _, extention := range t.Extention {
		enc.Pack(&extention)
	}
	return enc.GetBytes()
}

func (t *Transaction) Unpack(data []byte) (int, error) {
	var err error

	dec := NewDecoder(data)
	t.Expiration, err = dec.UnpackUint32()
	if err != nil {
		return 0, err
	}

	t.RefBlockNum, err = dec.UnpackUint16()
	if err != nil {
		return 0, err
	}

	t.RefBlockPrefix, err = dec.UnpackUint32()
	if err != nil {
		return 0, err
	}

	t.MaxNetUsageWords, err = dec.UnpackVarInt()
	if err != nil {
		return 0, err
	}

	t.MaxCpuUsageMs, err = dec.UnpackUint8()
	if err != nil {
		return 0, err
	}

	t.DelaySec, err = dec.UnpackVarInt()
	if err != nil {
		return 0, err
	}

	contextFreeActionLength, err := dec.UnpackVarInt()
	if err != nil {
		return 0, err
	}

	t.ContextFreeActions = make([]Action, contextFreeActionLength)
	for i := 0; i < int(contextFreeActionLength); i++ {
		_, err := dec.Unpack(&t.ContextFreeActions[i])
		if err != nil {
			return 0, err
		}
	}

	actionLength, err := dec.UnpackVarInt()
	if err != nil {
		return 0, err
	}

	t.Actions = make([]Action, actionLength)
	for i := 0; i < int(actionLength); i++ {
		_, err := dec.Unpack(&t.Actions[i])
		if err != nil {
			return 0, err
		}
	}

	extentionLength, err := dec.UnpackVarInt()
	if err != nil {
		return 0, err
	}
	t.Extention = make([]TransactionExtension, extentionLength)
	for i := 0; i < int(extentionLength); i++ {
		t.Extention[i].Type, err = dec.UnpackUint16()
		if err != nil {
			return 0, err
		}

		t.Extention[i].Data, err = dec.UnpackBytes()
		if err != nil {
			return 0, err
		}
	}
	return dec.Pos(), nil
}

func (t *Transaction) Print() {
	Print("{")
	Print(t.Expiration, t.RefBlockNum, t.RefBlockPrefix, t.MaxNetUsageWords, t.MaxCpuUsageMs, t.DelaySec)

	Print("[")
	for _, a := range t.ContextFreeActions {
		a.Print()
	}
	Print("]")

	Print("[")
	for _, a := range t.Actions {
		a.Print()
	}
	Print("]")
	Print("[")
	for _, e := range t.Extention {
		e.Print()
	}
	Print("]")
	Print("}")
	Println()
}
