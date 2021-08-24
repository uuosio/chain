package chain

type Name struct {
	N uint64
}

func NewName(s string) Name {
	return Name{N: S2N(s)}
}

func (a *Name) Pack() []byte {
	enc := NewEncoder(8)
	enc.WriteUint64(a.N)
	return enc.GetBytes()
}

func (a *Name) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	n, err := dec.UnpackUint64()
	if err != nil {
		return 0, err
	}
	a.N = n
	return 8, nil
}

func (a *Name) String() string {
	return N2S(a.N)
}
