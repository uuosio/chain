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

func (t *Checksum160) Pack() []byte {
	return t[:]
}

func (t *Checksum160) Unpack(data []byte) (int, error) {
	Check(len(data) >= t.Size(), "Unpack data overflow")
	copy(t[:], data)
	return t.Size(), nil
}

func (t *Checksum160) Size() int {
	return len(*t)
}

type Checksum256 [32]byte

func (t *Checksum256) Pack() []byte {
	return t[:]
}

func (t *Checksum256) Unpack(data []byte) (int, error) {
	Check(len(data) >= t.Size(), "Unpack data overflow")
	copy(t[:], data)
	return t.Size(), nil
}

func (t *Checksum256) Size() int {
	return len(*t)
}

type Checksum512 [64]byte

func (t *Checksum512) Pack() []byte {
	return t[:]
}

func (t *Checksum512) Unpack(data []byte) (int, error) {
	Check(len(data) >= t.Size(), "Unpack data overflow")
	copy(t[:], data)
	return t.Size(), nil
}

func (t *Checksum512) Size() int {
	return len(*t)
}

//Tests if the sha256 hash generated from data matches the provided checksum.
func AssertSha256(data []byte, hash Checksum256) {
	C.assert_sha256((*C.char)(unsafe.Pointer(&data[0])), uint32(len(data)), (*C.uint8_t)(unsafe.Pointer(&hash)))
}

//Tests if the sha1 hash generated from data matches the provided checksum.
func AssertSha1(data []byte, hash Checksum160) {
	C.assert_sha1((*C.char)(unsafe.Pointer(&data[0])), uint32(len(data)), (*C.uint8_t)(unsafe.Pointer(&hash)))
}

//Tests if the sha512 hash generated from data matches the provided checksum.
func AssertSha512(data []byte, hash Checksum512) {
	C.assert_sha512((*C.char)(unsafe.Pointer(&data[0])), uint32(len(data)), (*C.uint8_t)(unsafe.Pointer(&hash)))
}

//Tests if the ripemod160 hash generated from data matches the provided checksum.
func AssertRipemd160(data []byte, hash Checksum160) {
	C.assert_ripemd160((*C.char)(unsafe.Pointer(&data[0])), uint32(len(data)), (*C.uint8_t)(unsafe.Pointer(&hash)))
}

//Hashes data using sha256 and return hash value.
func Sha256(data []byte) Checksum256 {
	var hash Checksum256
	C.sha256((*C.char)(unsafe.Pointer(&data[0])), uint32(len(data)), (*C.uint8_t)(unsafe.Pointer(&hash)))
	return hash
}

//Hashes data using sha1 and return hash value.
func Sha1(data []byte) Checksum160 {
	var hash Checksum160
	C.sha1((*C.char)(unsafe.Pointer(&data[0])), uint32(len(data)), (*C.uint8_t)(unsafe.Pointer(&hash)))
	return hash
}

//Hashes data using sha512 and return hash value.
func Sha512(data []byte) Checksum512 {
	var hash Checksum512
	C.sha512((*C.char)(unsafe.Pointer(&data[0])), uint32(len(data)), (*C.uint8_t)(unsafe.Pointer(&hash)))
	return hash
}

//Hashes data using ripemd160 and return hash value.
func Ripemd160(data []byte) Checksum160 {
	var hash Checksum160
	C.ripemd160((*C.char)(unsafe.Pointer(&data[0])), uint32(len(data)), (*C.uint8_t)(unsafe.Pointer(&hash)))
	return hash
}

//Recover the public key from digest and signature
func RecoverKey(digest Checksum256, sig Signature) *PublicKey {
	//TODO: handle webauth signature
	var pub [128]byte //34
	_sig := sig.Pack()
	ret := C.recover_key((*C.uint8_t)(unsafe.Pointer(&digest)), (*C.char)(unsafe.Pointer(&_sig[0])), C.size_t(len(_sig)), (*C.char)(unsafe.Pointer(&pub[0])), C.size_t(len(pub)))
	_pub := &PublicKey{}
	_pub.Unpack(pub[:int(ret)])
	return _pub
}

//Tests a given public key with the generated key from digest and the signature
func AssertRecoverKey(digest Checksum256, sig Signature, pub PublicKey) {
	_sig := sig.Pack()
	_pub := pub.Pack()
	C.assert_recover_key((*C.uint8_t)(unsafe.Pointer(&digest)), (*C.char)(unsafe.Pointer(&_sig[0])), C.size_t(len(_sig)), (*C.char)(unsafe.Pointer(&_pub[0])), C.size_t(len(_pub)))
}

//TODO: implement Signature&PublicKey struct
type Signature struct {
	Type int // Signature type
	Data [65]byte
}

func (t *Signature) Pack() []byte {
	enc := NewEncoder(1 + len(t.Data))
	enc.WriteUint8(uint8(t.Type))
	enc.WriteBytes(t.Data[:])
	return enc.GetBytes()
}

func (t *Signature) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	n, err := dec.ReadUint8()
	if err != nil {
		return 0, err
	}
	t.Type = int(n)
	dec.Read(t.Data[:])
	return dec.Pos(), nil
}

func (t *Signature) Size() int {
	return 1 + len(t.Data)
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
