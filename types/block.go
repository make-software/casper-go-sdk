package types

import (
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/casper-go-sdk/types/keypair"
)

// Block represents a common object returned as result from RPC response unifying BlockV2 and BlockV1
// Block is inherited from BlockV2, BlockV1 should be matched to the view of BlockV2 to achieve backward compatibility
type Block struct {
	BlockV2
	// source BlockV1, nil if constructed from BlockV2
	OriginBlockV1 *BlockV1

	Proofs []Proof `json:"proofs"`
}

// NewBlockFromBlockWithSignatures construct Block from BlockWithSignatures
func NewBlockFromBlockWithSignatures(signBlock BlockWithSignatures) Block {
	if blockV1 := signBlock.Block.BlockV1; blockV1 != nil {
		block := NewBlockFromBlockV1(*blockV1)
		block.Proofs = signBlock.Proofs
		return block
	} else {
		return Block{
			BlockV2: *signBlock.Block.BlockV2,
			Proofs:  signBlock.Proofs,
		}
	}
}

// NewBlockFromBlockV1 construct Block from BlockV1
func NewBlockFromBlockV1(block BlockV1) Block {
	var eraEnd EraEndV2
	if block.Header.EraEnd != nil {
		rewards := make(map[string]clvalue.UInt512, len(block.Header.EraEnd.EraReport.Rewards))
		for _, reward := range block.Header.EraEnd.EraReport.Rewards {
			rewards[reward.Validator.ToHex()] = reward.Amount
		}

		eraEnd = EraEndV2{
			NextEraGasPrice:         1,
			Equivocators:            block.Header.EraEnd.EraReport.Equivocators,
			InactiveValidators:      block.Header.EraEnd.EraReport.InactiveValidators,
			NextEraValidatorWeights: block.Header.EraEnd.NextEraValidatorWeights,
			Rewards:                 rewards,
		}
	}

	mints := make([]TransactionHash, 0, len(block.Body.TransferHashes))
	for i := range block.Body.TransferHashes {
		mints = append(mints, TransactionHash{TransactionV1Hash: &block.Body.TransferHashes[i]})
	}

	// To achieve backward compatibility DeployHashes are transformed to Standard
	// Use OriginBlockV1 if you are looking for source data
	standard := make([]TransactionHash, 0, len(block.Body.DeployHashes))
	for i := range block.Body.DeployHashes {
		standard = append(standard, TransactionHash{TransactionV1Hash: &block.Body.DeployHashes[i]})
	}

	return Block{
		BlockV2: BlockV2{
			Hash: block.Hash,
			Header: BlockHeaderV2{
				AccumulatedSeed: block.Header.AccumulatedSeed,
				BodyHash:        block.Header.BodyHash,
				EraID:           block.Header.EraID,
				CurrentGasPrice: 1,
				Height:          block.Header.Height,
				ParentHash:      block.Header.ParentHash,
				ProtocolVersion: block.Header.ProtocolVersion,
				RandomBit:       block.Header.RandomBit,
				StateRootHash:   block.Header.StateRootHash,
				Timestamp:       block.Header.Timestamp,
				EraEnd:          &eraEnd,
			},
			Body: BlockBodyV2{
				Proposer: block.Body.Proposer,
				Mint:     mints,
				Standard: standard,
			},
		},
		OriginBlockV1: &block,
		Proofs:        block.Proofs,
	}
}

type BlockWithSignatures struct {
	Block  BlockWrapper `json:"block"`
	Proofs []Proof      `json:"proofs"`
}

type BlockWrapper struct {
	BlockV1 *BlockV1 `json:"Version1"`
	BlockV2 *BlockV2 `json:"Version2"`
}

type BlockV1 struct {
	Hash   key.Hash      `json:"hash"`
	Header BlockHeaderV1 `json:"header"`
	Body   BlockBodyV1   `json:"body"`
	Proofs []Proof       `json:"proofs"`
}

type BlockBodyV1 struct {
	// List of `Deploy` hashes included in the block
	DeployHashes []key.Hash `json:"deploy_hashes"`
	// Public key of the validator that proposed the block
	Proposer Proposer `json:"proposer"`
	// List of `TransferHash` hashes included in the block
	TransferHashes []key.Hash `json:"transfer_hashes"`
}

type BlockHeaderV1 struct {
	AccumulatedSeed *key.Hash `json:"accumulated_seed,omitempty"`
	BodyHash        key.Hash  `json:"body_hash"`
	EraID           uint32    `json:"era_id"`
	Height          uint64    `json:"height"`
	ParentHash      key.Hash  `json:"parent_hash"`
	ProtocolVersion string    `json:"protocol_version,omitempty"`
	RandomBit       bool      `json:"random_bit"`
	StateRootHash   key.Hash  `json:"state_root_hash"`
	Timestamp       Timestamp `json:"timestamp"`
	EraEnd          *EraEndV1 `json:"era_end"`
}

type BlockV2 struct {
	Hash   key.Hash      `json:"hash"`
	Header BlockHeaderV2 `json:"header"`
	Body   BlockBodyV2   `json:"body"`
}

// SingleBlockRewardedSignatures
// List of identifiers for finality signatures for a particular past block.
// That past block height is current_height - signature_rewards_max_delay, the latter being defined in the chainspec.
// We need to wait for a few blocks to pass (`signature_rewards_max_delay`) to store the finality signers because we need a bit of time to get the block finality.
type SingleBlockRewardedSignatures []uint16

type BlockBodyV2 struct {
	// Public key of the validator that proposed the block
	Proposer Proposer `json:"proposer"`
	// The hashes of the mint transactions within the block.
	Mint []TransactionHash `json:"mint"`
	// The hashes of the auction transactions within the block.
	Auction []TransactionHash `json:"auction"`
	// The hashes of the installer/upgrader transactions within the block
	InstallUpgrade []TransactionHash `json:"install_upgrade"`
	// The hashes of all other transactions within the block
	Standard []TransactionHash `json:"standard"`
	// List of identifiers for finality signatures for a particular past block
	RewardedSignatures []SingleBlockRewardedSignatures `json:"rewarded_signatures"`
}

type BlockHeaderV2 struct {
	AccumulatedSeed *key.Hash `json:"accumulated_seed,omitempty"`
	BodyHash        key.Hash  `json:"body_hash"`
	EraID           uint32    `json:"era_id"`
	CurrentGasPrice uint8     `json:"current_gas_price"`
	Height          uint64    `json:"height"`
	ParentHash      key.Hash  `json:"parent_hash"`
	ProtocolVersion string    `json:"protocol_version,omitempty"`
	RandomBit       bool      `json:"random_bit"`
	StateRootHash   key.Hash  `json:"state_root_hash"`
	Timestamp       Timestamp `json:"timestamp"`
	EraEnd          *EraEndV2 `json:"era_end"`
}

// Proof is a `BlockV1`'s finality signature.
type Proof struct {
	// Validator public key
	PublicKey keypair.PublicKey `json:"public_key"`
	// Validator signature
	Signature HexBytes `json:"signature"`
}
