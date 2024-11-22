package types

import (
	"errors"

	"github.com/make-software/casper-go-sdk/v2/types/key"
	"github.com/make-software/casper-go-sdk/v2/types/serialization"
	"github.com/make-software/casper-go-sdk/v2/types/serialization/encoding"
)

const (
	ByHashVariant   uint8  = 0
	ByHashHashIndex uint16 = 1

	ByNameVariant   uint8  = 1
	ByNameNameIndex uint16 = 1

	ByPackageHashVariant      uint8  = 2
	ByPackageHashAddrIndex    uint16 = 1
	ByPackageHashVersionIndex uint16 = 2

	ByPackageNameVariant      uint8  = 3
	ByPackageNameNameIndex    uint16 = 1
	ByPackageNameVersionIndex uint16 = 2
)

type TransactionInvocationTarget struct {
	// Hex-encoded entity address identifying the invocable entity.
	ByHash *key.Hash `json:"ByHash,omitempty"`
	// The alias identifying the invocable entity.
	ByName *string `json:"ByName,omitempty"`
	// The address and optional version identifying the package.
	ByPackageHash *ByPackageHashInvocationTarget `json:"ByPackageHash,omitempty"`
	// The alias and optional version identifying the package.
	ByPackageName *ByPackageNameInvocationTarget `json:"ByPackageName,omitempty"`
}

// ByPackageHashInvocationTarget The address and optional version identifying the package.
type ByPackageHashInvocationTarget struct {
	Addr    key.Hash `json:"addr"`
	Version *uint32  `json:"version"`
}

// ByPackageNameInvocationTarget The alias and optional version identifying the package.
type ByPackageNameInvocationTarget struct {
	Name    string  `json:"name"`
	Version *uint32 `json:"version"`
}

type TransactionInvocationTargetFromBytesDecoder struct{}

func (d *TransactionInvocationTargetFromBytesDecoder) FromBytes(bytes []byte) (*TransactionInvocationTarget, []byte, error) {
	envelope := &serialization.CallTableSerializationEnvelope{}
	binaryPayload, remainder, err := envelope.FromBytes(3, bytes)
	if err != nil {
		return nil, nil, err
	}

	window, err := binaryPayload.StartConsuming()
	if err != nil || window == nil {
		return nil, nil, serialization.ErrFormatting
	}

	if err = window.VerifyIndex(TagFieldIndex); err != nil {
		return nil, nil, err
	}

	tag, nextWindow, err := serialization.DeserializeAndMaybeNext[uint8](window, &encoding.U8FromBytesDecoder{})
	if err != nil {
		return nil, nil, err
	}

	switch tag {
	case ByHashVariant:
		if nextWindow == nil {
			return nil, nil, serialization.ErrFormatting
		}
		if err = nextWindow.VerifyIndex(ByHashHashIndex); err != nil {
			return nil, nil, err
		}

		hash, finalWindow, err := serialization.DeserializeAndMaybeNext(nextWindow, &key.HashFromBytesDecoder{})
		if err != nil {
			return nil, nil, err
		}
		if finalWindow != nil {
			return nil, nil, serialization.ErrFormatting
		}

		return &TransactionInvocationTarget{ByHash: &hash}, remainder, nil

	case ByNameVariant:
		if nextWindow == nil {
			return nil, nil, serialization.ErrFormatting
		}
		if err = nextWindow.VerifyIndex(ByNameNameIndex); err != nil {
			return nil, nil, err
		}

		name, finalWindow, err := serialization.DeserializeAndMaybeNext(nextWindow, &encoding.StringFromBytesDecoder{})
		if err != nil {
			return nil, nil, err
		}
		if finalWindow != nil {
			return nil, nil, serialization.ErrFormatting
		}

		return &TransactionInvocationTarget{ByName: &name}, remainder, nil

	case ByPackageHashVariant:
		if nextWindow == nil {
			return nil, nil, serialization.ErrFormatting
		}
		if err = nextWindow.VerifyIndex(ByPackageHashAddrIndex); err != nil {
			return nil, nil, err
		}

		bytesDecoder := encoding.SliceFromBytesDecoder[uint8, *encoding.U8FromBytesDecoder]{
			Decoder: &encoding.U8FromBytesDecoder{},
		}

		addr, nextWindow, err := serialization.DeserializeAndMaybeNext[[]uint8](nextWindow, &bytesDecoder)
		if err != nil {
			return nil, nil, err
		}

		if nextWindow == nil {
			return nil, nil, serialization.ErrFormatting
		}
		if err = nextWindow.VerifyIndex(ByPackageHashVersionIndex); err != nil {
			return nil, nil, err
		}

		decoder := encoding.OptionFromBytesDecoder[uint32, *encoding.U32FromBytesDecoder]{
			Decoder: &encoding.U32FromBytesDecoder{},
		}

		optionVersion, finalWindow, err := serialization.DeserializeAndMaybeNext[encoding.Option[uint32]](nextWindow, &decoder)
		if err != nil {
			return nil, nil, err
		}

		if finalWindow != nil {
			return nil, nil, serialization.ErrFormatting
		}

		return &TransactionInvocationTarget{ByPackageHash: &ByPackageHashInvocationTarget{
			Addr:    key.Hash(addr),
			Version: optionVersion.Some,
		}}, remainder, nil

	case ByPackageNameVariant:
		if nextWindow == nil {
			return nil, nil, serialization.ErrFormatting
		}
		if err = nextWindow.VerifyIndex(ByPackageNameNameIndex); err != nil {
			return nil, nil, err
		}

		name, nextWindow, err := serialization.DeserializeAndMaybeNext[string](nextWindow, &encoding.StringFromBytesDecoder{})
		if err != nil {
			return nil, nil, err
		}
		if nextWindow == nil {
			return nil, nil, serialization.ErrFormatting
		}
		if err = nextWindow.VerifyIndex(ByPackageNameVersionIndex); err != nil {
			return nil, nil, err
		}

		decoder := encoding.OptionFromBytesDecoder[uint32, *encoding.U32FromBytesDecoder]{
			Decoder: &encoding.U32FromBytesDecoder{},
		}
		optionVersion, finalWindow, err := serialization.DeserializeAndMaybeNext[encoding.Option[uint32]](nextWindow, &decoder)
		if err != nil {
			return nil, nil, err
		}
		if finalWindow != nil {
			return nil, nil, serialization.ErrFormatting
		}

		return &TransactionInvocationTarget{ByPackageName: &ByPackageNameInvocationTarget{
			Version: optionVersion.Some,
			Name:    name,
		}}, remainder, nil

	default:
		return nil, nil, serialization.ErrFormatting
	}
}

