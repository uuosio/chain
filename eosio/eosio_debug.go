//go:build !tinygo.wasm
// +build !tinygo.wasm

package eosio

import (
	"context"
	"encoding/binary"
	"unsafe"

	"github.com/uuosio/chaintester"
	"github.com/uuosio/chaintester/interfaces"
)

var ctx = context.Background()

func CheckError(err error) {
	if err != nil {
		panic(chaintester.NewAssertError(err))
	}
}

//action.h

//Read current action data
func ReadActionData() []byte {
	ret, err := chaintester.GetVMAPI().ReadActionData(ctx)
	CheckError(err)
	return ret
}

//Get the length of the current action's data field
func ActionDataSize() uint32 {
	size, err := chaintester.GetVMAPI().ActionDataSize(ctx)
	CheckError(err)
	return uint32(size)
}

//Add the specified account to set of accounts to be notified
func RequireRecipient(name uint64) {
	err := chaintester.GetVMAPI().RequireRecipient(ctx, to_raw_uint64(name))
	CheckError(err)
}

//Verifies that name exists in the set of provided auths on a action. Throws if not found.
func RequireAuth(name uint64) {
	err := chaintester.GetVMAPI().RequireAuth(ctx, to_raw_uint64(name))
	CheckError(err)
}

//Verifies that name has auth.
func HasAuth(name uint64) bool {
	ret, err := chaintester.GetVMAPI().HasAuth(ctx, to_raw_uint64(name))
	CheckError(err)
	return ret
}

//Verifies that name exists in the set of provided auths on a action. Throws if not found.
func RequireAuth2(name uint64, permission uint64) {
	err := chaintester.GetVMAPI().RequireAuth2(ctx, to_raw_uint64(name), to_raw_uint64(permission))
	CheckError(err)
}

//Verifies that name is an existing account.
func IsAccount(name uint64) bool {
	ret, err := chaintester.GetVMAPI().IsAccount(ctx, to_raw_uint64(name))
	CheckError(err)
	return ret
}

//Send an inline action in the context of this action's parent transaction
func SendInline(data []byte) {
	err := chaintester.GetVMAPI().SendInline(ctx, data)
	CheckError(err)
}

//Send an inline context free action in the context of this action's parent transaction
func SendContextFreeInline(data []byte) {
	err := chaintester.GetVMAPI().SendContextFreeInline(ctx, data)
	CheckError(err)
}

//Returns the time in microseconds from 1970 of the publication_time
func PublicationTime() uint64 {
	ret, err := chaintester.GetVMAPI().PublicationTime(ctx)
	CheckError(err)
	return from_raw_uint64(ret)
}

//Get the current receiver of the action
func CurrentReceiver() uint64 {
	ret, err := chaintester.GetVMAPI().CurrentReceiver(ctx)
	CheckError(err)
	return from_raw_uint64(ret)
}

//Set the action return value which will be included in the action_receipt
func SetActionReturnValue(return_value []byte) {
	panic("SetActionReturnValue not implemented")
}

//system.h
func Check(b bool, msg string) {
	if !b {
		EosioAssert(false, msg)
	}
}

//Aborts processing of this action and unwinds all pending changes if the test condition is true
func Assert(test bool, msg string) {
	if !test {
		EosioAssert(false, msg)
	}
}

//Aborts processing of this action and unwinds all pending changes if the test condition is true
func EosioAssert(test bool, msg string) {
	if !test {
		err := chaintester.GetVMAPI().EosioAssertMessage(ctx, false, []byte(msg))
		CheckError(err)
	}
}

func to_raw_uint64(value uint64) *interfaces.Uint64 {
	ret := &interfaces.Uint64{}
	ret.RawValue = make([]byte, 8)
	binary.LittleEndian.PutUint64(ret.RawValue, value)
	return ret
}

func from_raw_uint64(value *interfaces.Uint64) uint64 {
	return binary.LittleEndian.Uint64(value.RawValue)
}

//Aborts processing of this action and unwinds all pending changes if the test condition is true
func EosioAssertCode(test bool, code uint64) {
	if !test {
		err := chaintester.GetVMAPI().EosioAssertCode(ctx, false, to_raw_uint64(code))
		CheckError(err)
	}
}

