package chain

import (
	"github.com/uuosio/chain/eosio"
)

func Check(b bool, msg string) {
	if !b {
		eosio.EosioAssert(false, msg)
	}
}

//Aborts processing of this action and unwinds all pending changes if the test condition is true
func Assert(test bool, msg string) {
	if !test {
		eosio.EosioAssert(false, msg)
	}
}

//Aborts processing of this action and unwinds all pending changes if the test condition is true
func EosioAssert(test bool, msg string) {
	if !test {
		eosio.EosioAssert(test, msg)
	}
}

//Aborts processing of this action and unwinds all pending changes if the test condition is true
func EosioAssertCode(test bool, code uint64) {
	if !test {
		eosio.EosioAssertCode(test, code)
	}
}

//Returns the time in microseconds from 1970 of the current block
func CurrentTime() TimePoint {
	return TimePoint{uint64(eosio.CurrentTime())}
}

//Returns the time in microseconds from 1970 of the current block
func Now() TimePoint {
	return CurrentTime()
}

func NowSeconds() uint32 {
	t := CurrentTime().Elapsed / 1000000
	return uint32(t)
}

func CurrentTimeSeconds() uint32 {
	t := CurrentTime().Elapsed / 1000000
	return uint32(t)
}

//Check if specified protocol feature has been activated
func IsFeatureActivated(featureDigest [32]byte) bool {
	return eosio.IsFeatureActivated(featureDigest)
}

//Return name of account that sent current inline action
func GetSender() uint64 {
	return uint64(eosio.GetSender())
}

func Exit() {
	eosio.Exit()
}
