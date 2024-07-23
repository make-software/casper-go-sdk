package types

import (
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/key"
)

const (
	PricingModeClassicTag = iota
	PricingModeFixedTag
	PricingModeReservedTag
)

type PricingMode struct {
	// The original payment model, where the creator of the transaction specifies how much they will pay, at what gas price.
	Classic *ClassicMode `json:"Classic,omitempty"`
	// The cost of the transaction is determined by the cost table, per the transaction kind.
	Fixed *FixedMode `json:"Fixed,omitempty"`
	// The payment for this transaction was previously reserved, as proven by the receipt hash.
	Reserved *ReservedMode `json:"reserved,omitempty"`
}

func (d PricingMode) Bytes() []byte {
	result := make([]byte, 0, 2)
	if d.Classic != nil {
		result = append(result, PricingModeClassicTag)
		result = append(result, clvalue.NewCLUInt64(d.Classic.PaymentAmount).Bytes()...)
		result = append(result, d.Classic.GasPriceTolerance)
		if d.Classic.StandardPayment {
			result = append(result, 1)
		} else {
			result = append(result, 0)
		}
	} else if d.Fixed != nil {
		result = append(result, PricingModeFixedTag)
		result = append(result, d.Fixed.GasPriceTolerance)
	} else if d.Reserved != nil {
		result = append(result, PricingModeReservedTag)
		result = append(result, d.Reserved.Receipt.Bytes()...)
	}

	return result
}

type ClassicMode struct {
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
}

type ReservedMode struct {
	// Pre-paid receipt
	Receipt key.Hash `json:"receipt"`
}