//Returns the time in microseconds from 1970 of the current block
func CurrentTime() uint64 {
	ret, err := chaintester.GetVMAPI().CurrentTime(ctx)
	CheckError(err)
	return from_raw_uint64(ret)
}

//TODO:
//Check if specified protocol feature has been activated
func IsFeatureActivated(featureDigest [32]byte) bool {
	panic("IsFeatureActivated not implemented")
	return false
}

//Return name of account that sent current inline action
func GetSender() uint64 {
	ret, err := chaintester.GetVMAPI().GetSender(ctx)
	CheckError(err)
	return from_raw_uint64(ret)
}

func Exit() {
	err := chaintester.GetVMAPI().EosioExit(ctx, 0)
	CheckError(err)
}

//Tests if the sha256 hash generated from data matches the provided checksum.
func AssertSha256(data []byte, hash [32]byte) {
	err := chaintester.GetVMAPI().AssertSha256(ctx, data, hash[:])
	CheckError(err)
}

//Tests if the sha1 hash generated from data matches the provided checksum.
func AssertSha1(data []byte, hash [20]byte) {
	err := chaintester.GetVMAPI().AssertSha1(ctx, data, hash[:])
	CheckError(err)
}

//Tests if the sha512 hash generated from data matches the provided checksum.
func AssertSha512(data []byte, hash [64]byte) {
	err := chaintester.GetVMAPI().AssertSha512(ctx, data, hash[:])
	CheckError(err)
}

//Tests if the ripemod160 hash generated from data matches the provided checksum.
func AssertRipemd160(data []byte, hash [20]byte) {
	err := chaintester.GetVMAPI().AssertRipemd160(ctx, data, hash[:])
	CheckError(err)
}

//Hashes data using sha256 and return hash value.
func Sha256(data []byte) [32]byte {
	_ret, err := chaintester.GetVMAPI().Sha256(ctx, data[:])
	CheckError(err)
	var ret [32]byte
	copy(ret[:], _ret)
	return ret
}

//Hashes data using sha1 and return hash value.
func Sha1(data []byte) [20]byte {
	_ret, err := chaintester.GetVMAPI().Sha1(ctx, data[:])
	CheckError(err)
	var ret [20]byte
	copy(ret[:], _ret)
	return ret
}

//Hashes data using sha512 and return hash value.
func Sha512(data []byte) [64]byte {
	_ret, err := chaintester.GetVMAPI().Sha512(ctx, data[:])
	CheckError(err)
	var ret [64]byte
	copy(ret[:], _ret)
	return ret
}

//Hashes data using ripemd160 and return hash value.
func Ripemd160(data []byte) [20]byte {
	_ret, err := chaintester.GetVMAPI().Ripemd160(ctx, data[:])
	CheckError(err)
	var ret [20]byte
	copy(ret[:], _ret)
	return ret
}

//Recover the public key from digest and signature
func RecoverKey(digest [32]byte, sig []byte) []byte {
	ret, err := chaintester.GetVMAPI().RecoverKey(ctx, digest[:], sig)
	CheckError(err)
	return ret
}

//Tests a given public key with the generated key from digest and the signature
func AssertRecoverKey(digest [32]byte, sig []byte, pub []byte) {
	err := chaintester.GetVMAPI().AssertRecoverKey(ctx, digest[:], sig, pub)
	CheckError(err)
}

func DBStoreI64(scope uint64, table uint64, payer uint64, id uint64, data []byte) int32 {
	ret, err := chaintester.GetVMAPI().DbStoreI64(ctx, to_raw_uint64(scope), to_raw_uint64(table), to_raw_uint64(payer), to_raw_uint64(id), data)
	CheckError(err)
	return ret
}

// void db_update_i64(int32_t iterator, uint64_t payer, const char* data, uint32_t len);
func DBUpdateI64(iterator int32, payer uint64, data []byte) {
	err := chaintester.GetVMAPI().DbUpdateI64(ctx, iterator, to_raw_uint64(payer), data)
	CheckError(err)
}

// void db_remove_i64(int32_t iterator);
func DBRemoveI64(iterator int32) {
	err := chaintester.GetVMAPI().DbRemoveI64(ctx, iterator)
	CheckError(err)
}

// int32_t db_get_i64(int32_t iterator, const char* data, uint32_t len);
func DBGetI64(iterator int32) []byte {
	ret, err := chaintester.GetVMAPI().DbGetI64(ctx, iterator)
	CheckError(err)
	return ret
}

