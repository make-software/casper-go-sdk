package types

import "github.com/make-software/casper-go-sdk/types/key"

// DeployInfo provides information relating to the given Deploy.
type DeployInfo struct {
	// The `Deploy` hash.
	DeployHash key.Hash `json:"deploy_hash"`
	// `Account` identifier of the creator of the `Deploy`.
	From key.AccountHash `json:"from"`
	// `Gas` cost of executing the `Deploy`.
	Gas uint64 `json:"gas,string"`
	// `Source` purse used for payment of the `Deploy`.
	Source key.URef `json:"source"`
	// `Transfer` addresses performed by the `Deploy`.
	Transfers []key.TransferHash `json:"transfers"`
}
