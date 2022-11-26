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

func N(s string) Name {
	return Name{N: S2N(s)}
}

func (a *Name) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.WriteUint64(a.N)
	return enc.GetSize() - oldSize
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

type NameList []Name

func (a NameList) Len() int { return len(a) }
func (a NameList) Less(i, j int) bool {
	return a[i].N < a[j].N
}

func (a NameList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
