package types

import (
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

const (
	DefaultMinimumDelegationAmount = 500 * 1_000_000_000
	DefaultMaximumDelegationAmount = 1_000_000_000 * 1_000_000_000
)

type AuctionState struct {
	// All bids contained within a vector.
	Bids          []BidKindWrapper `json:"bids"`
	BlockHeight   uint64           `json:"block_height"`
	EraValidators []EraValidators  `json:"era_validators"`
	StateRootHash string           `json:"state_root_hash"`
}

// AuctionStateV1 is a data structure summarizing auction contract data (version V1).
type AuctionStateV1 struct {
	// All bids contained within a vector.
	Bids          []PublicKeyAndBid `json:"bids"`
	BlockHeight   uint64            `json:"block_height"`
	EraValidators []EraValidators   `json:"era_validators"`
	StateRootHash string            `json:"state_root_hash"`
}

type BidKindWrapper struct {
	PublicKey keypair.PublicKey `json:"public_key"`
	Bid       BidKind           `json:"bid"`
}

// AuctionStateV2 is a data structure summarizing auction contract data.
type AuctionStateV2 struct {
	// All bids contained within a vector.
	Bids          []BidKindWrapper `json:"bids"`
	BlockHeight   uint64           `json:"block_height"`
	EraValidators []EraValidators  `json:"era_validators"`
	StateRootHash string           `json:"state_root_hash"`
}

func NewAuctionStateFromV1(v1 AuctionStateV1) AuctionState {
	bids := make([]BidKindWrapper, 0, len(v1.Bids))
	for _, bid := range v1.Bids {
		bids = append(bids, BidKindWrapper{
			PublicKey: bid.PublicKey,
			Bid: BidKind{
				Validator: &ValidatorBid{
					ValidatorPublicKey:      bid.Bid.ValidatorPublicKey,
					BondingPurse:            bid.Bid.BondingPurse,
					DelegationRate:          bid.Bid.DelegationRate,
					Inactive:                bid.Bid.Inactive,
					StakedAmount:            bid.Bid.StakedAmount,
					MinimumDelegationAmount: DefaultMinimumDelegationAmount,
					MaximumDelegationAmount: DefaultMaximumDelegationAmount,
					VestingSchedule:         bid.Bid.VestingSchedule,
					ReservedSlots:           0,
				},
			},
		})

		for _, delegatorBid := range bid.Bid.Delegators {
			bids = append(bids, BidKindWrapper{
				PublicKey: bid.PublicKey,
				Bid: BidKind{
					Delegator: &Delegator{
						BondingPurse:       delegatorBid.BondingPurse,
						StakedAmount:       delegatorBid.StakedAmount,
						DelegatorKind:      delegatorBid.DelegatorKind,
						ValidatorPublicKey: delegatorBid.ValidatorPublicKey,
						VestingSchedule:    delegatorBid.VestingSchedule,
					},
				},
			})
		}
	}

	return AuctionState{
		Bids:          bids,
		BlockHeight:   v1.BlockHeight,
		EraValidators: v1.EraValidators,
		StateRootHash: v1.StateRootHash,
	}
}

func NewAuctionStateFromV2(v2 AuctionStateV2) AuctionState {
	bids := make([]BidKindWrapper, 0, len(v2.Bids))
	for _, bid := range v2.Bids {
		bids = append(bids, BidKindWrapper{
			PublicKey: bid.PublicKey,
			Bid:       bid.Bid,
		})
	}

	return AuctionState{
		Bids:          bids,
		BlockHeight:   v2.BlockHeight,
		EraValidators: v2.EraValidators,
		StateRootHash: v2.StateRootHash,
	}
}

// PublicKeyAndBid is an entry in a founding validator map representing a bid.
type PublicKeyAndBid struct {
	// Validator public key
	PublicKey keypair.PublicKey `json:"public_key"`
	Bid       Bid               `json:"bid"`
}

// EraValidators contains validators and weights for an Era.
type EraValidators struct {
	EraID uint64 `json:"era_id"`
	// List of the validator's weight in the Era
	ValidatorWeights []ValidatorWeightAuction `json:"validator_weights"`
}
