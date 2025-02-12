package types

import (
	"github.com/make-software/casper-go-sdk/v2/types/key"
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
	"github.com/make-software/casper-go-sdk/v2/types/serialization"
	"github.com/make-software/casper-go-sdk/v2/types/serialization/encoding"
)

const (
	TagFieldIndex         uint16 = 0
	PublicKeyVariantTag   uint8  = 0
	PublicKeyFieldIndex   uint16 = 1
	AccountHashVariantTag uint8  = 1
	AccountHashFieldIndex uint16 = 1
)

// InitiatorAddr the address of the initiator of a TransactionV1.
type InitiatorAddr struct {
	// The public key of the initiator
	PublicKey *keypair.PublicKey `json:"PublicKey,omitempty"`
	// The account hash derived from the public key of the initiator
	AccountHash *key.AccountHash `json:"AccountHash,omitempty"`
}

func (d InitiatorAddr) serializedFieldLengths() []int {
	if d.PublicKey != nil {
		return []int{
			encoding.U8SerializedLength,
			d.PublicKey.SerializedLength(),
		}
	} else {
		return []int{
			encoding.U8SerializedLength,
			key.ByteHashLen,
		}
	}
}

func (t *InitiatorAddr) SerializedLength() int {
	envelope := serialization.CallTableSerializationEnvelope{}
	return envelope.EstimateSize(t.serializedFieldLengths())
}

type InitiatorAddrFromBytesDecoder struct{}

func (addr *InitiatorAddrFromBytesDecoder) FromBytes(bytes []byte) (*InitiatorAddr, []byte, error) {
	envelope := &serialization.CallTableSerializationEnvelope{}
	binaryPayload, remainder, err := envelope.FromBytes(2, bytes)
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
	case PublicKeyVariantTag:
		if nextWindow == nil {
			return nil, nil, serialization.ErrFormatting
		}
		if err = nextWindow.VerifyIndex(PublicKeyFieldIndex); err != nil {
			return nil, nil, err
		}

		pubKey, finalWindow, err := serialization.DeserializeAndMaybeNext[keypair.PublicKey](nextWindow, &keypair.PublicKeyFromBytesDecoder{})
		if err != nil {
			return nil, nil, err
		}
		if finalWindow != nil {
			return nil, nil, serialization.ErrFormatting
		}

		return &InitiatorAddr{PublicKey: &pubKey}, remainder, nil
	case AccountHashVariantTag:
		if nextWindow == nil {
			return nil, nil, serialization.ErrFormatting
		}

		if err = nextWindow.VerifyIndex(AccountHashFieldIndex); err != nil {
			return nil, nil, err
		}

		decoder := encoding.SliceFromBytesDecoder[uint8, *encoding.U8FromBytesDecoder]{
			Decoder: &encoding.U8FromBytesDecoder{},
		}

		hash, finalWindow, err := serialization.DeserializeAndMaybeNext[[]uint8](nextWindow, &decoder)
		if err != nil {
			return nil, nil, err
		}
		if finalWindow != nil {
			return nil, nil, serialization.ErrFormatting
		}

		accountHash := key.AccountHash{
			Hash: key.Hash(hash),
		}
		return &InitiatorAddr{AccountHash: &accountHash}, remainder, nil
	default:
		return nil, nil, serialization.ErrFormatting
	}
}

func (d InitiatorAddr) Bytes() ([]byte, error) {
	builder, err := serialization.NewCallTableSerializationEnvelopeBuilder(d.serializedFieldLengths())
	if err != nil {
		return nil, err
	}

	if d.AccountHash != nil {
		if err = builder.AddField(TagFieldIndex, []byte{AccountHashVariantTag}); err != nil {
			return nil, err
		}

		if err = builder.AddField(AccountHashFieldIndex, d.AccountHash.Bytes()); err != nil {
			return nil, err
		}

	} else if d.PublicKey != nil {
		if err = builder.AddField(TagFieldIndex, []byte{PublicKeyVariantTag}); err != nil {
			return nil, err
		}

		if err = builder.AddField(PublicKeyFieldIndex, d.PublicKey.Bytes()); err != nil {
			return nil, err
		}
	}

	return builder.BinaryPayloadBytes()
}
