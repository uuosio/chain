package chain

type KeyWeight struct {
	Key    PublicKey
	Weight uint16
}

type PermissionLevel struct {
	Actor      Name
	Permission Name
}

type PermissionLevelWeight struct {
	Permission PermissionLevel
	Weight     uint16
}

type WaitWeight struct {
	WaitSec uint32
	Weight  uint16
}

type Authority struct {
	Threshold uint32
	Keys      []KeyWeight
	Accounts  []PermissionLevelWeight
	Waits     []WaitWeight
}

func NewPermissionLevel(actor Name, permission Name) *PermissionLevel {
	return &PermissionLevel{actor, permission}
}

func (t *PermissionLevel) Pack(enc *Encoder) int {
	size := enc.GetSize()
	enc.PackUint64(t.Actor.N)
	enc.PackUint64(t.Permission.N)
	return enc.GetSize() - size
}

func (t *PermissionLevel) Unpack(data []byte) int {
	dec := NewDecoder(data)
	t.Actor = dec.UnpackName()
	t.Permission = dec.UnpackName()
	return dec.Pos()
}

func (t *PermissionLevel) Size() int {
	size := 0
	size += 8 //Actor
	size += 8 //Permission
	return size
}

func (t *PermissionLevelWeight) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	t.Permission.Pack(enc)
	enc.PackUint16(t.Weight)
	return enc.GetSize() - oldSize
}

func (t *PermissionLevelWeight) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.UnpackI(&t.Permission)
	t.Weight = dec.UnpackUint16()
	return dec.Pos()
}

func (t *PermissionLevelWeight) Size() int {
	size := 0
	size += t.Permission.Size() //Permission
	size += 2                   //Weight
	return size
}

func (t *WaitWeight) Pack(enc *Encoder) int {
	size := enc.GetSize()
	enc.PackUint32(t.WaitSec)
	enc.PackUint16(t.Weight)
	return enc.GetSize() - size
}

func (t *WaitWeight) Unpack(data []byte) int {
	dec := NewDecoder(data)
	t.WaitSec = dec.UnpackUint32()
	t.Weight = dec.UnpackUint16()
	return dec.Pos()
}

func (t *WaitWeight) Size() int {
	size := 0
	size += 4 //WaitSec
	size += 2 //Weight
	return size
}

func (t *Authority) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()

	enc.PackUint32(t.Threshold)
	{
		enc.PackLength(len(t.Keys))
		for i := range t.Keys {
			t.Keys[i].Pack(enc)
		}
	}

	{
		enc.PackLength(len(t.Accounts))
		for i := range t.Accounts {
			t.Accounts[i].Pack(enc)
		}
	}

	{
		enc.PackLength(len(t.Waits))
		for i := range t.Waits {
			t.Waits[i].Pack(enc)
		}
	}

	return enc.GetSize() - oldSize
}

func (t *Authority) Unpack(data []byte) int {
	dec := NewDecoder(data)
	t.Threshold = dec.UnpackUint32()
	{
		length := dec.UnpackLength()
		t.Keys = make([]KeyWeight, length)
		for i := 0; i < length; i++ {
			dec.UnpackI(&t.Keys[i])
		}
	}

	{
		length := dec.UnpackLength()
		t.Accounts = make([]PermissionLevelWeight, length)
		for i := 0; i < length; i++ {
			dec.UnpackI(&t.Accounts[i])
		}
	}

	{
		length := dec.UnpackLength()
		t.Waits = make([]WaitWeight, length)
		for i := 0; i < length; i++ {
			dec.UnpackI(&t.Waits[i])
		}
	}

	return dec.Pos()
}

func (t *Authority) Size() int {
	size := 0
	size += 4 //Threshold
	size += PackedVarUint32Length(uint32(len(t.Keys)))

	for i := range t.Keys {
		size += t.Keys[i].Size()
	}
	size += PackedVarUint32Length(uint32(len(t.Accounts)))

	for i := range t.Accounts {
		size += t.Accounts[i].Size()
	}
	size += PackedVarUint32Length(uint32(len(t.Waits)))

	for i := range t.Waits {
		size += t.Waits[i].Size()
	}
	return size
}

func (t *KeyWeight) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	t.Key.Pack(enc)
	enc.PackUint16(t.Weight)
	return enc.GetSize() - oldSize
}

func (t *KeyWeight) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.UnpackI(&t.Key)
	t.Weight = dec.UnpackUint16()
	return dec.Pos()
}

func (t *KeyWeight) Size() int {
	size := 0
	size += t.Key.Size() //Key
	size += 2            //Weight
	return size
}
