package main

import (
	"github.com/uuosio/chain"
)

//table accounts
type Account struct {
	Balance chain.Asset //primary: t.Balance.Symbol.Code()
}

//table stat
type CurrencyStats struct {
	Supply    chain.Asset //primary: t.Supply.Symbol.Code()
	MaxSupply chain.Asset
	Issuer    chain.Name
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
func (token *Token) Create(issuer chain.Name, maximumSupply chain.Asset) {
	sym_code := maximumSupply.Symbol.Code()
	db := NewCurrencyStatsDB(token.code, chain.Name{sym_code})
	itr := db.Find(sym_code)
	chain.Check(!itr.IsOk(), "token with symbol already exists")

	stats := &CurrencyStats{}
	stats.Issuer = issuer
	stats.MaxSupply = maximumSupply
	stats.Supply.Symbol = maximumSupply.Symbol
	db.Store(stats, issuer)
}

//action issue
func (token *Token) Issue(to chain.Name, quantity chain.Asset, memo string) {
	sym_code := quantity.Symbol.Code()
	db := NewCurrencyStatsDB(token.code, chain.Name{sym_code})
	it, item := db.GetByKey(sym_code)
	chain.Check(it.IsOk(), "token with symbol does not exist, create token before issue")
	chain.Check(to == item.Issuer, "tokens can only be issued to issuer account")

	chain.RequireAuth(item.Issuer)
	chain.Check(quantity.IsValid(), "invalid quantity")
	chain.Check(quantity.Amount > 0, "must issue positive quantity")

	item.Supply.Add(&quantity)
	db.Store(item, chain.Name{0})
	// _data := t.Pack()
	// chain.Assert(Equal(_data, data), "bad issue")
}

//action retire
func (token *Token) Retire(quantity chain.Asset, memo string) {
}

//action transfer
func (token *Token) Transfer(from chain.Name, to chain.Name, quantity chain.Asset, memo string) {
}

//action open
func (token *Token) Open(owner chain.Name, symbol chain.Symbol, ramPayer chain.Name) {
}

//action close
func (token *Token) Close(owner chain.Name, symbol chain.Symbol) {
}
