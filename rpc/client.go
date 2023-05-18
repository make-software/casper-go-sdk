/*
Package rpc provides access to the exported methods of RPC Client and data structures where serialized response.
See details in [README.md]
*/

package rpc

import (
	"context"

	"github.com/make-software/casper-go-sdk/types"
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
}

// ClientInformational contains methods that return information from a node on a Casper network.
// The response should be identical, regardless of the node queried,
// as the information in question is objective and common to all nodes within a network.
type ClientInformational interface {
	// GetAccountBalance returns a purse's balance from a network.
	// The request takes in the formatted representation of a purse URef as a parameter.
	GetAccountBalance(ctx context.Context, stateRootHash, purseURef string) (StateGetBalanceResult, error)
	// GetDeploy retrieves a Deploy from a network. It requires a deploy_hash to query the Deploy.
	GetDeploy(ctx context.Context, hash string) (InfoGetDeployResult, error)
	// GetDictionaryItem returns an item from a Dictionary.
	// Every dictionary has a seed URef, findable by using a dictionary_identifier.
	// The address of a stored value is the blake2b hash of the seed URef and the byte representation of the dictionary key.
	GetDictionaryItem(ctx context.Context, stateRootHash, uref, key string) (StateGetItemResult, error)
	// TODO: No Documentation about this method in the Casper specification
	GetStateItem(ctx context.Context, stateRootHash, key string, path []string) (StateGetItemResult, error)

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

	// GetStateRootHashLatest returns a state root hash of the latest Block.
	GetStateRootHashLatest(ctx context.Context) (ChainGetStateRootHashResult, error)
	// GetStateRootHashByHash returns a state root hash of the latest Block the requested block hash.
	GetStateRootHashByHash(ctx context.Context, stateRootHash string) (ChainGetStateRootHashResult, error)
	// GetStateRootHashByHeight returns a state root hash of the latest Block the requested block height.
	GetStateRootHashByHeight(ctx context.Context, height uint64) (ChainGetStateRootHashResult, error)

	// GetStatus return the current status of a node on a Casper network.
	// The responses return information specific to the queried node, and as such, will vary.
	//TODO: maybe move to separate interface NodeInfoClient
	GetStatus(ctx context.Context) (InfoGetStatusResult, error)
	// GetPeers return a list of peers connected to the node on a Casper network.
	// The responses return information specific to the queried node, and as such, will vary.
	GetPeers(ctx context.Context) (InfoGetPeerResult, error)
}

// ClientTransactional contains the description of account_put_deploy,
// the only means by which users can send their compiled Wasm (as part of a Deploy) to a node on a Casper network.
type ClientTransactional interface {
	PutDeploy(ctx context.Context, deploy types.Deploy) (PutDeployResult, error)
}

// Client interface represent full RPC client that includes all possible queries.
//
//go:generate mockgen -destination=../tests/mocks/rpc_client_mock.go -package=mocks -source=./client.go Client
type Client interface {
	ClientPOS
	ClientInformational
	ClientTransactional
}

// Handler is responsible to implement interaction with underlying protocol.
type Handler interface {
	ProcessCall(ctx context.Context, params RpcRequest) (RpcResponse, error)
}