package key

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrInvalidByteCodeFormat = errors.New("invalid ByteCode format")
	ErrInvalidByteCodeKind   = errors.New("invalid ByteCodeKind")
)

type ByteCodeKind uint8

const (
	// EmptyKind Empty byte code
	EmptyKind ByteCodeKind = iota
	//V1CasperWasmKind  Byte code to be executed with the version 1 Casper execution engine.
	V1CasperWasmKind
)

func NewByteCodeKindFromByte(tag uint8) (ByteCodeKind, error) {
	kindTag := ByteCodeKind(tag)
	if kindTag != EmptyKind && kindTag != V1CasperWasmKind {
		return 0, ErrInvalidByteCodeKind
	}

	return kindTag, nil
}

const (
	EmptyPrefix  = "empty-"
	V1WasmPrefix = "v1-wasm-"
)

type ByteCode struct {
	V1CasperWasm *Hash

	isEmpty bool
}

func (h *ByteCode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	tmp, err := NewByteCode(s)
	if err != nil {
		return err
	}
	*h = tmp
	return nil
}

func (h ByteCode) IsEmpty() bool {
	return h.isEmpty
}

func (h ByteCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.ToPrefixedString())
}

func (h ByteCode) ToPrefixedString() string {
	switch {
	case h.V1CasperWasm != nil:
		return PrefixNameByteCode + V1WasmPrefix + h.V1CasperWasm.ToHex()
	case h.IsEmpty():
		emptyHash, _ := NewHashFromBytes(make([]byte, ByteHashLen))
		return PrefixNameByteCode + EmptyPrefix + emptyHash.ToHex()
	default:
		panic("Unexpected ByteCode type")
	}
}

func NewByteCode(source string) (ByteCode, error) {
	if strings.HasPrefix(source, V1WasmPrefix) {
		hexBytes, err := NewHash(strings.TrimPrefix(source, V1WasmPrefix))
		if err != nil {
			return ByteCode{}, err
		}
		return ByteCode{V1CasperWasm: &hexBytes}, nil
	} else if strings.HasPrefix(source, EmptyPrefix) {
		return ByteCode{isEmpty: true}, nil
	}

	return ByteCode{}, ErrInvalidByteCodeFormat
}

func NewByteCodeFromBuffer(buf *bytes.Buffer) (ByteCode, error) {
	tag, err := buf.ReadByte()
	if err != nil {
		return ByteCode{}, err
	}

	byteCodeKind, err := NewByteCodeKindFromByte(tag)
	if err != nil {
		return ByteCode{}, err
	}

	switch byteCodeKind {
	case EmptyKind:
		return ByteCode{isEmpty: true}, nil
	case V1CasperWasmKind:
		hash, err := NewByteHashFromBuffer(buf)
		if err != nil {
			return ByteCode{}, err
		}
		return ByteCode{V1CasperWasm: &hash}, nil
	}
	return ByteCode{}, ErrInvalidByteCodeFormat
}

func (h ByteCode) Bytes() []byte {
	switch {
	case h.V1CasperWasm != nil:
		res := make([]byte, 0)
		res = append(res, byte(V1CasperWasmKind))
		return append(res, h.V1CasperWasm.Bytes()...)
	case h.isEmpty:
		res := make([]byte, 0)
		res = append(res, byte(EmptyKind))
		return res
	default:
		panic("Unexpected ByteCode type")
	}
}
