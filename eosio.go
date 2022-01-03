//go:build eosio
// +build eosio

package chain

import (
	"runtime"
)

func GetApplyArgs() (Name, Name, Name) {
	receiver, code, action := runtime.GetApplyArgs()
	return Name{receiver}, Name{code}, Name{action}
}
