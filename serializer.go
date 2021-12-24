package chain

import (
	"encoding/binary"
	"math"
)

func PackVarInt32(v int32) []byte {
	i := 0
	buf := [8]byte{}
	val := uint32((v << 1) ^ (v >> 31))
	for {
		b := uint8(val) & uint8(0x7f)
		val >>= 7
		if val > 0 {
			b |= uint8(1 << 7)
		}
		buf[i] = b
		i += 1
		if val == 0 {
			break
		}
	}
	return buf[:i]
}

func PackedVarInt32Length(v int32) int {
	i := 0
	val := uint32((v << 1) ^ (v >> 31))
	for {
		b := uint8(val) & uint8(0x7f)
		val >>= 7
		if val > 0 {
			b |= uint8(1 << 7)
		}
		i += 1
		if val == 0 {
			break
		}
	}
	return i
}

func UnpackVarInt32(buf []byte) (int32, int) {
	v := uint32(0)
	by := int(0)
	length := 0
	for _, b := range buf {
		v |= uint32(b&0x7f) << by
		by += 7
		length += 1
		if b&0x80 == 0 {
			break
		}
	}
	v = (v >> 1) ^ (^(v & 1) + 1)
	return int32(v), length
}

func PackVarUint32(val uint32) []byte {
	result := make([]byte, 0, 5)
	for {
		b := byte(val & 0x7f)
		val >>= 7
		if val > 0 {
			b |= byte(1 << 7)
		}
		result = append(result, b)
		if val <= 0 {
			break
		}
	}
	return result
}

func UnpackVarUint32(val []byte) (v uint32, n int) {
	var by int = 0
	// if len(val) > 5 {
	// 	val = val[:5]
	// }
	n = 0
	for _, b := range val {
		v |= uint32(b&0x7f) << by
		by += 7
		n += 1
		if b&0x80 == 0 {
			break
		}
	}
	return
}

func PackedVarUint32Length(val uint32) int {
	n := 0
	for {
		b := byte(val & 0x7f)
		val >>= 7
		if val > 0 {
			b |= byte(1 << 7)
		}
		n += 1
		if val <= 0 {
			break
		}
	}
	return n
}

type Serializer interface {
	Pack() []byte
	Unpack([]byte) (int, uint64)
}

type Decoder struct {
	buf []byte
	pos int
}

type Unpacker interface {
	Unpack(data []byte) int
}

type PackedSize interface {
	Size() int
}

func NewDecoder(buf []byte) *Decoder {
	dec := &Decoder{}
	dec.buf = buf
	dec.pos = 0
	return dec
}

func (dec *Decoder) Pos() int {
	return dec.pos
}

func (dec *Decoder) checkPos(n int) {
	if dec.pos+n > len(dec.buf) {
		panic("checkPos: buffer overflow in Decoder")
	}
}

func (dec *Decoder) incPos(n int) {
	dec.pos += n
	if dec.pos > len(dec.buf) {
		panic("incPos: buffer overflow in Decoder")
	}
}

func (dec *Decoder) Read(b []byte) {
	dec.checkPos(len(b))
	copy(b[:], dec.buf[dec.pos:])
	dec.incPos(len(b))
}

func (dec *Decoder) ReadInt32() int32 {
	d := dec.ReadUint32()
	return int32(d)
}

func (dec *Decoder) ReadUint32() uint32 {
	var b [4]byte
	dec.Read(b[:])
	d := binary.LittleEndian.Uint32(b[:])
	return d
}

func (dec *Decoder) ReadInt16() int16 {
	n := dec.ReadUint16()
	return int16(n)
}

func (dec *Decoder) ReadUint16() uint16 {
	var b [2]byte
	dec.Read(b[:])
	d := binary.LittleEndian.Uint16(b[:])
	return d
}

func (dec *Decoder) ReadInt64() int64 {
	n := dec.ReadUint64()
	return int64(n)
}

