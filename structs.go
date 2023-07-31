package chain

import "encoding/binary"

type VarInt32 int32

func (t *VarInt32) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.PackVarInt32(int32(*t))
	return enc.GetSize() - oldSize
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

func (t *VarUint32) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.PackVarUint32(uint32(*t))
	return enc.GetSize() - oldSize
}

func (t *VarUint32) Unpack(data []byte) int {
	v, n := UnpackVarUint32(data)
	*t = VarUint32(v)
	return n
}

func (t *VarUint32) Size() int {
	return PackedVarUint32Length(uint32(*t))
}

type Uint256 [32]byte

func NewUint256(lo1, lo2, hi1, hi2 uint64) Uint256 {
	ret := Uint256{}
	binary.LittleEndian.PutUint64(ret[:8], lo1)
	binary.LittleEndian.PutUint64(ret[8:16], lo2)
	binary.LittleEndian.PutUint64(ret[16:24], hi1)
	binary.LittleEndian.PutUint64(ret[24:32], hi2)
	return ret
}

func NewUint256FromBytes(bs []byte) Uint256 {
	ret := Uint256{}
	Check(len(bs) >= 32, "bad size")
	copy(ret[:], bs)
	return ret
}

func (n *Uint256) SetUint64(v uint64) {
	tmp := Uint256{}
	copy(n[:], tmp[:]) //memset
	binary.LittleEndian.PutUint64(n[:8], v)
}

func (n *Uint256) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.Write(n[16:])
	enc.Write(n[:16])
	return enc.GetSize() - oldSize
}

func (n *Uint256) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Read(n[16:])
	dec.Read(n[:16])
	return 32
}

func (t *Uint256) Size() int {
	return 32
}

func (n *Uint256) Uint64() uint64 {
	return binary.LittleEndian.Uint64(n[:8])
}

type TimePoint struct {
	Elapsed uint64
}

func (t *TimePoint) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.PackUint64(t.Elapsed)
	return enc.GetSize() - oldSize
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

func (t *TimePointSec) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.PackUint32(t.UTCSeconds)
	return enc.GetSize() - oldSize
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

func (t *BlockTimestampType) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.PackUint32(t.Slot)
	return enc.GetSize() - oldSize
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
