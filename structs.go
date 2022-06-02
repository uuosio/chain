package chain

import "encoding/binary"

type VarInt32 int32

func (t *VarInt32) Pack() []byte {
	return PackVarInt32(int32(*t))
}

func (t *VarInt32) Unpack(data []byte) int {
	v, n := UnpackVarInt32(data)
	*t = VarInt32(v)
	return n
}

func (t *VarInt32) Size() int {
	return PackedVarInt32Length(int32(*t))
}

type VarUint32 uint32

func (t *VarUint32) Pack() []byte {
	return PackVarUint32(uint32(*t))
}

func (t *VarUint32) Unpack(data []byte) int {
	v, n := UnpackVarUint32(data)
	*t = VarUint32(v)
	return n
}

func (t *VarUint32) Size() int {
	return PackedVarUint32Length(uint32(*t))
}

type Uint256 [4]uint64

func NewUint256(lo1, lo2, hi1, hi2 uint64) *Uint256 {
	ret := Uint256{lo1, lo2, hi1, hi2}
	return &ret
}

func NewUint256FromBytes(bs []byte) *Uint256 {
	ret := &Uint256{}
	Check(len(bs) >= 32, "bad size")
	ret[0] = binary.LittleEndian.Uint64(bs[0:8])
	ret[1] = binary.LittleEndian.Uint64(bs[8:16])
	ret[2] = binary.LittleEndian.Uint64(bs[16:24])
	ret[3] = binary.LittleEndian.Uint64(bs[24:32])
	return ret
}

func (n *Uint256) SetUint64(v uint64) {
	tmp := Uint256{}
	copy(n[:], tmp[:]) //memset
	n[3] = 0
	n[2] = 0
	n[1] = 0
	n[0] = v
}

func (n *Uint256) Pack() []byte {
	enc := NewEncoder(32)
	for _, a := range n {
		enc.PackUint64(a)
	}
	return enc.GetBytes()
}

func (n *Uint256) Unpack(data []byte) int {
	dec := NewDecoder(data)
	for i := 0; i <= 3; i++ {
		n[i] = dec.UnpackUint64()
	}
	return 32
}

func (t *Uint256) Size() int {
	return 32
}

func (n *Uint256) Uint64() uint64 {
	return n[0]
}

type TimePoint struct {
	Elapsed uint64
}

func (t *TimePoint) Pack() []byte {
	enc := NewEncoder(t.Size())
	enc.PackUint64(t.Elapsed)
	return enc.GetBytes()
}

func (t *TimePoint) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Unpack(&t.Elapsed)
	return 8
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

func (t *TimePointSec) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Unpack(&t.UTCSeconds)
	return 4
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

func (t *BlockTimestampType) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Unpack(&t.Slot)
	return 4
}

func (t *BlockTimestampType) Size() int {
	return 4
}

type BinaryExtension struct {
	HasValue bool
}

type Optional struct {
	IsValid bool
}
