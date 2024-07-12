package types

import (
	"encoding/json"
	"errors"

	"github.com/make-software/casper-go-sdk/types/key"
)

// ExecutionInfo represents the result of executing a single deploy V2
type ExecutionInfo struct {
	BlockHash       key.Hash        `json:"block_hash"`
	BlockHeight     uint64          `json:"block_height"`
	ExecutionResult ExecutionResult `json:"execution_result"`
}

func ExecutionInfoFromV1(results []DeployExecutionResult, height *uint64) ExecutionInfo {
	if len(results) == 0 {
		return ExecutionInfo{}
	}

	var blockHeight uint64
	if height != nil {
		blockHeight = *height
	}

	result := results[0]
	return ExecutionInfo{
		BlockHash:       result.BlockHash,
		BlockHeight:     blockHeight,
		ExecutionResult: NewExecutionResultFromV1(result.Result),
	}
}

type ExecutionResult struct {
	Initiator    InitiatorAddr   `json:"initiator"`
	ErrorMessage *string         `json:"error_message"`
	Limit        uint64          `json:"limit,string"`
	Consumed     uint64          `json:"consumed,string"`
	Cost         uint64          `json:"cost,string"`
	Payment      json.RawMessage `json:"payment"`
	Transfers    []Transfer      `json:"transfers"`
	SizeEstimate uint64          `json:"size_estimate"`
	Effects      []TransformV2   `json:"effects"`

	originExecutionResultV1 *ExecutionResultV1
	originExecutionResultV2 *ExecutionResultV2
}

func (v *ExecutionResult) GetExecutionResultV1() *ExecutionResultV1 {
	return v.originExecutionResultV1
}

func (v *ExecutionResult) GetExecutionResultV2() *ExecutionResultV2 {
	return v.originExecutionResultV2
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
			Initiator:               versioned.ExecutionResultV2.Initiator,
			ErrorMessage:            versioned.ExecutionResultV2.ErrorMessage,
			Limit:                   versioned.ExecutionResultV2.Limit,
			Consumed:                versioned.ExecutionResultV2.Consumed,
			Cost:                    versioned.ExecutionResultV2.Cost,
			Payment:                 versioned.ExecutionResultV2.Payment,
			Transfers:               versioned.ExecutionResultV2.Transfers,
			SizeEstimate:            versioned.ExecutionResultV2.SizeEstimate,
			Effects:                 versioned.ExecutionResultV2.Effects,
			originExecutionResultV2: versioned.ExecutionResultV2,
		}
		return nil
	}

	if versioned.ExecutionResultV1 != nil {
		*v = NewExecutionResultFromV1(*versioned.ExecutionResultV1)
		return nil
	}

	return errors.New("incorrect RPC response structure")
}

func NewExecutionResultFromV1(v1 ExecutionResultV1) ExecutionResult {
	transforms := make([]TransformV2, 0)
	if v1.Success != nil {
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
				Amount: writeTransfer.Amount,
				TransactionHash: TransactionHash{
					Transaction: transform.Key.Hash,
				},
				From: InitiatorAddr{
					AccountHash: &writeTransfer.From,
				},
				Gas:    writeTransfer.Gas,
				ID:     id,
				Source: writeTransfer.Source,
				Target: writeTransfer.Target,
				To:     toHash,
			})
		}
		return ExecutionResult{
			Limit:                   0, // limit is unknown field for V1 Deploy
			Consumed:                v1.Success.Cost,
			Cost:                    0, // cost is unknown field for V1 Deploy
			Payment:                 nil,
			Transfers:               transfers,
			Effects:                 transforms,
			originExecutionResultV1: &v1,
		}
	}

	if v1.Failure != nil {
		for _, transform := range v1.Failure.Effect.Transforms {
			transforms = append(transforms, TransformV2{
				Key: transform.Key,
				// TODO: we should convert old Transform to new format
				Kind: TransformKind(transform.Transform),
			})
		}
	}

	return ExecutionResult{
		ErrorMessage:            &v1.Failure.ErrorMessage,
		Consumed:                v1.Failure.Cost,
		Effects:                 transforms,
		originExecutionResultV1: &v1,
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

type ExecutionResultV1 struct {
	Success *ExecutionResultStatusData `json:"Success,omitempty"`
	Failure *ExecutionResultStatusData `json:"Failure,omitempty"`
}

// DeployExecutionResult represents the result of executing a single deploy.
type DeployExecutionResult struct {
	BlockHash key.Hash          `json:"block_hash"`
	Result    ExecutionResultV1 `json:"result"`
}

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
		BlockHash:       result.BlockHash,
		BlockHeight:     blockHeight,
		ExecutionResult: NewExecutionResultFromV1(result.Result),
	}
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
