package helper

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/make-software/casper-go-sdk/v2/rpc"
)

type clientRetriesDecorator struct {
	handler       rpc.Handler
	retiesCount   int
	retriesPeriod time.Duration
}

func NewTestRetriesDecorator(
	handler rpc.Handler,
	retiesCount int,
	retriesSeconds time.Duration,
) rpc.Handler {
	return &clientRetriesDecorator{
		handler:       handler,
		retiesCount:   retiesCount,
		retriesPeriod: retriesSeconds * time.Second,
	}
}

func (c *clientRetriesDecorator) ProcessCall(ctx context.Context, params rpc.RpcRequest) (rpc.RpcResponse, error) {
	result, err := c.handler.ProcessCall(ctx, params)
	if err == nil {
		return result, nil
	}

	triesCount := 1
	ticker := time.NewTicker(c.retriesPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Context was closed")
			return result, err
		case <-ticker.C:
			result, err := c.handler.ProcessCall(ctx, params)
			var httpErr *rpc.HttpError
			if errors.As(err, &httpErr) {
				if triesCount >= c.retiesCount {
					return rpc.RpcResponse{}, err
				}
				if httpErr.StatusCode >= http.StatusBadGateway && httpErr.StatusCode <= http.StatusGatewayTimeout {
					triesCount += 1
					continue
				}
			}
			if err != nil {
				return rpc.RpcResponse{}, err
			}
			return result, nil
		}
	}
}
