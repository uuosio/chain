package main

import (
	"bytes"

	"github.com/uuosio/chain"
)

//wrap chain.Check
func check(cond bool, msg string) {
	chain.Check(cond, msg)
}

var gContractName = chain.NewName("hello")
var gActionName = chain.NewName("sayhello2")

//contract test
type SerializerTest struct {
	self   chain.Name
	code   chain.Name
	action chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *SerializerTest {
	return &SerializerTest{
		self:   receiver,
		code:   firstReceiver,
		action: action,
	}
}

//action test
func (c *SerializerTest) test(
	a0 bool,
	a1 int8,
	a2 uint8,
	a3 int16,
	a4 uint16,
	a5 int32,
	a6 uint32,
	a7 int64,
	a8 uint64,
	a9 chain.Int128,
	a10 chain.Uint128,
	a11 chain.VarInt32,
	a12 chain.VarUint32,
	a13 float32,
	a14 float64,
	a15 chain.Float128,
	a16 chain.TimePoint,
	a17 chain.TimePointSec,
	a18 chain.BlockTimestampType, //block_timestamp_type,
	a19 chain.Name,
	a20 []byte, //bytes,
	a21 string,
	a22 chain.Checksum160, //checksum160,
	a23 chain.Checksum256, //checksum256,
	a24 chain.Checksum512, //checksum512,
	a25 chain.PublicKey, //public_key,
	a26 chain.Signature, //signature,
	a27 chain.Symbol, //symbol,
	a28 chain.SymbolCode, //symbol_code,
	a29 chain.Asset,
	a30 chain.ExtendedAsset,
) {
	//h := chain.Sha256(data)
	check(a0 == true, "a0")
	check(a1 == int8(-1), "a1")
	check(a2 == uint8(0xff), "a1")
	check(a3 == int16(-1), "a3")
	check(a4 == uint16(0xffff), "a4")
	check(a5 == int32(-1), "a5")
	check(a6 == uint32(0xffffffff), "a6")
	check(a7 == int64(-1), "a7")
	check(a8 == uint64(0xffffffffffffffff), "a8")

	for i := range a9[:15] {
		check(a9[i] == byte(0xff), "a9")
	}
	check(a9[15] == byte(0x7f), "a9")

	for i := range a10 {
		check(a10[i] == byte(0xff), "a10")
	}
	check(a11 == chain.VarInt32(int32(-1)), "a11")
	check(a12 == chain.VarUint32(0xffffffff), "a12")
	check(a13 == 11.2233, "a13")
	check(a14 == 11.2233, "a14")
	chain.Println(a13, a14)

	check(a15[0] == byte(0x7f), "a15")
	for _, v := range a15[1:] {
		check(v == byte(0xff), "a15")
	}

	check(a16 == chain.TimePoint{1630642401 * 1000000}, "a16")
	check(a17 == chain.TimePointSec{1630642401}, "a17")
	check(a18 == chain.BlockTimestampType{193723200}, "a18")
	check(a19 == chain.NewName("helloworld"), "a19")
	check(bytes.Equal(a20, []byte("hello,world")), "a20")
	check(a21 == "hello,world", "a21")

	for _, a := range a22 {
		check(a == 0xaa, "a22")
	}
	for _, a := range a23 {
		check(a == 0xaa, "a23")
	}
	for _, a := range a24 {
		check(a == 0xaa, "a24")
	}
	chain.Println(a27.Value, chain.NewSymbol("EOS", 4).Value)
	chain.Println(a28.Value, chain.NewSymbolCode("EOS").Value)
	check(a27 == chain.NewSymbol("EOS", 4), "a27")
	check(a28 == chain.NewSymbolCode("EOS"), "a28")
	check(a29 == *chain.NewAsset(10000, chain.NewSymbol("EOS", 4)), "a29")
	check(a30 == *chain.NewExtendedAsset(*chain.NewAsset(10000, chain.NewSymbol("EOS", 4)), chain.NewName("eosio.token")), "a30")
}

//action testvarint
func (c *SerializerTest) testvarint(a1 chain.VarInt32, a2 chain.VarInt32) {
	chain.Println(int32(a1), int32(a2))
	check(a1 == chain.VarInt32(-1), "a1")
	check(a2 == chain.VarInt32(0x7fffffff), "a2")
}

//action testpack
func (c *SerializerTest) testpack() {
	a := chain.Asset{}
	enc := chain.NewEncoder(100)
	enc.Pack(a)
	chain.Check(len(enc.GetBytes()) == 16, "bad value")
}
