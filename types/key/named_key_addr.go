package key

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"strings"
)

// NamedKeyAddr NamedKey address
type NamedKeyAddr struct {
	// The address of the entity.
	BaseAddr EntityAddr
	//  The bytes of the name
	NameBytes [32]byte
}

func NewNamedKeyAddr(source string) (NamedKeyAddr, error) {
	nameBytesData := source[strings.LastIndex(source, "-")+1:]
	source = source[:strings.LastIndex(source, "-")]

	rawBytes, err := hex.DecodeString(nameBytesData)
	if err != nil {
		return NamedKeyAddr{}, err
	}
	nameBytes := [32]byte{}
	copy(nameBytes[:], rawBytes)

	baseAddr, err := NewEntityAddr(strings.TrimPrefix(source, PrefixNameAddressableEntity))
	if err != nil {
		return NamedKeyAddr{}, err
	}
	return NamedKeyAddr{
		BaseAddr:  baseAddr,
		NameBytes: nameBytes,
	}, nil
}

func (h *NamedKeyAddr) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, err := NewNamedKeyAddr(s)
	if err != nil {
		return err
	}
	*h = val
	return nil
}

func (h NamedKeyAddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.ToPrefixedString())
}

func (h NamedKeyAddr) ToPrefixedString() string {
	res := PrefixNameNamedKey
	res += h.BaseAddr.ToPrefixedString()
	res += "-" + hex.EncodeToString(h.NameBytes[:])
	return res
}

func (h NamedKeyAddr) Bytes() []byte {
	res := make([]byte, 0, ByteHashLen)
	res = append(res, h.BaseAddr.Bytes()...)
	res = append(res, h.NameBytes[:]...)
	return res
}

func NewNamedKeyAddrFromBuffer(buf *bytes.Buffer) (NamedKeyAddr, error) {
	entityAddr, err := NewEntityAddrFromBuffer(buf)
	if err != nil {
		return NamedKeyAddr{}, err
	}

	nameBytes, err := NewByteHashFromBuffer(buf)
	if err != nil {
		return NamedKeyAddr{}, err
	}

	return NamedKeyAddr{
		BaseAddr:  entityAddr,
		NameBytes: nameBytes,
	}, nil
}
