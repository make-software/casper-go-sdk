package types

import (
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/casper-go-sdk/types/keypair"
)

// Block in the network
type Block struct {
	Hash   key.Hash    `json:"hash"`
	Header BlockHeader `json:"header"`
	Body   BlockBody   `json:"body"`
	Proofs []Proof     `json:"proofs"`
}

type BlockBody struct {
	// List of `Deploy` hashes included in the block
	DeployHashes []key.Hash `json:"deploy_hashes"`
	// Public key of the validator that proposed the block
	Proposer HexBytes `json:"proposer"`
	// List of `TransferHash` hashes included in the block
	TransferHashes []key.TransferHash `json:"transfer_hashes"`
}

type BlockHeader struct {
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

// Proof is a `Block`'s finality signature.
type Proof struct {
	// Validator public key
	PublicKey keypair.PublicKey `json:"public_key"`
	// Validator signature
	Signature HexBytes `json:"signature"`
}
