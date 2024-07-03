package key

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

type UrefAccess = byte

var (
	ErrIncorrectUrefFormat = errors.New("incorrect uref format")
)

const (
	UrefAccessNone = iota
	UrefAccessRead
	UrefAccessWrite
	UrefAccessAdd
	UrefAccessReadWrite
	UrefAccessReadAdd
	UrefAccessAddWrite
	UrefAccessReadAddWrite
)

type URef struct {
	data   [ByteHashLen]byte
	access UrefAccess
}

func (v URef) Bytes() []byte {
	return append(v.data[:], []byte{v.access}...)
}

func (v URef) String() string {
	return v.ToPrefixedString()
}

func (v URef) ToPrefixedString() string {
	return PrefixNameURef + hex.EncodeToString(v.data[:]) + "-0" + hex.EncodeToString([]byte{v.access})
}

func (v *URef) SetAccess(access UrefAccess) {
	v.access = access
}

func (v *URef) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	tmp, err := NewURef(s)
	if err != nil {
		return err
	}
	*v = tmp
	return nil
}

func (v *URef) UnmarshalText(text []byte) error {
	tmp, err := NewURef(string(text))
	if err != nil {
		return err
	}
	*v = tmp
	return nil
}

func (v URef) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.ToPrefixedString())
}

func (v *URef) GobDecode(i []byte) error {
	tmp, err := NewURefFromBytes(i)
	if err != nil {
		return err
	}
	*v = tmp
	return nil
}

func (v URef) GobEncode() ([]byte, error) {
	return v.Bytes(), nil
}

func NewURef(source string) (res URef, err error) {
	parts := strings.Split(source, "-")
	if len(parts) != 3 {
		return res, ErrIncorrectUrefFormat
	}

	payloadInBytes, err := hex.DecodeString(parts[1])
	if err != nil {
		return res, err
	}
	accessInBytes, err := hex.DecodeString(strings.TrimPrefix(parts[2], "0"))
	if err != nil {
		return res, err
	}
	return NewURefFromBytes(append(payloadInBytes, accessInBytes...))
}

func NewURefFromBytes(source []byte) (URef, error) {
	return NewURefFromBuffer(bytes.NewBuffer(source))
}

func NewURefFromBuffer(buffer *bytes.Buffer) (res URef, err error) {
	var payload [32]byte
	_, err = buffer.Read(payload[:])
	if err != nil {
		return res, err
	}
	access, err := buffer.ReadByte()
	if err != nil {
		return res, err
	}
	return URef{
		access: access,
		data:   payload,
	}, nil
}

func (v URef) Value() (driver.Value, error) {
	return v.Bytes(), nil
}

func (v *URef) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan value type")
	}

	dst := make([]byte, len(data))
	copy(dst, data)

	tmp, err := NewURefFromBytes(dst)
	if err != nil {
		return err
	}
	*v = tmp
	return nil
}
