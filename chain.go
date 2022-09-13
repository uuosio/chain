package chain

import "github.com/uuosio/chain/eosio"

//Gets the set of active producers.
func GetActiveProducers() []Name {
	prods := eosio.GetActiveProducers()

	var _prods = make([]Name, 0, len(prods))
	for _, v := range prods {
		_prods = append(_prods, Name{v})
	}
	return _prods
}
