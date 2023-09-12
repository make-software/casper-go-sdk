package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/make-software/casper-go-sdk/types"
)

type SpeculativeClient struct {
	handler Handler
}

func NewSpeculativeClient(handler Handler) *SpeculativeClient {
	return &SpeculativeClient{handler: handler}
}

func (c SpeculativeClient) SpeculativeExec(ctx context.Context, deploy types.Deploy, identifier *BlockIdentifier) (SpeculativeExecResult, error) {
	var result SpeculativeExecResult
	request := DefaultRpcRequest(MethodSpeculativeExec, SpeculativeExecParams{
		Deploy:          deploy,
		BlockIdentifier: identifier,
	})
	if reqID := GetReqIdCtx(ctx); reqID != "0" {
		request.ID = NewIDFromString(reqID)
	}
	resp, err := c.handler.ProcessCall(ctx, request)
	if err != nil {
		return SpeculativeExecResult{}, err
	}

	if resp.Error != nil {
		return SpeculativeExecResult{}, fmt.Errorf("rpc call failed, details: %w", resp.Error)
	}

	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return SpeculativeExecResult{}, fmt.Errorf("%w, details: %s", ErrResultUnmarshal, err.Error())
	}

	return result, nil
}
