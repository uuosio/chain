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
	var pubKeys_ptr *C.char
	var perms_ptr *C.char
	if len(pubkeys_data) > 0 {
		pubKeys_ptr = (*C.char)(unsafe.Pointer(&pubkeys_data[0]))
	} else {
		pubKeys_ptr = nil
	}

	if len(perms_data) > 0 {
		perms_ptr = (*C.char)(unsafe.Pointer(&perms_data[0]))
	} else {
		perms_ptr = nil
	}
	ret := C.check_transaction_authorization(
		(*C.char)(unsafe.Pointer(&trx_data[0])), (C.uint32_t)(len(trx_data)),
		pubKeys_ptr, (C.uint32_t)(len(pubkeys_data)),
		perms_ptr, (C.uint32_t)(len(perms_data)),
	)
	return int32(ret)
}

//Checks if a permission is authorized by a provided delay and a provided set of keys and permissions
func CheckPermissionAuthorization(account Name, permission Name, pubkeys_data []byte, perms_data []byte, delay_us uint64) int32 {
	ret := C.check_permission_authorization(C.uint64_t(account.N), C.uint64_t(permission.N),
		(*C.char)(unsafe.Pointer(&pubkeys_data[0])), (C.uint32_t)(len(pubkeys_data)),
		(*C.char)(unsafe.Pointer(&perms_data[0])), (C.uint32_t)(len(perms_data)),
		C.uint64_t(delay_us),
	)
	return int32(ret)
}

//Returns the last used time of a permission
func GetPermissionLastUsed(account Name, permission Name) int64 {
	ret := C.get_permission_last_used(C.uint64_t(account.N), C.uint64_t(permission.N))
	return int64(ret)
}

//Returns the creation time of an account
func GetAccountCreationTime(account Name) int64 {
	ret := C.get_account_creation_time(C.uint64_t(account.N))
	return int64(ret)
}

type NewAccount struct {
	Creator Name
	Name    Name
	Owner   Authority
	Active  Authority
}

func CreateNewAccount(creator Name, newAccount Name) *NewAccount {
	a := &NewAccount{}
	a.Creator = creator
	a.Name = newAccount
	return a
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