// int32_t db_next_i64(int32_t iterator, uint64_t* primary);
func DBNextI64(iterator int32) (int32, uint64) {
	ret, err := chaintester.GetVMAPI().DbNextI64(ctx, iterator)
	CheckError(err)
	return ret.Iterator, from_raw_uint64(ret.Primary)
}

// int32_t db_previous_i64(int32_t iterator, uint64_t* primary);
func DBPreviousI64(iterator int32) (int32, uint64) {
	ret, err := chaintester.GetVMAPI().DbPreviousI64(ctx, iterator)
	CheckError(err)
	return ret.Iterator, from_raw_uint64(ret.Primary)
}

// int32_t db_find_i64(uint64_t code, uint64_t scope, uint64_t table, uint64_t id);
func DBFindI64(code uint64, scope uint64, table uint64, id uint64) int32 {
	ret, err := chaintester.GetVMAPI().DbFindI64(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), to_raw_uint64(id))
	CheckError(err)
	return ret
}

// int32_t db_lowerbound_i64(uint64_t code, uint64_t scope, uint64_t table, uint64_t id);
func DBLowerBoundI64(code uint64, scope uint64, table uint64, id uint64) int32 {
	ret, err := chaintester.GetVMAPI().DbLowerboundI64(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), to_raw_uint64(id))
	CheckError(err)
	return ret
}

// int32_t db_upperbound_i64(uint64_t code, uint64_t scope, uint64_t table, uint64_t id);
func DBUpperBoundI64(code uint64, scope uint64, table uint64, id uint64) int32 {
	ret, err := chaintester.GetVMAPI().DbUpperboundI64(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), to_raw_uint64(id))
	CheckError(err)
	return ret
}

// int32_t db_end_i64(uint64_t code, uint64_t scope, uint64_t table);
func DBEndI64(code uint64, scope uint64, table uint64) int32 {
	ret, err := chaintester.GetVMAPI().DbEndI64(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table))
	CheckError(err)
	return ret
}

// int32_t db_idx64_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const uint64_t* secondary);
func DBIdx64Store(scope uint64, table uint64, id uint64, secondary uint64, payer uint64) int32 {
	ret, err := chaintester.GetVMAPI().DbIdx64Store(ctx, to_raw_uint64(scope), to_raw_uint64(table), to_raw_uint64(payer), to_raw_uint64(id), to_raw_uint64(secondary))
	CheckError(err)
	return ret
}

// void db_idx64_update(int32_t iterator, uint64_t payer, const uint64_t* secondary);
func DBIdx64Update(it int32, secondary uint64, payer uint64) {
	err := chaintester.GetVMAPI().DbIdx64Update(ctx, it, to_raw_uint64(payer), to_raw_uint64(secondary))
	CheckError(err)
}

// void db_idx64_remove(int32_t iterator);
func DBIdx64Remove(it int32) {
	err := chaintester.GetVMAPI().DbIdx64Remove(ctx, it)
	CheckError(err)
}

// int32_t db_idx64_next(int32_t iterator, uint64_t* primary);
func DBIdx64Next(it int32) (int32, uint64) {
	ret, err := chaintester.GetVMAPI().DbIdx64Next(ctx, it)
	CheckError(err)
	return ret.Iterator, from_raw_uint64(ret.Primary)
}

// int32_t db_idx64_previous(int32_t iterator, uint64_t* primary);
func DBIdx64Previous(it int32) (int32, uint64) {
	ret, err := chaintester.GetVMAPI().DbIdx64Previous(ctx, it)
	CheckError(err)
	return ret.Iterator, from_raw_uint64(ret.Primary)
}

// int32_t db_idx64_find_primary(uint64_t code, uint64_t scope, uint64_t table, uint64_t* secondary, uint64_t primary);
func DBIdx64FindByPrimary(code uint64, scope uint64, table uint64, primary uint64) (int32, uint64) {
	ret, err := chaintester.GetVMAPI().DbIdx64FindPrimary(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), to_raw_uint64(primary))
	CheckError(err)
	return ret.Iterator, binary.LittleEndian.Uint64(ret.Secondary)
}

