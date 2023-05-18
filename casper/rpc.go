package casper

import (
	"github.com/make-software/casper-go-sdk/rpc"
)

type (
	RPCClient                   = rpc.Client
	RpcHandler                  = rpc.Handler
	RpcInformationalClient      = rpc.ClientInformational
	InfoGetDeployResult         = rpc.InfoGetDeployResult
	ChainGetBlockResult         = rpc.ChainGetBlockResult
	ChainGetEraInfoResult       = rpc.ChainGetEraInfoResult
	StateGetAuctionInfoResult   = rpc.StateGetAuctionInfoResult
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
