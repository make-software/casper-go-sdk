package types

import (
	"github.com/make-software/casper-go-sdk/v2/types/key"
	"github.com/make-software/casper-go-sdk/v2/types/serialization"
	"github.com/make-software/casper-go-sdk/v2/types/serialization/encoding"
)

const (
	PaymentLimitedPaymentAmountIndex uint16 = iota + 1
	PaymentLimitedGasPriceToleranceIndex
	PaymentLimitedStandardPaymentIndex
)

const (
	PaymentLimitedVariantTag uint8 = iota
	FixedVariantTag
	ReservedVariantTag
)

const (
	FixedGasPriceToleranceIndex           uint16 = 1
	FixedAdditionalComputationFactorIndex uint16 = 2
	ReservedReceiptIndex                  uint16 = 1
)

type PricingMode struct {
	// The original payment model, where the creator of the transaction specifies how much they will pay, at what gas price.
	Limited *LimitedMode `json:"Limited,omitempty"`
	// The cost of the transaction is determined by the cost table, per the transaction kind.
	Fixed *FixedMode `json:"Fixed,omitempty"`
	// The payment for this transaction was previously reserved, as proven by the receipt hash.
	Prepaid *PrepaidMode `json:"Prepaid,omitempty"`
}

func (d PricingMode) serializedFieldLengths() []int {
	switch {
	case d.Fixed != nil:
		return []int{
			encoding.U8SerializedLength,
			encoding.U8SerializedLength, // gas price tolerance
			encoding.U8SerializedLength, // additional computation factor
		}
	case d.Limited != nil:
		return []int{
			encoding.U8SerializedLength,
			encoding.U64SerializedLength,  // payment amount
			encoding.U8SerializedLength,   // gas price tolerance
			encoding.BoolSerializedLength, // standard_payment
		}
	case d.Prepaid != nil:
		return []int{
			encoding.U8SerializedLength,
			key.ByteHashLen,
		}
	default:
		return []int{}
	}
}

func (t *PricingMode) SerializedLength() int {
	envelope := serialization.CallTableSerializationEnvelope{}
	return envelope.EstimateSize(t.serializedFieldLengths())
}

type PricingModeFromBytesDecoder struct{}

func (addr *PricingModeFromBytesDecoder) FromBytes(bytes []byte) (*PricingMode, []byte, error) {
	envelope := &serialization.CallTableSerializationEnvelope{}
	binaryPayload, remainder, err := envelope.FromBytes(4, bytes)
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
	case PaymentLimitedVariantTag:
		if nextWindow == nil {
			return nil, nil, serialization.ErrFormatting
		}

		if err = nextWindow.VerifyIndex(PaymentLimitedPaymentAmountIndex); err != nil {
			return nil, nil, err
		}

		paymentAmount, nextWindow, err := serialization.DeserializeAndMaybeNext[uint64](nextWindow, &encoding.U64FromBytesDecoder{})
		if err != nil {
			return nil, nil, err
		}

		if nextWindow == nil {
			return nil, nil, serialization.ErrFormatting
		}

		if err = nextWindow.VerifyIndex(PaymentLimitedGasPriceToleranceIndex); err != nil {
			return nil, nil, err
		}

		gasPriceTolerance, nextWindow, err := serialization.DeserializeAndMaybeNext[uint8](nextWindow, &encoding.U8FromBytesDecoder{})
		if err != nil {
			return nil, nil, err
		}
		if nextWindow == nil {
			return nil, nil, serialization.ErrFormatting
		}

		if err = nextWindow.VerifyIndex(PaymentLimitedStandardPaymentIndex); err != nil {
			return nil, nil, err
		}

		standardPayment, _, err := serialization.DeserializeAndMaybeNext[bool](nextWindow, encoding.NewBoolFromBytesDecoder())
		if err != nil {
			return nil, nil, err
		}

		return &PricingMode{
			Limited: &LimitedMode{
				PaymentAmount:     paymentAmount,
				GasPriceTolerance: gasPriceTolerance,
				StandardPayment:   standardPayment,
			},
		}, remainder, nil

	case FixedVariantTag:
		if nextWindow == nil {
			return nil, nil, serialization.ErrFormatting
		}

		if err = nextWindow.VerifyIndex(FixedGasPriceToleranceIndex); err != nil {
			return nil, nil, err
		}
		gasPriceTolerance, nextWindow, err := serialization.DeserializeAndMaybeNext[uint8](nextWindow, &encoding.U8FromBytesDecoder{})
		if err != nil {
			return nil, nil, err
		}

		if nextWindow == nil {
			return nil, nil, serialization.ErrFormatting
		}
		if err = nextWindow.VerifyIndex(FixedAdditionalComputationFactorIndex); err != nil {
			return nil, nil, err
		}
		additionalComputationFactor, _, err := serialization.DeserializeAndMaybeNext[uint8](nextWindow, &encoding.U8FromBytesDecoder{})
		if err != nil {
			return nil, nil, err
		}

		return &PricingMode{
			Fixed: &FixedMode{
				GasPriceTolerance:           gasPriceTolerance,
				AdditionalComputationFactor: additionalComputationFactor,
			},
		}, remainder, nil

	case ReservedVariantTag:
		if nextWindow == nil {
			return nil, nil, serialization.ErrFormatting
		}

		if err = nextWindow.VerifyIndex(ReservedReceiptIndex); err != nil {
			return nil, nil, err
		}

		decoder := encoding.SliceFromBytesDecoder[uint8, *encoding.U8FromBytesDecoder]{
			Decoder: &encoding.U8FromBytesDecoder{},
		}

		receiptBytes, _, err := serialization.DeserializeAndMaybeNext[[]uint8](nextWindow, &decoder)
		if err != nil {
			return nil, nil, err
		}

		return &PricingMode{
			Prepaid: &PrepaidMode{
				Receipt: key.Hash(receiptBytes),
			},
		}, remainder, nil

	default:
		return nil, nil, serialization.ErrFormatting
	}
}

