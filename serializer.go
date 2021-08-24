package chain

import (
	"encoding/binary"
	"unsafe"
)

type Decoder struct {
	buf []byte
	pos int
}

type Unpacker interface {
	Unpack(data []byte) (int, error)
}

type StructSize interface {
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

func (dec *Decoder) Read(b []byte) error {
	dec.checkPos(len(b))
	copy(b[:], dec.buf[dec.pos:])
	dec.incPos(len(b))
	return nil
}

func (dec *Decoder) ReadInt() (int, error) {
	d, err := dec.ReadUint32()
	return int(d), err
}

func (dec *Decoder) ReadUint32() (uint32, error) {
	var b [4]byte
	err := dec.Read(b[:])
	if err != nil {
		return 0, err
	}
	d := binary.LittleEndian.Uint32(b[:])
	return d, nil
}

func (dec *Decoder) ReadInt16() (int16, error) {
	n, err := dec.ReadUint16()
	return int16(n), err
}

func (dec *Decoder) ReadUint16() (uint16, error) {
	var b [2]byte
	if err := dec.Read(b[:]); err != nil {
		return 0, err
	}
	d := binary.LittleEndian.Uint16(b[:])
	return d, nil
}

func (dec *Decoder) ReadInt32() (int32, error) {
	n, err := dec.ReadUint32()
	return int32(n), err
}

func (dec *Decoder) ReadInt64() (int64, error) {
	n, err := dec.ReadUint64()
	return int64(n), err
}

func (dec *Decoder) ReadUint64() (uint64, error) {
	var b [8]byte
	if err := dec.Read(b[:]); err != nil {
		return 0, err
	}
	d := binary.LittleEndian.Uint64(b[:])
	return d, nil
}

func (dec *Decoder) ReadBool() (bool, error) {
	var b [1]byte
	if err := dec.Read(b[:]); err != nil {
		return false, err
	}
	return b[0] == 1, nil
}

func (dec *Decoder) UnpackString() (string, error) {
	bb, err := dec.UnpackBytes()
	if err != nil {
		return "", err
	}
	return string(bb), nil
}

func (dec *Decoder) UnpackName() (Name, error) {
	n, err := dec.UnpackUint64()
	if err != nil {
		return Name{}, err
	}
	return Name{n}, nil
}

func (dec *Decoder) UnpackBytes() ([]byte, error) {
	length, err := dec.UnpackLength()
	if err != nil {
		return nil, err
	}
	dec.checkPos(length)
	buf := make([]byte, length)
	copy(buf, dec.buf[dec.pos:dec.pos+length])
	dec.incPos(length)
	return buf, nil
}

func (dec *Decoder) UnpackLength() (int, error) {
	n, v := UnpackUint32(dec.buf[dec.pos:])
	dec.incPos(n)
	return int(v), nil
}

func (dec *Decoder) UnpackVarInt() (uint32, error) {
	n, v := UnpackUint32(dec.buf[dec.pos:])
	dec.incPos(n)
	return v, nil
}

func (dec *Decoder) UnpackUint16() (uint16, error) {
	return dec.ReadUint16()
}

func (dec *Decoder) UnpackUint32() (uint32, error) {
	v, err := dec.ReadUint32()
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (dec *Decoder) UnpackUint64() (uint64, error) {
	v, err := dec.ReadUint64()
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (dec *Decoder) ReadUint8() (uint8, error) {
	var b [1]byte
	if err := dec.Read(b[:]); err != nil {
		return 0, err
	}
	return b[0], nil
}

func (dec *Decoder) UnpackUint8() (uint8, error) {
	return dec.ReadUint8()
}

func (dec *Decoder) UnpackAction() (*Action, error) {
	a := &Action{}
	n, err := a.Unpack(dec.buf[dec.pos:])
	if err != nil {
		return nil, err
	}
	dec.incPos(n)
	return a, nil
}

// Unpack supported type:
// Unpacker interface,
// *string, *[]byte,
// *uint8, *int16, *uint16, *int32, *uint32, *int64, *uint64, *bool
// *float64
// *Name

func (dec *Decoder) Unpack(i interface{}) (n int, err error) {
	switch v := i.(type) {
	case Unpacker:
		n, err := v.Unpack(dec.buf[dec.pos:])
		if err != nil {
			return 0, err
		}
		dec.incPos(n)
		return n, err
	case *string:
		n = dec.Pos()
		*v, err = dec.UnpackString()
		return dec.Pos() - n, err
	case *[]byte:
		n = dec.Pos()
		*v, err = dec.UnpackBytes()
		return dec.Pos() - n, err
	case *uint8:
		n, err := dec.ReadUint8()
		if err != nil {
			return 0, err
		}
		*v = n
		return 1, err
	case *int16:
		n, err := dec.ReadInt16()
		if err != nil {
			return 0, err
		}
		*v = n
		return 2, err
	case *uint16:
		n, err := dec.ReadUint16()
		if err != nil {
			return 0, err
		}
		*v = n
		return 2, err
	case *int32:
		n, err := dec.ReadInt32()
		if err != nil {
			return 0, err
		}
		*v = n
		return 4, err
	case *uint32:
		n, err := dec.ReadUint32()
		if err != nil {
			return 0, err
		}
		*v = n
		return 4, err
	case *int64:
		n, err := dec.ReadInt64()
		if err != nil {
			return 0, err
		}
		*v = n
		return 8, err
	case *uint64:
		n, err := dec.ReadUint64()
		if err != nil {
			return 0, err
		}
		*v = n
		return 8, err
	case *float64:
		n, err := dec.ReadUint64()
		if err != nil {
			return 0, err
		}
		*v = *(*float64)(unsafe.Pointer(&n))
		return 8, nil
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
		panic("unknown Unpack type")
	}
	return 0, err
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

func (enc *Encoder) Write(b []byte) {
	enc.buf = append(enc.buf, b[:]...)
}

// Pack supported types:
// Packer interface
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
	case uint8:
		enc.Write([]byte{v})
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
	default:
		// if DEBUG {
		// 	panic(fmt.Sprintf("Unknow Pack type <%v>", i))
		// }
		panic("Unknow Pack type")
		//		return errors.New("Unknow Pack type")
	}
	return nil
}

func (enc *Encoder) PackFloat32(f float32) {
	n := *(*uint32)(unsafe.Pointer(&f))
	enc.WriteUint32(n)
}

func (enc *Encoder) PackFloat64(f float64) {
	n := *(*uint64)(unsafe.Pointer(&f))
	enc.WriteUint64(n)
}

func (enc *Encoder) PackName(name Name) {
	enc.WriteUint64(name.N)
}

func (enc *Encoder) PackLength(n int) {
	enc.Write(PackUint32(uint32(n)))
}

func (enc *Encoder) PackBool(b bool) {
	if b {
		enc.Write([]byte{1})
	} else {
		enc.Write([]byte{0})
	}
}

func (enc *Encoder) PackUint8(d uint8) {
	enc.Write([]byte{d})
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

func (enc *Encoder) PackVarInt(n uint32) {
	enc.Write(PackUint32(uint32(n)))
}

func (enc *Encoder) PackString(s string) {
	enc.Write(PackUint32(uint32(len(s))))
	enc.Write([]byte(s))
}

func (enc *Encoder) PackBytes(v []byte) {
	enc.Write(PackUint32(uint32(len(v))))
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