// int32_t db_idx64_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const uint64_t* secondary, uint64_t* primary);
func DBIdx64Find(code uint64, scope uint64, table uint64, secondary uint64) (int32, uint64) {
	it, _, _secondary := DBIdx64Lowerbound(code, scope, table, secondary)
	if it >= 0 {
		if _secondary == secondary {
			return it, secondary
		}
	}
	return -1, 0
}

// int32_t db_idx64_lowerbound(uint64_t code, uint64_t scope, uint64_t table, uint64_t* secondary, uint64_t* primary);
func DBIdx64Lowerbound(code uint64, scope uint64, table uint64, secondary uint64) (int32, uint64, uint64) {
	var primary uint64 = 0
	ret, err := chaintester.GetVMAPI().DbIdx64Lowerbound(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), to_raw_uint64(secondary), to_raw_uint64(primary))
	CheckError(err)
	return ret.Iterator, binary.LittleEndian.Uint64(ret.Secondary), from_raw_uint64(ret.Primary)
}

// int32_t db_idx64_upperbound(uint64_t code, uint64_t scope, uint64_t table, uint64_t* secondary, uint64_t* primary);
func DBIdx64Upperbound(code uint64, scope uint64, table uint64, secondary uint64) (int32, uint64, uint64) {
	var primary uint64 = 0
	ret, err := chaintester.GetVMAPI().DbIdx64Upperbound(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), to_raw_uint64(secondary), to_raw_uint64(primary))
	CheckError(err)
	return ret.Iterator, binary.LittleEndian.Uint64(ret.Secondary), from_raw_uint64(ret.Primary)
}

// int32_t db_idx64_end(uint64_t code, uint64_t scope, uint64_t table);
func DBIdx64End(code uint64, scope uint64, table uint64) int32 {
	ret, err := chaintester.GetVMAPI().DbIdx64End(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table))
	CheckError(err)
	return ret
}

// int32_t db_idx128_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const uint128* secondary);
func DBIdx128Store(scope uint64, table uint64, id uint64, secondary [16]byte, payer uint64) int32 {
	it, err := chaintester.GetVMAPI().DbIdx128Store(ctx, to_raw_uint64(scope), to_raw_uint64(table), to_raw_uint64(payer), to_raw_uint64(id), secondary[:])
	CheckError(err)
	return it
}

// void db_idx128_update(int32_t iterator, uint64_t payer, const uint128* secondary);
func DBIdx128Update(it int32, secondary [16]byte, payer uint64) {
	err := chaintester.GetVMAPI().DbIdx128Update(ctx, it, to_raw_uint64(payer), secondary[:])
	CheckError(err)
}

// void db_idx128_remove(int32_t iterator);
func DBIdx128Remove(it int32) {
	err := chaintester.GetVMAPI().DbIdx128Remove(ctx, it)
	CheckError(err)
}

// int32_t db_idx128_next(int32_t iterator, uint64_t* primary);
func DBIdx128Next(it int32) (int32, uint64) {
	ret, err := chaintester.GetVMAPI().DbIdx128Next(ctx, it)
	CheckError(err)
	return ret.Iterator, from_raw_uint64(ret.Primary)
}

// int32_t db_idx128_previous(int32_t iterator, uint64_t* primary);
func DBIdx128Previous(it int32) (int32, uint64) {
	ret, err := chaintester.GetVMAPI().DbIdx128Previous(ctx, it)
	CheckError(err)
	return ret.Iterator, from_raw_uint64(ret.Primary)
}

// int32_t db_idx128_find_primary(uint64_t code, uint64_t scope, uint64_t table, uint128* secondary, uint64_t primary);
func DBIdx128FindByPrimary(code uint64, scope uint64, table uint64, primary uint64) (int32, [16]byte) {
	ret, err := chaintester.GetVMAPI().DbIdx128FindPrimary(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), to_raw_uint64(primary))
	CheckError(err)
	var _secondary [16]byte
	copy(_secondary[:], ret.Secondary)
	return ret.Iterator, _secondary
}

// int32_t db_idx128_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const uint128* secondary, uint64_t* primary);
func DBIdx128Find(code uint64, scope uint64, table uint64, secondary [16]byte) (int32, [16]byte, uint64) {
	it, _secondary, primary := DBIdx128Lowerbound(code, scope, table, secondary)
	if it >= 0 {
		if _secondary == secondary {
			return it, _secondary, primary
		}
	}
	return it, _secondary, 0
}

