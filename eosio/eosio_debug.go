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

//Read current action data
func ReadActionData() []byte {
	return nil
}

//Get the length of the current action's data field
func ActionDataSize() uint32 {
	return 0
}

//Add the specified account to set of accounts to be notified
func RequireRecipient(name uint64) {
}

func RequireRecipientEx(name uint64) {
}

//Verifies that name exists in the set of provided auths on a action. Throws if not found.
func RequireAuth(name uint64) {
}

//Verifies that name has auth.
func HasAuth(name uint64) bool {
	return false
}

//Verifies that name exists in the set of provided auths on a action. Throws if not found.
func RequireAuth2(name uint64, permission uint64) {
}

//Verifies that name is an existing account.
func IsAccount(name uint64) bool {
	return false
}

//Send an inline action in the context of this action's parent transaction
func SendInline(data []byte) {
	chaintester.GetVMAPI().SendInline(ctx, data)
}

//Send an inline context free action in the context of this action's parent transaction
func SendContextFreeInline(data []byte) {
}

//Returns the time in microseconds from 1970 of the publication_time
func PublicationTime() uint64 {
	return 0
}

//Get the current receiver of the action
func CurrentReceiver() uint64 {
	return 0
}

//Set the action return value which will be included in the action_receipt
func SetActionReturnValue(return_value []byte) {

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
