package chain

import "math/big"

const MAX_AMOUNT = (1 << 62) - 1

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
	_a := big.NewInt(a.Amount)
	_b := big.NewInt(b.Amount)
	_z := big.NewInt(0)
	_z.Mul(_a, _b)

	m := big.NewInt(MAX_AMOUNT)
	Check(m.Cmp(_z) >= 0, "multiplication overflow")

	m = big.NewInt(-MAX_AMOUNT)
	Check(_z.Cmp(m) >= 0, "multiplication underflow")
	a.Amount = _z.Int64()
	return a
}

func (a *Asset) Div(b *Asset) *Asset {
	Check(a.Symbol == b.Symbol, "Asset.Mul:Symbol not the same")
	Check(b.Amount != 0, "divide by zero")
	Check(!(a.Amount == int64(-9223372036854775808) && b.Amount == -1), "signed division overflow")
	a.Amount /= b.Amount
	return a
}

func (a *Asset) IsValid() bool {
	return isAmountWithInRange(a.Amount) && a.Symbol.IsValid()
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

func (t *Asset) Size() int {
	return 16
}

type ExtendedAsset struct {
	Quantity Asset
	Contract Name
}

func NewExtendedAsset(quantity Asset, contract Name) *ExtendedAsset {
	return &ExtendedAsset{quantity, contract}
}

func (t *ExtendedAsset) Pack() []byte {
	enc := NewEncoder(16 + 8)
	enc.Pack(&t.Quantity)
	enc.PackName(t.Contract)
	return enc.GetBytes()
}

func (t *ExtendedAsset) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	dec.Unpack(&t.Quantity)
	dec.Unpack(&t.Contract)
	return dec.Pos(), nil
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
