package main

import (
	"github.com/uuosio/chain"
)

func check(b bool, msg string) {
	chain.Check(b, msg)
}

//table accounts
type account struct {
	balance chain.Asset //primary: t.balance.Symbol.Code()
}

//table stat
type currency_stats struct {
	supply     chain.Asset //primary: t.supply.Symbol.Code()
	max_supply chain.Asset
	issuer     chain.Name
}

//contract token
type Token struct {
	receiver      chain.Name
	firstReceiver chain.Name
	action        chain.Name
}

func NewAccountTable(code chain.Name, scope chain.Name) *accountTable {
	return NewaccountTable(code, scope)
}

func NewCurrencyStatsTable(code chain.Name, scope chain.Name) *currency_statsTable {
	return Newcurrency_statsTable(code, scope)
}

func NewContract(receiver, firstReceiver, action chain.Name) *Token {
	return &Token{receiver, firstReceiver, action}
}

//action create
func (token *Token) Create(issuer chain.Name, maximum_supply chain.Asset) {
	chain.RequireAuth(token.receiver)
	check(maximum_supply.Symbol.IsValid(), "invalid symbol name")
	check(maximum_supply.IsValid(), "invalid supply")
	check(maximum_supply.Amount > 0, "max_supply must be positive")

	sym_code := maximum_supply.Symbol.Code()
	db := NewCurrencyStatsTable(token.receiver, chain.Name{sym_code})
	itr := db.Find(sym_code)
	check(!itr.IsOk(), "token with symbol already exists")

	stats := &currency_stats{}
	stats.supply.Symbol = maximum_supply.Symbol
	stats.max_supply = maximum_supply
	stats.issuer = issuer
	db.Store(stats, token.receiver)
}

//action issue
func (token *Token) Issue(to chain.Name, quantity chain.Asset, memo string) {
	check(quantity.Symbol.IsValid(), "invalid symbol name")
	check(len(memo) <= 256, "memo has more than 256 bytes")

	sym_code := quantity.Symbol.Code()
	db := NewCurrencyStatsTable(token.receiver, chain.Name{sym_code})
	it, item := db.GetByKey(sym_code)
	check(it.IsOk(), "token with symbol does not exist, create token before issue")
	check(to == item.issuer, "tokens can only be issued to issuer account")

	chain.RequireAuth(item.issuer)
	check(quantity.IsValid(), "invalid quantity")
	check(quantity.Amount > 0, "must issue positive quantity")

	check(quantity.Symbol == item.supply.Symbol, "symbol precision mismatch")
	check(quantity.Amount <= item.max_supply.Amount-item.supply.Amount, "quantity exceeds available supply")

	item.supply.Add(&quantity)
	db.Update(it, item, item.issuer)

	token.addBalance(to, quantity, to)
}

//action retire
func (token *Token) Retire(quantity chain.Asset, memo string) {
	check(quantity.Symbol.IsValid(), "invalid symbol name")
	check(len(memo) <= 256, "memo has more than 256 bytes")
	stats := NewCurrencyStatsTable(token.receiver, chain.Name{quantity.Symbol.Code()})
	it, item := stats.GetByKey(quantity.Symbol.Code())
	check(it.IsOk(), "token with symbol does not exist")
	chain.RequireAuth(item.issuer)
	check(quantity.IsValid(), "invalid quantity")
	check(quantity.Amount > 0, "must retire positive quantity")
	check(quantity.Symbol == item.supply.Symbol, "symbol precision mismatch")

	item.supply.Sub(&quantity)
	stats.Update(it, item, chain.SamePayer)
	token.subBalance(item.issuer, quantity)
}

//action transfer
func (token *Token) Transfer(from chain.Name, to chain.Name, quantity chain.Asset, memo string) {
	check(from != to, "cannot transfer to self")
	chain.RequireAuth(from)
	check(chain.IsAccount(to), "to account does not exist")

	stats := NewCurrencyStatsTable(token.receiver, chain.Name{quantity.Symbol.Code()})
	it, st := stats.GetByKey(quantity.Symbol.Code())
	check(it.IsOk(), "token with symbol does not exist")

	chain.RequireRecipient(from)
	chain.RequireRecipient(to)

	check(quantity.IsValid(), "invalid quantity")
	check(quantity.Amount > 0, "must transfer positive quantity")
	check(quantity.Symbol == st.supply.Symbol, "symbol precision mismatch")
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
	stats := NewCurrencyStatsTable(token.receiver, chain.Name{symbol.Code()})
	it, st := stats.GetByKey(symbol.Code())
	check(it.IsOk(), "symbol does not exist")
	check(st.supply.Symbol == symbol, "symbol precision mismatch")

	accountTable := NewAccountTable(token.receiver, owner)
	it = accountTable.Find(symbol.Code())
	if !it.IsOk() {
		account := &account{}
		account.balance = chain.Asset{0, symbol}
		accountTable.Store(account, ram_payer)
	}
}

//action close
func (token *Token) Close(owner chain.Name, symbol chain.Symbol) {
	chain.RequireAuth(owner)
	accountTable := NewAccountTable(token.receiver, owner)
	it, item := accountTable.GetByKey(symbol.Code())
	check(it.IsOk(), "Balance row already deleted or never existed. Action won't have any effect.")
	check(item.balance.Amount == 0, "Cannot close because the balance is not zero.")
	accountTable.Remove(it)
}

func (token *Token) subBalance(owner chain.Name, value chain.Asset) {
	accountTable := NewAccountTable(token.receiver, owner)
	it, from := accountTable.GetByKey(value.Symbol.Code())
	check(it.IsOk(), "no balance object found")
	check(from.balance.Amount >= value.Amount, "overdrawn balance")
	from.balance.Sub(&value)
	accountTable.Update(it, from, owner)
}

func (token *Token) addBalance(owner chain.Name, value chain.Asset, ramPayer chain.Name) {
	accountTable := NewAccountTable(token.receiver, owner)
	it, to := accountTable.GetByKey(value.Symbol.Code())
	if !it.IsOk() {
		account := &account{balance: value}
		accountTable.Store(account, ramPayer)
	} else {
		to.balance.Add(&value)
		accountTable.Update(it, to, chain.Name{0})
	}
}
