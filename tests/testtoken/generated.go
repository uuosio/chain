
package testtoken
import (
	"github.com/uuosio/chain"
    "github.com/uuosio/chain/database"
    "unsafe"
)

type create struct {
    issuer chain.Name
    maximumSupply chain.Asset
}


func (t *create) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.PackUint64(t.issuer.N)
	enc.Pack(&t.maximumSupply)
    return enc.GetBytes()
}

func (t *create) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	t.issuer.N = dec.UnpackUint64()
	dec.Unpack(&t.maximumSupply)
    return dec.Pos()
}

func (t *create) Size() int {
    size := 0
	size += 8 //issuer
	size += t.maximumSupply.Size() //maximumSupply
    return size
}

type issue struct {
    to chain.Name
    quantity chain.Asset
    memo string
}


func (t *issue) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.PackUint64(t.to.N)
	enc.Pack(&t.quantity)
	enc.PackString(t.memo)
    return enc.GetBytes()
}

func (t *issue) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	t.to.N = dec.UnpackUint64()
	dec.Unpack(&t.quantity)
	t.memo = dec.UnpackString()
    return dec.Pos()
}

func (t *issue) Size() int {
    size := 0
	size += 8 //to
	size += t.quantity.Size() //quantity
	size += chain.PackedVarUint32Length(uint32(len(t.memo))) + len(t.memo) //memo
    return size
}

type retire struct {
    quantity chain.Asset
    memo string
}


func (t *retire) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.Pack(&t.quantity)
	enc.PackString(t.memo)
    return enc.GetBytes()
}

func (t *retire) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	dec.Unpack(&t.quantity)
	t.memo = dec.UnpackString()
    return dec.Pos()
}

func (t *retire) Size() int {
    size := 0
	size += t.quantity.Size() //quantity
	size += chain.PackedVarUint32Length(uint32(len(t.memo))) + len(t.memo) //memo
    return size
}

type transfer struct {
    from chain.Name
    to chain.Name
    quantity chain.Asset
    memo string
}


func (t *transfer) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.PackUint64(t.from.N)
	enc.PackUint64(t.to.N)
	enc.Pack(&t.quantity)
	enc.PackString(t.memo)
    return enc.GetBytes()
}

func (t *transfer) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	t.from.N = dec.UnpackUint64()
	t.to.N = dec.UnpackUint64()
	dec.Unpack(&t.quantity)
	t.memo = dec.UnpackString()
    return dec.Pos()
}

func (t *transfer) Size() int {
    size := 0
	size += 8 //from
	size += 8 //to
	size += t.quantity.Size() //quantity
	size += chain.PackedVarUint32Length(uint32(len(t.memo))) + len(t.memo) //memo
    return size
}

type open struct {
    owner chain.Name
    symbol chain.Symbol
    ramPayer chain.Name
}


func (t *open) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.PackUint64(t.owner.N)
	enc.Pack(&t.symbol)
	enc.PackUint64(t.ramPayer.N)
    return enc.GetBytes()
}

func (t *open) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	t.owner.N = dec.UnpackUint64()
	dec.Unpack(&t.symbol)
	t.ramPayer.N = dec.UnpackUint64()
    return dec.Pos()
}

func (t *open) Size() int {
    size := 0
	size += 8 //owner
	size += 8 //symbol
	size += 8 //ramPayer
    return size
}

type close struct {
    owner chain.Name
    symbol chain.Symbol
}


func (t *close) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.PackUint64(t.owner.N)
	enc.Pack(&t.symbol)
    return enc.GetBytes()
}

func (t *close) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	t.owner.N = dec.UnpackUint64()
	dec.Unpack(&t.symbol)
    return dec.Pos()
}

func (t *close) Size() int {
    size := 0
	size += 8 //owner
	size += 8 //symbol
    return size
}


func (t *Account) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.Pack(&t.Balance)
    return enc.GetBytes()
}

func (t *Account) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	dec.Unpack(&t.Balance)
    return dec.Pos()
}

