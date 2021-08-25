package chain

/*
#include <stddef.h>
#include <stdint.h>

//fake typedef for call api from go
typedef uint8_t capi_checksum160;
typedef uint8_t capi_checksum256;
typedef uint8_t capi_checksum512;

void assert_sha256( const char* data, uint32_t length, const capi_checksum256* hash );
void assert_sha1( const char* data, uint32_t length, const capi_checksum160* hash );
void assert_sha512( const char* data, uint32_t length, const capi_checksum512* hash );
void assert_ripemd160( const char* data, uint32_t length, const capi_checksum160* hash );
void sha256( const char* data, uint32_t length, capi_checksum256* hash );
void sha1( const char* data, uint32_t length, capi_checksum160* hash );
void sha512( const char* data, uint32_t length, capi_checksum512* hash );
void ripemd160( const char* data, uint32_t length, capi_checksum160* hash );
int recover_key( const capi_checksum256* digest, const char* sig, size_t siglen, char* pub, size_t publen );
void assert_recover_key( const capi_checksum256* digest, const char* sig, size_t siglen, const char* pub, size_t publen );
*/
import "C"
import "unsafe"

type Checksum160 [20]byte
type Checksum256 [32]byte
type Checksum512 [64]byte

// void assert_sha256( const char* data, uint32_t length, const struct capi_checksum256* hash );
func AssertSha256(data []byte, hash Checksum256) {
	C.assert_sha256((*C.char)(unsafe.Pointer(&data[0])), uint32(len(data)), (*C.uint8_t)(unsafe.Pointer(&hash)))
}

// void assert_sha1( const char* data, uint32_t length, const struct capi_checksum160* hash );
func AssertSha1(data []byte, hash Checksum160) {
	C.assert_sha1((*C.char)(unsafe.Pointer(&data[0])), uint32(len(data)), (*C.uint8_t)(unsafe.Pointer(&hash)))
}

// void assert_sha512( const char* data, uint32_t length, const struct capi_checksum512* hash );
func AssertSha512(data []byte, hash Checksum512) {
	C.assert_sha512((*C.char)(unsafe.Pointer(&data[0])), uint32(len(data)), (*C.uint8_t)(unsafe.Pointer(&hash)))
}

// void assert_ripemd160( const char* data, uint32_t length, const struct capi_checksum160* hash );
func AssertRipemd160(data []byte, hash Checksum160) {
	C.assert_ripemd160((*C.char)(unsafe.Pointer(&data[0])), uint32(len(data)), (*C.uint8_t)(unsafe.Pointer(&hash)))
}

// void sha256( const char* data, uint32_t length, struct capi_checksum256* hash );
func Sha256(data []byte) Checksum256 {
	var hash Checksum256
	C.sha256((*C.char)(unsafe.Pointer(&data[0])), uint32(len(data)), (*C.uint8_t)(unsafe.Pointer(&hash)))
	return hash
}

// void sha1( const char* data, uint32_t length, struct capi_checksum160* hash );
func Sha1(data []byte) Checksum160 {
	var hash Checksum160
	C.sha1((*C.char)(unsafe.Pointer(&data[0])), uint32(len(data)), (*C.uint8_t)(unsafe.Pointer(&hash)))
	return hash
}

// void sha512( const char* data, uint32_t length, struct capi_checksum512* hash );
func Sha512(data []byte) Checksum512 {
	var hash Checksum512
	C.sha512((*C.char)(unsafe.Pointer(&data[0])), uint32(len(data)), (*C.uint8_t)(unsafe.Pointer(&hash)))
	return hash
}

// void ripemd160( const char* data, uint32_t length, struct capi_checksum160* hash );
func Ripemd160(data []byte) Checksum160 {
	var hash Checksum160
	C.ripemd160((*C.char)(unsafe.Pointer(&data[0])), uint32(len(data)), (*C.uint8_t)(unsafe.Pointer(&hash)))
	return hash
}

// int recover_key( const struct capi_checksum256* digest, const char* sig, size_t siglen, char* pub, size_t publen );
func RecoverKey(digest Checksum256, sig []byte) []byte {
	//TODO: handle webauth signature
	var pub [128]byte //34
	ret := C.recover_key((*C.uint8_t)(unsafe.Pointer(&digest)), (*C.char)(unsafe.Pointer(&sig[0])), C.size_t(len(sig)), (*C.char)(unsafe.Pointer(&pub[0])), C.size_t(len(pub)))
	return pub[:int(ret)]
}

// void assert_recover_key( const struct capi_checksum256* digest, const char* sig, size_t siglen, const char* pub, size_t publen );
func AssertRecoverKey(digest Checksum256, sig []byte, pub []byte) {
	C.assert_recover_key((*C.uint8_t)(unsafe.Pointer(&digest)), (*C.char)(unsafe.Pointer(&sig[0])), C.size_t(len(sig)), (*C.char)(unsafe.Pointer(&pub[0])), C.size_t(len(pub)))
}

//TODO: implement Signature&PublicKey struct
type Signature struct {
	Type int // Signature type
	Data []byte
}

func (t *Signature) Pack() []byte {
	enc := NewEncoder(5 + len(t.Data))
	enc.WriteBytes(t.Data[:])
	return enc.GetBytes()
}

func (t *Signature) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	dec.Read(t.Data[:])
	return dec.Pos(), nil
}

func (t *Signature) Size() int {
	return PackedSizeLength(uint32(len(t.Data))) + len(t.Data)
}

type PublicKey struct {
	Type int // Public key type
	Data [33]byte
}

func (t *PublicKey) Pack() []byte {
	enc := NewEncoder(34)
	enc.WriteUint8(uint8(t.Type))
	enc.WriteBytes(t.Data[:])
	return enc.GetBytes()
}

func (t *PublicKey) Unpack(data []byte) (int, error) {
	// var err error
	dec := NewDecoder(data)
	_type, err := dec.ReadUint8()
	if err != nil {
		return 0, err
	}

	t.Type = int(_type)
	dec.Read(t.Data[:])
	return dec.Pos(), nil
}

func (t *PublicKey) Size() int {
	return 34
}