func (dec *Decoder) ReadUint64() uint64 {
	var b [8]byte
	dec.Read(b[:])
	d := binary.LittleEndian.Uint64(b[:])
	return d
}

func (dec *Decoder) ReadFloat32() float32 {
	n := dec.ReadUint32()
	return math.Float32frombits(n)
}

func (dec *Decoder) ReadFloat64() float64 {
	n := dec.ReadUint64()
	return math.Float64frombits(n)
}

func (dec *Decoder) ReadBool() bool {
	var b [1]byte
	dec.Read(b[:])
	return b[0] == 1
}

func (dec *Decoder) UnpackBool() bool {
	return dec.ReadBool()
}

func (dec *Decoder) UnpackString() string {
	bb := dec.UnpackBytes()
	return string(bb)
}

func (dec *Decoder) UnpackName() Name {
	n := dec.UnpackUint64()
	return Name{n}
}

func (dec *Decoder) UnpackBytes() []byte {
	length := dec.UnpackLength()
	buf := make([]byte, length)
	dec.Read(buf)
	return buf
}

func (dec *Decoder) UnpackLength() int {
	v, n := UnpackVarUint32(dec.buf[dec.pos:])
	dec.incPos(n)
	return int(v)
}

func (dec *Decoder) UnpackVarInt32() int32 {
	v, n := UnpackVarInt32(dec.buf[dec.pos:])
	dec.incPos(n)
	return v
}

func (dec *Decoder) UnpackVarUint32() uint32 {
	v, n := UnpackVarUint32(dec.buf[dec.pos:])
	dec.incPos(n)
	return v
}

func (dec *Decoder) UnpackInt16() int16 {
	v := dec.ReadInt16()
	return v
}

func (dec *Decoder) UnpackUint16() uint16 {
	return dec.ReadUint16()
}

func (dec *Decoder) UnpackInt32() int32 {
	v := dec.ReadInt32()
	return int32(v)
}

func (dec *Decoder) UnpackUint32() uint32 {
	v := dec.ReadUint32()
	return v
}

func (dec *Decoder) UnpackFloat32() float32 {
	v := dec.ReadUint32()
	return math.Float32frombits(v)
}

func (dec *Decoder) UnpackFloat64() float64 {
	v := dec.ReadUint64()
	return math.Float64frombits(v)
}

func (dec *Decoder) UnpackInt64() int64 {
	v := dec.ReadInt64()
	return v
}

func (dec *Decoder) UnpackUint64() uint64 {
	v := dec.ReadUint64()
	return v
}

func (dec *Decoder) ReadInt8() int8 {
	var b [1]byte
	dec.Read(b[:])
	return int8(b[0])
}

func (dec *Decoder) ReadUint8() uint8 {
	var b [1]byte
	dec.Read(b[:])
	return b[0]
}

func (dec *Decoder) UnpackInt8() int8 {
	return dec.ReadInt8()
}

func (dec *Decoder) UnpackUint8() uint8 {
	return dec.ReadUint8()
}

func (dec *Decoder) UnpackAction() *Action {
	a := &Action{}
	n := a.Unpack(dec.buf[dec.pos:])
	dec.incPos(n)
	return a
}

func (dec *Decoder) UnpackI(unpacker Unpacker) {
	n := unpacker.Unpack(dec.buf[dec.pos:])
	dec.incPos(n)
}

// Unpack supported type:
// Unpacker, interface,
// *string, *[]byte,
// *uint8, *int16, *uint16, *int32, *uint32, *int64, *uint64, *bool
// *float64
// *Name