func (t *TransactionInvocationTarget) Bytes() ([]byte, error) {
	builder, err := serialization.NewCallTableSerializationEnvelopeBuilder(t.serializedFieldLengths())
	if err != nil {
		return nil, err
	}

	switch {
	case t.ByHash != nil:
		if err = builder.AddField(TagFieldIndex, []byte{ByHashVariant}); err != nil {
			return nil, err
		}

		byHashBytes, _ := encoding.NewStringToBytesEncoder(t.ByHash.String()).Bytes()
		if err = builder.AddField(ByHashHashIndex, byHashBytes); err != nil {
			return nil, err
		}
	case t.ByName != nil:
		if err = builder.AddField(TagFieldIndex, []byte{ByNameVariant}); err != nil {
			return nil, err
		}

		byNameBytes, _ := encoding.NewStringToBytesEncoder(*t.ByName).Bytes()
		if err = builder.AddField(ByNameNameIndex, byNameBytes); err != nil {
			return nil, err
		}

	case t.ByPackageHash != nil:
		if err = builder.AddField(TagFieldIndex, []byte{ByPackageHashVariant}); err != nil {
			return nil, err
		}

		byPackageBytes, _ := encoding.NewStringToBytesEncoder(t.ByPackageHash.Addr.String()).Bytes()
		if err = builder.AddField(ByPackageHashAddrIndex, byPackageBytes); err != nil {
			return nil, err
		}

		versionBytes, _ := encoding.NewU32ToBytesEncoder(*t.ByPackageHash.Version).Bytes()
		if err = builder.AddField(ByPackageHashVersionIndex, versionBytes); err != nil {
			return nil, err
		}

	case t.ByPackageName != nil:
		if err = builder.AddField(TagFieldIndex, []byte{ByPackageNameVariant}); err != nil {
			return nil, err
		}
		if err = builder.AddField(ByPackageNameNameIndex, []byte(t.ByPackageName.Name)); err != nil {
			return nil, err
		}

		versionBytes, _ := encoding.NewU32ToBytesEncoder(*t.ByPackageName.Version).Bytes()
		if err = builder.AddField(ByPackageNameVersionIndex, versionBytes); err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("unknown transaction invocation target")
	}

	return builder.BinaryPayloadBytes()
}

func (t *TransactionInvocationTarget) SerializedLength() int {
	envelope := serialization.CallTableSerializationEnvelope{}
	return envelope.EstimateSize(t.serializedFieldLengths())
}

func (t *TransactionInvocationTarget) serializedFieldLengths() []int {
	switch {
	case t.ByHash != nil:
		return []int{
			encoding.U8SerializedLength,
		}
	case t.ByName != nil:
		return []int{}
	case t.ByPackageHash != nil:
		return []int{}
	case t.ByPackageName != nil:
		return []int{}
	default:
		return []int{}
	}
}
