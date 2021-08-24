package chain

type Template struct{}

func (t *Template) Pack() []byte {
	// enc := chain.NewEncoder(10)
	// return enc.GetBytes()
	return nil
}

func (t *Template) Unpack(data []byte) (int, error) {
	// dec := chain.NewDecoder(data)
	// return dec.Pos(), nil
	return 0, nil
}

func (t *Template) Print() {
}