func (dec *Decoder) Unpack(i interface{}) int {
	switch v := i.(type) {
	case Unpacker:
		n := v.Unpack(dec.buf[dec.pos:])
		dec.incPos(n)
		return n
	case *string:
		n := dec.Pos()
		*v = dec.UnpackString()
		return dec.Pos() - n
	case *[]byte:
		n := dec.Pos()
		*v = dec.UnpackBytes()
		return dec.Pos() - n
	case *bool:
		n := dec.ReadBool()
		*v = n
		return 1
	case *int8:
		n := dec.ReadInt8()
		*v = n
		return 1
	case *uint8:
		n := dec.ReadUint8()
		*v = n
		return 1
	case *int16:
		n := dec.ReadInt16()
		*v = n
		return 2
	case *uint16:
		n := dec.ReadUint16()
		*v = n
		return 2
	case *int32:
		n := dec.ReadInt32()
		*v = n
		return 4
	case *uint32:
		n := dec.ReadUint32()
		*v = n
		return 4
	case *int64:
		n := dec.ReadInt64()
		*v = n
		return 8
	case *uint64:
		n := dec.ReadUint64()
		*v = n
		return 8
	case *float32:
		n := dec.ReadFloat32()
		*v = n
		return 4
	case *float64:
		n := dec.ReadFloat64()
		*v = n
		return 8
	// Name struct implemented Unpacker interface
	// case *Name:
	// 	n, err := dec.UnpackUint64()
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// 	v.N = n
	// 	return 8, nil
	default:
		// if DEBUG {
		// 	panic(fmt.Sprintf("unknown Unpack type <%v>", i))
		// }
		panic("unknown type in Unpack")
	}
	return 0
}

type Encoder struct {
	buf []byte
}

type Packer interface {
	Pack() []byte
}

func NewEncoder(initSize int) *Encoder {
	ret := &Encoder{}
	ret.buf = make([]byte, 0, initSize)
	return ret
}

func (enc *Encoder) Reset() {
	enc.buf = enc.buf[:0]
}

func (enc *Encoder) Bytes() []byte {
	return enc.buf
}

func (enc *Encoder) Write(b []byte) {
	enc.buf = append(enc.buf, b...)
}

func (enc *Encoder) WriteByte(b byte) {
	enc.buf = append(enc.buf, b)
}

// Pack supported types:
// Packer
// string, bytes
// byte, uint16, int32, uint32, int64, uint64, float64
// Name
func (enc *Encoder) Pack(i interface{}) error {
	switch v := i.(type) {
	case Packer:
		enc.Write(v.Pack())
	case string:
		enc.PackString(v)
	case []byte:
		enc.PackBytes(v)
	case bool:
		if v {
			enc.Write([]byte{1})
		} else {
			enc.Write([]byte{0})
		}
	case int8:
		enc.WriteByte(byte(v))
	case uint8:
		enc.WriteByte(byte(v))
	case int16:
		enc.WriteInt16(v)
	case uint16:
		enc.WriteUint16(v)
	case int32:
		enc.WriteInt32(v)
	case uint32:
		enc.WriteUint32(v)
	case int64:
		enc.WriteInt64(v)
	case uint64:
		enc.WriteUint64(v)
	// case Uint128:
	// 	enc.PackBytes(v[:])
	// case Uint256:
	// 	enc.PackBytes(v[:])
	case float32:
		enc.PackFloat32(v)
	case float64:
		enc.PackFloat64(v)
	case Name:
		enc.WriteUint64(v.N)
	case Asset:
		enc.Write(v.Pack())
	default:
		// if DEBUG {
		// 	panic(fmt.Sprintf("Unknow Pack type <%v>", i))
		// }
		panic("Unknown Pack type")
		//		return errors.New("Unknow Pack type")
	}
	return nil
}

func (enc *Encoder) PackFloat32(f float32) {
	n := math.Float32bits(f)
	enc.WriteUint32(n)
}

func (enc *Encoder) PackFloat64(f float64) {
	n := math.Float64bits(f)
	enc.WriteUint64(n)
}

func (enc *Encoder) PackName(name Name) {
	enc.WriteUint64(name.N)
}

