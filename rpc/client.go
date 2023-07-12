/*
Package rpc provides access to the exported methods of RPC Client and data structures where serialized response.
See details in [README.md]
*/

package rpc

import (
	"context"

	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/keypair"
)

// ClientPOS contains methods pertain to the Proof-of-Stake functionality of a Casper network.
// They return information related to auctions, bids and validators.
// This information is necessary for users involved with node operations and validation.
type ClientPOS interface {
	// GetAuctionInfoLatest returns the types.ValidatorBid and types.EraValidators from the most recent Block.
	GetAuctionInfoLatest(ctx context.Context) (StateGetAuctionInfoResult, error)
	// GetAuctionInfoByHash returns the types.ValidatorBid and types.EraValidators of either a specific Block by hash
	GetAuctionInfoByHash(ctx context.Context, blockHash string) (StateGetAuctionInfoResult, error)
	// GetAuctionInfoByHeight returns the types.ValidatorBid and types.EraValidators of either a specific Block by height
	GetAuctionInfoByHeight(ctx context.Context, height uint64) (StateGetAuctionInfoResult, error)

	// GetEraInfoLatest returns an EraInfo from the network.
	// Only the last Block in an era, known as a switch block, will contain an era_summary.
	// This method return information about the latest block in the chain, it may not be the last block in the era.
	GetEraInfoLatest(ctx context.Context) (ChainGetEraInfoResult, error)
	// GetEraInfoByBlockHeight returns an EraInfo from the network.
	// Only the last Block in an era, known as a switch block, will contain an era_summary.
	// Querying by block height.
	GetEraInfoByBlockHeight(ctx context.Context, height uint64) (ChainGetEraInfoResult, error)
	// GetEraInfoByBlockHash returns an EraInfo from the network.
	// Only the last Block in an era, known as a switch block, will contain an era_summary.
	// Querying by block hash.
	GetEraInfoByBlockHash(ctx context.Context, hash string) (ChainGetEraInfoResult, error)

	// GetValidatorChangesInfo returns status changes of active validators. Listed changes occurred during the EraId
	// contained within the response itself. A validator may show more than one change in a single era.
	GetValidatorChangesInfo(ctx context.Context) (InfoGetValidatorChangesResult, error)
}

