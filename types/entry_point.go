package types

import (
	"encoding/json"
)

// EntryPointType defines whether the code runs in the contract's or the session context.
type EntryPointType string

const (
	EntryPointTypeSession  EntryPointType = "Session"
	EntryPointTypeContract EntryPointType = "Contract"
)

// EntryPoint is a type signature of a method.
// Order of arguments matter since can be referenced by index as well as name.
type EntryPoint struct {
	// Access control options for a contract entry point
	Access EntryPointAccess `json:"access"`
	// List of input parameters to the method.
	// Order of arguments matter since can be referenced by index as well as name.
	Args []EntryPointArg `json:"args"`
	// Context of method execution
	EntryPointType EntryPointType `json:"entry_point_type"`
	// Name of the entry point
	Name string `json:"name"`
	// Returned type
	Ret CLTypeRaw `json:"ret"`
}

type EntryPointArg struct {
	Name   string    `json:"name"`
	CLType CLTypeRaw `json:"cl_type"`
}

// EntryPointAccess is access control options for a contract entry point (method).
// TODO: to implement
type EntryPointAccess struct {
	json.RawMessage
}
