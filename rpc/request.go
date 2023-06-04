package rpc

import (
	"context"

	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/keypair"
)

var ApiVersion = "2.0"

type CtxRequestID string

const RequestIDKey CtxRequestID = "RequestID"

func WithRequestId(ctx context.Context, requestID int) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func GetReqIdCtx(ctx context.Context) int {
	value := ctx.Value(RequestIDKey)
	if value == nil {
		return 0
	}
	return value.(int)
}

// Method is represented a name of the RPC endpoint
type Method string

const (
	MethodGetDeploy           Method = "info_get_deploy"
	MethodGetStateItem        Method = "state_get_item"
	MethodQueryGlobalState    Method = "query_global_state"
	MethodGetDictionaryItem   Method = "state_get_dictionary_item"
	MethodGetStateBalance     Method = "state_get_balance"
	MethodGetStateAccount     Method = "state_get_account_info"
	MethodGetEraInfo          Method = "chain_get_era_info_by_switch_block"
	MethodGetBlock            Method = "chain_get_block"
	MethodGetBlockTransfers   Method = "chain_get_block_transfers"
	MethodGetEraSummary       Method = "chain_get_era_summary"
	MethodGetAuctionInfo      Method = "state_get_auction_info"
	MethodGetValidatorChanges Method = "info_get_validator_changes"
	MethodGetStateRootHash    Method = "chain_get_state_root_hash"
	MethodGetStatus           Method = "info_get_status"
	MethodGetPeers            Method = "info_get_peers"
	MethodPutDeploy           Method = "account_put_deploy"
)

// RpcRequest is a wrapper struct for an RPC call method that can be serialized to JSON.
type RpcRequest struct {
	// Version of the RPC protocol in use
	Version string `json:"jsonrpc"`
	// Id of the RPC request that can be correlated with the equivalent Id in the RPC response
	//TODO: ID doesn't work from the Node side (always return 1 in response)
	ID     int         `json:"id"`
	Method Method      `json:"method"`
	Params interface{} `json:"params"`
}

func DefaultRpcRequest(method Method, params interface{}) RpcRequest {
	return RpcRequest{
		Version: ApiVersion,
		ID:      1,
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
	StateIdentifier ParamQueryGlobalStateID `json:"state_identifier"`
	Key             string                  `json:"key"`
	Path            []string                `json:"path,omitempty"`
}

type ParamQueryGlobalStateID struct {
	StateRootHash string `json:"StateRootHash,omitempty"`
	BlockHash     string `json:"BlockHash,omitempty"`
}

func NewQueryGlobalStateParam(key string, path []string, id ParamQueryGlobalStateID) ParamQueryGlobalState {
	return ParamQueryGlobalState{StateIdentifier: id, Key: key, Path: path}
}

type ParamGetAccountInfoBalance struct {
	PublicKey keypair.PublicKey `json:"public_key"`
	ParamBlockIdentifier
}

type PutDeployRequest struct {
	Deploy types.Deploy `json:"deploy"`
}

type BlockIdentifier struct {
	Hash   string `json:"Hash,omitempty"`
	Height uint64 `json:"Height,omitempty"`
}

type ParamBlockIdentifier struct {
	BlockIdentifier BlockIdentifier `json:"block_identifier"`
}

func NewParamBlockByHeight(height uint64) ParamBlockIdentifier {
	return ParamBlockIdentifier{BlockIdentifier: BlockIdentifier{Height: height}}
}

func NewParamBlockByHash(hash string) ParamBlockIdentifier {
	return ParamBlockIdentifier{BlockIdentifier: BlockIdentifier{Hash: hash}}
}

func NewParamStateDictionaryItem(stateRootHash, uref, key string) map[string]interface{} {
	return map[string]interface{}{
		"state_root_hash": stateRootHash,
		"dictionary_identifier": map[string]interface{}{
			"URef": map[string]string{
				"dictionary_item_key": key,
				"seed_uref":           uref,
			},
		},
	}
}
