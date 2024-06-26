package types

import (
	"encoding/json"
	"errors"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/casper-go-sdk/types/keypair"
)

// Block represents a common object returned as result from RPC response unifying BlockV2 and BlockV1
// Block is inherited from BlockV2, BlockV1 should be matched to the view of BlockV2 to achieve backward compatibility
type Block struct {
	Hash               key.Hash                        `json:"hash"`
	Height             uint64                          `json:"height"`
	StateRootHash      key.Hash                        `json:"state_root_hash"`
	ParentHash         key.Hash                        `json:"parent_hash"`
	EraID              uint32                          `json:"era_id"`
	Timestamp          Timestamp                       `json:"timestamp"`
	AccumulatedSeed    *key.Hash                       `json:"accumulated_seed,omitempty"`
	RandomBit          bool                            `json:"random_bit"`
	CurrentGasPrice    uint8                           `json:"current_gas_price"`
	Proposer           Proposer                        `json:"proposer"`
	ProtocolVersion    string                          `json:"protocol_version,omitempty"`
	EraEnd             *EraEnd                         `json:"era_end"`
	Transactions       BlockTransactions               `json:"transactions"`
	RewardedSignatures []SingleBlockRewardedSignatures `json:"rewarded_signatures"`
	Proofs             []Proof                         `json:"proofs"`

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

// NewBlockFromBlockWithSignatures construct Block from BlockWithSignatures
func NewBlockFromBlockWithSignatures(signBlock BlockWithSignatures) Block {
	if blockV1 := signBlock.Block.BlockV1; blockV1 != nil {
		block := NewBlockFromBlockV1(*blockV1)
		block.Proofs = signBlock.Proofs
		return block
	} else {
		return Block{
			Hash:               signBlock.Block.BlockV2.Hash,
			Height:             signBlock.Block.BlockV2.Header.Height,
			StateRootHash:      signBlock.Block.BlockV2.Header.StateRootHash,
			ParentHash:         signBlock.Block.BlockV2.Header.ParentHash,
			EraID:              signBlock.Block.BlockV2.Header.EraID,
			Timestamp:          signBlock.Block.BlockV2.Header.Timestamp,
			AccumulatedSeed:    signBlock.Block.BlockV2.Header.AccumulatedSeed,
			RandomBit:          signBlock.Block.BlockV2.Header.RandomBit,
			CurrentGasPrice:    signBlock.Block.BlockV2.Header.CurrentGasPrice,
			Proposer:           signBlock.Block.BlockV2.Header.Proposer,
			ProtocolVersion:    signBlock.Block.BlockV2.Header.ProtocolVersion,
			EraEnd:             NewEraEndFromV2(signBlock.Block.BlockV2.Header.EraEnd),
			Transactions:       signBlock.Block.BlockV2.Body.Transactions,
			RewardedSignatures: signBlock.Block.BlockV2.Body.RewardedSignatures,
			Proofs:             signBlock.Proofs,
			originBlockV2:      signBlock.Block.BlockV2,
		}
	}
}

// NewBlockFromBlockV1 construct Block from BlockV1
func NewBlockFromBlockV1(blockV1 BlockV1) Block {
	blockTransactions := make(BlockTransactions, 0)
	for i := range blockV1.Body.TransferHashes {
		blockTransactions = append(blockTransactions, BlockTransaction{
			Category: BlockTransactionCategoryMint,
			Version:  BlockTransactionDeploy,
			Hash:     blockV1.Body.TransferHashes[i],
		})
	}

	for i := range blockV1.Body.DeployHashes {
		blockTransactions = append(blockTransactions, BlockTransaction{
			Category: BlockTransactionCategoryLarge,
			Version:  BlockTransactionDeploy,
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
	if header.EraEnd != nil {
		rewards := make(map[string][]clvalue.UInt512, len(header.EraEnd.EraReport.Rewards))
		for _, reward := range header.EraEnd.EraReport.Rewards {
			list := rewards[reward.Validator.ToHex()]
			list = append(list, reward.Amount)
			rewards[reward.Validator.ToHex()] = list
		}
	}

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
type SingleBlockRewardedSignatures []uint16

type BlockTransactionCategory uint

const (
	BlockTransactionCategoryMint BlockTransactionCategory = iota
	BlockTransactionCategoryAuction
	BlockTransactionCategoryInstallUpgrade
	BlockTransactionCategoryLarge
	BlockTransactionCategoryMedium
	BlockTransactionCategorySmall
)

type BlockTransactionVersion uint

const (
	BlockTransactionVersionV1 BlockTransactionVersion = iota
	BlockTransactionDeploy
)

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
	res = append(res, getBlockTransactionsFromTransactionHashes(source.Mint, BlockTransactionCategoryMint)...)
	res = append(res, getBlockTransactionsFromTransactionHashes(source.Auction, BlockTransactionCategoryAuction)...)
	res = append(res, getBlockTransactionsFromTransactionHashes(source.InstallUpgrade, BlockTransactionCategoryInstallUpgrade)...)
	res = append(res, getBlockTransactionsFromTransactionHashes(source.Large, BlockTransactionCategoryLarge)...)
	res = append(res, getBlockTransactionsFromTransactionHashes(source.Medium, BlockTransactionCategoryMedium)...)
	res = append(res, getBlockTransactionsFromTransactionHashes(source.Small, BlockTransactionCategorySmall)...)
	*t = res
	return nil
}

type BlockTransaction struct {
	Category BlockTransactionCategory
	Version  BlockTransactionVersion
	Hash     key.Hash
}

type BlockBodyV2 struct {
	// Map of transactions mapping categories to a list of transaction hashes.
	Transactions BlockTransactions `json:"transactions"`
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
	Proposer        Proposer  `json:"proposer"`
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

func getBlockTransactionsFromTransactionHashes(hashes []TransactionHash, category BlockTransactionCategory) BlockTransactions {
	if len(hashes) == 0 {
		return nil
	}

	res := make(BlockTransactions, 0)
	for i := range hashes {
		blockTransaction := BlockTransaction{
			Category: category,
		}
		if hashes[i].TransactionV1Hash != nil {
			blockTransaction.Hash = *hashes[i].TransactionV1Hash
			blockTransaction.Version = BlockTransactionVersionV1
		} else {
			blockTransaction.Hash = *hashes[i].Deploy
			blockTransaction.Version = BlockTransactionDeploy
		}

		res = append(res, blockTransaction)
	}

	return res
}