// int32_t db_idx128_lowerbound(uint64_t code, uint64_t scope, uint64_t table, uint128* secondary, uint64_t* primary);
func DBIdx128Lowerbound(code uint64, scope uint64, table uint64, secondary [16]byte) (int32, [16]byte, uint64) {
	ret, err := chaintester.GetVMAPI().DbIdx128Lowerbound(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), secondary[:], to_raw_uint64(0))
	CheckError(err)
	copy(secondary[:], ret.Secondary)
	return ret.Iterator, secondary, from_raw_uint64(ret.Primary)
}

// int32_t db_idx128_upperbound(uint64_t code, uint64_t scope, uint64_t table, uint128* secondary, uint64_t* primary);
func DBIdx128Upperbound(code uint64, scope uint64, table uint64, secondary [16]byte) (int32, [16]byte, uint64) {
	ret, err := chaintester.GetVMAPI().DbIdx128Upperbound(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), secondary[:], to_raw_uint64(0))
	CheckError(err)
	copy(secondary[:], ret.Secondary)
	return ret.Iterator, secondary, from_raw_uint64(ret.Primary)
}

// int32_t db_idx128_end(uint64_t code, uint64_t scope, uint64_t table);
func DBIdx128End(code uint64, scope uint64, table uint64) int32 {
	it, err := chaintester.GetVMAPI().DbIdx128End(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table))
	CheckError(err)
	return it
}

// int32_t db_idx256_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const uint128* data, uint32_t data_len );
func DBIdx256Store(scope uint64, table uint64, id uint64, secondary [32]byte, payer uint64) int32 {
	it, err := chaintester.GetVMAPI().DbIdx256Store(ctx, to_raw_uint64(scope), to_raw_uint64(table), to_raw_uint64(payer), to_raw_uint64(id), secondary[:])
	CheckError(err)
	return it
}

// void db_idx256_update(int32_t iterator, uint64_t payer, const uint128* data, uint32_t data_len);
func DBIdx256Update(it int32, secondary [32]byte, payer uint64) {
	err := chaintester.GetVMAPI().DbIdx256Update(ctx, it, to_raw_uint64(payer), secondary[:])
	CheckError(err)
}

// void db_idx256_remove(int32_t iterator);
func DBIdx256Remove(it int32) {
	err := chaintester.GetVMAPI().DbIdx256Remove(ctx, it)
	CheckError(err)
}

// int32_t db_idx256_next(int32_t iterator, uint64_t* primary);
func DBIdx256Next(it int32) (int32, uint64) {
	ret, err := chaintester.GetVMAPI().DbIdx256Next(ctx, it)
	CheckError(err)
	return ret.Iterator, from_raw_uint64(ret.Primary)
}

// int32_t db_idx256_previous(int32_t iterator, uint64_t* primary);
func DBIdx256Previous(it int32) (int32, uint64) {
	ret, err := chaintester.GetVMAPI().DbIdx256Previous(ctx, it)
	CheckError(err)
	return ret.Iterator, from_raw_uint64(ret.Primary)
}

// int32_t db_idx256_find_primary(uint64_t code, uint64_t scope, uint64_t table, uint128* data, uint32_t data_len, uint64_t primary);
func DBIdx256FindByPrimary(code uint64, scope uint64, table uint64, primary uint64) (int32, [32]byte) {
	ret, err := chaintester.GetVMAPI().DbIdx256FindPrimary(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), to_raw_uint64(primary))
	CheckError(err)
	var _secondary [32]byte
	copy(_secondary[:], ret.Secondary)
	return ret.Iterator, _secondary
}

// int32_t db_idx256_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const uint128* data, uint32_t data_len, uint64_t* primary);
func DBIdx256Find(code uint64, scope uint64, table uint64, secondary [32]byte) (int32, [32]byte, uint64) {
	it, _secondary, primary := DBIdx256Lowerbound(code, scope, table, secondary)
	if it >= 0 {
		if _secondary == secondary {
			return it, _secondary, primary
		}
	}
	return it, _secondary, 0
}

