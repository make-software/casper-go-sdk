package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/make-software/casper-go-sdk/types"
)

var ErrResultUnmarshal = errors.New("failed to unmarshal rpc result")

// client implements a Client interface. The main responsibility is to provide a useful set of methods to interact with
// the external RPC server. client declares the methods' signatures, builds RpcRequest from signature params,
// and serializes RpcResponse to the corresponding data structures.
// Most interaction work with RPC delegates to the Handler.
type client struct {
	handler Handler
}

// NewClient is a constructor for client that suppose to configure Handler
// examples of usage can be found here [Test_ConfigurableClient_GetDeploy]
func NewClient(handler Handler) Client {
	return &client{handler: handler}
}

func (c *client) GetDeploy(ctx context.Context, hash string) (InfoGetDeployResult, error) {
	var result InfoGetDeployResult
	return result, c.processRequest(ctx, MethodGetDeploy, map[string]string{
		"deploy_hash": hash,
	}, &result)
}

func (c *client) GetStateItem(ctx context.Context, stateRootHash, key string, path []string) (StateGetItemResult, error) {
	var result StateGetItemResult
	return result, c.processRequest(ctx, MethodGetStateItem, ParamStateRootHash{
		StateRootHash: stateRootHash,
		Key:           key,
		Path:          path,
	}, &result)
}

func (c *client) QueryGlobalStateByBlockHash(ctx context.Context, blockHash, key string, path []string) (QueryGlobalStateResult, error) {
	var result QueryGlobalStateResult
	return result, c.processRequest(ctx, MethodQueryGlobalState, NewQueryGlobalStateParam(key, path, ParamQueryGlobalStateID{
		BlockHash: blockHash,
	}), &result)
}

func (c *client) QueryGlobalStateByStateHash(ctx context.Context, stateRootHash, key string, path []string) (QueryGlobalStateResult, error) {
	var result QueryGlobalStateResult
	return result, c.processRequest(ctx, MethodQueryGlobalState, NewQueryGlobalStateParam(key, path, ParamQueryGlobalStateID{
		StateRootHash: stateRootHash,
	}), &result)
}

func (c *client) GetDictionaryItem(ctx context.Context, stateRootHash, uref, key string) (StateGetDictionaryResult, error) {
	var result StateGetDictionaryResult
	return result, c.processRequest(ctx, MethodGetDictionaryItem,
		NewParamStateDictionaryItem(stateRootHash, uref, key), &result)
}

func (c *client) GetAccountBalance(ctx context.Context, stateRootHash, purseURef string) (StateGetBalanceResult, error) {
	var result StateGetBalanceResult
	return result, c.processRequest(ctx, MethodGetStateBalance, map[string]string{
		"state_root_hash": stateRootHash,
		"purse_uref":      purseURef,
	}, &result)
}

func (c *client) GetEraInfoLatest(ctx context.Context) (ChainGetEraInfoResult, error) {
	var result ChainGetEraInfoResult
	return result, c.processRequest(ctx, MethodGetEraInfo, nil, &result)
}

func (c *client) GetEraInfoByBlockHeight(ctx context.Context, height uint64) (ChainGetEraInfoResult, error) {
	var result ChainGetEraInfoResult
	return result, c.processRequest(ctx, MethodGetEraInfo, NewParamBlockByHeight(height), &result)
}

func (c *client) GetEraInfoByBlockHash(ctx context.Context, hash string) (ChainGetEraInfoResult, error) {
	var result ChainGetEraInfoResult
	return result, c.processRequest(ctx, MethodGetEraInfo, NewParamBlockByHash(hash), &result)
}

func (c *client) GetBlockLatest(ctx context.Context) (ChainGetBlockResult, error) {
	var result ChainGetBlockResult
	return result, c.processRequest(ctx, MethodGetBlock, nil, &result)
}

func (c *client) GetBlockByHash(ctx context.Context, hash string) (ChainGetBlockResult, error) {
	var result ChainGetBlockResult
	return result, c.processRequest(ctx, MethodGetBlock, NewParamBlockByHash(hash), &result)
}

func (c *client) GetBlockByHeight(ctx context.Context, height uint64) (ChainGetBlockResult, error) {
	var result ChainGetBlockResult
	return result, c.processRequest(ctx, MethodGetBlock, NewParamBlockByHeight(height), &result)
}

