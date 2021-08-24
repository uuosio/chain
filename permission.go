package chain

/*
#include <stdint.h>
typedef uint64_t capi_name;

int32_t
check_transaction_authorization( const char* trx_data,     uint32_t trx_size,
								const char* pubkeys_data, uint32_t pubkeys_size,
								const char* perms_data,   uint32_t perms_size
							);
int32_t
check_permission_authorization( capi_name account,
								capi_name permission,
								const char* pubkeys_data, uint32_t pubkeys_size,
								const char* perms_data,   uint32_t perms_size,
								uint64_t delay_us
							);
int64_t get_permission_last_used( capi_name account, capi_name permission );
int64_t get_account_creation_time( capi_name account );

*/
import "C"
import "unsafe"

// int32_t
// check_transaction_authorization( const char* trx_data,     uint32_t trx_size,
// 								const char* pubkeys_data, uint32_t pubkeys_size,
// 								const char* perms_data,   uint32_t perms_size
// 							);
func CheckTransactionAuthorization(trx_data []byte, pubkeys_data []byte, perms_data []byte) int32 {
	return C.check_transaction_authorization(
		(*C.char)(unsafe.Pointer(&trx_data[0])), (C.uint32_t)(len(trx_data)),
		(*C.char)(unsafe.Pointer(&pubkeys_data[0])), (C.uint32_t)(len(pubkeys_data)),
		(*C.char)(unsafe.Pointer(&perms_data[0])), (C.uint32_t)(len(perms_data)),
	)
}

// int32_t
// check_permission_authorization( capi_name account,
// 								capi_name permission,
// 								const char* pubkeys_data, uint32_t pubkeys_size,
// 								const char* perms_data,   uint32_t perms_size,
// 								uint64_t delay_us
// 							);
func CheckPermissionAuthorization(account Name, permission Name, pubkeys_data []byte, perms_data []byte, delay_us uint64) int32 {
	return C.check_permission_authorization(account.N, permission.N,
		(*C.char)(unsafe.Pointer(&pubkeys_data[0])), (C.uint32_t)(len(pubkeys_data)),
		(*C.char)(unsafe.Pointer(&perms_data[0])), (C.uint32_t)(len(perms_data)),
		delay_us,
	)
}

// int64_t get_permission_last_used( capi_name account, capi_name permission );
func GetPermissionLastUsed(account Name, permission Name) int64 {
	return C.get_permission_last_used(account.N, permission.N)
}

// int64_t get_account_creation_time( capi_name account );
func GetAccountCreationTime(account Name) int64 {
	return C.get_account_creation_time(account.N)
}

type PermissionLevel struct {
	Actor      Name
	Permission Name
}

func (t *PermissionLevel) Pack() []byte {
	enc := NewEncoder(16)
	enc.Pack(&t.Actor)
	enc.Pack(&t.Permission)
	return enc.GetBytes()
}

func (t *PermissionLevel) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	dec.Unpack(&t.Actor)
	dec.Unpack(&t.Permission)
	return dec.Pos(), nil
}

type PermissionLevelWeight struct {
	Permission PermissionLevel
	Weight     uint16
}

func (t *PermissionLevelWeight) Pack() []byte {
	enc := NewEncoder(16 + 2)
	enc.Pack(&t.Permission)
	enc.Pack(t.Weight)
	return enc.GetBytes()
}

func (t *PermissionLevelWeight) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	dec.Unpack(&t.Permission)
	dec.Unpack(&t.Weight)
	return dec.Pos(), nil
}

type KeyWeight struct {
	Key    PublicKey
	Weight uint16
}

func (t *KeyWeight) Pack() []byte {
	enc := NewEncoder(34 + 2)
	enc.Pack(&t.Key)
	enc.Pack(t.Weight)
	return enc.GetBytes()
}

func (t *KeyWeight) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	dec.Unpack(&t.Key)
	dec.Unpack(&t.Weight)
	return dec.Pos(), nil
}

type WaitWeight struct {
	WaitSec uint32
	Weight  uint16
}

func (t *WaitWeight) Pack() []byte {
	enc := NewEncoder(4 + 2)
	enc.Pack(t.WaitSec)
	enc.Pack(t.Weight)
	return enc.GetBytes()
}

func (t *WaitWeight) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	dec.Unpack(&t.WaitSec)
	dec.Unpack(&t.Weight)
	return dec.Pos(), nil
}

type Authority struct {
	Threshold uint32
	Keys      []KeyWeight
	Accounts  []PermissionLevelWeight
	Waits     []WaitWeight
}

func (t *Authority) Pack() []byte {
	enc := NewEncoder(4 + 2)
	enc.Pack(t.Threshold)

	enc.PackLength(len(t.Keys))
	for _, v := range t.Keys {
		enc.Pack(&v)
	}

	enc.PackLength(len(t.Accounts))
	for _, v := range t.Accounts {
		enc.Pack(&v)
	}

	enc.PackLength(len(t.Waits))
	for _, v := range t.Waits {
		enc.Pack(&v)
	}

	return enc.GetBytes()
}

func (t *Authority) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)

	dec.Unpack(&t.Threshold)

	length, _ := dec.UnpackLength()
	t.Keys = make([]KeyWeight, 0, length)
	for i := 0; i < length; i++ {
		v := KeyWeight{}
		dec.Unpack(&v)
		t.Keys = append(t.Keys, v)
	}

	length, _ = dec.UnpackLength()
	t.Accounts = make([]PermissionLevelWeight, 0, length)
	for i := 0; i < length; i++ {
		v := PermissionLevelWeight{}
		dec.Unpack(&v)
		t.Accounts = append(t.Accounts, v)
	}

	length, _ = dec.UnpackLength()
	t.Waits = make([]WaitWeight, 0, length)
	for i := 0; i < length; i++ {
		v := WaitWeight{}
		dec.Unpack(&v)
		t.Waits = append(t.Waits, v)
	}

	return dec.Pos(), nil
}

type NewAccount struct {
	Creator Name
	Name    Name
	Owner   Authority
	Active  Authority
}

func (t *NewAccount) Pack() []byte {
	enc := NewEncoder(int(unsafe.Sizeof(*t)))
	enc.Pack(&t.Creator)
	enc.Pack(&t.Name)
	enc.Pack(&t.Owner)
	enc.Pack(&t.Active)
	return nil
}

func (t *NewAccount) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	dec.Unpack(&t.Creator)
	dec.Unpack(&t.Name)
	dec.Unpack(&t.Owner)
	dec.Unpack(&t.Active)
	return 0, nil
}

type UpdateAuth struct {
	Account    Name
	Permission Name
	Parent     Name
	Auth       Authority
}

func (t *UpdateAuth) Pack() []byte {
	enc := NewEncoder(int(unsafe.Sizeof(*t)))
	enc.Pack(&t.Account)
	enc.Pack(&t.Permission)
	enc.Pack(&t.Parent)
	enc.Pack(&t.Auth)
	return enc.GetBytes()
}

func (t *UpdateAuth) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	dec.Unpack(&t.Account)
	dec.Unpack(&t.Permission)
	dec.Unpack(&t.Parent)
	dec.Unpack(&t.Auth)
	return dec.Pos(), nil
}
