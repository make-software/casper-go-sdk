package rpc

import (
	"context"
	"strconv"
	"strings"

	"github.com/make-software/casper-go-sdk/v2/types"
	"github.com/make-software/casper-go-sdk/v2/types/key"
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

var ApiVersion = "2.0"

type CtxRequestID string

const RequestIDKey CtxRequestID = "RequestID"

func WithRequestId(ctx context.Context, requestID int) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func GetReqIdCtx(ctx context.Context) string {
	value := ctx.Value(RequestIDKey)
	if value == nil {
		return "0"
	}
	return strconv.Itoa(value.(int))
}

// Method is represented a name of the RPC endpoint
type Method string

const (
	MethodGetDeploy           Method = "info_get_deploy"
	MethodGetTransaction      Method = "info_get_transaction"
	MethodGetStateItem        Method = "state_get_item"
	MethodQueryGlobalState    Method = "query_global_state"
	MethodGetDictionaryItem   Method = "state_get_dictionary_item"
	MethodGetStateBalance     Method = "state_get_balance"
	MethodGetStateAccount     Method = "state_get_account_info"
	MethodGetStatePackage     Method = "state_get_package"
	MethodGetStateEntity      Method = "state_get_entity"
	MethodGetEraInfo          Method = "chain_get_era_info_by_switch_block"
	MethodGetBlock            Method = "chain_get_block"
	MethodGetBlockTransfers   Method = "chain_get_block_transfers"
	MethodGetEraSummary       Method = "chain_get_era_summary"
	MethodGetAuctionInfo      Method = "state_get_auction_info"
	MethodGetAuctionInfoV2    Method = "state_get_auction_info_v2"
	MethodGetValidatorChanges Method = "info_get_validator_changes"
	MethodGetStateRootHash    Method = "chain_get_state_root_hash"
	MethodGetStatus           Method = "info_get_status"
	MethodGetReward           Method = "info_get_reward"
	MethodGetPeers            Method = "info_get_peers"
	MethodPutDeploy           Method = "account_put_deploy"
	MethodPutTransaction      Method = "account_put_transaction"
	MethodSpeculativeExec     Method = "speculative_exec"
	MethodQueryBalance        Method = "query_balance"
	MethodQueryBalanceDetails Method = "query_balance_details"
	MethodInfoGetChainspec    Method = "info_get_chainspec"
)

// RpcRequest is a wrapper struct for an RPC call method that can be serialized to JSON.
type RpcRequest struct {
	// Version of the RPC protocol in use
	Version string `json:"jsonrpc"`
	// Id of the RPC request that can be correlated with the equivalent Id in the RPC response
	ID     *IDValue    `json:"id,omitempty"`
	Method Method      `json:"method"`
	Params interface{} `json:"params"`
}

func DefaultRpcRequest(method Method, params interface{}) RpcRequest {
	return RpcRequest{
		Version: ApiVersion,
		ID:      NewIDFromString("1"),
		Method:  method,
		Params:  params,
	}
}

type ParamStateRootHash struct {
	StateRootHash string   `json:"state_root_hash"`
	Key           string   `json:"key"`
	Path          []string `json:"path,omitempty"`
}

type ParamQueryGlobalState struct {
	StateIdentifier *ParamQueryGlobalStateID `json:"state_identifier,omitempty"`
	Key             string                   `json:"key"`
	Path            []string                 `json:"path,omitempty"`
}

type ParamQueryGlobalStateID struct {
	StateRootHash string  `json:"StateRootHash,omitempty"`
	BlockHash     string  `json:"BlockHash,omitempty"`
	BlockHeight   *uint64 `json:"BlockHeight,omitempty"`
}

type ParamTransactionHash struct {
	TransactionHash    types.TransactionHash `json:"transaction_hash,omitempty"`
	FinalizedApprovals *bool                 `json:"finalized_approvals,omitempty"`
}

func NewQueryGlobalStateParam(key string, path []string, id *ParamQueryGlobalStateID) ParamQueryGlobalState {
	return ParamQueryGlobalState{StateIdentifier: id, Key: key, Path: path}
}

type ParamGetAccountInfoBalance struct {
	AccountIdentifier string `json:"account_identifier"`
	ParamBlockIdentifier
}

type ParamGetStateEntity struct {
	EntityIdentifier EntityIdentifier `json:"entity_identifier"`
	BlockIdentifier  *BlockIdentifier `json:"block_identifier,omitempty"`
}

type AccountIdentifier struct {
	AccountHash *key.AccountHash
	PublicKey   *keypair.PublicKey
}

// EntityIdentifier Identifier of an addressable entity.
type EntityIdentifier struct {
	// The account hash of an account.
	AccountHash *key.AccountHash `json:"AccountHash,omitempty"`
	// The public key of an account.
	PublicKey *keypair.PublicKey `json:"PublicKey,omitempty"`
	// The address of an addressable entity.
	EntityAddr *key.EntityAddr `json:"EntityAddr,omitempty"`
}

func NewEntityIdentifierFromAccountHash(accountHash key.AccountHash) EntityIdentifier {
	return EntityIdentifier{
		AccountHash: &accountHash,
	}
}

func NewEntityIdentifierFromPublicKey(pubKey keypair.PublicKey) EntityIdentifier {
	return EntityIdentifier{
		PublicKey: &pubKey,
	}
}

func NewEntityIdentifierFromEntityAddr(entityAddr key.EntityAddr) EntityIdentifier {
	return EntityIdentifier{
		EntityAddr: &entityAddr,
	}
}

type PutDeployRequest struct {
	Deploy types.Deploy `json:"deploy"`
}

type PutTransactionRequest struct {
	Transaction types.TransactionWrapper `json:"transaction"`
}

type BlockIdentifier struct {
	Hash   *string `json:"Hash,omitempty"`
	Height *uint64 `json:"Height,omitempty"`
}

type GlobalStateIdentifier struct {
	BlockHash   *string `json:"BlockHash,omitempty"`
	BlockHeight *uint64 `json:"BlockHeight,omitempty"`
	StateRoot   *string `json:"StateRootHash,omitempty"`
}

type EraIdentifier struct {
	Block *BlockIdentifier `json:"Block,omitempty"`
	Era   *uint64          `json:"Era,omitempty"`
}

type ParamBlockIdentifier struct {
	BlockIdentifier *BlockIdentifier `json:"block_identifier"`
}

func NewParamBlockByHeight(height uint64) ParamBlockIdentifier {
	return ParamBlockIdentifier{BlockIdentifier: &BlockIdentifier{Height: &height}}
}

func NewParamBlockByHash(hash string) ParamBlockIdentifier {
	return ParamBlockIdentifier{BlockIdentifier: &BlockIdentifier{Hash: &hash}}
}

type ParamDictionaryIdentifier struct {
	AccountNamedKey  *AccountNamedKey                           `json:"AccountNamedKey,omitempty"`
	ContractNamedKey *ParamDictionaryIdentifierContractNamedKey `json:"ContractNamedKey,omitempty"`
	URef             *ParamDictionaryIdentifierURef             `json:"URef,omitempty"`
	Dictionary       *string                                    `json:"Dictionary,omitempty"`
}

type AccountNamedKey struct {
	Key               string `json:"key"`
	DictionaryName    string `json:"dictionary_name"`
	DictionaryItemKey string `json:"dictionary_item_key"`
}

type ParamDictionaryIdentifierContractNamedKey struct {
	Key               string `json:"key"`
	DictionaryName    string `json:"dictionary_name"`
	DictionaryItemKey string `json:"dictionary_item_key"`
}

type ParamDictionaryIdentifierURef struct {
	DictionaryItemKey string `json:"dictionary_item_key"`
	SeedUref          string `json:"seed_uref"`
}

type SpeculativeExecParams struct {
	Deploy          types.Deploy     `json:"deploy"`
	BlockIdentifier *BlockIdentifier `json:"block_identifier,omitempty"`
}

type ParamStateGetPackage struct {
	PackageIdentifier PackageIdentifier `json:"package_identifier"`
	ParamBlockIdentifier
}

type PackageIdentifier struct {
	PackageAddr         *string `json:"PackageAddr,omitempty"`
	ContractPackageHash *string `json:"ContractPackageHash,omitempty"`
}

func NewPackageIdentifierFromHash(hash string) PackageIdentifier {
	if strings.HasPrefix(hash, "package-") {
		return PackageIdentifier{
			PackageAddr: &hash,
		}
	}
	return PackageIdentifier{
		ContractPackageHash: &hash,
	}
}

type PurseIdentifier struct {
	MainPurseUnderPublicKey   *keypair.PublicKey `json:"main_purse_under_public_key,omitempty"`
	MainPurseUnderAccountHash *key.AccountHash   `json:"main_purse_under_account_hash,omitempty"`
	MainPurseUnderEntityAddr  *key.EntityAddr    `json:"main_purse_under_entity_addr,omitempty"`
	PurseUref                 *key.URef          `json:"purse_uref,omitempty"`
}

func NewPurseIdentifierFromPublicKey(pubKey keypair.PublicKey) PurseIdentifier {
	return PurseIdentifier{
		MainPurseUnderPublicKey: &pubKey,
	}
}

func NewPurseIdentifierFromAccountHash(accountHash key.AccountHash) PurseIdentifier {
	return PurseIdentifier{
		MainPurseUnderAccountHash: &accountHash,
	}
}

func NewPurseIdentifierFromEntityAddr(entityAddr key.EntityAddr) PurseIdentifier {
	return PurseIdentifier{
		MainPurseUnderEntityAddr: &entityAddr,
	}
}

func NewPurseIdentifierFromUref(uref key.URef) PurseIdentifier {
	return PurseIdentifier{
		PurseUref: &uref,
	}
}

type QueryBalanceRequest struct {
	PurseIdentifier PurseIdentifier        `json:"purse_identifier"`
	StateIdentifier *GlobalStateIdentifier `json:"state_identifier,omitempty"`
}

type QueryBalanceDetailsRequest struct {
	PurseIdentifier PurseIdentifier        `json:"purse_identifier"`
	StateIdentifier *GlobalStateIdentifier `json:"state_identifier,omitempty"`
}

type InfoGetRewardRequest struct {
	Validator     keypair.PublicKey  `json:"validator"`
	Delegator     *keypair.PublicKey `json:"delegator,omitempty"`
	EraIdentifier *EraIdentifier     `json:"era_identifier,omitempty"`
}
