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
