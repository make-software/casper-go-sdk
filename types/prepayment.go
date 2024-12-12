package types

import "github.com/make-software/casper-go-sdk/v2/types/key"

// PrepaymentKind Container for bytes recording location, type and data for a gas pre payment
type PrepaymentKind struct {
	Receipt        key.Hash `json:"receipt"`
	PrepaymentData HexBytes `json:"prepayment_data"`
	PrepaymentKind uint8    `json:"prepayment_kind"`
}