func (p *PricingMode) Bytes() ([]byte, error) {
	builder, err := serialization.NewCallTableSerializationEnvelopeBuilder(p.serializedFieldLengths())
	if err != nil {
		return nil, err
	}

	switch {
	case p.Limited != nil:
		if err = builder.AddField(TagFieldIndex, []byte{PaymentLimitedVariantTag}); err != nil {
			return nil, err
		}

		encodedPaymentAmount, _ := encoding.NewU64ToBytesEncoder(p.Limited.PaymentAmount).Bytes()
		if err = builder.AddField(PaymentLimitedPaymentAmountIndex, encodedPaymentAmount); err != nil {
			return nil, err
		}

		if err = builder.AddField(PaymentLimitedGasPriceToleranceIndex, []byte{p.Limited.GasPriceTolerance}); err != nil {
			return nil, err
		}

		var standardPaymentByte byte
		if p.Limited.StandardPayment {
			standardPaymentByte = 1
		}

		if err = builder.AddField(PaymentLimitedStandardPaymentIndex, []byte{standardPaymentByte}); err != nil {
			return nil, err
		}

	case p.Fixed != nil:
		if err = builder.AddField(TagFieldIndex, []byte{FixedVariantTag}); err != nil {
			return nil, err
		}

		if err = builder.AddField(FixedGasPriceToleranceIndex, []byte{p.Fixed.GasPriceTolerance}); err != nil {
			return nil, err
		}

		if err = builder.AddField(FixedAdditionalComputationFactorIndex, []byte{p.Fixed.AdditionalComputationFactor}); err != nil {
			return nil, err
		}

	case p.Prepaid != nil:
		if err = builder.AddField(TagFieldIndex, []byte{ReservedVariantTag}); err != nil {
			return nil, err
		}

		receiptBytes, _ := encoding.NewBytesToBytesEncoder(p.Prepaid.Receipt.Bytes()).Bytes()
		if err = builder.AddField(ReservedReceiptIndex, receiptBytes); err != nil {
			return nil, err
		}
	default:
		return nil, serialization.ErrFormatting
	}

	return builder.BinaryPayloadBytes()
}

type LimitedMode struct {
	// User-specified gas_price tolerance (minimum 1). This is interpreted to mean "do not include this transaction in a block if the current gas price is greater than this number"
	GasPriceTolerance uint8 `json:"gas_price_tolerance"`
	// User-specified payment amount.
	PaymentAmount uint64 `json:"payment_amount"`
	// Standard payment.
	StandardPayment bool `json:"standard_payment"`
}

type FixedMode struct {
	// User-specified gas_price tolerance (minimum 1). This is interpreted to mean "do not include this transaction in a block if the current gas price is greater than this number"
	GasPriceTolerance uint8 `json:"gas_price_tolerance"`
	/// User-specified additional computation factor (minimum 0). If "0" is provided,
	///  no additional logic is applied to the computation limit. Each value above "0"
	///  tells the node that it needs to treat the transaction as if it uses more gas
	///  than it's serialized size indicates. Each "1" will increase the "wasm lane"
	///  size bucket for this transaction by 1. So if the size of the transaction
	///  indicates bucket "0" and "additional_computation_factor = 2", the transaction
	///  will be treated as a "2".
	AdditionalComputationFactor uint8 `json:"additional_computation_factor"`
}

type PrepaidMode struct {
	// Pre-paid receipt
	Receipt key.Hash `json:"receipt"`
}
