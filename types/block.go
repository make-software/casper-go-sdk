package types

import (
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/casper-go-sdk/types/keypair"
)

type BlockWithSignatures struct {
	Block  Block   `json:"block"`
	Proofs []Proof `json:"proofs"`
}

type Block struct {
	BlockV1 *BlockV1 `json:"Version1"`
	BlockV2 *BlockV2 `json:"Version2"`
}

type BlockV1 struct {
	Hash   key.Hash      `json:"hash"`
	Header BlockHeaderV1 `json:"header"`
	Body   BlockBodyV1   `json:"body"`
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
	EraEnd          *EraEnd   `json:"era_end"`
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
