package main

import (
	"chain"
	"chain/database"
)

type Symbol struct {
	Value uint64
}

func (a *Symbol) Code() uint64 {
	return a.Value >> 8
}

type Asset struct {
	Amount int64
	Symbol Symbol
}

func (a *Asset) IsValid() bool {
	return true
}

func (a *Asset) Add(b Asset) {
	a.Amount += b.Amount
}

type Account struct {
	balance Asset
}

func (a *Account) GetPrimary() uint64 {
	return a.balance.Symbol.Value >> 8
}

type CurrencyStats struct {
	Supply    Asset
	MaxSupply Asset
	Issuer    chain.Name
}

func (stats *CurrencyStats) GetPrimary() uint64 {
	return stats.Supply.Symbol.Value >> 8
}

func (stats *CurrencyStats) Print() {
	chain.Print(stats.Issuer, stats.MaxSupply.Amount)
}

type AccountDB struct {
	database.DBI64
}

func NewAccountDB(code chain.Name, scope chain.Name, table chain.Name) *AccountDB {
	v := &AccountDB{}
	v.Init(code, scope, table, accountUnpacker)
	return v
}

func accountUnpacker(data []byte) (database.DBValue, error) {
	a := &Account{}
	_, err := a.Unpack(data)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (db *AccountDB) Get(iterator database.Iterator) *Account {
	data := db.GetRawByIterator(iterator)
	if len(data) <= 0 {
		return nil
	}

	_data := &Account{}
	_data.Unpack(data)
	return _data
}

type CurrencyStatsDB struct {
	database.DBI64
}

func statsUnpacker(data []byte) (database.DBValue, error) {
	a := &CurrencyStats{}
	_, err := a.Unpack(data)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func NewCurrencyStatsDB(code chain.Name, scope chain.Name, table chain.Name) *CurrencyStatsDB {
	v := &CurrencyStatsDB{}
	v.Init(code, scope, table, statsUnpacker)
	return v
}

func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

//contract token
type Token struct {
	receiver chain.Name
	code     chain.Name
	action   chain.Name
	// statsDb CurrencyStatsDB
	// accountDB AccountDB
}

func NewContract(receiver chain.Name, firstReceiver chain.Name, action chain.Name) *Token {
	return &Token{receiver, firstReceiver, action}
}

//action create
func (token *Token) Create(issuer chain.Name, maximumSupply Asset) {
	sym_code := maximumSupply.Symbol.Code()
	db := NewCurrencyStatsDB(token.code, chain.Name{sym_code}, chain.NewName("stat"))
	itr := db.Find(sym_code)
	chain.Check(!itr.IsOk(), "token with symbol already exists")

	stats := CurrencyStats{}
	stats.Issuer = issuer
	stats.MaxSupply = maximumSupply
	stats.Supply.Symbol = maximumSupply.Symbol
	db.Set(sym_code, stats.Pack(), issuer)
}

//action issue
func (token *Token) Issue(to chain.Name, quantity Asset, memo string) {
	sym_code := quantity.Symbol.Code()
	db := NewCurrencyStatsDB(token.code, chain.Name{sym_code}, chain.NewName("stat"))
	_item, err := db.Get(sym_code)
	chain.Check(err == nil, "token with symbol does not exist, create token before issue")
	item, ok := _item.(*CurrencyStats)
	chain.Check(ok, "bad item type")
	chain.Check(to == item.Issuer, "tokens can only be issued to issuer account")

	chain.RequireAuth(item.Issuer)
	chain.Check(quantity.IsValid(), "invalid quantity")
	chain.Check(quantity.Amount > 0, "must issue positive quantity")

	item.Supply.Add(quantity)
	db.Set(sym_code, item.Pack(), chain.Name{0})
	// _data := t.Pack()
	// chain.Assert(Equal(_data, data), "bad issue")
}

//action retire
func (token *Token) Retire(quantity Asset, memo string) {
}

//action transfer
func (token *Token) Transfer(from chain.Name, to chain.Name, quantity Asset, memo string) {
}

//action open
func (token *Token) Open(owner chain.Name, symbol Symbol, ramPayer chain.Name) {
}

//action close
func (token *Token) Close(owner chain.Name, symbol Symbol) {
}
