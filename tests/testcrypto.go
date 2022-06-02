package main

import (
	"github.com/uuosio/chain/hex"

	"github.com/uuosio/chain"
	"github.com/uuosio/chain/logger"
)

var gContractName = chain.NewName("hello")
var gActionName = chain.NewName("sayhello2")

//contract test
type CryptoTest struct {
	self   chain.Name
	code   chain.Name
	action chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *CryptoTest {
	return &CryptoTest{
		self:   receiver,
		code:   firstReceiver,
		action: action,
	}
}

//action testhash
func (c *CryptoTest) TestHash() {
	producers := chain.GetActiveProducers()
	logger.Println(producers)
	hash := chain.Sha1([]byte("hello"))
	logger.Println(hex.EncodeToString(hash[:]))

	helloworld := []byte("hello,world")
	{
		hash := chain.Sha256(helloworld)
		chain.AssertSha256(helloworld, hash)
	}

	{
		hash := chain.Sha1(helloworld)
		chain.AssertSha1(helloworld, hash)
	}

	{
		hash := chain.Sha512(helloworld)
		chain.AssertSha512(helloworld, hash)
	}

	{
		hash := chain.Ripemd160(helloworld)
		chain.AssertRipemd160(helloworld, hash)
	}
}

//action testrecover
func (c *CryptoTest) TestRecover(data []byte, sig *chain.Signature, pub *chain.PublicKey) {
	hash := chain.Sha256([]byte("hello,world"))
	pub2 := chain.RecoverKey(hash, sig)
	chain.Check(*pub == *pub2, "bad recovery")
}
