//go:build !tinygo.wasm
// +build !tinygo.wasm

package eosio

import (
	"unsafe"

	"github.com/uuosio/chaintester"
)

//Prints string
func Prints(str string) {
	err := chaintester.GetVMAPI().Prints(ctx, str)
	CheckError(err)
}

//Prints value as a 64 bit signed integer
func Printi(value int64) {
	err := chaintester.GetVMAPI().Printi(ctx, value)
	CheckError(err)
}

//Prints value as a 64 bit unsigned integer
func PrintUi(value uint64) {
	err := chaintester.GetVMAPI().Printui(ctx, to_raw_uint64(value))
	CheckError(err)
}

func Printui(value uint64) {
	err := chaintester.GetVMAPI().Printui(ctx, to_raw_uint64(value))
	CheckError(err)
}

//Prints value as a 128 bit signed integer
func PrintI128(value [16]byte) {
	err := chaintester.GetVMAPI().Printi128(ctx, value[:])
	CheckError(err)
}

//Prints value as a 128 bit unsigned integer
func PrintUi128(value [16]byte) {
	err := chaintester.GetVMAPI().Printui128(ctx, value[:])
	CheckError(err)
}

//Prints value as single-precision floating point number
func PrintSf(value float32) {
	rawValue := (*[4]byte)(unsafe.Pointer(&value))
	err := chaintester.GetVMAPI().Printsf(ctx, rawValue[:])
	CheckError(err)
}

//Prints value as double-precision floating point number
func PrintDf(value float64) {
	rawValue := (*[8]byte)(unsafe.Pointer(&value))
	err := chaintester.GetVMAPI().Printdf(ctx, rawValue[:])
	CheckError(err)
}

//Prints value as quadruple-precision floating point number
func PrintQf(value [16]byte) {
	err := chaintester.GetVMAPI().Printqf(ctx, value[:])
	CheckError(err)
}

//Prints a 64 bit names as base32 encoded string
func PrintN(name uint64) {
	err := chaintester.GetVMAPI().Printn(ctx, to_raw_uint64(name))
	CheckError(err)
}

//Prints hexidecimal data of length datalen
func PrintHex(data []byte) {
	err := chaintester.GetVMAPI().Printhex(ctx, data)
	CheckError(err)
}
