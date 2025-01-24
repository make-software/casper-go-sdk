package helper

import (
	"context"
	"log"
	"time"

	"github.com/make-software/casper-go-sdk/v2/rpc"
)

type rpcCallLoggerDecorator struct {
	handler rpc.Handler
}

func NewRPCLoggerDecorator(handler rpc.Handler) rpc.Handler {
	return &rpcCallLoggerDecorator{handler: handler}
}

func (m rpcCallLoggerDecorator) ProcessCall(ctx context.Context, params rpc.RpcRequest) (rpc.RpcResponse, error) {
	result, err := m.handler.ProcessCall(ctx, params)
	if err != nil {
		log.Printf("RPC call failed (method %s): \n, timestamp: %s", params.Method, time.Now().UTC())
		return rpc.RpcResponse{}, err
	}
	log.Printf("RPC call success (method %s): \n, timestamp: %s", params.Method, time.Now().UTC())
	return result, err
}
