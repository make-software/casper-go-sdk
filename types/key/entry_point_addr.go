package key

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrInvalidEntryPointTag    = errors.New("invalid EntryPointTag")
	ErrInvalidEntryPointFormat = errors.New("invalid EntryPoint format")
)

type EntryPointTag uint8

const (
	V1EntryPoint EntryPointTag = iota
	V2EntryPoint
)

func NewEntryPointTagFromByte(tag uint8) (EntryPointTag, error) {
	entryPointTag := EntryPointTag(tag)
	if entryPointTag != V1EntryPoint && entryPointTag != V2EntryPoint {
		return 0, ErrInvalidEntryPointTag
	}

	return V2EntryPoint, nil
}

const SelectorBytesLen = 8

const (
	V1Prefix = "v1-"
	V2Prefix = "v2-"
)

// VmCasperV1 The address for a V1 Entrypoint.
type VmCasperV1 struct {
	// The addr of the entity.
	EntityAddr EntityAddr
	// The 32 byte hash of the name of the entry point
	NameBytes [32]byte
}

// VmCasperV2 The address for a V2 entrypoint
type VmCasperV2 struct {
	// The addr of the entity.
	EntityAddr EntityAddr
	// The selector.
	Selector uint32
}

type EntryPointAddr struct {
	VmCasperV1 *VmCasperV1
	VmCasperV2 *VmCasperV2
}

func (h *EntryPointAddr) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	tmp, err := NewEntryPointAddr(s)
	if err != nil {
		return err
	}
	*h = tmp
	return nil
}

func (h EntryPointAddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.ToPrefixedString())
}

func (h EntryPointAddr) ToPrefixedString() string {
	switch {
	case h.VmCasperV1 != nil:
		return PrefixEntryPoint + V1Prefix + h.VmCasperV1.EntityAddr.ToPrefixedString() + "-" + hex.EncodeToString(h.VmCasperV1.NameBytes[:])
	case h.VmCasperV2 != nil:
		selector := make([]byte, 0)
		binary.LittleEndian.PutUint32(selector, h.VmCasperV2.Selector)
		return PrefixEntryPoint + V2Prefix + h.VmCasperV2.EntityAddr.ToPrefixedString() + "-" + hex.EncodeToString(selector)
	default:
		panic("Unexpected ByteCode type")
	}
}

func NewEntryPointAddr(source string) (EntryPointAddr, error) {
	lastIndex := strings.LastIndex(source, "-")
	if lastIndex == -1 {
		return EntryPointAddr{}, ErrInvalidByteCodeFormat
	}
	prefix := source[:lastIndex]
	data := source[lastIndex+1:]

	var (
		baseAddr EntityAddr
		err      error
	)

	rawBytes, err := hex.DecodeString(data)
	if err != nil {
		return EntryPointAddr{}, err
	}

	switch {
	case strings.HasPrefix(prefix, V1Prefix):
		prefix = strings.TrimPrefix(prefix, V1Prefix)
		var nameBytes [32]byte
		copy(nameBytes[:], rawBytes)

		baseAddr, err = NewEntityAddr(strings.TrimPrefix(prefix, PrefixNameAddressableEntity))
		if err != nil {
			return EntryPointAddr{}, err
		}

		return EntryPointAddr{
			VmCasperV1: &VmCasperV1{
				EntityAddr: baseAddr,
				NameBytes:  nameBytes,
			},
		}, nil
	case strings.HasPrefix(prefix, V2Prefix):
		prefix = strings.TrimPrefix(prefix, V2Prefix)
		baseAddr, err = NewEntityAddr(strings.TrimPrefix(prefix, PrefixNameAddressableEntity))
		if err != nil {
			return EntryPointAddr{}, err
		}

		selector := binary.LittleEndian.Uint32(rawBytes)

		return EntryPointAddr{
			VmCasperV2: &VmCasperV2{
				EntityAddr: baseAddr,
				Selector:   selector,
			},
		}, nil
	default:
		return EntryPointAddr{}, ErrInvalidEntryPointFormat
	}
}

func NewEntryPointAddrFromBuffer(buf *bytes.Buffer) (EntryPointAddr, error) {
	tag, err := buf.ReadByte()
	if err != nil {
		return EntryPointAddr{}, err
	}

	byteCodeKind, err := NewEntryPointTagFromByte(tag)
	if err != nil {
		return EntryPointAddr{}, err
	}

	entityAddr, err := NewEntityAddrFromBuffer(buf)
	if err != nil {
		return EntryPointAddr{}, err
	}

	switch byteCodeKind {
	case V1EntryPoint:
		nameBytes, err := NewByteHashFromBuffer(buf)
		if err != nil {
			return EntryPointAddr{}, err
		}

		return EntryPointAddr{
			VmCasperV1: &VmCasperV1{
				EntityAddr: entityAddr,
				NameBytes:  nameBytes,
			},
		}, nil
	case V2EntryPoint:
		return EntryPointAddr{
			VmCasperV2: &VmCasperV2{
				EntityAddr: entityAddr,
				Selector:   binary.LittleEndian.Uint32(buf.Next(SelectorBytesLen)),
			},
		}, nil
	}
	return EntryPointAddr{}, ErrInvalidEntryPointFormat
}

func (h EntryPointAddr) Bytes() []byte {
	switch {
	case h.VmCasperV1 != nil:
		res := make([]byte, 0)
		res = append(res, byte(V1EntryPoint))
		res = append(res, h.VmCasperV1.EntityAddr.Bytes()...)
		return append(res, h.VmCasperV1.NameBytes[:]...)
	case h.VmCasperV2 != nil:
		res := make([]byte, 0)
		res = append(res, byte(V2EntryPoint))
		res = append(res, h.VmCasperV2.EntityAddr.Bytes()...)
		binary.LittleEndian.PutUint32(res, h.VmCasperV2.Selector)
		return res
	default:
		panic("Unexpected EntryPointAddr type")
	}
}