func (t *Account) Size() int {
    size := 0
	size += t.Balance.Size() //Balance
    return size
}


func (t *CurrencyStats) Pack() []byte {
    enc := chain.NewEncoder(t.Size())
	enc.Pack(&t.Supply)
	enc.Pack(&t.MaxSupply)
	enc.PackUint64(t.Issuer.N)
    return enc.GetBytes()
}

func (t *CurrencyStats) Unpack(data []byte) int {
    dec := chain.NewDecoder(data)
	dec.Unpack(&t.Supply)
	dec.Unpack(&t.MaxSupply)
	t.Issuer.N = dec.UnpackUint64()
    return dec.Pos()
}

func (t *CurrencyStats) Size() int {
    size := 0
	size += t.Supply.Size() //Supply
	size += t.MaxSupply.Size() //MaxSupply
	size += 8 //Issuer
    return size
}

var (
	AccountSecondaryTypes = []int{
	}
)

func AccountTableNameToIndex(indexName string) int {
	switch indexName {
	default:
		panic("unknow indexName")
	}
}

func AccountUnpacker(buf []byte) database.MultiIndexValue {
	v := &Account{}
	v.Unpack(buf)
	return v
}

func (t *Account) GetSecondaryValue(index int) interface{} {
	switch index {
		default:
			panic("index out of bound")
	}
}

func (t *Account) SetSecondaryValue(index int, v interface{}) {
	switch index {
	default:
		panic("unknown index")
	}
}

func (t *Account) GetPrimary() uint64 {
    return t.Balance.Symbol.Code()
}

type AccountTable struct {
	database.MultiIndexInterface
}

func (mi *AccountTable) Store(v *Account, payer chain.Name) {
	mi.MultiIndexInterface.Store(v, payer)
}

func (mi *AccountTable) GetByKey(id uint64) (*database.Iterator, *Account) {
	it, data := mi.MultiIndexInterface.GetByKey(id)
	if !it.IsOk() {
		return it, nil
	}
	return it, data.(*Account)
}

func (mi *AccountTable) GetByIterator(it *database.Iterator) *Account {
	data := mi.MultiIndexInterface.GetByIterator(it)
	return data.(*Account)
}

func (mi *AccountTable) Update(it *database.Iterator, v *Account, payer chain.Name) {
	mi.MultiIndexInterface.Update(it, v, payer)
}

func NewAccountTable(code chain.Name, optionalScope ...chain.Name) *AccountTable {
	var scope chain.Name
	if len(optionalScope) > 0 {
		scope = optionalScope[0]
	} else {
		scope = chain.Name{N: 0}
	}
	table := chain.Name{N:uint64(3607749779137757184)} //table name: Account
	if table.N&uint64(0x0f) != 0 {
		// Limit table names to 12 characters so that the last character (4 bits) can be used to distinguish between the secondary indices.
		panic("NewMultiIndex:Invalid multi-index table name ")
	}

	mi := &database.MultiIndex{}
	mi.SetTable(code, scope, table)
	mi.Table = database.NewTableI64(code, scope, table, func(data []byte) database.TableValue {
		return mi.Unpack(data)
	})
	mi.IdxTableNameToIndex = AccountTableNameToIndex
	mi.IndexTypes = AccountSecondaryTypes
	mi.IDXTables = make([]database.SecondaryTable, len(AccountSecondaryTypes))
	mi.Unpack = AccountUnpacker
	return &AccountTable{mi}
}

var (
	CurrencyStatsSecondaryTypes = []int{
	}
)

func CurrencyStatsTableNameToIndex(indexName string) int {
	switch indexName {
	default:
		panic("unknow indexName")
	}
}

func CurrencyStatsUnpacker(buf []byte) database.MultiIndexValue {
	v := &CurrencyStats{}
	v.Unpack(buf)
	return v
}

func (t *CurrencyStats) GetSecondaryValue(index int) interface{} {
	switch index {
		default:
			panic("index out of bound")
	}
}

