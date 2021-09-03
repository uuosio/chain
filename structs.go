package chain

import (
	"encoding/binary"
)

type VarInt32 uint32

func (t *VarInt32) Pack() []byte {
	return PackVarUint32(uint32(*t))
}

func (t *VarInt32) Unpack(data []byte) (int, error) {
	v, n := UnpackVarUint32(data)
	*t = VarInt32(v)
	return n, nil
}

func (t *VarInt32) Size() int {
	return PackedVarInt32Length(int32(*t))
}

type VarUint32 uint32

func (t *VarUint32) Pack() []byte {
	return PackVarUint32(uint32(*t))
}

func (t *VarUint32) Unpack(data []byte) (int, error) {
	v, n := UnpackVarUint32(data)
	*t = VarUint32(v)
	return n, nil
}

func (t *VarUint32) Size() int {
	return PackedVarUint32Length(uint32(*t))
}

type Int128 [16]byte

func (n *Int128) Pack() []byte {
	return n[:]
}

func (n *Int128) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	if err := dec.Read(n[:]); err != nil {
		return 0, err
	}
	return 16, nil
}

func (t *Int128) Size() int {
	return 16
}

type Uint128 [16]byte

func (n *Uint128) Pack() []byte {
	return n[:]
}

func (n *Uint128) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	if err := dec.Read(n[:]); err != nil {
		return 0, err
	}
	return 16, nil
}

func (t *Uint128) Size() int {
	return 16
}

func (n *Uint128) SetUint64(v uint64) {
	tmp := Uint128{}
	copy(n[:], tmp[:]) //memset
	binary.LittleEndian.PutUint64(n[:], v)
}

func (n *Uint128) Uint64() uint64 {
	return binary.LittleEndian.Uint64(n[:])
}

type Uint256 [32]uint8

func (n *Uint256) Pack() []byte {
	return n[:]
}

func (n *Uint256) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	if err := dec.Read(n[:]); err != nil {
		return 0, err
	}
	return 32, nil
}

func (t *Uint256) Size() int {
	return 32
}

func (n *Uint256) SetUint64(v uint64) {
	tmp := Uint256{}
	copy(n[:], tmp[:]) //memset
	binary.LittleEndian.PutUint64(n[:], v)
}

func (n *Uint256) Uint64() uint64 {
	return binary.LittleEndian.Uint64(n[:])
}

type Float128 [16]byte

func (n *Float128) Pack() []byte {
	return n[:]
}

func (n *Float128) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	if err := dec.Read(n[:]); err != nil {
		return 0, err
	}
	return 16, nil
}

func (t *Float128) Size() int {
	return 16
}

type TimePoint struct {
	Elapsed uint64
}

func (t *TimePoint) Pack() []byte {
	enc := NewEncoder(t.Size())
	enc.PackUint64(t.Elapsed)
	return enc.GetBytes()
}

func (t *TimePoint) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	dec.Unpack(&t.Elapsed)
	return 8, nil
}

func (t *TimePoint) Size() int {
	return 8
}

type TimePointSec struct {
	UTCSeconds uint32
}

func (t *TimePointSec) Pack() []byte {
	enc := NewEncoder(t.Size())
	enc.PackUint32(t.UTCSeconds)
	return enc.GetBytes()
}

func (t *TimePointSec) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	dec.Unpack(&t.UTCSeconds)
	return 4, nil
}

func (t *TimePointSec) Size() int {
	return 4
}

type BlockTimestampType struct {
	Slot uint32
}

func (t *BlockTimestampType) Pack() []byte {
	enc := NewEncoder(t.Size())
	enc.PackUint32(t.Slot)
	return enc.GetBytes()
}

func (t *BlockTimestampType) Unpack(data []byte) (int, error) {
	dec := NewDecoder(data)
	dec.Unpack(&t.Slot)
	return 4, nil
}

func (t *BlockTimestampType) Size() int {
	return 4
}
