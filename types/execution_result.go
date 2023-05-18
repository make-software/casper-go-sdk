package types

import "github.com/make-software/casper-go-sdk/types/key"

// DeployExecutionResult represents the result of executing a single deploy.
type DeployExecutionResult struct {
	BlockHash key.Hash              `json:"block_hash"`
	Result    ExecutionResultStatus `json:"result"`
}

type ExecutionResultStatus struct {
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
