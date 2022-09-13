//go:build !tinygo.wasm
// +build !tinygo.wasm

package chain

func GetApplyArgs() (Name, Name, Name) {
	panic("should not be called!")
	return Name{0}, Name{0}, Name{0}
}
