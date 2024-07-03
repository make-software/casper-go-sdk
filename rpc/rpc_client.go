package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/keypair"
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
	resp, err := c.processRequest(ctx, MethodGetDeploy, map[string]string{
		"deploy_hash": hash,
	}, &result)
	if err != nil {
		return InfoGetDeployResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetDeployFinalizedApproval(ctx context.Context, hash string) (InfoGetDeployResult, error) {
	var result InfoGetDeployResult

	resp, err := c.processRequest(ctx, MethodGetDeploy, map[string]interface{}{
		"deploy_hash":         hash,
		"finalized_approvals": true,
	}, &result)
	if err != nil {
		return InfoGetDeployResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetStateItem(ctx context.Context, stateRootHash *string, key string, path []string) (StateGetItemResult, error) {
	if stateRootHash == nil {
		latestHashResult, err := c.GetStateRootHashLatest(ctx)
		if err != nil {
			return StateGetItemResult{}, err
		}
		latestHashString := latestHashResult.StateRootHash.String()
		stateRootHash = &latestHashString
	}

	var result StateGetItemResult
	resp, err := c.processRequest(ctx, MethodGetStateItem, ParamStateRootHash{
		StateRootHash: *stateRootHash,
		Key:           key,
		Path:          path,
	}, &result)
	if err != nil {
		return StateGetItemResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) QueryGlobalStateByBlockHash(ctx context.Context, blockHash, key string, path []string) (QueryGlobalStateResult, error) {
	var result QueryGlobalStateResult
	resp, err := c.processRequest(ctx, MethodQueryGlobalState, NewQueryGlobalStateParam(key, path, &ParamQueryGlobalStateID{
		BlockHash: blockHash,
	}), &result)
	if err != nil {
		return QueryGlobalStateResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) QueryGlobalStateByBlockHeight(ctx context.Context, blockHeight uint64, key string, path []string) (QueryGlobalStateResult, error) {
	var result QueryGlobalStateResult
	resp, err := c.processRequest(ctx, MethodQueryGlobalState, NewQueryGlobalStateParam(key, path, &ParamQueryGlobalStateID{
		BlockHeight: &blockHeight,
	}), &result)
	if err != nil {
		return QueryGlobalStateResult{}, nil
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) QueryGlobalStateByStateHash(ctx context.Context, stateRootHash *string, key string, path []string) (QueryGlobalStateResult, error) {
	var result QueryGlobalStateResult
	if stateRootHash == nil {
		resp, err := c.processRequest(ctx, MethodQueryGlobalState, NewQueryGlobalStateParam(key, path, nil), &result)
		if err != nil {
			return QueryGlobalStateResult{}, err
		}

		result.rawJSON = resp.Result
		return result, nil
	}

	resp, err := c.processRequest(ctx, MethodQueryGlobalState, NewQueryGlobalStateParam(key, path, &ParamQueryGlobalStateID{
		StateRootHash: *stateRootHash,
	}), &result)
	if err != nil {
		return QueryGlobalStateResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetAccountInfoByBlochHash(ctx context.Context, blockHash string, pub keypair.PublicKey) (StateGetAccountInfo, error) {
	var result StateGetAccountInfo

	resp, err := c.processRequest(ctx, MethodGetStateAccount, ParamGetAccountInfoBalance{AccountIdentifier: pub.String(), ParamBlockIdentifier: NewParamBlockByHash(blockHash)}, &result)
	if err != nil {
		return StateGetAccountInfo{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetAccountInfoByBlochHeight(ctx context.Context, blockHeight uint64, pub keypair.PublicKey) (StateGetAccountInfo, error) {
	var result StateGetAccountInfo
	resp, err := c.processRequest(ctx, MethodGetStateAccount, ParamGetAccountInfoBalance{AccountIdentifier: pub.String(), ParamBlockIdentifier: NewParamBlockByHeight(blockHeight)}, &result)
	if err != nil {
		return StateGetAccountInfo{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetAccountInfo(ctx context.Context, blockIdentifier *ParamBlockIdentifier, accountIdentifier AccountIdentifier) (StateGetAccountInfo, error) {
	if blockIdentifier == nil {
		blockIdentifier = &ParamBlockIdentifier{}
	}
	var accountParam string
	if accountIdentifier.AccountHash != nil {
		accountParam = accountIdentifier.AccountHash.ToPrefixedString()
	} else if accountIdentifier.PublicKey != nil {
		accountParam = accountIdentifier.PublicKey.String()
	} else {
		return StateGetAccountInfo{}, fmt.Errorf("account identifier is empty")
	}

	var result StateGetAccountInfo
	resp, err := c.processRequest(ctx, MethodGetStateAccount, ParamGetAccountInfoBalance{AccountIdentifier: accountParam, ParamBlockIdentifier: *blockIdentifier}, &result)
	if err != nil {
		return StateGetAccountInfo{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetDictionaryItem(ctx context.Context, stateRootHash *string, uref, key string) (StateGetDictionaryResult, error) {
	return c.GetDictionaryItemByIdentifier(ctx, stateRootHash, ParamDictionaryIdentifier{
		URef: &ParamDictionaryIdentifierURef{
			DictionaryItemKey: key,
			SeedUref:          uref,
		},
	})
}

func (c *client) GetDictionaryItemByIdentifier(ctx context.Context, stateRootHash *string, identifier ParamDictionaryIdentifier) (StateGetDictionaryResult, error) {
	if stateRootHash == nil {
		latestHashResult, err := c.GetStateRootHashLatest(ctx)
		if err != nil {
			return StateGetDictionaryResult{}, err
		}
		latestHashString := latestHashResult.StateRootHash.String()
		stateRootHash = &latestHashString
	}
	var result StateGetDictionaryResult
	resp, err := c.processRequest(ctx, MethodGetDictionaryItem, map[string]interface{}{
		"state_root_hash":       *stateRootHash,
		"dictionary_identifier": identifier,
	}, &result)
	if err != nil {
		return StateGetDictionaryResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetBalance(ctx context.Context, purseURef string, stateRootHash *string) (StateGetBalanceResult, error) {
	if stateRootHash == nil {
		latestHashResult, err := c.GetStateRootHashLatest(ctx)
		if err != nil {
			return StateGetBalanceResult{}, err
		}
		latestHashString := latestHashResult.StateRootHash.String()
		stateRootHash = &latestHashString
	}
	var result StateGetBalanceResult
	resp, err := c.processRequest(ctx, MethodGetStateBalance, map[string]string{
		"state_root_hash": *stateRootHash,
		"purse_uref":      purseURef,
	}, &result)
	if err != nil {
		return StateGetBalanceResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetEraInfoLatest(ctx context.Context) (ChainGetEraInfoResult, error) {
	var result ChainGetEraInfoResult

	resp, err := c.processRequest(ctx, MethodGetEraInfo, nil, &result)
	if err != nil {
		return ChainGetEraInfoResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetEraInfoByBlockHeight(ctx context.Context, height uint64) (ChainGetEraInfoResult, error) {
	var result ChainGetEraInfoResult

	resp, err := c.processRequest(ctx, MethodGetEraInfo, NewParamBlockByHeight(height), &result)
	if err != nil {
		return ChainGetEraInfoResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetEraInfoByBlockHash(ctx context.Context, hash string) (ChainGetEraInfoResult, error) {
	var result ChainGetEraInfoResult

	resp, err := c.processRequest(ctx, MethodGetEraInfo, NewParamBlockByHash(hash), &result)
	if err != nil {
		return ChainGetEraInfoResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetBlockLatest(ctx context.Context) (ChainGetBlockResult, error) {
	var result chainGetBlockResultV1Compatible

	resp, err := c.processRequest(ctx, MethodGetBlock, nil, &result)
	if err != nil {
		return ChainGetBlockResult{}, err
	}

	blockResult, err := newChainGetBlockResultFromV1Compatible(result, resp.Result)
	if err != nil {
		return ChainGetBlockResult{}, err
	}

	blockResult.rawJSON = resp.Result
	return blockResult, nil
}

func (c *client) GetBlockByHash(ctx context.Context, hash string) (ChainGetBlockResult, error) {
	var result chainGetBlockResultV1Compatible

	resp, err := c.processRequest(ctx, MethodGetBlock, NewParamBlockByHash(hash), &result)
	if err != nil {
		return ChainGetBlockResult{}, err
	}

	blockResult, err := newChainGetBlockResultFromV1Compatible(result, resp.Result)
	if err != nil {
		return ChainGetBlockResult{}, err
	}

	blockResult.rawJSON = resp.Result
	return blockResult, nil
}

func (c *client) GetBlockByHeight(ctx context.Context, height uint64) (ChainGetBlockResult, error) {
	var result chainGetBlockResultV1Compatible

	resp, err := c.processRequest(ctx, MethodGetBlock, NewParamBlockByHeight(height), &result)
	if err != nil {
		return ChainGetBlockResult{}, err
	}

	blockResult, err := newChainGetBlockResultFromV1Compatible(result, resp.Result)
	if err != nil {
		return ChainGetBlockResult{}, err
	}

	blockResult.rawJSON = resp.Result
	return blockResult, nil
}

func (c *client) GetBlockTransfersLatest(ctx context.Context) (ChainGetBlockTransfersResult, error) {
	var result ChainGetBlockTransfersResult

	resp, err := c.processRequest(ctx, MethodGetBlockTransfers, nil, &result)
	if err != nil {
		return ChainGetBlockTransfersResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetBlockTransfersByHash(ctx context.Context, blockHash string) (ChainGetBlockTransfersResult, error) {
	var result ChainGetBlockTransfersResult

	resp, err := c.processRequest(ctx, MethodGetBlockTransfers, NewParamBlockByHash(blockHash), &result)
	if err != nil {
		return ChainGetBlockTransfersResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetBlockTransfersByHeight(ctx context.Context, height uint64) (ChainGetBlockTransfersResult, error) {
	var result ChainGetBlockTransfersResult

	resp, err := c.processRequest(ctx, MethodGetBlockTransfers, NewParamBlockByHeight(height), &result)
	if err != nil {
		return ChainGetBlockTransfersResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetEraSummaryLatest(ctx context.Context) (ChainGetEraSummaryResult, error) {
	var result ChainGetEraSummaryResult

	resp, err := c.processRequest(ctx, MethodGetEraSummary, nil, &result)
	if err != nil {
		return ChainGetEraSummaryResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetEraSummaryByHash(ctx context.Context, blockHash string) (ChainGetEraSummaryResult, error) {
	var result ChainGetEraSummaryResult

	resp, err := c.processRequest(ctx, MethodGetEraSummary, NewParamBlockByHash(blockHash), &result)
	if err != nil {
		return ChainGetEraSummaryResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetEraSummaryByHeight(ctx context.Context, height uint64) (ChainGetEraSummaryResult, error) {
	var result ChainGetEraSummaryResult

	resp, err := c.processRequest(ctx, MethodGetEraSummary, NewParamBlockByHeight(height), &result)
	if err != nil {
		return ChainGetEraSummaryResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetAuctionInfoLatest(ctx context.Context) (StateGetAuctionInfoResult, error) {
	var result StateGetAuctionInfoResult

	resp, err := c.processRequest(ctx, MethodGetAuctionInfo, nil, &result)
	if err != nil {
		return StateGetAuctionInfoResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetAuctionInfoByHash(ctx context.Context, blockHash string) (StateGetAuctionInfoResult, error) {
	var result StateGetAuctionInfoResult
	resp, err := c.processRequest(ctx, MethodGetAuctionInfo, NewParamBlockByHash(blockHash), &result)
	if err != nil {
		return StateGetAuctionInfoResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetAuctionInfoByHeight(ctx context.Context, height uint64) (StateGetAuctionInfoResult, error) {
	var result StateGetAuctionInfoResult

	resp, err := c.processRequest(ctx, MethodGetAuctionInfo, NewParamBlockByHeight(height), &result)
	if err != nil {
		return StateGetAuctionInfoResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetStateRootHashLatest(ctx context.Context) (ChainGetStateRootHashResult, error) {
	var result ChainGetStateRootHashResult

	resp, err := c.processRequest(ctx, MethodGetStateRootHash, nil, &result)
	if err != nil {
		return ChainGetStateRootHashResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetStateRootHashByHash(ctx context.Context, blockHash string) (ChainGetStateRootHashResult, error) {
	var result ChainGetStateRootHashResult

	resp, err := c.processRequest(ctx, MethodGetStateRootHash, NewParamBlockByHash(blockHash), &result)
	if err != nil {
		return ChainGetStateRootHashResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetStateRootHashByHeight(ctx context.Context, height uint64) (ChainGetStateRootHashResult, error) {
	var result ChainGetStateRootHashResult
	resp, err := c.processRequest(ctx, MethodGetStateRootHash, NewParamBlockByHeight(height), &result)
	if err != nil {
		return ChainGetStateRootHashResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetValidatorChangesInfo(ctx context.Context) (InfoGetValidatorChangesResult, error) {
	var result InfoGetValidatorChangesResult
	resp, err := c.processRequest(ctx, MethodGetValidatorChanges, nil, &result)
	if err != nil {
		return InfoGetValidatorChangesResult{}, nil
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetStatus(ctx context.Context) (InfoGetStatusResult, error) {
	var result InfoGetStatusResult

	resp, err := c.processRequest(ctx, MethodGetStatus, nil, &result)
	if err != nil {
		return InfoGetStatusResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetPeers(ctx context.Context) (InfoGetPeerResult, error) {
	var result InfoGetPeerResult

	resp, err := c.processRequest(ctx, MethodGetPeers, nil, &result)
	if err != nil {
		return InfoGetPeerResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) PutDeploy(ctx context.Context, deploy types.Deploy) (PutDeployResult, error) {
	var result PutDeployResult

	resp, err := c.processRequest(ctx, MethodPutDeploy, PutDeployRequest{Deploy: deploy}, &result)
	if err != nil {
		return PutDeployResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) QueryLatestBalance(ctx context.Context, identifier PurseIdentifier) (QueryBalanceResult, error) {
	var result QueryBalanceResult

	resp, err := c.processRequest(ctx, MethodQueryBalance, QueryBalanceRequest{PurseIdentifier: identifier}, &result)
	if err != nil {
		return QueryBalanceResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) QueryBalanceByBlockHeight(ctx context.Context, purseIdentifier PurseIdentifier, height uint64) (QueryBalanceResult, error) {
	var result QueryBalanceResult

	resp, err := c.processRequest(ctx, MethodQueryBalance, QueryBalanceRequest{PurseIdentifier: purseIdentifier, StateIdentifier: &GlobalStateIdentifier{
		BlockHeight: &height,
	}}, &result)
	if err != nil {
		return QueryBalanceResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) QueryBalanceByBlockHash(ctx context.Context, purseIdentifier PurseIdentifier, blockHash string) (QueryBalanceResult, error) {
	var result QueryBalanceResult

	resp, err := c.processRequest(ctx, MethodQueryBalance, QueryBalanceRequest{PurseIdentifier: purseIdentifier, StateIdentifier: &GlobalStateIdentifier{
		BlockHash: &blockHash,
	}}, &result)
	if err != nil {
		return QueryBalanceResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) QueryBalanceByStateRootHash(ctx context.Context, purseIdentifier PurseIdentifier, stateRootHash string) (QueryBalanceResult, error) {
	var result QueryBalanceResult

	resp, err := c.processRequest(ctx, MethodQueryBalance, QueryBalanceRequest{PurseIdentifier: purseIdentifier, StateIdentifier: &GlobalStateIdentifier{
		StateRoot: &stateRootHash,
	}}, &result)
	if err != nil {
		return QueryBalanceResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) QueryLatestBalanceDetails(ctx context.Context, purseIdentifier PurseIdentifier) (QueryBalanceDetailsResult, error) {
	var result QueryBalanceDetailsResult

	resp, err := c.processRequest(ctx, MethodQueryBalanceDetails, QueryBalanceDetailsRequest{PurseIdentifier: purseIdentifier}, &result)
	if err != nil {
		return QueryBalanceDetailsResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) QueryBalanceDetailsByStateRootHash(ctx context.Context, purseIdentifier PurseIdentifier, stateRootHash string) (QueryBalanceDetailsResult, error) {
	var result QueryBalanceDetailsResult

	resp, err := c.processRequest(ctx, MethodQueryBalanceDetails, QueryBalanceDetailsRequest{purseIdentifier, &GlobalStateIdentifier{
		StateRoot: &stateRootHash,
	}}, &result)
	if err != nil {
		return QueryBalanceDetailsResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) QueryBalanceDetailsByBlockHeight(ctx context.Context, purseIdentifier PurseIdentifier, height uint64) (QueryBalanceDetailsResult, error) {
	var result QueryBalanceDetailsResult

	resp, err := c.processRequest(ctx, MethodQueryBalanceDetails, QueryBalanceDetailsRequest{purseIdentifier, &GlobalStateIdentifier{
		BlockHeight: &height,
	}}, &result)
	if err != nil {
		return QueryBalanceDetailsResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) QueryBalanceDetailsByBlockHash(ctx context.Context, purseIdentifier PurseIdentifier, blockHash string) (QueryBalanceDetailsResult, error) {
	var result QueryBalanceDetailsResult

	resp, err := c.processRequest(ctx, MethodQueryBalanceDetails, QueryBalanceDetailsRequest{purseIdentifier, &GlobalStateIdentifier{
		BlockHash: &blockHash,
	}}, &result)
	if err != nil {
		return QueryBalanceDetailsResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) GetChainspec(ctx context.Context) (InfoGetChainspecResult, error) {
	var result InfoGetChainspecResult

	resp, err := c.processRequest(ctx, MethodInfoGetChainspec, nil, &result)
	if err != nil {
		return InfoGetChainspecResult{}, err
	}

	result.rawJSON = resp.Result
	return result, nil
}

func (c *client) processRequest(ctx context.Context, method Method, params interface{}, result any) (RpcResponse, error) {
	request := DefaultRpcRequest(method, params)
	if reqID := GetReqIdCtx(ctx); reqID != "0" {
		request.ID = NewIDFromString(reqID)
	}
	resp, err := c.handler.ProcessCall(ctx, request)
	if err != nil {
		return resp, err
	}

	if resp.Error != nil {
		return resp, fmt.Errorf("rpc call failed, details: %w", resp.Error)
	}

	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return resp, fmt.Errorf("%w, details: %s", ErrResultUnmarshal, err.Error())
	}

	return resp, nil
}
