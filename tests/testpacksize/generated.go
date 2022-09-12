
package testpacksize
import (
	"github.com/uuosio/chain"
    "github.com/uuosio/chain/database"
    "unsafe"
)


func (t *B) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.PackUint32(t.a)
	enc.PackUint64(t.b)
    return enc.GetBytes()
}

func (t *B) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	t.a = dec.UnpackUint32()
	t.b = dec.UnpackUint64()
    return dec.Pos()
}

func (t *B) Size() int {
    size := 0
	size += 4 //a
	size += 8 //b
    return size
}


func (t *A) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.PackUint32(t.a)
	enc.PackString(t.b)
	enc.PackBytes(t.c)
	
	{
		enc.PackLength(len(t.e))
		for i := range t.e {
			enc.PackString(t.e[i])
		}
	}
	enc.Pack(&t.g)
	
	{
		enc.PackLength(len(t.h))
		for i := range t.h {
			enc.Pack(&t.h[i])
		}
	}
    return enc.GetBytes()
}

func (t *A) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	t.a = dec.UnpackUint32()
	t.b = dec.UnpackString()
	t.c = dec.UnpackBytes()
	
	{
		length := dec.UnpackLength()
		t.e = make([]string, length)
		for i:=0; i<length; i++ {
		t.e[i] = dec.UnpackString()
		}
	}
	dec.UnpackI(&t.g)
	
	{
		length := dec.UnpackLength()
		t.h = make([]B, length)
		for i:=0; i<length; i++ {
		dec.UnpackI(&t.h[i])
		}
	}
    return dec.Pos()
}

func (t *A) Size() int {
    size := 0
	size += 4 //a
	size += chain.PackedVarUint32Length(uint32(len(t.b))) + len(t.b) //b
	size += chain.PackedVarUint32Length(uint32(len(t.c)))
    size += len(t.c)
	size += chain.PackedVarUint32Length(uint32(len(t.e)))
    for i := range t.e {
	 size += chain.PackedVarUint32Length(uint32(len(t.e[i]))) + len(t.e[i])
}
	size += t.g.Size() //g
	size += chain.PackedVarUint32Length(uint32(len(t.h)))
    
    for i := range t.h {
        size += t.h[i].Size()
    }
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