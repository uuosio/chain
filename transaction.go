package chain

import (
	"github.com/uuosio/chain/eosio"
)

// void send_deferred(const uint128_t* sender_id, capi_name payer, const char *serialized_transaction, size_t size, uint32_t replace_existing);
func SendDeferred(senderID Uint128, payer Name, transaction []byte, replaceExisting bool) {
	eosio.SendDeferred(senderID, payer.N, transaction, replaceExisting)
}

// int cancel_deferred(const uint128_t* sender_id);
func CancelDeferred(senderID Uint128) int32 {
	return eosio.CancelDeferred(senderID)
}

// size_t read_transaction(char *buffer, size_t size);
func ReadTransaction() []byte {
	return eosio.ReadTransaction()
}

// __attribute__((eosio_wasm_import))
// size_t transaction_size( void );
func TransactionSize() int32 {
	return eosio.TransactionSize()
}

// int tapos_block_num( void );
func TaposBlockNum() int32 {
	return eosio.TaposBlockNum()
}

// int tapos_block_prefix( void );
func TaposBlockPrefix() int32 {
	return eosio.TaposBlockPrefix()
}

// uint32_t expiration( void );
func Expiration() uint32 {
	return eosio.Expiration()
}

// int get_action( uint32_t type, uint32_t index, char* buff, size_t size );
func GetAction(_type uint32, index uint32) []byte {
	return eosio.GetAction(_type, index)
}

// int get_context_free_data( uint32_t index, char* buff, size_t size );
func GetContextFreeData(index uint32) []byte {
	return eosio.GetContextFreeData(index)
}

type TransactionExtension struct {
	Type uint16
	Data []byte
}

func (a *TransactionExtension) Size() int {
	return 2 + 5 + len(a.Data)
}

func (t *TransactionExtension) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.PackUint16(t.Type)
	enc.PackBytes(t.Data)
	return enc.GetSize() - oldSize
}

func (t *TransactionExtension) Unpack(data []byte) int {
	dec := NewDecoder(data)
	t.Type = dec.UnpackUint16()

	t.Data = dec.UnpackBytes()
	return dec.Pos()
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
	MaxNetUsageWords   VarUint32
	MaxCpuUsageMs      uint8
	DelaySec           VarUint32 //unsigned_int
	ContextFreeActions []*Action
	Actions            []*Action
	Extention          []*TransactionExtension
}

func NewTransaction(delaySec int) *Transaction {
	t := &Transaction{}
	t.Expiration = CurrentTimeSeconds() + uint32(60*60)
	t.RefBlockNum = uint16(TaposBlockNum())
	t.RefBlockPrefix = uint32(TaposBlockPrefix())
	t.MaxNetUsageWords = VarUint32(0)
	t.MaxCpuUsageMs = uint8(0)
	t.DelaySec = VarUint32(delaySec)
	return t
}

type TransactionCache struct {
	sendId      *Uint128
	payer       Name
	transaction []byte
	replace     bool
}

var gTransactionCache []*TransactionCache

func AddTransactionCache(sendId *Uint128, payer Name, transaction []byte, replace bool) {
	if gTransactionCache == nil {
		gTransactionCache = make([]*TransactionCache, 0, 2)
	}
	gTransactionCache = append(gTransactionCache, &TransactionCache{sendId, payer, transaction, replace})
}

func (t *Transaction) Send(senderId Uint128, replaceExisting bool, payer Name) {
	SendDeferred(senderId, payer, EncoderPack(t), replaceExisting)
}

func (t *Transaction) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	//4+2+4
	enc.PackUint32(t.Expiration)
	enc.PackUint16(t.RefBlockNum)
	enc.PackUint32(t.RefBlockPrefix)
	t.MaxNetUsageWords.Pack(enc)
	enc.PackUint8(t.MaxCpuUsageMs)
	t.DelaySec.Pack(enc)

	enc.PackLength(len(t.ContextFreeActions))
	for _, action := range t.ContextFreeActions {
		action.Pack(enc)
	}

	enc.PackLength(len(t.Actions))
	for _, action := range t.Actions {
		action.Pack(enc)
	}

	enc.PackLength(len(t.Extention))
	for _, extention := range t.Extention {
		extention.Pack(enc)
	}
	return enc.GetSize() - oldSize
}

func (t *Transaction) Size() int {
	initSize := 4 + 2 + 4 + 5 + 1 + 5

	initSize += PackedVarUint32Length(uint32(len(t.ContextFreeActions))) // Max varint size
	for _, action := range t.ContextFreeActions {
		initSize += action.Size()
	}

	initSize += PackedVarUint32Length(uint32(len(t.Actions)))
	for _, action := range t.Actions {
		initSize += action.Size()
	}

	initSize += PackedVarUint32Length(uint32(len(t.Extention)))
	for _, extention := range t.Extention {
		initSize += extention.Size()
	}
	return initSize
}

func (t *Transaction) Unpack(data []byte) int {

	dec := NewDecoder(data)
	t.Expiration = dec.UnpackUint32()

	t.RefBlockNum = dec.UnpackUint16()

	t.RefBlockPrefix = dec.UnpackUint32()

	dec.Unpack(&t.MaxNetUsageWords)

	t.MaxCpuUsageMs = dec.UnpackUint8()

	dec.Unpack(&t.DelaySec)

	contextFreeActionLength := dec.UnpackVarUint32()

	t.ContextFreeActions = make([]*Action, contextFreeActionLength)
	for i := 0; i < int(contextFreeActionLength); i++ {
		action := &Action{}
		dec.Unpack(action)
		t.ContextFreeActions[i] = action
	}

	actionLength := dec.UnpackVarUint32()

	t.Actions = make([]*Action, actionLength)
	for i := 0; i < int(actionLength); i++ {
		action := &Action{}
		dec.Unpack(action)
		t.Actions[i] = action
	}

	extentionLength := dec.UnpackVarUint32()
	t.Extention = make([]*TransactionExtension, extentionLength)
	for i := 0; i < int(extentionLength); i++ {
		extention := &TransactionExtension{}
		extention.Type = dec.UnpackUint16()
		extention.Data = dec.UnpackBytes()
		t.Extention[i] = extention
	}
	return dec.Pos()
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
