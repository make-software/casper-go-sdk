package types

import (
	"encoding/json"
	"errors"

	"github.com/make-software/casper-go-sdk/v2/types/key"
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

// Block represents a common object returned as result from RPC response unifying BlockV2 and BlockV1
type Block struct {
	Hash                key.Hash                        `json:"hash"`
	Height              uint64                          `json:"height"`
	StateRootHash       key.Hash                        `json:"state_root_hash"`
	LastSwitchBlockHash key.Hash                        `json:"last_switch_block_hash"`
	ParentHash          key.Hash                        `json:"parent_hash"`
	EraID               uint32                          `json:"era_id"`
	Timestamp           Timestamp                       `json:"timestamp"`
	AccumulatedSeed     *key.Hash                       `json:"accumulated_seed,omitempty"`
	RandomBit           bool                            `json:"random_bit"`
	CurrentGasPrice     uint8                           `json:"current_gas_price"`
	Proposer            Proposer                        `json:"proposer"`
	ProtocolVersion     string                          `json:"protocol_version,omitempty"`
	EraEnd              *EraEnd                         `json:"era_end"`
	Transactions        BlockTransactions               `json:"transactions"`
	RewardedSignatures  []SingleBlockRewardedSignatures `json:"rewarded_signatures"`
	Proofs              []Proof                         `json:"proofs"`

	// source BlockV1, nil if constructed from BlockV2
	originBlockV1 *BlockV1
	// source BlockV2, nil if constructed from BlockV1
	originBlockV2 *BlockV2
}

func (b Block) GetBlockV1() *BlockV1 {
	return b.originBlockV1
}

func (b Block) GetBlockV2() *BlockV2 {
	return b.originBlockV2
}

// NewBlockFromBlockWrapper construct Block from BlockWithSignatures
func NewBlockFromBlockWrapper(blockWrapper BlockWrapper, proofs []Proof) Block {
	if blockV1 := blockWrapper.BlockV1; blockV1 != nil {
		block := NewBlockFromBlockV1(*blockV1)
		block.Proofs = proofs
		return block
	} else {
		blockV2 := blockWrapper.BlockV2
		return Block{
			Hash:                blockV2.Hash,
			Height:              blockV2.Header.Height,
			StateRootHash:       blockV2.Header.StateRootHash,
			LastSwitchBlockHash: blockV2.Header.LastSwitchBlockHash,
			ParentHash:          blockV2.Header.ParentHash,
			EraID:               blockV2.Header.EraID,
			Timestamp:           blockV2.Header.Timestamp,
			AccumulatedSeed:     blockV2.Header.AccumulatedSeed,
			RandomBit:           blockV2.Header.RandomBit,
			CurrentGasPrice:     blockV2.Header.CurrentGasPrice,
			Proposer:            blockV2.Header.Proposer,
			ProtocolVersion:     blockV2.Header.ProtocolVersion,
			EraEnd:              NewEraEndFromV2(blockV2.Header.EraEnd),
			Transactions:        blockV2.Body.Transactions,
			RewardedSignatures:  blockV2.Body.RewardedSignatures,
			Proofs:              proofs,
			originBlockV2:       blockV2,
		}
	}
}

// NewBlockFromBlockV1 construct Block from BlockV1
func NewBlockFromBlockV1(blockV1 BlockV1) Block {
	blockTransactions := make(BlockTransactions, 0)
	for i := range blockV1.Body.TransferHashes {
		blockTransactions = append(blockTransactions, BlockTransaction{
			Category: TransactionCategoryMint,
			Version:  TransactionDeploy,
			Hash:     blockV1.Body.TransferHashes[i],
		})
	}

	for i := range blockV1.Body.DeployHashes {
		blockTransactions = append(blockTransactions, BlockTransaction{
			Category: TransactionCategoryLarge,
			Version:  TransactionDeploy,
			Hash:     blockV1.Body.DeployHashes[i],
		})
	}

	return Block{
		Hash:            blockV1.Hash,
		Height:          blockV1.Header.Height,
		StateRootHash:   blockV1.Header.StateRootHash,
		ParentHash:      blockV1.Header.ParentHash,
		EraID:           blockV1.Header.EraID,
		Timestamp:       blockV1.Header.Timestamp,
		AccumulatedSeed: blockV1.Header.AccumulatedSeed,
		RandomBit:       blockV1.Header.RandomBit,
		CurrentGasPrice: 1,
		Proposer:        blockV1.Body.Proposer,
		ProtocolVersion: blockV1.Header.ProtocolVersion,
		EraEnd:          NewEraEndFromV1(blockV1.Header.EraEnd),
		Transactions:    blockTransactions,
		Proofs:          blockV1.Proofs,
		originBlockV1:   &blockV1,
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

type BlockHeaderWrapper struct {
	BlockHeaderV1 *BlockHeaderV1 `json:"Version1"`
	BlockHeaderV2 *BlockHeaderV2 `json:"Version2"`
}

type BlockHeader struct {
	AccumulatedSeed *key.Hash `json:"accumulated_seed,omitempty"`
	BodyHash        key.Hash  `json:"body_hash"`
	EraID           uint32    `json:"era_id"`
	CurrentGasPrice uint8     `json:"current_gas_price"`
	Height          uint64    `json:"height"`
	ParentHash      key.Hash  `json:"parent_hash"`
	Proposer        Proposer  `json:"proposer"`
	ProtocolVersion string    `json:"protocol_version,omitempty"`
	RandomBit       bool      `json:"random_bit"`
	StateRootHash   key.Hash  `json:"state_root_hash"`
	Timestamp       Timestamp `json:"timestamp"`
	EraEnd          *EraEnd   `json:"era_end"`

	// source OriginBlockHeaderV1, nil if constructed from BlockHeaderV2
	originBlockHeaderV1 *BlockHeaderV1
	// source OriginBlockHeaderV2, nil if constructed from BlockHeaderV1
	originBlockHeaderV2 *BlockHeaderV2
}

func (b BlockHeader) GetBlockHeaderV1() *BlockHeaderV1 {
	return b.originBlockHeaderV1
}

func (b BlockHeader) GetBlockHeaderV2() *BlockHeaderV2 {
	return b.originBlockHeaderV2
}

func NewBlockHeaderFromV1(header BlockHeaderV1) BlockHeader {
	return BlockHeader{
		AccumulatedSeed:     header.AccumulatedSeed,
		BodyHash:            header.BodyHash,
		EraID:               header.EraID,
		CurrentGasPrice:     1,
		Height:              header.Height,
		ParentHash:          header.ParentHash,
		ProtocolVersion:     header.ProtocolVersion,
		RandomBit:           header.RandomBit,
		StateRootHash:       header.StateRootHash,
		Timestamp:           header.Timestamp,
		EraEnd:              NewEraEndFromV1(header.EraEnd),
		originBlockHeaderV1: &header,
	}
}

func NewBlockHeaderFromV2(header BlockHeaderV2) BlockHeader {
	return BlockHeader{
		AccumulatedSeed:     header.AccumulatedSeed,
		BodyHash:            header.BodyHash,
		EraID:               header.EraID,
		CurrentGasPrice:     header.CurrentGasPrice,
		Height:              header.Height,
		ParentHash:          header.ParentHash,
		Proposer:            header.Proposer,
		ProtocolVersion:     header.ProtocolVersion,
		RandomBit:           header.RandomBit,
		StateRootHash:       header.StateRootHash,
		Timestamp:           header.Timestamp,
		EraEnd:              NewEraEndFromV2(header.EraEnd),
		originBlockHeaderV2: &header,
	}
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
type SingleBlockRewardedSignatures []uint8

type BlockTransactions []BlockTransaction

func (t *BlockTransactions) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}

	source := struct {
		Mint           []TransactionHash `json:"0,omitempty"`
		Auction        []TransactionHash `json:"1,omitempty"`
		InstallUpgrade []TransactionHash `json:"2,omitempty"`
		Large          []TransactionHash `json:"3,omitempty"`
		Medium         []TransactionHash `json:"4,omitempty"`
		Small          []TransactionHash `json:"5,omitempty"`
	}{}

	if err := json.Unmarshal(data, &source); err != nil {
		return err
	}

	res := make(BlockTransactions, 0)
	res = append(res, getBlockTransactionsFromTransactionHashes(source.Mint, TransactionCategoryMint)...)
	res = append(res, getBlockTransactionsFromTransactionHashes(source.Auction, TransactionCategoryAuction)...)
	res = append(res, getBlockTransactionsFromTransactionHashes(source.InstallUpgrade, TransactionCategoryInstallUpgrade)...)
	res = append(res, getBlockTransactionsFromTransactionHashes(source.Large, TransactionCategoryLarge)...)
	res = append(res, getBlockTransactionsFromTransactionHashes(source.Medium, TransactionCategoryMedium)...)
	res = append(res, getBlockTransactionsFromTransactionHashes(source.Small, TransactionCategorySmall)...)
	*t = res
	return nil
}

type BlockTransaction struct {
	Category TransactionCategory
	Version  TransactionVersion
	Hash     key.Hash
}

type BlockBodyV2 struct {
	// Map of transactions mapping categories to a list of transaction hashes.
	Transactions BlockTransactions `json:"transactions"`
	// List of identifiers for finality signatures for a particular past block
	RewardedSignatures []SingleBlockRewardedSignatures `json:"rewarded_signatures"`
}

type BlockHeaderV2 struct {
	AccumulatedSeed     *key.Hash `json:"accumulated_seed,omitempty"`
	BodyHash            key.Hash  `json:"body_hash"`
	EraID               uint32    `json:"era_id"`
	CurrentGasPrice     uint8     `json:"current_gas_price"`
	Height              uint64    `json:"height"`
	ParentHash          key.Hash  `json:"parent_hash"`
	Proposer            Proposer  `json:"proposer"`
	ProtocolVersion     string    `json:"protocol_version,omitempty"`
	RandomBit           bool      `json:"random_bit"`
	StateRootHash       key.Hash  `json:"state_root_hash"`
	LastSwitchBlockHash key.Hash  `json:"last_switch_block_hash"`
	Timestamp           Timestamp `json:"timestamp"`
	EraEnd              *EraEndV2 `json:"era_end"`
}

// Proof is a `BlockV1`'s finality signature.
type Proof struct {
	// Validator public key
	PublicKey keypair.PublicKey `json:"public_key"`
	// Validator signature
	Signature HexBytes `json:"signature"`
}

func getBlockTransactionsFromTransactionHashes(hashes []TransactionHash, category TransactionCategory) BlockTransactions {
	if len(hashes) == 0 {
		return nil
	}

	res := make(BlockTransactions, 0)
	for i := range hashes {
		blockTransaction := BlockTransaction{
			Category: category,
		}
		if hashes[i].TransactionV1 != nil {
			blockTransaction.Hash = *hashes[i].TransactionV1
			blockTransaction.Version = TransactionVersionV1
		} else {
			blockTransaction.Hash = *hashes[i].Deploy
			blockTransaction.Version = TransactionDeploy
		}

		res = append(res, blockTransaction)
	}

	return res
}
