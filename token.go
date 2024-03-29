package chain

const MAX_AMOUNT = (1 << 62) - 1

type SymbolCode struct {
	Value uint64
}

func NewSymbolCode(sym string) SymbolCode {
	n := SymbolCode{}
	for i := range sym {
		n.Value <<= 8
		n.Value |= uint64(sym[len(sym)-i-1])
	}
	return n
}

func (a *SymbolCode) IsValid() bool {
	sym := a.Value

	if sym == 0 {
		return false
	}

	if sym>>56 != 0 {
		return false
	}

	for i := 0; i < 7; i++ {
		c := byte(sym & 0xFF)
		if !('A' <= c && c <= 'Z') {
			return false
		}
		sym >>= 8
		if sym&0xFF != 0 {
			continue
		}
		for ; i < 7; i++ {
			sym >>= 8
			if sym&0xFF != 0 {
				return false
			}
		}
	}
	return true
}

func (a *SymbolCode) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.PackUint64(a.Value)
	return enc.GetSize() - oldSize
}

func (a *SymbolCode) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Unpack(&a.Value)
	return dec.Pos()
}

func (t *SymbolCode) Size() int {
	return 8
}

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

func (a *Symbol) Precision() uint64 {
	return a.Value & 0xff
}

func (a *Symbol) IsValid() bool {
	sym := a.Code()
	for i := 0; i < 7; i++ {
		c := byte(sym & 0xFF)
		if !('A' <= c && c <= 'Z') {
			return false
		}
		sym >>= 8
		if sym&0xFF != 0 {
			continue
		}
		for ; i < 7; i++ {
			sym >>= 8
			if sym&0xFF != 0 {
				return false
			}
		}
	}
	return true
}

func (a *Symbol) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.PackUint64(a.Value)
	return enc.GetSize() - oldSize
}

func (a *Symbol) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Unpack(&a.Value)
	return dec.Pos()
}

func (t *Symbol) Size() int {
	return 8
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

func isAmountWithInRange(amount int64) bool {
	return -MAX_AMOUNT <= amount && amount <= MAX_AMOUNT
}

func NewAsset(amount int64, symbol Symbol) *Asset {
	Check(symbol.IsValid(), "bad symbol")
	a := &Asset{amount, symbol}
	Check(isAmountWithInRange(amount), "magnitude of asset amount must be less than 2^62")
	return a
}

func (a *Asset) Add(b *Asset) *Asset {
	Check(a.Symbol == b.Symbol, "Asset.Add:Symbol not the same")
	a.Amount += b.Amount
	Check(-MAX_AMOUNT <= a.Amount, "addition underflow")
	Check(a.Amount <= MAX_AMOUNT, "addition overflow")
	return a
}

func (a *Asset) Sub(b *Asset) *Asset {
	Check(a.Symbol == b.Symbol, "Asset.Sub:Symbol not the same")
	a.Amount -= b.Amount
	Check(a.Amount >= -MAX_AMOUNT, "subtraction underflow")
	Check(a.Amount <= MAX_AMOUNT, "subtraction overflow")
	return a
}

func (a *Asset) Mul(b *Asset) *Asset {
	Check(a.Symbol == b.Symbol, "Asset.Mul:Symbol not the same")
	_a := NewInt128FromInt64(a.Amount)
	_b := NewInt128FromInt64(b.Amount)
	_a.Mul(&_b)

	m := NewInt128FromInt64(MAX_AMOUNT)
	Check(m.Cmp(&_a) >= 0, "multiplication overflow")

	m = NewInt128FromInt64(-MAX_AMOUNT)
	Check(_a.Cmp(&m) >= 0, "multiplication underflow")
	a.Amount = _a.Int64()
	return a
}

func (a *Asset) Div(b *Asset) *Asset {
	Check(a.Symbol == b.Symbol, "Asset.Mul:Symbol not the same")
	Check(b.Amount != 0, "divide by zero")
	Check(b.Amount > 0, "divide by negative value")
	// Check(!(a.Amount == int64(-9223372036854775808) && b.Amount == -1), "signed division overflow")
	a.Amount /= b.Amount
	return a
}

func (a *Asset) IsValid() bool {
	return isAmountWithInRange(a.Amount) && a.Symbol.IsValid()
}

func (a *Asset) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.WriteUint64(uint64(a.Amount))
	enc.WriteUint64(a.Symbol.Value)
	return enc.GetSize() - oldSize
}

func (a *Asset) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Unpack(&a.Amount)
	dec.Unpack(&a.Symbol)
	return 16
}

func (t *Asset) Size() int {
	return 16
}

type ExtendedAsset struct {
	Quantity Asset
	Contract Name
}

func NewExtendedAsset(quantity *Asset, contract Name) *ExtendedAsset {
	return &ExtendedAsset{*quantity, contract}
}

func (t *ExtendedAsset) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	t.Quantity.Pack(enc)
	t.Contract.Pack(enc)
	return enc.GetSize() - oldSize
}

func (t *ExtendedAsset) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Unpack(&t.Quantity)
	dec.Unpack(&t.Contract)
	return dec.Pos()
}

func (t *ExtendedAsset) Size() int {
	return 16 + 8
}

type Transfer struct {
	From     Name
	To       Name
	Quantity Asset
	Memo     string
}

func (a *Transfer) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.PackName(a.From)
	enc.PackName(a.To)
	a.Quantity.Pack(enc)
	enc.PackString(a.Memo)
	return enc.GetSize() - oldSize
}

func (a *Transfer) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Unpack(&a.From)
	dec.Unpack(&a.To)
	dec.Unpack(&a.Quantity)
	dec.Unpack(&a.Memo)
	return dec.Pos()
}

func (t *Transfer) Size() int {
	size := 0
	size += 8                                                        //From
	size += 8                                                        //To
	size += t.Quantity.Size()                                        //Quantity
	size += PackedVarUint32Length(uint32(len(t.Memo))) + len(t.Memo) //Memo
	return size
}
