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

//Checks if a transaction is authorized by a provided set of keys and permissions
func CheckTransactionAuthorization(trx_data []byte, pubkeys_data []byte, perms_data []byte) int32 {
	return C.check_transaction_authorization(
		(*C.char)(unsafe.Pointer(&trx_data[0])), (C.uint32_t)(len(trx_data)),
		(*C.char)(unsafe.Pointer(&pubkeys_data[0])), (C.uint32_t)(len(pubkeys_data)),
		(*C.char)(unsafe.Pointer(&perms_data[0])), (C.uint32_t)(len(perms_data)),
	)
}

//Checks if a permission is authorized by a provided delay and a provided set of keys and permissions
func CheckPermissionAuthorization(account Name, permission Name, pubkeys_data []byte, perms_data []byte, delay_us uint64) int32 {
	return C.check_permission_authorization(account.N, permission.N,
		(*C.char)(unsafe.Pointer(&pubkeys_data[0])), (C.uint32_t)(len(pubkeys_data)),
		(*C.char)(unsafe.Pointer(&perms_data[0])), (C.uint32_t)(len(perms_data)),
		delay_us,
	)
}

//Returns the last used time of a permission
func GetPermissionLastUsed(account Name, permission Name) int64 {
	return C.get_permission_last_used(account.N, permission.N)
}

//Returns the creation time of an account
func GetAccountCreationTime(account Name) int64 {
	return C.get_account_creation_time(account.N)
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

func (t *NewAccount) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Unpack(&t.Creator)
	dec.Unpack(&t.Name)
	dec.Unpack(&t.Owner)
	dec.Unpack(&t.Active)
	return 0
}
