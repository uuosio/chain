package main

import (
	"chain"
	"chain/logger"
)

type B struct {
	a int
	b uint64
}

type A struct {
	a int
	b string
	c []byte
	d [][]byte
	e []string
	// f [][][]byte
	g B
	h []B
}

func main() {
	//128 bytes
	text := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	a := A{}
	a.a = 1
	a.b = text
	a.c = []byte(text)
	a.d = [][]byte{[]byte(text)}
	a.e = []string{text}
	a.g = B{}
	a.h = []B{B{}}

	size := 4                 //a.a int
	size += 2 + len(a.b)      //a.b string
	size += 2 + len(a.c)      //a.c []byte
	size += 1 + 2 + len(text) //a.d [][]byte
	size += 1 + 2 + len(text) //a.e []string
	size += 12                //a.g B
	size += 1 + 12            //a.h []B
	// chain.Check(a.Size() == size, "bad size")
	logger.Println(a.Size(), size)
	logger.Println("Hello,world!")
	chain.Check(true, "")
}
