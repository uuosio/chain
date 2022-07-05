package database

/*
#include "../structs.h"
*/
import "C"

type IdxTable struct {
	dbIndex int
	code    C.uint64_t
	scope   C.uint64_t
	table   C.uint64_t
}

func (db *IdxTable) GetIndex() int {
	return db.dbIndex
}

func (db *IdxTable) GetTable() (uint64, uint64, uint64) {
	return uint64(db.code), uint64(db.scope), uint64(db.table)
}