func (t *CurrencyStats) SetSecondaryValue(index int, v interface{}) {
	switch index {
	default:
		panic("unknown index")
	}
}

func (t *CurrencyStats) GetPrimary() uint64 {
    return t.Supply.Symbol.Code()
}

type CurrencyStatsTable struct {
	database.MultiIndexInterface
}

func (mi *CurrencyStatsTable) Store(v *CurrencyStats, payer chain.Name) {
	mi.MultiIndexInterface.Store(v, payer)
}

func (mi *CurrencyStatsTable) GetByKey(id uint64) (*database.Iterator, *CurrencyStats) {
	it, data := mi.MultiIndexInterface.GetByKey(id)
	if !it.IsOk() {
		return it, nil
	}
	return it, data.(*CurrencyStats)
}

func (mi *CurrencyStatsTable) GetByIterator(it *database.Iterator) *CurrencyStats {
	data := mi.MultiIndexInterface.GetByIterator(it)
	return data.(*CurrencyStats)
}

func (mi *CurrencyStatsTable) Update(it *database.Iterator, v *CurrencyStats, payer chain.Name) {
	mi.MultiIndexInterface.Update(it, v, payer)
}

func NewCurrencyStatsTable(code chain.Name, optionalScope ...chain.Name) *CurrencyStatsTable {
	var scope chain.Name
	if len(optionalScope) > 0 {
		scope = optionalScope[0]
	} else {
		scope = chain.Name{N: 0}
	}
	table := chain.Name{N:uint64(14289235522390851584)} //table name: CurrencyStats
	if table.N&uint64(0x0f) != 0 {
		// Limit table names to 12 characters so that the last character (4 bits) can be used to distinguish between the secondary indices.
		panic("NewMultiIndex:Invalid multi-index table name ")
	}

	mi := &database.MultiIndex{}
	mi.SetTable(code, scope, table)
	mi.Table = database.NewTableI64(code, scope, table, func(data []byte) database.TableValue {
		return mi.Unpack(data)
	})
	mi.IdxTableNameToIndex = CurrencyStatsTableNameToIndex
	mi.IndexTypes = CurrencyStatsSecondaryTypes
	mi.IDXTables = make([]database.SecondaryTable, len(CurrencyStatsSecondaryTypes))
	mi.Unpack = CurrencyStatsUnpacker
	return &CurrencyStatsTable{mi}
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


func main() {
	receiver, firstReceiver, action := chain.GetApplyArgs()
	ContractApply(receiver.N, firstReceiver.N, action.N)
}

func ContractApply(_receiver, _firstReceiver, _action uint64) {
	receiver := chain.Name{_receiver}
	firstReceiver := chain.Name{_firstReceiver}
	action := chain.Name{_action}

	contract := NewContract(receiver, firstReceiver, action)
	if contract == nil {
		return
	}
	data := chain.ReadActionData()
	
	//Fix data declared but not used error
	if false {
		println(len(data))
	}

    if receiver == firstReceiver {
        switch action.N {
        case uint64(5031766152489992192): //create
            t := create{}
            t.Unpack(data)
            contract.Create(t.issuer, t.maximumSupply)
        case uint64(8516769789752901632): //issue
            t := issue{}
            t.Unpack(data)
            contract.Issue(t.to, t.quantity, t.memo)
        case uint64(13453074143696125952): //retire
            t := retire{}
            t.Unpack(data)
            contract.Retire(t.quantity, t.memo)
        case uint64(14829575313431724032): //transfer
            t := transfer{}
            t.Unpack(data)
            contract.Transfer(t.from, t.to, t.quantity, t.memo)
        case uint64(11913481165836648448): //open
            t := open{}
            t.Unpack(data)
            contract.Open(t.owner, t.symbol, t.ramPayer)
        case uint64(4929617502180212736): //close
            t := close{}
            t.Unpack(data)
            contract.Close(t.owner, t.symbol)
        }
    }
    if receiver != firstReceiver {
        switch action.N {
        }
    }
}
