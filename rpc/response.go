package rpc

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/casper-go-sdk/types/keypair"
)

// RpcResponse is a wrapper struct for an RPC Response. For a successful response the Result property
// contains the returned data as a JSON object. If an error occurs Error property contains a description of an error.
type RpcResponse struct {
	Version string          `json:"jsonrpc"`
	Id      *IDValue        `json:"id,omitempty"`
	Result  json.RawMessage `json:"result"`
	Error   *RpcError       `json:"error,omitempty"`
}

type StateGetAuctionInfoResult struct {
	Version      string             `json:"api_version"`
	AuctionState types.AuctionState `json:"auction_state"`

	rawJSON json.RawMessage
}

func (b StateGetAuctionInfoResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type StateGetBalanceResult struct {
	ApiVersion   string          `json:"api_version"`
	BalanceValue clvalue.UInt512 `json:"balance_value"`

	rawJSON json.RawMessage
}

func (b StateGetBalanceResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type StateGetAccountInfo struct {
	ApiVersion string        `json:"api_version"`
	Account    types.Account `json:"account"`

	rawJSON json.RawMessage
}

func (b StateGetAccountInfo) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type ChainGetBlockResult struct {
	APIVersion string `json:"api_version"`
	Block      types.Block

	rawJSON json.RawMessage
}

func (b ChainGetBlockResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type chainGetBlockResultV1Compatible struct {
	APIVersion          string                     `json:"api_version"`
	BlockWithSignatures *types.BlockWithSignatures `json:"block_with_signatures"`
	BlockV1             *types.BlockV1             `json:"block"`
}

func newChainGetBlockResultFromV1Compatible(result chainGetBlockResultV1Compatible, rawJSON json.RawMessage) (ChainGetBlockResult, error) {
	if result.BlockV1 != nil {
		return ChainGetBlockResult{
			APIVersion: result.APIVersion,
			Block:      types.NewBlockFromBlockV1(*result.BlockV1),
			rawJSON:    rawJSON,
		}, nil
	}

	if result.BlockWithSignatures != nil {
		return ChainGetBlockResult{
			APIVersion: result.APIVersion,
			Block:      types.NewBlockFromBlockWithSignatures(*result.BlockWithSignatures),
			rawJSON:    rawJSON,
		}, nil
	}
	return ChainGetBlockResult{}, errors.New("incorrect RPC response structure")
}

type ChainGetBlockTransfersResult struct {
	Version   string           `json:"api_version"`
	BlockHash string           `json:"block_hash"`
	Transfers []types.Transfer `json:"transfers"`

	rawJSON json.RawMessage
}

func (b ChainGetBlockTransfersResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type ChainGetEraSummaryResult struct {
	Version    string           `json:"api_version"`
	EraSummary types.EraSummary `json:"era_summary"`

	rawJSON json.RawMessage
}

func (b ChainGetEraSummaryResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type InfoGetDeployResult struct {
	ApiVersion       string                        `json:"api_version"`
	Deploy           types.Deploy                  `json:"deploy"`
	ExecutionResults []types.DeployExecutionResult `json:"execution_results"`
	BlockHash        *key.Hash                     `json:"block_hash,omitempty"`
	BlockHeight      *uint64                       `json:"block_height,omitempty"`

	rawJSON json.RawMessage
}

func (b InfoGetDeployResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type ChainGetEraInfoResult struct {
	Version    string           `json:"api_version"`
	EraSummary types.EraSummary `json:"era_summary"`

	rawJSON json.RawMessage
}

func (b ChainGetEraInfoResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type StateGetItemResult struct {
	StoredValue types.StoredValue `json:"stored_value"`
	//MerkleProof is a construction created using a merkle trie that allows verification of the associated hashes.
	MerkleProof json.RawMessage `json:"merkle_proof"`

	rawJSON json.RawMessage
}

func (b StateGetItemResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type StateGetDictionaryResult struct {
	ApiVersion    string            `json:"api_version"`
	DictionaryKey string            `json:"dictionary_key"`
	StoredValue   types.StoredValue `json:"stored_value"`
	MerkleProof   json.RawMessage   `json:"merkle_proof"`

	rawJSON json.RawMessage
}

func (b StateGetDictionaryResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type QueryGlobalStateResult struct {
	ApiVersion  string            `json:"api_version"`
	BlockHeader types.BlockHeader `json:"block_header,omitempty"`
	StoredValue types.StoredValue `json:"stored_value"`
	//MerkleProof is a construction created using a merkle trie that allows verification of the associated hashes.
	MerkleProof json.RawMessage `json:"merkle_proof"`

	rawJSON json.RawMessage
}

func (b QueryGlobalStateResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type InfoGetPeerResult struct {
	ApiVersion string     `json:"api_version"`
	Peers      []NodePeer `json:"peers"`

	rawJSON json.RawMessage
}

func (b InfoGetPeerResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type NodePeer struct {
	NodeID  string `json:"node_id"`
	Address string `json:"address"`
}

type ChainGetStateRootHashResult struct {
	Version       string   `json:"api_version"`
	StateRootHash key.Hash `json:"state_root_hash"`

	rawJSON json.RawMessage
}

func (b ChainGetStateRootHashResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type ValidatorState string

const (
	// ValidatorStateAdded means that the validator has been added to the set.
	ValidatorStateAdded ValidatorState = "Added"
	// ValidatorStateRemoved means that the validator has been removed from the set.
	ValidatorStateRemoved ValidatorState = "Removed"
	// ValidatorStateBanned means that the validator has been banned in the current era.
	ValidatorStateBanned ValidatorState = "Banned"
	// ValidatorStateCannotPropose means that the validator cannot propose a Block.
	ValidatorStateCannotPropose ValidatorState = "CannotPropose"
	// ValidatorStateSeenAsFaulty means that the validator has performed questionable activity.
	ValidatorStateSeenAsFaulty ValidatorState = "SeenAsFaulty"
)

type StatusChanges struct {
	EraID          uint64         `json:"era_id"`
	ValidatorState ValidatorState `json:"validator_change"`

	rawJSON json.RawMessage
}

func (b StatusChanges) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type ValidatorChanges struct {
	PublicKey     keypair.PublicKey `json:"public_key"`
	StatusChanges []StatusChanges   `json:"status_changes"`

	rawJSON json.RawMessage
}

func (b ValidatorChanges) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type InfoGetValidatorChangesResult struct {
	APIVersion string             `json:"api_version"`
	Changes    []ValidatorChanges `json:"changes"`

	rawJSON json.RawMessage
}

func (b InfoGetValidatorChangesResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type InfoGetStatusResult struct {
	// The RPC API version.
	APIVersion string `json:"api_version"`
	// The compiled node version.
	BuildVersion string `json:"build_version"`
	// The chainspec name, used to identify the currently connected network.
	ChainSpecName string `json:"chainspec_name"`
	// The minimal info of the last block from the linear chain.
	LastAddedBlockInfo types.MinimalBlockInfo `json:"last_added_block_info"`
	// Information about the next scheduled upgrade.
	NextUpgrade NodeNextUpgrade `json:"next_upgrade,omitempty"`
	// Node public signing key.
	OutPublicSigningKey string `json:"our_public_signing_key"`
	// The list of node ID and network address of each connected peer.
	Peers []NodePeer `json:"peers"`
	// The next round length if this node is a validator.
	RoundLength string `json:"round_length"`
	// The state root hash used at the start of the current session.
	StartingStateRootHash string `json:"starting_state_root_hash"`
	// Time that passed since the node has started. format "2months 20days 22h 3m 21s 512ms"
	Uptime string `json:"uptime"`
	// Indicating the node's current operating mode
	ReactorState string `json:"reactor_state"`
	// Indicating the time the node last made progress
	LastProgress types.Timestamp `json:"last_progress"`
	// The hash of the latest switch block
	LatestSwitchBlockHash key.Hash `json:"latest_switch_block_hash"`
	// Indicating the highest contiguous sequence of the block chain for which the node has complete data
	AvailableBlockRange struct {
		Low  uint64 `json:"low"`
		High uint64 `json:"high"`
	} `json:"available_block_range"`
	// Indicating the state of the block synchronizer component
	BlockSync struct {
		Historical string `json:"historical,omitempty"`
		Forward    string `json:"forward,omitempty"`
	} `json:"block_sync"`

	rawJSON json.RawMessage
}

func (b PutDeployResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

// NodeNextUpgrade contains the information about the next protocol upgrade.
type NodeNextUpgrade struct {
	//The first era to which the associated protocol version applies.
	ActivationPoint uint64 `json:"activation_point"`
	// The protocol version of the next upgrade
	ProtocolVersion string `json:"protocol_version"`
}

type PutDeployResult struct {
	ApiVersion string   `json:"api_version"`
	DeployHash key.Hash `json:"deploy_hash"`

	rawJSON json.RawMessage
}

func (b InfoGetStatusResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type SpeculativeExecResult struct {
	ApiVersion      string                      `json:"api_version"`
	BlockHash       key.Hash                    `json:"block_hash"`
	ExecutionResult types.ExecutionResultStatus `json:"execution_result"`

	rawJSON json.RawMessage
}

func (b SpeculativeExecResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type QueryBalanceResult struct {
	ApiVersion string          `json:"api_version"`
	Balance    clvalue.UInt512 `json:"balance"`
	rawJSON    json.RawMessage
}

func (b QueryBalanceResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type QueryBalanceDetailsResult struct {
	APIVersion        string                 `json:"api_version"`
	TotalBalance      clvalue.UInt512        `json:"total_balance"`
	AvailableBalance  clvalue.UInt512        `json:"available_balance"`
	TotalBalanceProof string                 `json:"total_balance_proof"`
	Holds             []BalanceHoldWithProof `json:"holds"`
}

// BalanceHoldWithProof The block time at which the hold was created.
type BalanceHoldWithProof struct {
	Time   types.BlockTime `json:"time"`
	Amount clvalue.UInt512 `json:"amount"`
	Proof  string          `json:"proof"`
}

type InfoGetChainspecResult struct {
	ApiVersion     string `json:"api_version"`
	ChainspecBytes struct {
		ChainspecBytes            string `json:"chainspec_bytes,omitempty"`
		MaybeGenesisAccountsBytes string `json:"maybe_genesis_accounts_bytes,omitempty"`
		MaybeGlobalStateBytes     string `json:"maybe_global_state_bytes,omitempty"`
	} `json:"chainspec_bytes"`
	rawJSON json.RawMessage
}

func (b InfoGetChainspecResult) GetRawJSON() json.RawMessage {
	return b.rawJSON
}

type queryGlobalStateResultV1Compatible struct {
	ApiVersion  string              `json:"api_version"`
	BlockHeader types.BlockHeaderV1 `json:"block_header,omitempty"`
	StoredValue types.StoredValue   `json:"stored_value"`
	//MerkleProof is a construction created using a merkle trie that allows verification of the associated hashes.
	MerkleProof json.RawMessage `json:"merkle_proof"`
}

// UnmarshalJSON handle the backward compatibility logic with V1
func (h *QueryGlobalStateResult) UnmarshalJSON(bytes []byte) error {
	// Check the API version
	version := struct {
		ApiVersion string `json:"api_version"`
	}{}

	if err := json.Unmarshal(bytes, &version); err != nil {
		return err
	}

	// handle V1 version
	if strings.HasPrefix(version.ApiVersion, "1") {
		var v1Compatible queryGlobalStateResultV1Compatible
		if err := json.Unmarshal(bytes, &v1Compatible); err != nil {
			return err
		}
		*h = QueryGlobalStateResult{
			ApiVersion:  v1Compatible.ApiVersion,
			BlockHeader: types.NewBlockHeaderFromV1(v1Compatible.BlockHeader),
			StoredValue: v1Compatible.StoredValue,
			MerkleProof: v1Compatible.MerkleProof,
		}
		return nil
	}

	var result struct {
		ApiVersion  string                   `json:"api_version"`
		BlockHeader types.BlockHeaderWrapper `json:"block_header,omitempty"`
		StoredValue types.StoredValue        `json:"stored_value"`
		MerkleProof json.RawMessage          `json:"merkle_proof"`
	}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	if result.BlockHeader.BlockHeaderV2 == nil {
		return errors.New("incorrect RPC response structure")
	}

	*h = QueryGlobalStateResult{
		ApiVersion:  result.ApiVersion,
		BlockHeader: types.NewBlockHeaderFromV2(*result.BlockHeader.BlockHeaderV2),
		StoredValue: result.StoredValue,
		MerkleProof: result.MerkleProof,
	}
	return nil
}
