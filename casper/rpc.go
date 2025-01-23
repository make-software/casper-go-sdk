package casper

import (
	"github.com/make-software/casper-go-sdk/v2/rpc"
)

type (
	RPCClient                   = rpc.Client
	RpcHandler                  = rpc.Handler
	RpcInformationalClient      = rpc.ClientInformational
	InfoGetDeployResult         = rpc.InfoGetDeployResult
	InfoGetTransactionResult    = rpc.InfoGetTransactionResult
	ChainGetBlockResult         = rpc.ChainGetBlockResult
	ChainGetEraInfoResult       = rpc.ChainGetEraInfoResult
	StateGetEntity              = rpc.StateGetEntityResult
	StateGetAuctionInfoV1Result = rpc.StateGetAuctionInfoV1Result
	StateGetAuctionInfoV2Result = rpc.StateGetAuctionInfoV2Result
	StateGetItemResult          = rpc.StateGetItemResult
	InfoGetStatusResult         = rpc.InfoGetStatusResult
	NodePeer                    = rpc.NodePeer
	ChainGetStateRootHashResult = rpc.ChainGetStateRootHashResult
	HttpError                   = rpc.HttpError
)

var (
	NewRPCClient  = rpc.NewClient
	NewRPCHandler = rpc.NewHttpHandler
)
