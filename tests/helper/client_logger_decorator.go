package helper

import (
	"context"
	"encoding/json"
	"log"

	"github.com/make-software/casper-go-sdk/rpc"
)

type clientLoggerDecorator struct {
	handler rpc.Handler
}

func NewTestLoggerDecorator(handler rpc.Handler) rpc.Handler {
	return &clientLoggerDecorator{handler: handler}
}

func (m clientLoggerDecorator) ProcessCall(ctx context.Context, params rpc.RpcRequest) (rpc.RpcResponse, error) {
	paramsStr, _ := json.Marshal(params)
	log.Print("Parsed query params: " + string(paramsStr))
	result, err := m.handler.ProcessCall(ctx, params)
	if err != nil {
		log.Print("Error was occurred: " + string(paramsStr))
		return rpc.RpcResponse{}, err
	}
	resultStr, _ := json.Marshal(result)
	log.Print("Parsed rpc response: " + string(resultStr))
	return result, err
}
