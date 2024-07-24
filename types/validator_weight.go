package types

import (
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

type ValidatorWeightEraEnd struct {
	Validator keypair.PublicKey `json:"validator"`
	Weight    clvalue.UInt512   `json:"weight"`
}

type ValidatorWeightAuction struct {
	Validator keypair.PublicKey `json:"public_key"`
	Weight    clvalue.UInt512   `json:"weight"`
}
