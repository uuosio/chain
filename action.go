package chain

import (
	"github.com/uuosio/chain/eosio"
)

var gActionCache [][]byte
var gNotifyCache []Name

// Read current action data
func ReadActionData() []byte {
	return eosio.ReadActionData()
}

// Get the length of the current action's data field
func ActionDataSize() uint32 {
	return eosio.ActionDataSize()
}

// Add the specified account to set of accounts to be notified
func RequireRecipient(name Name) {
	eosio.RequireRecipient(name.N)
}

// Verifies that name exists in the set of provided auths on a action. Throws if not found.
func RequireAuth(name Name) {
	eosio.RequireAuth(name.N)
}

// Verifies that name has auth.
func HasAuth(name Name) bool {
	return eosio.HasAuth(name.N)
}

// Verifies that name exists in the set of provided auths on a action. Throws if not found.
func RequireAuth2(name Name, permission Name) {
	eosio.RequireAuth2(name.N, permission.N)
}

// Verifies that name is an existing account.
func IsAccount(name Name) bool {
	return eosio.IsAccount(name.N)
}

// Send an inline action in the context of this action's parent transaction
func SendInline(data []byte) {
	eosio.SendInline(data)
}

// Send an inline context free action in the context of this action's parent transaction
func SendContextFreeInline(data []byte) {
	eosio.SendContextFreeInline(data)
}

// Returns the time in microseconds from 1970 of the publication_time
func PublicationTime() uint64 {
	return eosio.PublicationTime()
}

// Get the current receiver of the action
func CurrentReceiver() Name {
	n := eosio.CurrentReceiver()
	return Name{n}
}

// Set the action return value which will be included in the action_receipt
func SetActionReturnValue(return_value []byte) {
	eosio.SetActionReturnValue(return_value)
}

type Action struct {
	Account       Name
	Name          Name
	Authorization []*PermissionLevel
	Data          []byte
}

func NewAction(perm *PermissionLevel, account Name, name Name, args ...interface{}) *Action {
	return NewActionEx([]*PermissionLevel{perm}, account, name, args...)
}

func NewActionEx(perms []*PermissionLevel, account Name, name Name, args ...interface{}) *Action {
	a := &Action{}
	a.Account = account
	a.Name = name

	if perms != nil {
		for _, perm := range perms {
			a.Authorization = append(a.Authorization, perm)
		}
	}

	if len(args) == 0 {
		a.Data = []byte{}
		return a
	}

	size := 0
	for _, v := range args {
		n := CalcPackedSize(v)
		size += n
	}
	enc := NewEncoder(size)
	for _, arg := range args {
		enc.Pack(arg)
	}
	a.Data = enc.GetBytes()

	return a
}

func (a *Action) Size() int {
	return 8 + 8 + 5 + len(a.Authorization)*8 + 5 + len(a.Data)
}

func (a *Action) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()
	enc.PackName(a.Account)
	enc.PackName(a.Name)
	enc.PackLength(len(a.Authorization))
	for _, v := range a.Authorization {
		v.Pack(enc)
	}
	enc.PackBytes(a.Data)
	return enc.GetSize() - oldSize
}

func (a *Action) Unpack(b []byte) int {
	dec := NewDecoder(b)
	a.Account = dec.UnpackName()
	a.Name = dec.UnpackName()
	length := dec.UnpackLength()
	a.Authorization = make([]*PermissionLevel, length)
	for i := 0; i < length; i++ {
		a.Authorization[i] = new(PermissionLevel)
		dec.Unpack(a.Authorization[i])
	}
	dec.Unpack(&a.Data)
	return dec.Pos()
}

func (a *Action) Print() {
	Print("{")
	Print(a.Account, a.Name)
	Print("[")
	for _, v := range a.Authorization {
		Print("[", v.Actor, v.Permission, "]")
	}
	Print("]")
	Print(a.Data)
	Print("}")
}

func (a *Action) AddPermission(actor Name, permission Name) {
	a.Authorization = append(a.Authorization, &PermissionLevel{actor, permission})
}

func (a *Action) Send() {
	data := EncoderPack(a)
	SendInline(data)
}

// send action directly, no cache
func (a *Action) SendEx() {
	data := EncoderPack(a)
	SendInline(data)
}

type GetCodeHashResult struct {
	StructVersion VarUint32
	CodeSequence  uint64
	CodeHash      Checksum256
	VMType        uint8
	VMVersion     uint8
}

func (t *GetCodeHashResult) Pack(enc *Encoder) int {
	oldSize := enc.GetSize()

	t.StructVersion.Pack(enc)
	enc.PackUint64(t.CodeSequence)
	t.CodeHash.Pack(enc)
	enc.PackUint8(t.VMType)
	enc.PackUint8(t.VMVersion)

	return enc.GetSize() - oldSize
}

func (t *GetCodeHashResult) Unpack(data []byte) int {
	dec := NewDecoder(data)
	dec.Unpack(&t.StructVersion)
	dec.Unpack(&t.CodeSequence)
	dec.Unpack(&t.CodeHash)
	dec.Unpack(&t.VMType)
	dec.Unpack(&t.VMVersion)
	return dec.Pos()
}

func (t *GetCodeHashResult) Size() int {
	return t.StructVersion.Size() + 8 + 32 + 1 + 1
}

func GetCodeHash(account Name) Checksum256 {
	data := eosio.GetCodeHash(account.N)
	ret := GetCodeHashResult{}
	ret.Unpack(data)
	return ret.CodeHash
}
