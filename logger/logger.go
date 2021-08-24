package logger

import "chain"

type Printable interface {
	Print()
}

type ExtendedPrintFunc func(value interface{})

var (
	ExtendedPrint ExtendedPrintFunc
)

func Print(args ...interface{}) {
	chain.Print(args...)
}

func Println(args ...interface{}) {
	chain.Println(args...)
}

// const DEBUG = false

// comment out Printf for reduce build binary size
// func Printf(format string, a ...interface{}) {
// 	if !chain.DEBUG {
// 		return
// 	}
// 	s := fmt.Sprintf(format, a...)
// 	chain.Prints(s)
// }

func Fatal(err interface{}) {
	switch v := err.(type) {
	case error:
		panic(v.Error())
	case string:
		panic(v)
	default:
		break
	}
}