// int32_t db_idx256_lowerbound(uint64_t code, uint64_t scope, uint64_t table, uint128* data, uint32_t data_len, uint64_t* primary);
func DBIdx256Lowerbound(code uint64, scope uint64, table uint64, secondary [32]byte) (int32, [32]byte, uint64) {
	ret, err := chaintester.GetVMAPI().DbIdx256Lowerbound(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), secondary[:], to_raw_uint64(0))
	CheckError(err)
	copy(secondary[:], ret.Secondary)
	return ret.Iterator, secondary, from_raw_uint64(ret.Primary)
}

// int32_t db_idx256_upperbound(uint64_t code, uint64_t scope, uint64_t table, uint128* data, uint32_t data_len, uint64_t* primary);
func DBIdx256Upperbound(code uint64, scope uint64, table uint64, secondary [32]byte) (int32, [32]byte, uint64) {
	ret, err := chaintester.GetVMAPI().DbIdx256Upperbound(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), secondary[:], to_raw_uint64(0))
	CheckError(err)
	copy(secondary[:], ret.Secondary)
	return ret.Iterator, secondary, from_raw_uint64(ret.Primary)
}

// int32_t db_idx256_end(uint64_t code, uint64_t scope, uint64_t table);
func DBIdx256End(code uint64, scope uint64, table uint64) int32 {
	it, err := chaintester.GetVMAPI().DbIdx256End(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table))
	CheckError(err)
	return it
}

// int32_t db_idx_double_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const Double* secondary);
func DBIdxFloat64Store(scope uint64, table uint64, id uint64, secondary float64, payer uint64) int32 {
	_secondary := (*[8]byte)(unsafe.Pointer(&secondary))
	ret, err := chaintester.GetVMAPI().DbIdxDoubleStore(ctx, to_raw_uint64(scope), to_raw_uint64(table), to_raw_uint64(payer), to_raw_uint64(id), _secondary[:])
	CheckError(err)
	return ret
}

// void db_idx_double_update(int32_t iterator, uint64_t payer, const Double* secondary);
func DBIdxFloat64Update(it int32, secondary float64, payer uint64) {
	_secondary := (*[8]byte)(unsafe.Pointer(&secondary))
	err := chaintester.GetVMAPI().DbIdxDoubleUpdate(ctx, it, to_raw_uint64(payer), _secondary[:])
	CheckError(err)
}

// void db_idx_double_remove(int32_t iterator);
func DBIdxFloat64Remove(it int32) {
	err := chaintester.GetVMAPI().DbIdxDoubleRemove(ctx, it)
	CheckError(err)
}

// int32_t db_idx_double_next(int32_t iterator, uint64_t* primary);
func DBIdxFloat64Next(it int32) (int32, uint64) {
	ret, err := chaintester.GetVMAPI().DbIdxDoubleNext(ctx, it)
	CheckError(err)
	return ret.Iterator, from_raw_uint64(ret.Primary)
}

// int32_t db_idx_double_previous(int32_t iterator, uint64_t* primary);
func DBIdxFloat64Previous(it int32) (int32, uint64) {
	ret, err := chaintester.GetVMAPI().DbIdxDoublePrevious(ctx, it)
	CheckError(err)
	return ret.Iterator, from_raw_uint64(ret.Primary)
}

// int32_t db_idx_double_find_primary(uint64_t code, uint64_t scope, uint64_t table, Double* secondary, uint64_t primary);
func DBIdxFloat64FindByPrimary(code uint64, scope uint64, table uint64, primary uint64) (int32, float64) {
	ret, err := chaintester.GetVMAPI().DbIdxDoubleFindPrimary(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), to_raw_uint64(primary))
	CheckError(err)
	secondary := float64(0.0)
	_secondary := (*[8]byte)(unsafe.Pointer(&secondary))
	copy(_secondary[:], ret.Secondary)
	return ret.Iterator, secondary
}

// int32_t db_idx_double_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const Double* secondary, uint64_t* primary);
func DBIdxFloat64Find(code uint64, scope uint64, table uint64, secondary float64) (int32, uint64) {
	it, _secondary, primary := DBIdxFloat64Lowerbound(code, scope, table, secondary)
	if it >= 0 {
		if _secondary == secondary {
			return it, primary
		}
	}
	return -1, 0
}

