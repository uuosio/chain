//go:build !tinygo.wasm
// +build !tinygo.wasm

package eosio

import "C"
import (
	"context"
	"encoding/binary"

	"github.com/learnforpractice/chaintester"
	"github.com/learnforpractice/chaintester/interfaces"
)

var ctx = context.Background()

type AssertError struct {
	Err error
}

func (err *AssertError) Error() string {
	return err.Err.Error()
}

func NewAssertError(err error) *AssertError {
	return &AssertError{err}
}

func CheckError(err error) {
	if err != nil {
		panic(NewAssertError(err))
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
