package types

import (
	"encoding/json"
	"errors"
	"github.com/make-software/casper-go-sdk/types/key"
)

// DeployExecutionInfo represents the result of executing a single deploy V2
type DeployExecutionInfo struct {
	BlockHash       key.Hash        `json:"block_hash"`
	BlockHeight     uint64          `json:"block_height"`
	ExecutionResult ExecutionResult `json:"execution_result"`
}

func DeployExecutionInfoFromV1(results []DeployExecutionResult, height *uint64) DeployExecutionInfo {
	if len(results) == 0 {
		return DeployExecutionInfo{}
	}

	var blockHeight uint64
	if height != nil {
		blockHeight = *height
	}

	result := results[0]
	return DeployExecutionInfo{
		BlockHash:   result.BlockHash,
		BlockHeight: blockHeight,
		ExecutionResult: ExecutionResult{
			ExecutionResultV2:       NewExecutionResultV2FromV1(result.Result),
			OriginExecutionResultV1: &result.Result,
		},
	}
}

type ExecutionResult struct {
	ExecutionResultV2
	OriginExecutionResultV1 *ExecutionResultV1
}

func (v *ExecutionResult) UnmarshalJSON(data []byte) error {
	var versioned = struct {
		ExecutionResultV2 *ExecutionResultV2 `json:"Version2"`
		ExecutionResultV1 *ExecutionResultV1 `json:"Version1"`
	}{}
	if err := json.Unmarshal(data, &versioned); err != nil {
		return err
	}

	if versioned.ExecutionResultV2 != nil {
		*v = ExecutionResult{
			ExecutionResultV2: *versioned.ExecutionResultV2,
		}
		return nil
	}

	if versioned.ExecutionResultV1 != nil {
		*v = ExecutionResult{
			ExecutionResultV2:       NewExecutionResultV2FromV1(*versioned.ExecutionResultV1),
			OriginExecutionResultV1: versioned.ExecutionResultV1,
		}
		return nil
	}

	return errors.New("incorrect RPC response structure")
}

func NewExecutionResultV2FromV1(v1 ExecutionResultV1) ExecutionResultV2 {
	if v1.Success != nil {
		transforms := make([]TransformV2, 0)
		for _, transform := range v1.Success.Effect.Transforms {
			transforms = append(transforms, TransformV2{
				Key:  transform.Key,
				Kind: TransformKind(transform.Transform),
			})
		}

		transfers := make([]Transfer, 0)
		for _, transform := range v1.Success.Effect.Transforms {
			writeTransfer, err := transform.Transform.ParseAsWriteTransfer()
			if err != nil {
				continue
			}

			var toHash *key.AccountHash
			if writeTransfer.To != nil {
				toHash = writeTransfer.To
			}

			var id uint64
			if writeTransfer.ID != nil {
				id = *writeTransfer.ID
			}

			transfers = append(transfers, Transfer{
				TransferV2: TransferV2{
					Amount: writeTransfer.Amount,
					TransactionHash: TransactionHash{
						TransactionV1Hash: transform.Key.Hash,
					},
					From: InitiatorAddr{
						AccountHash: &writeTransfer.From,
					},
					Gas:    writeTransfer.Gas,
					ID:     id,
					Source: writeTransfer.Source,
					Target: writeTransfer.Target,
					To:     toHash,
				},
			})
		}
		return ExecutionResultV2{
			Cost:      v1.Success.Cost,
			Transfers: transfers,
			Effects:   transforms,
		}
	}

	return ExecutionResultV2{
		ErrorMessage: &v1.Failure.ErrorMessage,
		Cost:         v1.Failure.Cost,
	}
}

// ExecutionResultV2 represents the result of executing a single deploy for V2 version
type ExecutionResultV2 struct {
	Initiator    InitiatorAddr   `json:"initiator"`
	ErrorMessage *string         `json:"error_message"`
	Limit        uint64          `json:"limit,string"`
	Consumed     uint64          `json:"consumed,string"`
	Cost         uint64          `json:"cost,string"`
	Payment      json.RawMessage `json:"payment"`
	Transfers    []Transfer      `json:"transfers"`
	SizeEstimate uint64          `json:"size_estimate"`
	Effects      []TransformV2   `json:"effects"`
}

// DeployExecutionResult represents the result of executing a single deploy.
type DeployExecutionResult struct {
	BlockHash key.Hash          `json:"block_hash"`
	Result    ExecutionResultV1 `json:"result"`
}

type ExecutionResultV1 struct {
	Success *ExecutionResultStatusData `json:"Success,omitempty"`
	Failure *ExecutionResultStatusData `json:"Failure,omitempty"`
}

type ExecutionResultStatusData struct {
	Effect Effect `json:"effect"`
	// A record of `Transfers` performed while executing the `deploy`.
	Transfers    []key.TransferHash `json:"transfers"`
	Cost         uint64             `json:"cost,string"`
	ErrorMessage string             `json:"error_message"`
}

type Effect struct {
	Operations []Operation    `json:"operations"`
	Transforms []TransformKey `json:"transforms"`
}

type Operation struct {
	Key  key.Key `json:"key"`
	Kind string  `json:"kind"`
}
