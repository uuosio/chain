package main

import (
	"github.com/uuosio/chain"
	"github.com/uuosio/chain/logger"
)

func main() {
	a, b, c := chain.GetResourceLimits(chain.NewName("hello"))
	logger.Println(a, b, c)
}
