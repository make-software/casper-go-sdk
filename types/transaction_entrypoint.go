package types

import (
	"encoding/json"
	"errors"

	"github.com/make-software/casper-go-sdk/types/clvalue"
)

const (
	TransactionEntryPointCustom             = "Custom"
	TransactionEntryPointTransfer           = "Transfer"
	TransactionEntryPointAddBid             = "AddBid"
	TransactionEntryPointWithdrawBid        = "WithdrawBid"
	TransactionEntryPointDelegate           = "Delegate"
	TransactionEntryPointUndelegate         = "Undelegate"
	TransactionEntryPointRedelegate         = "Redelegate"
	TransactionEntryPointActivateBid        = "ActivateBid"
	TransactionEntryPointChangeBidPublicKey = "ChangeBidPublicKey"
	TransactionEntryPointCall               = "Call"
)

const (
	TransactionEntryPointCustomTag = iota
	TransactionEntryPointTransferTag
	TransactionEntryPointAddBidTag
	TransactionEntryPointWithdrawBidTag
	TransactionEntryPointDelegateTag
	TransactionEntryPointUndelegateTag
	TransactionEntryPointRedelegateTag
	TransactionEntryPointActivateBidTag
	TransactionEntryPointChangeBidPublicKeyTag
	TransactionEntryPointCallTag
)

type TransactionEntryPoint struct {
	Custom *struct {
		Type string
	}
	// The `transfer` native entry point, used to transfer `Motes` from a source purse to a target purse.
	Transfer *struct{}
	// The `add_bid` native entry point, used to create or top off a bid purse.
	AddBid *struct{}
	// The `withdraw_bid` native entry point, used to decrease a stake.
	WithdrawBid *struct{}
	// The `delegate` native entry point, used to add a new delegator or increase an existing delegator's stake.
	Delegate *struct{}
	// The `undelegate` native entry point, used to reduce a delegator's stake or remove the delegator if the remaining stake is 0.
	Undelegate *struct{}
	// The `redelegate` native entry point, used to reduce a delegator's stake or remove the delegator if
	// the remaining stake is 0, and after the unbonding delay, automatically delegate to a new validator.
	Redelegate *struct{}
	// The `activate_bid` native entry point, used to used to reactivate an inactive bid.
	ActivateBid *struct{}
	// The `change_bid_public_key` native entry point, used to change a bid's public key.
	ChangeBidPublicKey *struct{}
	// Used to call entry point call() in session transactions
	Call *struct{}
}

func (t *TransactionEntryPoint) Bytes() []byte {
	result := make([]byte, 0, 2)
	result = append(result, t.Tag())

	if t.Custom != nil {
		result = append(result, clvalue.NewCLString(t.Custom.Type).Bytes()...)
	}
	return result
}

func (t *TransactionEntryPoint) Tag() byte {
	switch {
	case t.Custom != nil:
		return TransactionEntryPointCustomTag
	case t.Transfer != nil:
		return TransactionEntryPointTransferTag
	case t.AddBid != nil:
		return TransactionEntryPointAddBidTag
	case t.WithdrawBid != nil:
		return TransactionEntryPointWithdrawBidTag
	case t.Delegate != nil:
		return TransactionEntryPointDelegateTag
	case t.Undelegate != nil:
		return TransactionEntryPointUndelegateTag
	case t.Redelegate != nil:
		return TransactionEntryPointRedelegateTag
	case t.ActivateBid != nil:
		return TransactionEntryPointActivateBidTag
	case t.ChangeBidPublicKey != nil:
		return TransactionEntryPointChangeBidPublicKeyTag
	case t.Call != nil:
		return TransactionEntryPointCallTag
	default:
		return 0
	}
}

func (t *TransactionEntryPoint) UnmarshalJSON(data []byte) error {
	var custom struct {
		Custom string `json:"Custom"`
	}
	if err := json.Unmarshal(data, &custom); err == nil {
		*t = TransactionEntryPoint{
			Custom: &struct{ Type string }{Type: custom.Custom},
		}
		return nil
	}

	var key string
	if err := json.Unmarshal(data, &key); err != nil {
		return err
	}

	var entryPoint TransactionEntryPoint
	switch key {
	case TransactionEntryPointTransfer:
		entryPoint.Transfer = &struct{}{}
	case TransactionEntryPointAddBid:
		entryPoint.AddBid = &struct{}{}
	case TransactionEntryPointWithdrawBid:
		entryPoint.WithdrawBid = &struct{}{}
	case TransactionEntryPointDelegate:
		entryPoint.Delegate = &struct{}{}
	case TransactionEntryPointUndelegate:
		entryPoint.Undelegate = &struct{}{}
	case TransactionEntryPointRedelegate:
		entryPoint.Redelegate = &struct{}{}
	case TransactionEntryPointActivateBid:
		entryPoint.ActivateBid = &struct{}{}
	case TransactionEntryPointChangeBidPublicKey:
		entryPoint.ChangeBidPublicKey = &struct{}{}
	case TransactionEntryPointCall:
		entryPoint.Call = &struct{}{}
	}

	*t = entryPoint
	return nil
}

func (t TransactionEntryPoint) MarshalJSON() ([]byte, error) {
	if t.Custom != nil {
		return json.Marshal(struct {
			Custom string `json:"Custom"`
		}{
			Custom: t.Custom.Type,
		})
	}

	switch {
	case t.Transfer != nil:
		return json.Marshal(TransactionEntryPointTransfer)
	case t.AddBid != nil:
		return json.Marshal(TransactionEntryPointAddBid)
	case t.WithdrawBid != nil:
		return json.Marshal(TransactionEntryPointWithdrawBid)
	case t.Delegate != nil:
		return json.Marshal(TransactionEntryPointDelegate)
	case t.Undelegate != nil:
		return json.Marshal(TransactionEntryPointUndelegate)
	case t.Redelegate != nil:
		return json.Marshal(TransactionEntryPointRedelegate)
	case t.ActivateBid != nil:
		return json.Marshal(TransactionEntryPointActivateBid)
	case t.ChangeBidPublicKey != nil:
		return json.Marshal(TransactionEntryPointChangeBidPublicKey)
	case t.Call != nil:
		return json.Marshal(TransactionEntryPointCall)
	default:
		return nil, errors.New("unknown entry point")
	}
}
