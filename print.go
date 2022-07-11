package chain

import (
	"github.com/uuosio/chain/eosio"
)

//Prints string
func Prints(str string) {
	eosio.Prints(str)
}

//Prints value as a 64 bit signed integer
func Printi(value int64) {
	eosio.Printi(value)
}

//Prints value as a 64 bit unsigned integer
func PrintUi(value uint64) {
	eosio.PrintUi(value)
}

func Printui(value uint64) {
	eosio.Printui(value)
}

//Prints value as a 128 bit signed integer
func PrintI128(value [16]byte) {
	eosio.PrintI128(value)
}

//Prints value as a 128 bit unsigned integer
func PrintUi128(value [16]byte) {
	eosio.PrintUi128(value)
}

//Prints value as single-precision floating point number
func PrintSf(value float32) {
	eosio.PrintSf(value)
}

//Prints value as double-precision floating point number
func PrintDf(value float64) {
	eosio.PrintDf(value)
}

//Prints value as quadruple-precision floating point number
func PrintQf(value Float128) {
	eosio.PrintQf(value)
}

//Prints a 64 bit names as base32 encoded string
func PrintN(name uint64) {
	eosio.PrintN(name)
}

//Prints hexidecimal data of length datalen
func PrintHex(data []byte) {
	eosio.PrintHex(data)
}

type Printable interface {
	Print()
}

type ExtendedPrintFunc func(value interface{})

var (
	ExtendedPrint ExtendedPrintFunc
)

func PrintVariant(variant interface{}) {
	switch v := variant.(type) {
	case nil:
		Prints("nil")
	case Printable:
		v.Print()
	case error:
		Prints(v.Error())
	case string:
		Prints(v)
	case bool:
		if v {
			Prints("true")
		} else {
			Prints("false")
		}
	case int8:
		Printi(int64(v))
	case uint8:
		PrintUi(uint64(v))
	case int16:
		Printi(int64(v))
	case uint16:
		PrintUi(uint64(v))
	case int32:
		Printi(int64(v))
	case uint32:
		PrintUi(uint64(v))
	case int:
		Printi(int64(v))
	case int64:
		Printi(v)
	case uint64:
		PrintUi(v)
	case [16]byte:
		PrintI128(v)
	// case [16]byte:
	// 	PrintUi128(v)
	case float32:
		PrintSf(v)
	case float64:
		PrintDf(v)
	case Float128:
		PrintQf(v)
	// case [16]byte:
	// 	PrintQf(v)
	case Name:
		PrintN(v.N)
	case []Name:
		for _, n := range v {
			PrintN(n.N)
		}
	case Symbol:
		v.Print()
	case []byte:
		PrintHex(v)
	default:
		// if DEBUG {
		// 	s := fmt.Sprintf("%v", v)
		// 	Prints(s)
		// }
		Prints("<unprintable>")
	}
}

func PrintNoEndSpace(args ...interface{}) {
	for _, arg := range args {
		PrintVariant(arg)
	}
}

func Print(args ...interface{}) {
	for i, v := range args {
		PrintVariant(v)
		if i < len(args)-1 {
			PrintVariant(" ")
		}
	}
}

func Println(args ...interface{}) {
	Print(args...)
	Print("\n")
}
