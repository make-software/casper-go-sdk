/*
Package rpc provides access to the exported methods of RPC Client and data structures where serialized response.
See details in [README.md]
*/

package rpc

import (
	"context"

	"github.com/make-software/casper-go-sdk/v2/types"
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

// ClientPOS contains methods pertain to the Proof-of-Stake functionality of a Casper network.
// They return information related to auctions, bids and validators.
// This information is necessary for users involved with node operations and validation.
type ClientPOS interface {
	// GetLatestAuctionInfo returns the types.ValidatorBid and types.EraValidators from the most recent Block.
	// RPC: state_get_auction_info_v2 with fallback on state_get_auction_info
	GetLatestAuctionInfo(ctx context.Context) (StateGetAuctionInfoResult, error)
	// GetAuctionInfoByHash returns the types.ValidatorBid and types.EraValidators of either a specific Block by hash
	// RPC: state_get_auction_info_v2 with fallback on state_get_auction_info
	GetAuctionInfoByHash(ctx context.Context, blockHash string) (StateGetAuctionInfoResult, error)
	// GetAuctionInfoByHeight returns the types.ValidatorBid and types.EraValidators of either a specific Block by height
	// RPC: state_get_auction_info_v2 with fallback on state_get_auction_info
	GetAuctionInfoByHeight(ctx context.Context, height uint64) (StateGetAuctionInfoResult, error)
	// GetLatestAuctionInfoV1 returns the types.ValidatorBid and types.EraValidators from the most recent Block.
	// RPC: state_get_auction_info
	GetLatestAuctionInfoV1(ctx context.Context) (StateGetAuctionInfoV1Result, error)
	// GetAuctionInfoV1ByHash returns the types.ValidatorBid and types.EraValidators of either a specific Block by hash
	// RPC: state_get_auction_info
	GetAuctionInfoV1ByHash(ctx context.Context, blockHash string) (StateGetAuctionInfoV1Result, error)
	// GetAuctionInfoV1ByHeight returns the types.ValidatorBid and types.EraValidators of either a specific Block by height
	// RPC: state_get_auction_info
	GetAuctionInfoV1ByHeight(ctx context.Context, height uint64) (StateGetAuctionInfoV1Result, error)
	// GetLatestAuctionInfoV2 returns the types.ValidatorBid and types.EraValidators from the most recent Block.
	// RPC: state_get_auction_info_v2
	GetLatestAuctionInfoV2(ctx context.Context) (StateGetAuctionInfoV2Result, error)
	// GetAuctionInfoV2ByHash returns the types.ValidatorBid and types.EraValidators of either a specific Block by hash
	// RPC: state_get_auction_info_v2
	GetAuctionInfoV2ByHash(ctx context.Context, blockHash string) (StateGetAuctionInfoV2Result, error)
	// GetAuctionInfoV2ByHeight returns the types.ValidatorBid and types.EraValidators of either a specific Block by height
	// RPC: state_get_auction_info_v2
	GetAuctionInfoV2ByHeight(ctx context.Context, height uint64) (StateGetAuctionInfoV2Result, error)
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
	// GetLatestBalance returns a purse's balance from a network.
	// The request takes in the formatted representation of a purse URef as a parameter.
	// The client will make an additional RPC call to retrieve the latest stateRootHash.
	GetLatestBalance(ctx context.Context, purseURef string) (StateGetBalanceResult, error)
	// GetBalanceByStateRootHash returns a purse's balance and state root hash from a network.
	GetBalanceByStateRootHash(ctx context.Context, purseURef string, stateRootHash string) (StateGetBalanceResult, error)
	// GetDeploy retrieves a Deploy from a network. It requires a deploy_hash to query the Deploy.
	GetDeploy(ctx context.Context, hash string) (InfoGetDeployResult, error)
	// GetDeployFinalizedApproval returns Deploy with the finalized approvals substituted.
	GetDeployFinalizedApproval(ctx context.Context, hash string) (InfoGetDeployResult, error)
	// GetTransactionByTransactionHash returns a Transaction from the network
	GetTransactionByTransactionHash(ctx context.Context, transactionHash string) (InfoGetTransactionResult, error)
	// GetTransactionByDeployHash returns a Deploy as Transaction from the network
	GetTransactionByDeployHash(ctx context.Context, deployHash string) (InfoGetTransactionResult, error)
	// GetTransactionFinalizedApprovalByTransactionHash return the Transaction with the finalized approvals substituted.
	GetTransactionFinalizedApprovalByTransactionHash(ctx context.Context, transactionHash string) (InfoGetTransactionResult, error)
	// GetTransactionFinalizedApprovalByDeployHash return the Deploy as Transaction with the finalized approvals substituted.
	GetTransactionFinalizedApprovalByDeployHash(ctx context.Context, deployHash string) (InfoGetTransactionResult, error)
	// GetDictionaryItem returns an item from a Dictionary.
	// The address of a stored value is the blake2b hash of the seed URef and the byte representation of the dictionary key.
	// If the param stateRootHash is nil, the client will make an additional RPC call to retrieve the latest stateRootHash.
	GetDictionaryItem(ctx context.Context, stateRootHash *string, uref, key string) (StateGetDictionaryResult, error)
	// GetDictionaryItemByIdentifier returns an item from a Dictionary.
	// Every dictionary has a seed URef, findable by using a dictionary_identifier.
	GetDictionaryItemByIdentifier(ctx context.Context, stateRootHash *string, identifier ParamDictionaryIdentifier) (StateGetDictionaryResult, error)
	// GetStateItem allows to get item from the global state
	// If the param stateRootHash is nil, the client will make an additional RPC call to retrieve the latest stateRootHash.
	// Deprecated: use QueryGlobalStateByStateHash instead
	GetStateItem(ctx context.Context, stateRootHash *string, key string, path []string) (StateGetItemResult, error)

	// QueryLatestGlobalState allows for you to query for the latest value stored under certain keys in global state.
	QueryLatestGlobalState(ctx context.Context, key string, path []string) (QueryGlobalStateResult, error)
	// QueryGlobalStateByBlockHash allows for you to query for a value stored under certain keys in global state.
	QueryGlobalStateByBlockHash(ctx context.Context, blockHash, key string, path []string) (QueryGlobalStateResult, error)
	// QueryGlobalStateByBlockHeight allows for you to query for a value stored under certain keys in global state.
	QueryGlobalStateByBlockHeight(ctx context.Context, blockHeight uint64, key string, path []string) (QueryGlobalStateResult, error)
	// QueryGlobalStateByStateHash allows for you to query for a value stored under certain keys in global state.
	// If the param stateRootHash is nil, the client will make an additional RPC call to retrieve the latest stateRootHash.
	QueryGlobalStateByStateHash(ctx context.Context, stateRootHash *string, key string, path []string) (QueryGlobalStateResult, error)

	// GetAccountInfoByBlockHash returns a JSON representation of an Account from the network.
	// The blockHash must refer to a  Block after the Account's creation, or the method will return an empty response.
	GetAccountInfoByBlockHash(ctx context.Context, blockHash string, pub keypair.PublicKey) (StateGetAccountInfo, error)
	// GetAccountInfoByBlockHeight returns a JSON representation of an Account from the network.
	// The blockHeight must refer to a Block after the Account's creation, or the method will return an empty response.
	GetAccountInfoByBlockHeight(ctx context.Context, blockHeight uint64, pub keypair.PublicKey) (StateGetAccountInfo, error)
	// GetAccountInfo returns a JSON representation of an Account from the network.
	// This is the most generic interface.
	GetAccountInfo(ctx context.Context, blockIdentifier *ParamBlockIdentifier, accountIdentifier AccountIdentifier) (StateGetAccountInfo, error)
	// GetPackageByBlockHeight returns a Package from the network by BlockHeight
	// The blockHeight must refer to a Block after the Package's creation, or the method will return an empty response.
	GetPackageByBlockHeight(ctx context.Context, packageHash string, blockHeight uint64) (StateGetPackage, error)
	// GetPackageByBlockHash returns a Package from the network by BlockHash
	// The blockHash must refer to a Block after the Package's creation, or the method will return an empty response.
	GetPackageByBlockHash(ctx context.Context, packageHash string, blockHash string) (StateGetPackage, error)
	// GetPackage returns a Package from the network
	// This is the most generic interface.
	GetPackage(ctx context.Context, packageIdentifier PackageIdentifier, blockIdentifier *ParamBlockIdentifier) (StateGetPackage, error)

	// GetLatestEntity returns latest AddressableEntity from the network.
	GetLatestEntity(ctx context.Context, entityIdentifier EntityIdentifier) (StateGetEntityResult, error)
	// GetEntityByBlockHash returns an AddressableEntity by block hash from the network.
	GetEntityByBlockHash(ctx context.Context, entityIdentifier EntityIdentifier, hash string) (StateGetEntityResult, error)
	// GetEntityByBlockHeight returns an AddressableEntity by block height from the network.
	GetEntityByBlockHeight(ctx context.Context, entityIdentifier EntityIdentifier, height uint64) (StateGetEntityResult, error)

	// GetLatestBlock returns the latest types.Block from the network.
	GetLatestBlock(ctx context.Context) (ChainGetBlockResult, error)
	// GetBlockByHash returns the types.Block from the network the requested block hash.
	GetBlockByHash(ctx context.Context, hash string) (ChainGetBlockResult, error)
	// GetBlockByHeight returns the types.Block from the network the requested block height.
	GetBlockByHeight(ctx context.Context, height uint64) (ChainGetBlockResult, error)

	// GetLatestBlockTransfers returns all native transfers within a lasted Block from a network.
	GetLatestBlockTransfers(ctx context.Context) (ChainGetBlockTransfersResult, error)
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

	// QueryLatestBalance queries for balances under a given PurseIdentifier
	QueryLatestBalance(ctx context.Context, identifier PurseIdentifier) (QueryBalanceResult, error)
	// QueryBalanceByBlockHeight query for balance information using a purse identifier and block height
	QueryBalanceByBlockHeight(ctx context.Context, purseIdentifier PurseIdentifier, height uint64) (QueryBalanceResult, error)
	// QueryBalanceByBlockHash query for balance information using a purse identifier and block hash
	QueryBalanceByBlockHash(ctx context.Context, purseIdentifier PurseIdentifier, blockHash string) (QueryBalanceResult, error)
	// QueryBalanceByStateRootHash query for full balance information using a purse identifier and state root hash
	QueryBalanceByStateRootHash(ctx context.Context, purseIdentifier PurseIdentifier, stateRootHash string) (QueryBalanceResult, error)

	// QueryLatestBalanceDetails query for full balance information using a purse identifier
	QueryLatestBalanceDetails(ctx context.Context, purseIdentifier PurseIdentifier) (QueryBalanceDetailsResult, error)
	// QueryBalanceDetailsByBlockHeight query for full balance information using a purse identifier and block height
	QueryBalanceDetailsByBlockHeight(ctx context.Context, purseIdentifier PurseIdentifier, height uint64) (QueryBalanceDetailsResult, error)
	// QueryBalanceDetailsByBlockHash query for full balance information using a purse identifier and block hash
	QueryBalanceDetailsByBlockHash(ctx context.Context, purseIdentifier PurseIdentifier, blockHash string) (QueryBalanceDetailsResult, error)
	// QueryBalanceDetailsByStateRootHash query for full balance information using a purse identifier and state root hash
	QueryBalanceDetailsByStateRootHash(ctx context.Context, purseIdentifier PurseIdentifier, stateRootHash string) (QueryBalanceDetailsResult, error)

	// GetChainspec returns the raw bytes of the chainspec.toml, accounts.toml and global_state.toml files as read at node startup.
	GetChainspec(ctx context.Context) (InfoGetChainspecResult, error)

	// GetLatestValidatorReward returns the latest reward for a given validator
	GetLatestValidatorReward(ctx context.Context, validator keypair.PublicKey) (InfoGetRewardResult, error)
	// GetValidatorRewardByEraID returns the reward for a given era and a validator
	GetValidatorRewardByEraID(ctx context.Context, validator keypair.PublicKey, eraID uint64) (InfoGetRewardResult, error)
	// GetValidatorRewardByBlockHash returns the reward for a given block hash and a validator
	GetValidatorRewardByBlockHash(ctx context.Context, validator keypair.PublicKey, blockHash string) (InfoGetRewardResult, error)
	// GetValidatorRewardByBlockHeight returns the reward for a given block height and a validator
	GetValidatorRewardByBlockHeight(ctx context.Context, validator keypair.PublicKey, height uint64) (InfoGetRewardResult, error)
	// GetLatestDelegatorReward returns the latest delegator reward for a given validator
	GetLatestDelegatorReward(ctx context.Context, validator, delegator keypair.PublicKey) (InfoGetRewardResult, error)
	// GetDelegatorRewardByEraID returns the delegator reward for a given era and a validator
	GetDelegatorRewardByEraID(ctx context.Context, validator, delegator keypair.PublicKey, eraID uint64) (InfoGetRewardResult, error)
	// GetDelegatorRewardByBlockHash returns the delegator reward for a given block hash and a validator
	GetDelegatorRewardByBlockHash(ctx context.Context, validator, delegator keypair.PublicKey, blockHash string) (InfoGetRewardResult, error)
	// GetDelegatorRewardByBlockHeight returns the delegator reward for a given block height and a validator
	GetDelegatorRewardByBlockHeight(ctx context.Context, validator, delegator keypair.PublicKey, height uint64) (InfoGetRewardResult, error)
}

// ClientTransactional contains the description of account_put_deploy, account_put_transaction
// the only means by which users can send their compiled Wasm (as part of a Deploy or TransactionV1) to a node on a Casper network.
type ClientTransactional interface {
	PutDeploy(ctx context.Context, deploy types.Deploy) (PutDeployResult, error)
	PutTransactionV1(ctx context.Context, transaction types.TransactionV1) (PutTransactionResult, error)
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
