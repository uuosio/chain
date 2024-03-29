package chain

type Template struct{}

func (t *Template) Pack(enc *Encoder) int {
	// enc := chain.NewEncoder(10)
	// return enc.GetBytes()
	return 0
}

func (t *Template) Unpack(data []byte) int {
	// dec := chain.NewDecoder(data)
	// return dec.Pos()
	return 0
}

func (t *Template) Size() int {
	return 0
}

func (t *Template) Print() {
}

//table template
type TableTemplate struct {
	Id      uint64
	Account Name
}

var (
	TableTemplateSecondaryTypes = []int{}
)

func (t *TableTemplate) GetPrimary() uint64 {
	return t.Id
}

func (t *TableTemplate) GetSecondaryValue(index int) interface{} {
	switch index {
	case 0:
		return t.Account.N
	default:
		panic("unknown index")
	}
}

func (t *TableTemplate) SetSecondaryValue(index int, v interface{}) {
	switch index {
	case 0:
		t.Account = v.(Name)
	default:
		panic("unknown index")
	}
}