// int32_t db_idx_double_lowerbound(uint64_t code, uint64_t scope, uint64_t table, Double* secondary, uint64_t* primary);
func DBIdxFloat64Lowerbound(code uint64, scope uint64, table uint64, secondary float64) (int32, float64, uint64) {
	_secondary := (*[8]byte)(unsafe.Pointer(&secondary))
	var primary uint64 = 0
	ret, err := chaintester.GetVMAPI().DbIdxDoubleLowerbound(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), _secondary[:], to_raw_uint64(primary))
	CheckError(err)
	copy(_secondary[:], ret.Secondary)
	return ret.Iterator, secondary, from_raw_uint64(ret.Primary)
}

// int32_t db_idx_double_upperbound(uint64_t code, uint64_t scope, uint64_t table, Double* secondary, uint64_t* primary);
func DBIdxFloat64Upperbound(code uint64, scope uint64, table uint64, secondary float64) (int32, float64, uint64) {
	_secondary := (*[8]byte)(unsafe.Pointer(&secondary))
	var primary uint64 = 0
	ret, err := chaintester.GetVMAPI().DbIdxDoubleUpperbound(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), _secondary[:], to_raw_uint64(primary))
	CheckError(err)
	copy(_secondary[:], ret.Secondary)
	return ret.Iterator, secondary, from_raw_uint64(ret.Primary)
}

// int32_t db_idx_double_end(uint64_t code, uint64_t scope, uint64_t table);
func DBIdxFloat64End(code uint64, scope uint64, table uint64) int32 {
	ret, err := chaintester.GetVMAPI().DbIdxDoubleEnd(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table))
	CheckError(err)
	return ret
}

// int32_t db_idx_long_double_lowerbound(uint64_t code, uint64_t scope, uint64_t table, float128_t* secondary, uint64_t* primary);
func DBIdxFloat128Store(scope uint64, table uint64, id uint64, secondary [16]byte, payer uint64) int32 {
	ret, err := chaintester.GetVMAPI().DbIdxLongDoubleStore(ctx, to_raw_uint64(scope), to_raw_uint64(table), to_raw_uint64(payer), to_raw_uint64(id), secondary[:])
	CheckError(err)
	return ret
}

// int32_t db_idx_long_double_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const float128_t* secondary, uint64_t* primary);
func DBIdxFloat128Update(it int32, secondary [16]byte, payer uint64) {
	err := chaintester.GetVMAPI().DbIdxLongDoubleUpdate(ctx, it, to_raw_uint64(payer), secondary[:])
	CheckError(err)
}

// int32_t db_idx_long_double_find_primary(uint64_t code, uint64_t scope, uint64_t table, float128_t* secondary, uint64_t primary);
func DBIdxFloat128Remove(it int32) {
	err := chaintester.GetVMAPI().DbIdxLongDoubleRemove(ctx, it)
	CheckError(err)
}

// int32_t db_idx_long_double_previous(int32_t iterator, uint64_t* primary);
func DBIdxFloat128Next(it int32) (int32, uint64) {
	ret, err := chaintester.GetVMAPI().DbIdxLongDoubleNext(ctx, it)
	CheckError(err)
	return ret.Iterator, from_raw_uint64(ret.Primary)
}

// int32_t db_idx_long_double_next(int32_t iterator, uint64_t* primary);
func DBIdxFloat128Previous(it int32) (int32, uint64) {
	ret, err := chaintester.GetVMAPI().DbIdxLongDoublePrevious(ctx, it)
	CheckError(err)
	return ret.Iterator, from_raw_uint64(ret.Primary)
}

// void db_idx_long_double_remove(int32_t iterator);
func DBIdxFloat128FindByPrimary(code uint64, scope uint64, table uint64, primary uint64) (int32, [16]byte) {
	ret, err := chaintester.GetVMAPI().DbIdxLongDoubleFindPrimary(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), to_raw_uint64(primary))
	CheckError(err)
	secondary := [16]byte{}
	copy(secondary[:], ret.Secondary)
	return ret.Iterator, secondary
}

// int32_t db_idx_long_double_upperbound(uint64_t code, uint64_t scope, uint64_t table, float128_t* secondary, uint64_t* primary);
func DBIdxFloat128Find(code uint64, scope uint64, table uint64, secondary [16]byte) (int32, uint64) {
	it, _secondary, primary := DBIdxFloat128Lowerbound(code, scope, table, secondary)
	if it >= 0 {
		if _secondary == secondary {
			return it, primary
		}
	}
	return -1, 0
}

