package rpc

import (
	"encoding/json"

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
}

type StateGetBalanceResult struct {
	ApiVersion   string          `json:"api_version"`
	BalanceValue clvalue.UInt512 `json:"balance_value"`
}

type StateGetAccountInfo struct {
	ApiVersion string        `json:"api_version"`
	Account    types.Account `json:"account"`
}

type ChainGetBlockResult struct {
	Version string      `json:"version"`
	Block   types.Block `json:"block"`
}

type ChainGetBlockTransfersResult struct {
	Version   string           `json:"api_version"`
	BlockHash string           `json:"block_hash"`
	Transfers []types.Transfer `json:"transfers"`
}

type ChainGetEraSummaryResult struct {
	Version    string           `json:"api_version"`
	EraSummary types.EraSummary `json:"era_summary"`
}

type InfoGetDeployResult struct {
	ApiVersion       string                        `json:"api_version"`
	Deploy           types.Deploy                  `json:"deploy"`
	ExecutionResults []types.DeployExecutionResult `json:"execution_results"`
	BlockHash        *key.Hash                     `json:"block_hash,omitempty"`
	BlockHeight      *uint64                       `json:"block_height,omitempty"`
}

type ChainGetEraInfoResult struct {
	Version    string           `json:"api_version"`
	EraSummary types.EraSummary `json:"era_summary"`
}

type StateGetItemResult struct {
	StoredValue types.StoredValue `json:"stored_value"`
	//MerkleProof is a construction created using a merkle trie that allows verification of the associated hashes.
	MerkleProof json.RawMessage `json:"merkle_proof"`
}

type StateGetDictionaryResult struct {
	ApiVersion    string            `json:"api_version"`
	DictionaryKey string            `json:"dictionary_key"`
	StoredValue   types.StoredValue `json:"stored_value"`
	MerkleProof   json.RawMessage   `json:"merkle_proof"`
}

type QueryGlobalStateResult struct {
	ApiVersion  string            `json:"api_version"`
	BlockHeader types.BlockHeader `json:"block_header,omitempty"`
	StoredValue types.StoredValue `json:"stored_value"`
	//MerkleProof is a construction created using a merkle trie that allows verification of the associated hashes.
	MerkleProof json.RawMessage `json:"merkle_proof"`
}

type InfoGetPeerResult struct {
	ApiVersion string     `json:"api_version"`
	Peers      []NodePeer `json:"peers"`
}

type NodePeer struct {
	NodeID  string `json:"node_id"`
	Address string `json:"address"`
}

type ChainGetStateRootHashResult struct {
	Version       string   `json:"api_version"`
	StateRootHash key.Hash `json:"state_root_hash"`
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
}

type ValidatorChanges struct {
	PublicKey     keypair.PublicKey `json:"public_key"`
	StatusChanges []StatusChanges   `json:"status_changes"`
}

type InfoGetValidatorChangesResult struct {
	APIVersion string             `json:"api_version"`
	Changes    []ValidatorChanges `json:"changes"`
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
}

type SpeculativeExecResult struct {
	ApiVersion      string                      `json:"api_version"`
	DeployHash      key.Hash                    `json:"deploy_hash"`
	ExecutionResult types.ExecutionResultStatus `json:"execution_result"`
}

type QueryBalanceResult struct {
	ApiVersion string          `json:"api_version"`
	Balance    clvalue.UInt512 `json:"balance"`
}

type InfoGetChainspecResult struct {
	ApiVersion     string `json:"api_version"`
	ChainspecBytes struct {
		ChainspecBytes            string `json:"chainspec_bytes,omitempty"`
		MaybeGenesisAccountsBytes string `json:"maybe_genesis_accounts_bytes,omitempty"`
		MaybeGlobalStateBytes     string `json:"maybe_global_state_bytes,omitempty"`
	} `json:"chainspec_bytes"`
}
