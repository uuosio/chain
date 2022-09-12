
package testdb
import (
	"github.com/uuosio/chain"
    "github.com/uuosio/chain/database"
    "unsafe"
)


func (t *MyData) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.PackUint64(t.primary)
	enc.PackUint32(t.n)
    return enc.GetBytes()
}

func (t *MyData) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	t.primary = dec.UnpackUint64()
	t.n = dec.UnpackUint32()
    return dec.Pos()
}

func (t *MyData) Size() int {
    size := 0
	size += 8 //primary
	size += 4 //n
    return size
}


//eliminate unused package errors
func dummy() {
	if false {
		v := 0;
		n := unsafe.Sizeof(v);
		chain.Printui(uint64(n));
		chain.Printui(database.IDX64);
	}
}