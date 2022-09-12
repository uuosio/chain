
package testfloat128
import (
	"github.com/uuosio/chain"
    "github.com/uuosio/chain/database"
    "unsafe"
)


//eliminate unused package errors
func dummy() {
	if false {
		v := 0;
		n := unsafe.Sizeof(v);
		chain.Printui(uint64(n));
		chain.Printui(database.IDX64);
	}
}