package main

import (
	"chain"
	"chain/logger"
	"encoding/hex"
)

func main() {
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

	{
		signature := []byte("\x00\x1f*\x9ef\xe8\xfb\xbf\xf1\x1c_\x88\xa5\xa5\xcb\xacz\x07\x90xk\x06\xbd?\xb2\\\xa6g\xcc,*\xe5G\x01\x04[g6o\x1e\xdefuU\xb4\xf2\xeaJ\xab\x01\x08\t\x0e={~\x81\xed\x1b\x0b\x06\xf5o\xa6\x19\xeb")
		hash := chain.Sha256([]byte("hello,world"))
		pub := chain.RecoverKey(hash, signature)
		logger.Println(hex.EncodeToString(pub))
	}
}
