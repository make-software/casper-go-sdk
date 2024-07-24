package types

import (
	"github.com/make-software/casper-go-sdk/v2/types/key"
)

// Account representing a user's account, stored in a global state.
type Account struct {
	// AccountHash is an Account's identity key
	AccountHash key.AccountHash `json:"account_hash"`
	// TODO: Is it could be any type of keys or certain types?
	NamedKeys NamedKeys `json:"named_keys"`
	// Purse that can hold Casper tokens
	MainPurse key.URef `json:"main_purse"`
	// Set of public keys allowed to provide signatures on deploys for the account
	AssociatedKeys []AssociatedKey `json:"associated_keys"`
	// Thresholds that have to be met when executing an action of a certain type.
	ActionThresholds ActionThresholds `json:"action_thresholds"`
}

// AssociatedKey is allowed to provide signatures on deploys for the account
type AssociatedKey struct {
	// AccountHash is an Account's identity key
	AccountHash key.AccountHash `json:"account_hash"`
	// Weight of the associated key
	Weight uint64 `json:"weight"`
}

// ActionThresholds have to be met when executing an action of a certain type.
type ActionThresholds struct {
	// Threshold that has to be met for a deployment action
	Deployment uint64 `json:"deployment"`
	// Threshold that has to be met for a key management action
	KeyManagement uint64 `json:"key_management"`
}
