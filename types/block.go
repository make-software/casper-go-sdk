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
		rewards := make(map[string][]clvalue.UInt512, len(block.Header.EraEnd.EraReport.Rewards))
		for _, reward := range block.Header.EraEnd.EraReport.Rewards {
			list := rewards[reward.Validator.ToHex()]
			list = append(list, reward.Amount)
			rewards[reward.Validator.ToHex()] = list
		}

		eraEnd = EraEndV2{
			NextEraGasPrice:         1,
			Equivocators:            block.Header.EraEnd.EraReport.Equivocators,
			InactiveValidators:      block.Header.EraEnd.EraReport.InactiveValidators,
			NextEraValidatorWeights: block.Header.EraEnd.NextEraValidatorWeights,
			Rewards:                 rewards,
		}
	}

	blockTransactions := make(BlockTransactions, 0)
	for i := range block.Body.TransferHashes {
		blockTransactions = append(blockTransactions, BlockTransaction{
			Category: BlockTransactionCategoryMint,
			Version:  BlockTransactionDeploy,
			Hash:     block.Body.TransferHashes[i],
		})
	}

	for i := range block.Body.DeployHashes {
		blockTransactions = append(blockTransactions, BlockTransaction{
			Category: BlockTransactionCategoryLarge,
			Version:  BlockTransactionDeploy,
			Hash:     block.Body.DeployHashes[i],
		})
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
				Proposer:        block.Body.Proposer,
			},
			Body: BlockBodyV2{
				Transactions: blockTransactions,
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

type BlockHeaderWrapper struct {
	BlockHeaderV1 *BlockHeaderV1 `json:"Version1"`
	BlockHeaderV2 *BlockHeaderV2 `json:"Version2"`
}

type BlockHeader struct {
	BlockHeaderV2

	// source OriginBlockHeaderV1, nil if constructed from BlockHeaderV2
	OriginBlockHeaderV1 *BlockHeaderV1
}

func NewBlockHeaderFromV1(header BlockHeaderV1) BlockHeader {
	var eraEnd *EraEndV2
	if header.EraEnd != nil {
		rewards := make(map[string][]clvalue.UInt512, len(header.EraEnd.EraReport.Rewards))
		for _, reward := range header.EraEnd.EraReport.Rewards {
			list := rewards[reward.Validator.ToHex()]
			list = append(list, reward.Amount)
			rewards[reward.Validator.ToHex()] = list
		}

		eraEnd = &EraEndV2{
			Equivocators:            header.EraEnd.EraReport.Equivocators,
			InactiveValidators:      header.EraEnd.EraReport.InactiveValidators,
			NextEraValidatorWeights: header.EraEnd.NextEraValidatorWeights,
			Rewards:                 rewards,
			NextEraGasPrice:         1,
		}
	}

	return BlockHeader{
		BlockHeaderV2: BlockHeaderV2{
			AccumulatedSeed: header.AccumulatedSeed,
			BodyHash:        header.BodyHash,
			EraID:           header.EraID,
			CurrentGasPrice: 1,
			Height:          header.Height,
			ParentHash:      header.ParentHash,
			ProtocolVersion: header.ProtocolVersion,
			RandomBit:       header.RandomBit,
			StateRootHash:   header.StateRootHash,
			Timestamp:       header.Timestamp,
			EraEnd:          eraEnd,
		},
		OriginBlockHeaderV1: &header,
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
