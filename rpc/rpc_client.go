package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/key"
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
	return result, c.processRequest(ctx, MethodGetDeploy, map[string]string{
		"deploy_hash": hash,
	}, &result)
}

func (c *client) GetDeployFinalizedApproval(ctx context.Context, hash string) (InfoGetDeployResult, error) {
	var result InfoGetDeployResult
	return result, c.processRequest(ctx, MethodGetDeploy, map[string]interface{}{
		"deploy_hash":         hash,
		"finalized_approvals": true,
	}, &result)
}

func (c *client) GetTransaction(ctx context.Context, transactionHash string) (InfoGetTransactionResult, error) {
	hash, err := key.NewHash(transactionHash)
	if err != nil {
		return InfoGetTransactionResult{}, err
	}

	var result infoGetTransactionResultV1Compatible
	c.processRequest(ctx, MethodGetTransaction, ParamTransactionHash{
		TransactionHash: types.TransactionHash{
			TransactionV1Hash: &hash,
		},
	}, &result)
	if err != nil {
		return InfoGetTransactionResult{}, err
	}

	return newInfoGetTransactionResultFromV1Compatible(result)
}

func (c *client) GetTransactionFinalizedApproval(ctx context.Context, transactionHash string) (InfoGetTransactionResult, error) {
	hash, err := key.NewHash(transactionHash)
	if err != nil {
		return InfoGetTransactionResult{}, err
	}

	var result infoGetTransactionResultV1Compatible
	c.processRequest(ctx, MethodGetTransaction, ParamTransactionHash{
		TransactionHash: types.TransactionHash{
			TransactionV1Hash: &hash,
		},
		FinalizedApprovals: &[]bool{true}[0],
	}, &result)
	if err != nil {
		return InfoGetTransactionResult{}, err
	}

	return newInfoGetTransactionResultFromV1Compatible(result)
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
	return result, c.processRequest(ctx, MethodGetStateItem, ParamStateRootHash{
		StateRootHash: *stateRootHash,
		Key:           key,
		Path:          path,
	}, &result)
}

func (c *client) QueryGlobalStateByBlockHash(ctx context.Context, blockHash, key string, path []string) (QueryGlobalStateResult, error) {
	var result QueryGlobalStateResult
	return result, c.processRequest(ctx, MethodQueryGlobalState, NewQueryGlobalStateParam(key, path, &ParamQueryGlobalStateID{
		BlockHash: blockHash,
	}), &result)
}

func (c *client) QueryGlobalStateByBlockHeight(ctx context.Context, blockHeight uint64, key string, path []string) (QueryGlobalStateResult, error) {
	var result QueryGlobalStateResult
	return result, c.processRequest(ctx, MethodQueryGlobalState, NewQueryGlobalStateParam(key, path, &ParamQueryGlobalStateID{
		BlockHeight: &blockHeight,
	}), &result)
}

func (c *client) QueryGlobalStateByStateHash(ctx context.Context, stateRootHash *string, key string, path []string) (QueryGlobalStateResult, error) {
	var result QueryGlobalStateResult
	if stateRootHash == nil {
		return result, c.processRequest(ctx, MethodQueryGlobalState, NewQueryGlobalStateParam(key, path, nil), &result)
	}
	return result, c.processRequest(ctx, MethodQueryGlobalState, NewQueryGlobalStateParam(key, path, &ParamQueryGlobalStateID{
		StateRootHash: *stateRootHash,
	}), &result)
}

func (c *client) GetAccountInfoByBlochHash(ctx context.Context, blockHash string, pub keypair.PublicKey) (StateGetAccountInfo, error) {
	var result StateGetAccountInfo
	return result, c.processRequest(ctx, MethodGetStateAccount, ParamGetAccountInfoBalance{AccountIdentifier: pub.String(), ParamBlockIdentifier: NewParamBlockByHash(blockHash)}, &result)
}

func (c *client) GetAccountInfoByBlochHeight(ctx context.Context, blockHeight uint64, pub keypair.PublicKey) (StateGetAccountInfo, error) {
	var result StateGetAccountInfo
	return result, c.processRequest(ctx, MethodGetStateAccount, ParamGetAccountInfoBalance{AccountIdentifier: pub.String(), ParamBlockIdentifier: NewParamBlockByHeight(blockHeight)}, &result)
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
	return result, c.processRequest(ctx, MethodGetStateAccount, ParamGetAccountInfoBalance{AccountIdentifier: accountParam, ParamBlockIdentifier: *blockIdentifier}, &result)
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
	return result, c.processRequest(ctx, MethodGetDictionaryItem, map[string]interface{}{
		"state_root_hash":       *stateRootHash,
		"dictionary_identifier": identifier,
	}, &result)
}

func (c *client) GetAccountBalance(ctx context.Context, stateRootHash *string, purseURef string) (StateGetBalanceResult, error) {
	if stateRootHash == nil {
		latestHashResult, err := c.GetStateRootHashLatest(ctx)
		if err != nil {
			return StateGetBalanceResult{}, err
		}
		latestHashString := latestHashResult.StateRootHash.String()
		stateRootHash = &latestHashString
	}
	var result StateGetBalanceResult
	return result, c.processRequest(ctx, MethodGetStateBalance, map[string]string{
		"state_root_hash": *stateRootHash,
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
	var result chainGetBlockResultV1Compatible
	if err := c.processRequest(ctx, MethodGetBlock, nil, &result); err != nil {
		return ChainGetBlockResult{}, err
	}

	return newChainGetBlockResultFromV1Compatible(result)
}

func (c *client) GetBlockByHash(ctx context.Context, hash string) (ChainGetBlockResult, error) {
	var result chainGetBlockResultV1Compatible
	if err := c.processRequest(ctx, MethodGetBlock, NewParamBlockByHash(hash), &result); err != nil {
		return ChainGetBlockResult{}, err
	}

	return newChainGetBlockResultFromV1Compatible(result)
}

func (c *client) GetBlockByHeight(ctx context.Context, height uint64) (ChainGetBlockResult, error) {
	var result chainGetBlockResultV1Compatible
	if err := c.processRequest(ctx, MethodGetBlock, NewParamBlockByHeight(height), &result); err != nil {
		return ChainGetBlockResult{}, err
	}

	return newChainGetBlockResultFromV1Compatible(result)
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

func (c *client) GetValidatorChangesInfo(ctx context.Context) (InfoGetValidatorChangesResult, error) {
	var result InfoGetValidatorChangesResult
	return result, c.processRequest(ctx, MethodGetValidatorChanges, nil, &result)
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

func (c *client) QueryBalance(ctx context.Context, identifier PurseIdentifier) (QueryBalanceResult, error) {
	var result QueryBalanceResult
	return result, c.processRequest(ctx, MethodQueryBalance, QueryBalanceRequest{identifier}, &result)
}

func (c *client) GetChainspec(ctx context.Context) (InfoGetChainspecResult, error) {
	var result InfoGetChainspecResult
	return result, c.processRequest(ctx, MethodInfoGetChainspec, nil, &result)
}

func (c *client) processRequest(ctx context.Context, method Method, params interface{}, result any) error {
	request := DefaultRpcRequest(method, params)
	if reqID := GetReqIdCtx(ctx); reqID != "0" {
		request.ID = NewIDFromString(reqID)
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
