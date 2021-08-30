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
	receiver      chain.Name
	firstReceiver chain.Name
	action        chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *Token {
	return &Token{receiver, firstReceiver, action}
}

//action create
func (token *Token) Create(issuer chain.Name, maximumSupply chain.Asset) {
	sym_code := maximumSupply.Symbol.Code()
	db := NewCurrencyStatsDB(token.receiver, chain.Name{sym_code})
	itr := db.Find(sym_code)
	check(!itr.IsOk(), "token with symbol already exists")

	stats := &CurrencyStats{}
	stats.Issuer = issuer
	stats.MaxSupply = maximumSupply
	stats.Supply.Symbol = maximumSupply.Symbol
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

	item.Supply.Add(quantity)
	db.Store(item, chain.Name{0})
	// _data := t.Pack()
	// chain.Assert(Equal(_data, data), "bad issue")
}

//action retire
func (token *Token) Retire(quantity chain.Asset, memo string) {
	// auto sym = quantity.symbol;
	// check( sym.is_valid(), "invalid symbol name" );
	// check( memo.size() <= 256, "memo has more than 256 bytes" );

	// stats statstable( get_self(), sym.code().raw() );
	// auto existing = statstable.find( sym.code().raw() );
	// check( existing != statstable.end(), "token with symbol does not exist" );
	// const auto& st = *existing;

	// require_auth( st.issuer );
	// check( quantity.is_valid(), "invalid quantity" );
	// check( quantity.amount > 0, "must retire positive quantity" );

	// check( quantity.symbol == st.supply.symbol, "symbol precision mismatch" );

	// statstable.modify( st, same_payer, [&]( auto& s ) {
	//    s.supply -= quantity;
	// });

	// sub_balance( st.issuer, quantity );
	check(quantity.Symbol.IsValid(), "invalid symbol name")
	check(len(memo) <= 256, "memo has more than 256 bytes")
	stats := NewCurrencyStatsDB(token.receiver, token.receiver)
	it, item := stats.Get(quantity.Symbol.Code())
	check(it.IsOk(), "token with symbol does not exist")
	chain.RequireAuth(item.Issuer)
	check(quantity.IsValid(), "invalid quantity")
	check(quantity.Amount > 0, "must retire positive quantity")
	check(quantity.Symbol == item.Supply.Symbol, "symbol precision mismatch")

	item.Supply.Sub(quantity)
	stats.Store(item, chain.Name{0})
}

//action transfer
func (token *Token) Transfer(from chain.Name, to chain.Name, quantity chain.Asset, memo string) {
	// check( from != to, "cannot transfer to self" );
	// require_auth( from );
	// check( is_account( to ), "to account does not exist");
	// auto sym = quantity.symbol.code();
	// stats statstable( get_self(), sym.raw() );
	// const auto& st = statstable.get( sym.raw() );

	// require_recipient( from );
	// require_recipient( to );

	// check( quantity.is_valid(), "invalid quantity" );
	// check( quantity.amount > 0, "must transfer positive quantity" );
	// check( quantity.symbol == st.supply.symbol, "symbol precision mismatch" );
	// check( memo.size() <= 256, "memo has more than 256 bytes" );

	// auto payer = has_auth( to ) ? to : from;

	// sub_balance( from, quantity );
	// add_balance( to, quantity, payer );
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
func (token *Token) Open(owner chain.Name, symbol chain.Symbol, ramPayer chain.Name) {
	// require_auth( ram_payer );

	// check( is_account( owner ), "owner account does not exist" );

	// auto sym_code_raw = symbol.code().raw();
	// stats statstable( get_self(), sym_code_raw );
	// const auto& st = statstable.get( sym_code_raw, "symbol does not exist" );
	// check( st.supply.symbol == symbol, "symbol precision mismatch" );

	// accounts acnts( get_self(), owner.value );
	// auto it = acnts.find( sym_code_raw );
	// if( it == acnts.end() ) {
	// 	acnts.emplace( ram_payer, [&]( auto& a ){
	// 		a.balance = asset{0, symbol};
	// 	});
	// }
	chain.RequireAuth(ramPayer)
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
		accountDB.Store(account, ramPayer)
	}
}

//action close
func (token *Token) Close(owner chain.Name, symbol chain.Symbol) {
	// require_auth( owner );
	// accounts acnts( get_self(), owner.value );
	// auto it = acnts.find( symbol.code().raw() );
	// check( it != acnts.end(), "Balance row already deleted or never existed. Action won't have any effect." );
	// check( it->balance.amount == 0, "Cannot close because the balance is not zero." );
	// acnts.erase( it );
	chain.RequireAuth(owner)
	accountDB := NewAccountDB(token.receiver, owner)
	it, item := accountDB.Get(symbol.Code())
	check(it.IsOk(), "Balance row already deleted or never existed. Action won't have any effect.")
	check(item.Balance.Amount == 0, "Cannot close because the balance is not zero.")
	accountDB.Remove(it)
}

func (token *Token) subBalance(owner chain.Name, value chain.Asset) {
	// accounts from_acnts( get_self(), owner.value );
	// const auto& from = from_acnts.get( value.symbol.code().raw(), "no balance object found" );
	// check( from.balance.amount >= value.amount, "overdrawn balance" );

	// from_acnts.modify( from, owner, [&]( auto& a ) {
	// 	  a.balance -= value;
	//    });

	accountDB := NewAccountDB(token.receiver, owner)
	it, from := accountDB.Get(value.Symbol.Code())
	check(it.IsOk(), "no balance object found")
	check(from.Balance.Amount >= value.Amount, "overdrawn balance")
	from.Balance.Sub(value)
	accountDB.Store(from, owner)
}

func (token *Token) addBalance(owner chain.Name, value chain.Asset, ramPayer chain.Name) {
	//    accounts to_acnts( get_self(), owner.value );
	//    auto to = to_acnts.find( value.symbol.code().raw() );
	//    if( to == to_acnts.end() ) {
	//       to_acnts.emplace( ram_payer, [&]( auto& a ){
	//         a.balance = value;
	//       });
	//    } else {
	//       to_acnts.modify( to, same_payer, [&]( auto& a ) {
	//         a.balance += value;
	//       });
	//    }
	accountDB := NewAccountDB(token.receiver, owner)
	it, to := accountDB.Get(value.Symbol.Code())
	if !it.IsOk() {
		account := &Account{Balance: value}
		accountDB.Store(account, owner)
	} else {
		to.Balance.Add(value)
		accountDB.Update(it, to, owner)
	}
}