func (enc *Encoder) PackLength(n int) {
	enc.Write(PackVarUint32(uint32(n)))
}

func (enc *Encoder) PackVarUint32(n uint32) {
	enc.Write(PackVarUint32(uint32(n)))
}

func (enc *Encoder) PackBool(b bool) {
	if b {
		enc.WriteByte(byte(1))
	} else {
		enc.WriteByte(byte(0))
	}
}

func (enc *Encoder) PackInt8(d int8) {
	enc.WriteByte(byte(d))
}

func (enc *Encoder) PackUint8(d uint8) {
	enc.WriteByte(byte(d))
}

func (enc *Encoder) PackInt16(d int16) {
	enc.WriteUint16(uint16(d))
}

func (enc *Encoder) PackUint16(d uint16) {
	enc.WriteUint16(d)
}

func (enc *Encoder) PackInt32(d int32) {
	enc.WriteUint32(uint32(d))
}

func (enc *Encoder) PackUint32(d uint32) {
	enc.WriteUint32(d)
}

func (enc *Encoder) PackInt64(d int64) {
	b := [8]byte{}
	binary.LittleEndian.PutUint64(b[:], uint64(d))
	enc.Write(b[:])
}

func (enc *Encoder) PackUint64(d uint64) {
	b := [8]byte{}
	binary.LittleEndian.PutUint64(b[:], d)
	enc.Write(b[:])
}

func (enc *Encoder) PackVarInt32(n int32) {
	enc.Write(PackVarInt32(int32(n)))
}

func (enc *Encoder) PackString(s string) {
	enc.Write(PackVarUint32(uint32(len(s))))
	enc.Write([]byte(s))
}

func (enc *Encoder) PackBytes(v []byte) {
	enc.Write(PackVarUint32(uint32(len(v))))
	enc.Write([]byte(v))
}

func (enc *Encoder) WriteBytes(v []byte) {
	enc.Write(v)
}

func (enc *Encoder) WriteInt(d int) {
	enc.WriteUint32(uint32(d))
}

func (enc *Encoder) WriteInt32(d int32) {
	enc.WriteUint32(uint32(d))
}

func (enc *Encoder) WriteUint32(d uint32) {
	b := [4]byte{}
	binary.LittleEndian.PutUint32(b[:], d)
	enc.Write(b[:])
}

func (enc *Encoder) WriteUint8(d uint8) {
	b := [1]byte{d}
	enc.Write(b[:])
}

func (enc *Encoder) WriteInt16(d int16) {
	enc.WriteUint16(uint16(d))
}

func (enc *Encoder) WriteUint16(d uint16) {
	b := [2]byte{}
	binary.LittleEndian.PutUint16(b[:], d)
	enc.Write(b[:])
}

func (enc *Encoder) WriteInt64(d int64) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(d))
	enc.Write(b)
}

func (enc *Encoder) WriteUint64(d uint64) {
	b := [8]byte{}
	binary.LittleEndian.PutUint64(b[:], d)
	enc.Write(b[:])
}

func (enc *Encoder) GetBytes() []byte {
	return enc.buf
}

func CalcPackedSize(i interface{}) int {
	switch v := i.(type) {
	case PackedSize:
		return v.Size()
	case string:
		return PackedVarUint32Length(uint32(len(v))) + len(v)
	case []byte:
		return PackedVarUint32Length(uint32(len(v))) + len(v)
	case bool:
		return 1
	case uint8:
		return 1
	case int16:
		return 2
	case uint16:
		return 2
	case int32:
		return 4
	case uint32:
		return 4
	case int64:
		return 8
	case uint64:
		return 8
	case Uint128:
		return 16
	case Float128:
		return 16
	case Uint256:
		return 32
	case float32:
		return 4
	case float64:
		return 8
	case Name:
		return 8
	case Asset:
		return 16
	default:
		panic("Unknow type in CalcPackedSize")
	}
	return 0
}
