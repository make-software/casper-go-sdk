package types

import (
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/keypair"
)

type ValidatorWeightEraEnd struct {
	Validator keypair.PublicKey `json:"validator"`
	Weight    clvalue.UInt512   `json:"weight"`
}

type ValidatorWeightAuction struct {
	Validator keypair.PublicKey `json:"public_key"`
	Weight    clvalue.UInt512   `json:"weight"`
}