// void db_idx_long_double_update(int32_t iterator, uint64_t payer, const float128_t* secondary);
func DBIdxFloat128Lowerbound(code uint64, scope uint64, table uint64, secondary [16]byte) (int32, [16]byte, uint64) {
	var primary uint64 = 0
	ret, err := chaintester.GetVMAPI().DbIdxLongDoubleLowerbound(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), secondary[:], to_raw_uint64(primary))
	CheckError(err)
	copy(secondary[:], ret.Secondary)
	return ret.Iterator, secondary, from_raw_uint64(ret.Primary)
}

// int32_t db_idx_long_double_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const float128_t* secondary);
func DBIdxFloat128Upperbound(code uint64, scope uint64, table uint64, secondary [16]byte) (int32, [16]byte, uint64) {
	var primary uint64 = 0
	ret, err := chaintester.GetVMAPI().DbIdxLongDoubleUpperbound(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table), secondary[:], to_raw_uint64(primary))
	CheckError(err)
	copy(secondary[:], ret.Secondary)
	return ret.Iterator, secondary, from_raw_uint64(ret.Primary)
}

// int32_t db_idx_long_double_end(uint64_t code, uint64_t scope, uint64_t table);
func DBIdxFloat128End(code uint64, scope uint64, table uint64) int32 {
	ret, err := chaintester.GetVMAPI().DbIdxLongDoubleEnd(ctx, to_raw_uint64(code), to_raw_uint64(scope), to_raw_uint64(table))
	CheckError(err)
	return ret
}

// void send_deferred(const uint128_t* sender_id, capi_name payer, const char *serialized_transaction, size_t size, uint32_t replace_existing);
func SendDeferred(senderID [16]byte, payer uint64, transaction []byte, replaceExisting bool) {
	var _replaceExisting int32
	if replaceExisting {
		_replaceExisting = 1
	} else {
		_replaceExisting = 0
	}
	err := chaintester.GetVMAPI().SendDeferred(ctx, senderID[:], to_raw_uint64(payer), transaction, _replaceExisting)
	CheckError(err)
}

// int cancel_deferred(const uint128_t* sender_id);
func CancelDeferred(senderID [16]byte) int32 {
	ret, err := chaintester.GetVMAPI().CancelDeferred(ctx, senderID[:])
	CheckError(err)
	return ret
}

// size_t read_transaction(char *buffer, size_t size);
func ReadTransaction() []byte {
	ret, err := chaintester.GetVMAPI().ReadTransaction(ctx)
	CheckError(err)
	return ret
}

// __attribute__((eosio_wasm_import))
// size_t transaction_size( void );
func TransactionSize() int32 {
	ret, err := chaintester.GetVMAPI().TransactionSize(ctx)
	CheckError(err)
	return ret
}

// int tapos_block_num( void );
func TaposBlockNum() int32 {
	ret, err := chaintester.GetVMAPI().TaposBlockNum(ctx)
	CheckError(err)
	return ret
}

// int tapos_block_prefix( void );
func TaposBlockPrefix() int32 {
	ret, err := chaintester.GetVMAPI().TaposBlockPrefix(ctx)
	CheckError(err)
	return ret
}

// uint32_t expiration( void );
func Expiration() uint32 {
	ret, err := chaintester.GetVMAPI().Expiration((ctx))
	CheckError(err)
	return uint32(ret)
}

// int get_action( uint32_t type, uint32_t index, char* buff, size_t size );
func GetAction(_type uint32, index uint32) []byte {
	ret, err := chaintester.GetVMAPI().GetAction(ctx, int32(_type), int32(index))
	CheckError(err)
	return ret
}

// int get_context_free_data( uint32_t index, char* buff, size_t size );
func GetContextFreeData(index uint32) []byte {
	ret, err := chaintester.GetVMAPI().GetContextFreeData(ctx, int32(index))
	CheckError(err)
	return ret
}

func GetActiveProducers() []uint64 {
	ret, err := chaintester.GetVMAPI().GetActiveProducers(ctx)
	producers := make([]uint64, 0, len(ret)/8)
	for i := 0; i < len(ret)/8; i += 8 {
		producers = append(producers, binary.LittleEndian.Uint64(ret[i:i+8]))
	}
	CheckError(err)
	return producers
}
