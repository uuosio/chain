package chain

type PermissionLevel struct {
	Actor      Name
	Permission Name
}

func (t *PermissionLevel) Pack() []byte {
	enc := NewEncoder(16)
	enc.Pack(&t.Actor)
	enc.Pack(&t.Permission)
	return enc.GetBytes()
}

func (t *PermissionLevel) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Unpack(&t.Actor)
	dec.Unpack(&t.Permission)
	return dec.Pos()
}

func (t *PermissionLevel) Size() int {
	return 16
}

type PermissionLevelWeight struct {
	Permission PermissionLevel
	Weight     uint16
}

func (t *PermissionLevelWeight) Pack() []byte {
	enc := NewEncoder(16 + 2)
	enc.Pack(&t.Permission)
	enc.Pack(t.Weight)
	return enc.GetBytes()
}

func (t *PermissionLevelWeight) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Unpack(&t.Permission)
	dec.Unpack(&t.Weight)
	return dec.Pos()
}

func (t *PermissionLevelWeight) Size() int {
	size := 0
	size += t.Permission.Size()
	size += 2
	return size
}

type KeyWeight struct {
	Key    PublicKey
	Weight uint16
}

func (t *KeyWeight) Pack() []byte {
	enc := NewEncoder(t.Size())
	enc.Pack(&t.Key)
	enc.PackUint16(t.Weight)
	return enc.GetBytes()
}

func (t *KeyWeight) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.UnpackI(&t.Key)
	{
		v := dec.UnpackUint16()
		t.Weight = v
	}
	return dec.Pos()
}

func (t *KeyWeight) Size() int {
	size := 0
	size += t.Key.Size()
	size += 2
	return size
}

type WaitWeight struct {
	WaitSec uint32
	Weight  uint16
}

func (t *WaitWeight) Pack() []byte {
	enc := NewEncoder(t.Size())
	enc.PackUint32(t.WaitSec)
	enc.PackUint16(t.Weight)
	return enc.GetBytes()
}

func (t *WaitWeight) Unpack(data []byte) int {
	dec := NewDecoder(data)
	{
		v := dec.UnpackUint32()
		t.WaitSec = v
	}
	{
		v := dec.UnpackUint16()
		t.Weight = v
	}
	return dec.Pos()
}

func (t *WaitWeight) Size() int {
	size := 0
	size += 4
	size += 2
	return size
}

type Authority struct {
	Threshold uint32
	Keys      []KeyWeight
	Accounts  []PermissionLevelWeight
	Waits     []WaitWeight
}

func (t *Authority) SetThreshold(threshold uint32) {
	t.Threshold = threshold
}

func (t *Authority) InitAccountWeight(size int) {
	t.Accounts = make([]PermissionLevelWeight, 0, size)
}

func (t *Authority) AddAccountWeight(perm PermissionLevel, weight uint16) {
	t.Accounts = append(t.Accounts, PermissionLevelWeight{perm, weight})
}

func (t *Authority) InitKeyWeight(size int) {
	t.Accounts = make([]PermissionLevelWeight, 0, size)
}

func (t *Authority) AddKeyWeight(pub PublicKey, weight uint16) {
	t.Keys = append(t.Keys, KeyWeight{pub, weight})
}

func (t *Authority) Pack() []byte {
	enc := NewEncoder(t.Size())
	enc.PackUint32(t.Threshold)
	{
		enc.PackLength(len(t.Keys))
		for i := range t.Keys {
			enc.Pack(&t.Keys[i])
		}
	}
	{
		enc.PackLength(len(t.Accounts))
		for i := range t.Accounts {
			enc.Pack(&t.Accounts[i])
		}
	}
	{
		enc.PackLength(len(t.Waits))
		for i := range t.Waits {
			enc.Pack(&t.Waits[i])
		}
	}
	return enc.GetBytes()
}

func (t *Authority) Unpack(data []byte) int {
	dec := NewDecoder(data)
	{
		v := dec.UnpackUint32()
		t.Threshold = v
	}
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
	size += 4
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

type UpdateAuth struct {
	Account    Name
	Permission Name
	Parent     Name
	Auth       Authority
}

func NewUpdateAuth(account Name, permission Name, parent Name) *UpdateAuth {
	auth := &UpdateAuth{}
	auth.Account = account
	auth.Permission = permission
	auth.Parent = parent
	return auth
}

func (t *UpdateAuth) Pack() []byte {
	enc := NewEncoder(t.Size())
	enc.PackUint64(t.Account.N)
	enc.PackUint64(t.Permission.N)
	enc.PackUint64(t.Parent.N)
	enc.Pack(&t.Auth)
	return enc.GetBytes()
}

func (t *UpdateAuth) Unpack(data []byte) int {
	dec := NewDecoder(data)
	{
		v := dec.UnpackName()
		t.Account = v
	}
	{
		v := dec.UnpackName()
		t.Permission = v
	}
	{
		v := dec.UnpackName()
		t.Parent = v
	}
	dec.UnpackI(&t.Auth)
	return dec.Pos()
}

func (t *UpdateAuth) Size() int {
	size := 0
	size += 8
	size += 8
	size += 8
	size += t.Auth.Size()
	return size
}
