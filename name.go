package chain

type Name struct {
	N uint64
}

//called in compiler
func newname(n uint64) Name {
	return Name{N: uint64(n)}
}

func NewName(s string) Name {
	return Name{N: S2N(s)}
}

func (a *Name) Pack() []byte {
	enc := NewEncoder(8)
	enc.WriteUint64(a.N)
	return enc.GetBytes()
}

func (a *Name) Unpack(data []byte) int {
	dec := NewDecoder(data)
	n := dec.UnpackUint64()
	a.N = n
	return 8
}

func (t *Name) Size() int {
	return 8
}

func (a *Name) String() string {
	return N2S(a.N)
}
