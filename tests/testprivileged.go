package main

import (
	"chain"
	"chain/logger"
)

func main() {
	a, b, c := chain.GetResourceLimits(chain.NewName("hello"))
	logger.Println(a, b, c)
}