func (c *client) GetBlockTransfersLatest(ctx context.Context) (ChainGetBlockTransfersResult, error) {
	var result ChainGetBlockTransfersResult
	return result, c.processRequest(ctx, MethodGetBlockTransfers, nil, &result)
}

func (c *client) GetBlockTransfersByHash(ctx context.Context, blockHash string) (ChainGetBlockTransfersResult, error) {
	var result ChainGetBlockTransfersResult
	return result, c.processRequest(ctx, MethodGetBlockTransfers, NewParamBlockByHash(blockHash), &result)
}

func (c *client) GetBlockTransfersByHeight(ctx context.Context, height uint64) (ChainGetBlockTransfersResult, error) {
	var result ChainGetBlockTransfersResult
	return result, c.processRequest(ctx, MethodGetBlockTransfers, NewParamBlockByHeight(height), &result)
}

func (c *client) GetEraSummaryLatest(ctx context.Context) (ChainGetEraSummaryResult, error) {
	var result ChainGetEraSummaryResult
	return result, c.processRequest(ctx, MethodGetEraSummary, nil, &result)
}

func (c *client) GetEraSummaryByHash(ctx context.Context, blockHash string) (ChainGetEraSummaryResult, error) {
	var result ChainGetEraSummaryResult
	return result, c.processRequest(ctx, MethodGetEraSummary, NewParamBlockByHash(blockHash), &result)
}

func (c *client) GetEraSummaryByHeight(ctx context.Context, height uint64) (ChainGetEraSummaryResult, error) {
	var result ChainGetEraSummaryResult
	return result, c.processRequest(ctx, MethodGetEraSummary, NewParamBlockByHeight(height), &result)
}

func (c *client) GetAuctionInfoLatest(ctx context.Context) (StateGetAuctionInfoResult, error) {
	var result StateGetAuctionInfoResult
	return result, c.processRequest(ctx, MethodGetAuctionInfo, nil, &result)
}

func (c *client) GetAuctionInfoByHash(ctx context.Context, blockHash string) (StateGetAuctionInfoResult, error) {
	var result StateGetAuctionInfoResult
	return result, c.processRequest(ctx, MethodGetAuctionInfo, NewParamBlockByHash(blockHash), &result)
}

func (c *client) GetAuctionInfoByHeight(ctx context.Context, height uint64) (StateGetAuctionInfoResult, error) {
	var result StateGetAuctionInfoResult
	return result, c.processRequest(ctx, MethodGetAuctionInfo, NewParamBlockByHeight(height), &result)
}

func (c *client) GetStateRootHashLatest(ctx context.Context) (ChainGetStateRootHashResult, error) {
	var result ChainGetStateRootHashResult
	return result, c.processRequest(ctx, MethodGetStateRootHash, nil, &result)
}

func (c *client) GetStateRootHashByHash(ctx context.Context, blockHash string) (ChainGetStateRootHashResult, error) {
	var result ChainGetStateRootHashResult
	return result, c.processRequest(ctx, MethodGetStateRootHash, NewParamBlockByHash(blockHash), &result)
}

func (c *client) GetStateRootHashByHeight(ctx context.Context, height uint64) (ChainGetStateRootHashResult, error) {
	var result ChainGetStateRootHashResult
	return result, c.processRequest(ctx, MethodGetStateRootHash, NewParamBlockByHeight(height), &result)
}

func (c *client) GetStatus(ctx context.Context) (InfoGetStatusResult, error) {
	var result InfoGetStatusResult
	return result, c.processRequest(ctx, MethodGetStatus, nil, &result)
}

func (c *client) GetPeers(ctx context.Context) (InfoGetPeerResult, error) {
	var result InfoGetPeerResult
	return result, c.processRequest(ctx, MethodGetPeers, nil, &result)
}

func (c *client) PutDeploy(ctx context.Context, deploy types.Deploy) (PutDeployResult, error) {
	var result PutDeployResult
	return result, c.processRequest(ctx, MethodPutDeploy, PutDeployRequest{Deploy: deploy}, &result)
}

func (c *client) processRequest(ctx context.Context, method Method, params interface{}, result any) error {
	request := DefaultRpcRequest(method, params)
	if reqID := GetReqIdCtx(ctx); reqID != 0 {
		request.ID = reqID
	}
	resp, err := c.handler.ProcessCall(ctx, request)
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return fmt.Errorf("rpc call failed, details: %w", resp.Error)
	}

	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return fmt.Errorf("%w, details: %s", ErrResultUnmarshal, err.Error())
	}

	return nil
}
