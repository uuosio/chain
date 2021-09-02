package main

import (
	"github.com/uuosio/chain"
)

func check(b bool, msg string) {
	chain.Check(b, msg)
}

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
	receiver      chain.Name
	firstReceiver chain.Name
	action        chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *Token {
	return &Token{receiver, firstReceiver, action}
}

//action create
func (token *Token) Create(issuer chain.Name, maximum_supply chain.Asset) {
	check(maximum_supply.Symbol.IsValid(), "invalid symbol name")
	check(maximum_supply.IsValid(), "invalid supply")
	check(maximum_supply.Amount > 0, "max-supply must be positive")

	sym_code := maximum_supply.Symbol.Code()
	db := NewCurrencyStatsDB(token.receiver, chain.Name{sym_code})
	itr := db.Find(sym_code)
	check(!itr.IsOk(), "token with symbol already exists")

	stats := &CurrencyStats{}
	stats.Supply.Symbol = maximum_supply.Symbol
	stats.MaxSupply = maximum_supply
	stats.Issuer = issuer
	db.Store(stats, issuer)
}

//action issue
func (token *Token) Issue(to chain.Name, quantity chain.Asset, memo string) {
	sym_code := quantity.Symbol.Code()
	db := NewCurrencyStatsDB(token.receiver, chain.Name{sym_code})
	it, item := db.Get(sym_code)
	check(it.IsOk(), "token with symbol does not exist, create token before issue")
	check(to == item.Issuer, "tokens can only be issued to issuer account")

	chain.RequireAuth(item.Issuer)
	check(quantity.IsValid(), "invalid quantity")
	check(quantity.Amount > 0, "must issue positive quantity")

	item.Supply.Add(&quantity)
	db.Update(it, item, item.Issuer)

	token.addBalance(to, quantity, to)
}

//action retire
func (token *Token) Retire(quantity chain.Asset, memo string) {
	check(quantity.Symbol.IsValid(), "invalid symbol name")
	check(len(memo) <= 256, "memo has more than 256 bytes")
	stats := NewCurrencyStatsDB(token.receiver, chain.Name{quantity.Symbol.Code()})
	it, item := stats.Get(quantity.Symbol.Code())
	check(it.IsOk(), "token with symbol does not exist")
	chain.RequireAuth(item.Issuer)
	check(quantity.IsValid(), "invalid quantity")
	check(quantity.Amount > 0, "must retire positive quantity")
	check(quantity.Symbol == item.Supply.Symbol, "symbol precision mismatch")

	item.Supply.Sub(&quantity)
	stats.Update(it, item, chain.SamePayer)
	token.subBalance(item.Issuer, quantity)
}

//action transfer
func (token *Token) Transfer(from chain.Name, to chain.Name, quantity chain.Asset, memo string) {
	check(from != to, "cannot transfer to self")
	chain.RequireAuth(from)
	check(chain.IsAccount(to), "to account does not exist")

	stats := NewCurrencyStatsDB(token.receiver, chain.Name{quantity.Symbol.Code()})
	it, st := stats.Get(quantity.Symbol.Code())
	check(it.IsOk(), "token with symbol does not exist")

	chain.RequireRecipient(from)
	chain.RequireRecipient(to)

	check(quantity.IsValid(), "invalid quantity")
	check(quantity.Amount > 0, "must transfer positive quantity")
	check(quantity.Symbol == st.Supply.Symbol, "symbol precision mismatch")
	check(len(memo) <= 256, "memo has more than 256 bytes")

	payer := chain.Name{0}
	if chain.HasAuth(to) {
		payer = to
	} else {
		payer = from
	}

	token.subBalance(from, quantity)
	token.addBalance(to, quantity, payer)
}

//action open
func (token *Token) Open(owner chain.Name, symbol chain.Symbol, ram_payer chain.Name) {
	chain.RequireAuth(ram_payer)
	check(chain.IsAccount(owner), "owner account does not exist")
	stats := NewCurrencyStatsDB(token.receiver, chain.Name{symbol.Code()})
	it, st := stats.Get(symbol.Code())
	check(it.IsOk(), "symbol does not exist")
	check(st.Supply.Symbol == symbol, "symbol precision mismatch")

	accountDB := NewAccountDB(token.receiver, owner)
	it = accountDB.Find(symbol.Code())
	if !it.IsOk() {
		account := &Account{}
		account.Balance = chain.Asset{0, symbol}
		accountDB.Store(account, ram_payer)
	}
}

//action close
func (token *Token) Close(owner chain.Name, symbol chain.Symbol) {
	chain.RequireAuth(owner)
	accountDB := NewAccountDB(token.receiver, owner)
	it, item := accountDB.Get(symbol.Code())
	check(it.IsOk(), "Balance row already deleted or never existed. Action won't have any effect.")
	check(item.Balance.Amount == 0, "Cannot close because the balance is not zero.")
	accountDB.Remove(it)
}

func (token *Token) subBalance(owner chain.Name, value chain.Asset) {
	accountDB := NewAccountDB(token.receiver, owner)
	it, from := accountDB.Get(value.Symbol.Code())
	check(it.IsOk(), "no balance object found")
	check(from.Balance.Amount >= value.Amount, "overdrawn balance")
	from.Balance.Sub(&value)
	accountDB.Update(it, from, owner)
}

func (token *Token) addBalance(owner chain.Name, value chain.Asset, ramPayer chain.Name) {
	accountDB := NewAccountDB(token.receiver, owner)
	it, to := accountDB.Get(value.Symbol.Code())
	if !it.IsOk() {
		account := &Account{Balance: value}
		accountDB.Store(account, owner)
	} else {
		to.Balance.Add(&value)
		accountDB.Update(it, to, owner)
	}
}
