package chain

import (
	"github.com/uuosio/chain/eosio"
)

type Checksum160 [20]byte

func (t *Checksum160) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.WriteBytes(t[:])
	return enc.GetSize() - oldSize
}

func (t *Checksum160) Unpack(data []byte) int {
	Check(len(data) >= t.Size(), "Unpack data overflow")
	copy(t[:], data)
	return t.Size()
}

func (t *Checksum160) Size() int {
	return len(*t)
}

type Checksum256 [32]byte

func (t *Checksum256) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.WriteBytes(t[:])
	return enc.GetSize() - oldSize
}

func (t *Checksum256) Unpack(data []byte) int {
	Check(len(data) >= t.Size(), "Unpack data overflow")
	copy(t[:], data)
	return t.Size()
}

func (t *Checksum256) Size() int {
	return len(*t)
}

type Checksum512 [64]byte

func (t *Checksum512) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.WriteBytes(t[:])
	return enc.GetSize() - oldSize
}

func (t *Checksum512) Unpack(data []byte) int {
	Check(len(data) >= t.Size(), "Unpack data overflow")
	copy(t[:], data)
	return t.Size()
}

func (t *Checksum512) Size() int {
	return len(*t)
}

//Tests if the sha256 hash generated from data matches the provided checksum.
func AssertSha256(data []byte, hash Checksum256) {
	eosio.AssertSha256(data, hash)
}

//Tests if the sha1 hash generated from data matches the provided checksum.
func AssertSha1(data []byte, hash Checksum160) {
	eosio.AssertSha1(data, hash)
}

//Tests if the sha512 hash generated from data matches the provided checksum.
func AssertSha512(data []byte, hash Checksum512) {
	eosio.AssertSha512(data, hash)
}

//Tests if the ripemod160 hash generated from data matches the provided checksum.
func AssertRipemd160(data []byte, hash Checksum160) {
	eosio.AssertRipemd160(data, hash)
}

//Hashes data using sha256 and return hash value.
func Sha256(data []byte) Checksum256 {
	return eosio.Sha256(data)
}

//Hashes data using sha1 and return hash value.
func Sha1(data []byte) Checksum160 {
	return eosio.Sha1(data)
}

//Hashes data using sha512 and return hash value.
func Sha512(data []byte) Checksum512 {
	return eosio.Sha512(data)
}

//Hashes data using ripemd160 and return hash value.
func Ripemd160(data []byte) Checksum160 {
	return eosio.Ripemd160(data)
}

//Recover the public key from digest and signature
func RecoverKey(digest Checksum256, sig *Signature) *PublicKey {
	_sig := EncoderPack(sig)
	pub := eosio.RecoverKey(digest, _sig)
	_pub := &PublicKey{}
	_pub.Unpack(pub[:])
	return _pub
}

//Tests a given public key with the generated key from digest and the signature
func AssertRecoverKey(digest Checksum256, sig Signature, pub PublicKey) {
	_sig := EncoderPack(&sig)
	_pub := EncoderPack(&pub)
	eosio.AssertRecoverKey(digest, _sig, _pub)
}

type Signature struct {
	Type uint8 // Signature type
	Data [65]byte
}

func (t *Signature) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.WriteUint8(t.Type)
	enc.WriteBytes(t.Data[:])
	return enc.GetSize() - oldSize
}

func (t *Signature) Unpack(data []byte) int {
	dec := NewDecoder(data)
	n := dec.ReadUint8()
	t.Type = n
	dec.Read(t.Data[:])
	return dec.Pos()
}

func (t *Signature) Size() int {
	return 1 + len(t.Data)
}

type PublicKey struct {
	Type uint8 // Public key type
	Data [33]byte
}

func (t *PublicKey) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.WriteUint8(t.Type)
	enc.WriteBytes(t.Data[:])
	return enc.GetSize() - oldSize
}

func (t *PublicKey) Unpack(data []byte) int {
	// var err error
	dec := NewDecoder(data)
	_type := dec.ReadUint8()

	t.Type = _type
	dec.Read(t.Data[:])
	return dec.Pos()
}

func (t *PublicKey) Size() int {
	return 34
}

func less(a, b []byte) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] < b[i] {
			return true
		}
		if a[i] > b[i] {
			return false
		}
	}
	return len(a) < len(b)
}

type PublicKeyList []PublicKey

func (a PublicKeyList) Len() int { return len(a) }
func (a PublicKeyList) Less(i, j int) bool {
	return less(a[i].Data[:], a[j].Data[:])
}

func (a PublicKeyList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
