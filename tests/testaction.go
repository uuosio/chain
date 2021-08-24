package main

import (
	"chain"
)

var gContractName = chain.NewName("hello")
var gActionName = chain.NewName("sayhello2")

/*
   def get_secondary_values(self):
       return (self.a, self.b, self.c, self.d)

   @staticmethod
   def get_secondary_indexes():
       return (db.idx64, db.idx128, db.idx256, db.idx_double)
*/

type uint128 [16]byte
type uint256 [32]byte

func (n *uint128) Pack() []byte {
	return n[:]
}

func IsEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	// for i := 0; i < 10; i++ {
	// 	fmt.Printf("hello,world! %d\n", i*i)
	// }
	// return

	{
		// n := gContractName
		// chain.Prints(chain.N2S(n))
	}
	// a := ReadActionData()
	// prints(string(a))

	_, _, action := chain.GetApplyArgs()
	if action == gActionName {
		a := chain.ReadActionData()
		chain.Println("sayhello2 received:", string(a))
		return
	} else {
		chain.Prints("hello,world\n")
	}

	{
		a := chain.Action{
			chain.NewName("hello"), //gContractName,
			gActionName,
			[]chain.PermissionLevel{{gContractName, chain.ActiveName}},
			[]byte("hello,world"),
		}
		a.Send()
	}

	{
		a := chain.Action{}
		a.Account = gContractName
		a.Name = gActionName
		a.AddPermission(gContractName, chain.ActiveName)
		a.Data = []byte("hello,worldddd")
		a.Send()
	}
	// n := chain.Name{}
	// b := n.Unpack([]byte{})
	//	C.printui(b.N)
}
