package types

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/make-software/casper-go-sdk/v2/types/serialization"
	"github.com/make-software/casper-go-sdk/v2/types/serialization/encoding"
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
	TransactionEntryPointAddReservations    = "AddReservations"
	TransactionEntryPointCancelReservations = "CancelReservations"
	TransactionEntryPointCall               = "Call"
)

const (
	TransactionEntryPointCallTag = iota
	TransactionEntryPointCustomTag
	TransactionEntryPointTransferTag
	TransactionEntryPointAddBidTag
	TransactionEntryPointWithdrawBidTag
	TransactionEntryPointDelegateTag
	TransactionEntryPointUndelegateTag
	TransactionEntryPointRedelegateTag
	TransactionEntryPointActivateBidTag
	TransactionEntryPointChangeBidPublicKeyTag
	TransactionEntryPointAddReservationsTag
	TransactionEntryCancelReservationsTag
)

const CustomCustomIndex uint16 = 1

type TransactionEntryPoint struct {
	Custom *string
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
	// The `add_reservations` native entry point, used to add delegator to validator's reserve list.
	AddReservations *struct{}
	// The `cancel_reservations` native entry point, used to remove delegator from validator's reserve list.
	CancelReservations *struct{}
	// Used to call entry point call() in session transactions
	Call *struct{}
}

func (t *TransactionEntryPoint) SerializedLength() int {
	envelope := serialization.CallTableSerializationEnvelope{}
	return envelope.EstimateSize(t.serializedFieldLengths())
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
	case t.AddReservations != nil:
		return TransactionEntryPointAddReservationsTag
	case t.CancelReservations != nil:
		return TransactionEntryCancelReservationsTag
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
			Custom: &custom.Custom,
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
	case TransactionEntryPointAddReservations:
		entryPoint.AddReservations = &struct{}{}
	case TransactionEntryPointCancelReservations:
		entryPoint.CancelReservations = &struct{}{}
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
			Custom: *t.Custom,
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
	case t.AddReservations != nil:
		return json.Marshal(TransactionEntryPointAddReservations)
	case t.CancelReservations != nil:
		return json.Marshal(TransactionEntryPointCancelReservations)
	case t.Call != nil:
		return json.Marshal(TransactionEntryPointCall)
	default:
		return nil, errors.New("unknown entry point")
	}
}

func (t *TransactionEntryPoint) Bytes() ([]byte, error) {
	builder, err := serialization.NewCallTableSerializationEnvelopeBuilder(t.serializedFieldLengths())
	if err != nil {
		return nil, err
	}

	switch {
	case t.Call != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionEntryPointCallTag}); err != nil {
			return nil, err
		}

	case t.Custom != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionEntryPointCustomTag}); err != nil {
			return nil, err
		}

		customBytes, _ := encoding.NewStringToBytesEncoder(*t.Custom).Bytes()
		if err = builder.AddField(CustomCustomIndex, customBytes); err != nil {
			return nil, err
		}
	case t.Transfer != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionEntryPointTransferTag}); err != nil {
			return nil, err
		}
	case t.AddBid != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionEntryPointAddBidTag}); err != nil {
			return nil, err
		}
	case t.WithdrawBid != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionEntryPointWithdrawBidTag}); err != nil {
			return nil, err
		}
	case t.Delegate != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionEntryPointDelegateTag}); err != nil {
			return nil, err
		}
	case t.Undelegate != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionEntryPointUndelegateTag}); err != nil {
			return nil, err
		}
	case t.Redelegate != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionEntryPointRedelegateTag}); err != nil {
			return nil, err
		}
	case t.ActivateBid != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionEntryPointActivateBidTag}); err != nil {
			return nil, err
		}
	case t.ChangeBidPublicKey != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionEntryPointChangeBidPublicKeyTag}); err != nil {
			return nil, err
		}
	case t.AddReservations != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionEntryPointAddReservationsTag}); err != nil {
			return nil, err
		}
	case t.CancelReservations != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionEntryCancelReservationsTag}); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported variant")
	}
	return builder.BinaryPayloadBytes()
}

func (t *TransactionEntryPoint) serializedFieldLengths() []int {
	switch {
	case t.Custom != nil:
		return []int{
			encoding.U8SerializedLength,
			encoding.StringSerializedLength(*t.Custom),
		}
	default:
		return []int{
			encoding.U8SerializedLength,
		}
	}
}
