//go:build !tinygo.wasm
// +build !tinygo.wasm

package eosio

import (
	"github.com/learnforpractice/chaintester"
)

//Prints string
func Prints(str string) {
	chaintester.GetVMAPI().Prints(ctx, str)
}

//Prints value as a 64 bit signed integer
func Printi(value int64) {
}

//Prints value as a 64 bit unsigned integer
func PrintUi(value uint64) {
}

func Printui(value uint64) {
}

//Prints value as a 128 bit signed integer
func PrintI128(value [16]byte) {
}

//Prints value as a 128 bit unsigned integer
func PrintUi128(value [16]byte) {
}

//Prints value as single-precision floating point number
func PrintSf(value float32) {
}

//Prints value as double-precision floating point number
func PrintDf(value float64) {
}

//Prints value as quadruple-precision floating point number
func PrintQf(value [16]byte) {
}

//Prints a 64 bit names as base32 encoded string
func PrintN(name uint64) {
}

//Prints hexidecimal data of length datalen
func PrintHex(data []byte) {
}
