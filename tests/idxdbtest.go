package main

/*
#include <stdint.h>

int32_t db_idx64_store(uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, const uint64_t* secondary);
void db_idx64_update(int32_t iterator, uint64_t payer, const uint64_t* secondary);
void db_idx64_remove(int32_t iterator);
int32_t db_idx64_next(int32_t iterator, uint64_t* primary);
int32_t db_idx64_previous(int32_t iterator, uint64_t* primary);
int32_t db_idx64_find_primary(uint64_t code, uint64_t scope, uint64_t table, uint64_t* secondary, uint64_t primary);
int32_t db_idx64_find_secondary(uint64_t code, uint64_t scope, uint64_t table, const uint64_t* secondary, uint64_t* primary);
int32_t db_idx64_lowerbound(uint64_t code, uint64_t scope, uint64_t table, uint64_t* secondary, uint64_t* primary);
int32_t db_idx64_upperbound(uint64_t code, uint64_t scope, uint64_t table, uint64_t* secondary, uint64_t* primary);
int32_t db_idx64_end(uint64_t code, uint64_t scope, uint64_t table);
*/
import "C"
import (
	"chain"
	"chain/database"
	"fmt"
)

//func NewDBIdx64(code uint64, scope uint64, table uint64) *DBIdx64

func main() {
	code := chain.S2N("hello")
	scope := chain.S2N("helloo")
	table := chain.S2N("idxtablee5")
	payer := code

	db := database.NewDBIdx64(code, scope, table)
	if true {
		it := db.End()
		fmt.Println("++++It:", it)
		if it.IsEnd() {
			for {
				primary := uint64(0)
				fmt.Println("++++it.I", it.I)
				it, primary = db.Previous(it)
				// it = C.db_idx64_previous(it, (*C.uint64_t)(&primary))
				if true {
					fmt.Println("++++It:", it, primary)
				}
				if !it.IsOk() {
					break
				}
			}
		}
		return
	}
	{
		itr := db.End()
		fmt.Println("++++It:", itr)
		it := itr.I
		if it < -1 {
			for {
				primary := uint64(0)
				// it, primary := db.Previous(it)
				it = C.db_idx64_previous(it, (*C.uint64_t)(&primary))
				fmt.Println("++++It:", it, primary)
				if it < 0 {
					break
				}
			}
		}
		return
	}
	// fmt.Println(table, chain.N2S(table))
	// {
	// 	primary := uint64(555)
	// 	secondary := uint64(0)
	// 	it := C.db_idx64_find_primary(code, scope, table, &secondary, primary)
	// 	if it < 0 {
	// 		C.db_idx64_store(scope, table, payer, primary, &secondary)
	// 	}
	// 	fmt.Println(it)
	// }
	// return

	primary := uint64(11)
	it, secondary := db.Find(primary)
	fmt.Println("++++Find:", it, secondary)

	if it.IsOk() {
		db.Update(it, payer, secondary.(uint64)+1)
	} else {
		secondary := uint64(111)
		db.Store(payer, primary, secondary)
	}

	// secondary := uint64(66)
	// it, primary := db.FindSecondary(secondary)
	// fmt.Println("++++Find:", it, primary)

	// payer := code
	// if it.IsOk() {
	// 	// db.Update(it, payer, secondary+1)
	// } else {
	// 	// secondary := uint64(55)
	// 	primary := uint64(6)
	// 	db.Store(payer, primary, secondary)
	// }
	return

	// for it := db.End(); it.IsOk(); it, primary = db.Previous(it) {
	// 	fmt.Println("----Previous:", it, primary)
	// }

	// fmt.Printf("+++%v", it)
}