// ClientInformational contains methods that return information from a node on a Casper network.
// The response should be identical, regardless of the node queried,
// as the information in question is objective and common to all nodes within a network.
type ClientInformational interface {
	// GetAccountBalance returns a purse's balance from a network.
	// The request takes in the formatted representation of a purse URef as a parameter.
	// If the param stateRootHash is nil, the client will make an additional RPC call to retrieve the latest stateRootHash.
	GetAccountBalance(ctx context.Context, stateRootHash *string, purseURef string) (StateGetBalanceResult, error)
	// GetDeploy retrieves a Deploy from a network. It requires a deploy_hash to query the Deploy.
	GetDeploy(ctx context.Context, hash string) (InfoGetDeployResult, error)
	// GetDeployFinalizedApproval returns Deploy with the finalized approvals substituted.
	GetDeployFinalizedApproval(ctx context.Context, hash string) (InfoGetDeployResult, error)
	// GetDictionaryItem returns an item from a Dictionary.
	// Every dictionary has a seed URef, findable by using a dictionary_identifier.
	// The address of a stored value is the blake2b hash of the seed URef and the byte representation of the dictionary key.
	// If the param stateRootHash is nil, the client will make an additional RPC call to retrieve the latest stateRootHash.
	GetDictionaryItem(ctx context.Context, stateRootHash *string, uref, key string) (StateGetDictionaryResult, error)
	// GetStateItem allows to get item from the global state
	// If the param stateRootHash is nil, the client will make an additional RPC call to retrieve the latest stateRootHash.
	// Deprecated: use QueryGlobalStateByStateHash instead
	GetStateItem(ctx context.Context, stateRootHash *string, key string, path []string) (StateGetItemResult, error)

	// QueryGlobalStateByBlockHash allows for you to query for a value stored under certain keys in global state.
	QueryGlobalStateByBlockHash(ctx context.Context, blockHash, key string, path []string) (QueryGlobalStateResult, error)
	// QueryGlobalStateByStateHash allows for you to query for a value stored under certain keys in global state.
	// If the param stateRootHash is nil, the client will make an additional RPC call to retrieve the latest stateRootHash.
	QueryGlobalStateByStateHash(ctx context.Context, stateRootHash *string, key string, path []string) (QueryGlobalStateResult, error)

	// GetAccountInfoByBlochHash returns a JSON representation of an Account from the network.
	// The blockHash must refer to a  Block after the Account's creation, or the method will return an empty response.
	GetAccountInfoByBlochHash(ctx context.Context, blockHash string, pub keypair.PublicKey) (StateGetAccountInfo, error)
	// GetAccountInfoByBlochHeight returns a JSON representation of an Account from the network.
	// The blockHeight must refer to a Block after the Account's creation, or the method will return an empty response.
	GetAccountInfoByBlochHeight(ctx context.Context, blockHeight uint64, pub keypair.PublicKey) (StateGetAccountInfo, error)

	// GetBlockLatest returns the latest types.Block from the network.
	GetBlockLatest(ctx context.Context) (ChainGetBlockResult, error)
	// GetBlockByHash returns the types.Block from the network the requested block hash.
	GetBlockByHash(ctx context.Context, hash string) (ChainGetBlockResult, error)
	// GetBlockByHeight returns the types.Block from the network the requested block height.
	GetBlockByHeight(ctx context.Context, height uint64) (ChainGetBlockResult, error)

	// GetBlockTransfersLatest returns all native transfers within a lasted Block from a network.
	GetBlockTransfersLatest(ctx context.Context) (ChainGetBlockTransfersResult, error)
	// GetBlockTransfersByHash returns all native transfers within a given Block from a network the requested block hash.
	GetBlockTransfersByHash(ctx context.Context, blockHash string) (ChainGetBlockTransfersResult, error)
	// GetBlockTransfersByHeight returns all native transfers within a given Block from a network the requested block height.
	GetBlockTransfersByHeight(ctx context.Context, height uint64) (ChainGetBlockTransfersResult, error)

	// GetEraSummaryLatest returns the era summary at a latest Block.
	GetEraSummaryLatest(ctx context.Context) (ChainGetEraSummaryResult, error)
	// GetEraSummaryByHash returns the era summary at a Block by hash.
	GetEraSummaryByHash(ctx context.Context, blockHash string) (ChainGetEraSummaryResult, error)
	// GetEraSummaryByHeight returns the era summary at a Block by height.
	GetEraSummaryByHeight(ctx context.Context, height uint64) (ChainGetEraSummaryResult, error)

	// GetStateRootHashLatest returns a state root hash of the latest Block.
	GetStateRootHashLatest(ctx context.Context) (ChainGetStateRootHashResult, error)
	// GetStateRootHashByHash returns a state root hash of the latest Block the requested block hash.
	GetStateRootHashByHash(ctx context.Context, blockHash string) (ChainGetStateRootHashResult, error)
	// GetStateRootHashByHeight returns a state root hash of the latest Block the requested block height.
	GetStateRootHashByHeight(ctx context.Context, height uint64) (ChainGetStateRootHashResult, error)

	// GetStatus return the current status of a node on a Casper network.
	// The responses return information specific to the queried node, and as such, will vary.
	GetStatus(ctx context.Context) (InfoGetStatusResult, error)
	// GetPeers return a list of peers connected to the node on a Casper network.
	// The responses return information specific to the queried node, and as such, will vary.
	GetPeers(ctx context.Context) (InfoGetPeerResult, error)
	// QueryBalance queries for balances under a given PurseIdentifier
	QueryBalance(ctx context.Context, identifier PurseIdentifier) (QueryBalanceResult, error)
	// GetChainspec returns the raw bytes of the chainspec.toml, accounts.toml and global_state.toml files as read at node startup.
	GetChainspec(ctx context.Context) (InfoGetChainspecResult, error)
}

// ClientTransactional contains the description of account_put_deploy,
// the only means by which users can send their compiled Wasm (as part of a Deploy) to a node on a Casper network.
type ClientTransactional interface {
	PutDeploy(ctx context.Context, deploy types.Deploy) (PutDeployResult, error)
}

// Client interface represent full RPC client that includes all possible queries.
type Client interface {
	ClientPOS
	ClientInformational
	ClientTransactional
}

// Handler is responsible to implement interaction with underlying protocol.
type Handler interface {
	ProcessCall(ctx context.Context, params RpcRequest) (RpcResponse, error)
}
