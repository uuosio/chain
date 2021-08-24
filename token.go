package chain

type Symbol struct {
	Value uint64
}

func NewSymbol(name string, precision int) Symbol {
	Check(len(name) <= 7, "bad symbol name")
	value := uint64(0)
	for i := range name {
		v := name[len(name)-1-i]
		value |= uint64(v)
		value <<= 8
	}
	value |= uint64(precision)
	return Symbol{value}
}

func (a *Symbol) Code() uint64 {
	return a.Value >> 8
}

func (a *Symbol) Pack() []byte {
	enc := NewEncoder(8)
	enc.Pack(a.Value)
	return enc.GetBytes()
}

func (a *Symbol) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	dec.Unpack(&a.Value)
	return dec.Pos(), nil
}

func (a *Symbol) Print() {
	buf := [7]byte{}
	n := 0
	value := a.Value
	for {
		value >>= 8
		if value <= 0 {
			break
		}
		buf[n] = byte(value)
		n += 1
	}
	PrintNoEndSpace(a.Value&0xff, ",", string(buf[:n]))
}

type Asset struct {
	Amount int64
	Symbol Symbol
}

func (a *Asset) Pack() []byte {
	enc := NewEncoder(16)
	enc.WriteUint64(uint64(a.Amount))
	enc.WriteUint64(a.Symbol.Value)
	return enc.GetBytes()
}

func (a *Asset) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	dec.Unpack(&a.Amount)
	dec.Unpack(&a.Symbol)
	return 16, nil
}

type Transfer struct {
	From     Name
	To       Name
	Quantity Asset
	Memo     string
}

func (a *Transfer) Pack() []byte {
	enc := NewEncoder(8 + 8 + 16 + len(a.Memo) + 5)
	enc.Pack(&a.From)
	enc.Pack(&a.To)
	enc.Pack(&a.Quantity)
	enc.PackString(a.Memo)
	return enc.GetBytes()
}

func (a *Transfer) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	dec.Unpack(&a.From)
	dec.Unpack(&a.To)
	dec.Unpack(&a.Quantity)
	dec.Unpack(&a.Memo)
	return dec.Pos(), nil
}
