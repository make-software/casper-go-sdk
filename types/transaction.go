package types

import (
	"encoding/json"
	"strconv"

	"github.com/make-software/casper-go-sdk/types/key"
)

type Transaction struct {
	TransactionV1
	// source BlockV1, nil if constructed from BlockV2
	OriginDeployV1 *Deploy
}

func NewTransactionFromDeploy(deploy Deploy) Transaction {
	var (
		paymentAmount         uint64
		transactionEntryPoint TransactionEntryPoint
	)

	if deploy.Session.Transfer != nil {
		transactionEntryPoint.Transfer = &struct{}{}
	} else {
		transactionEntryPoint.Custom = &struct {
			Type string
		}{
			Type: "Call",
		}
	}

	if args := deploy.Payment.Args(); args != nil {
		argument, err := args.Find("amount")
		if err == nil {
			parsed, err := argument.Parsed()
			if err == nil {
				var amount string
				json.Unmarshal(parsed, &amount)
				paymentAmount, _ = strconv.ParseUint(amount, 10, 64)
			}
		}
	}

	// Use StandardPayment as true only for payments without explicit `payment amount`
	var standardPayment = paymentAmount == 0
	return Transaction{
		TransactionV1: TransactionV1{
			TransactionV1Hash: &deploy.Hash,
			TransactionV1Header: TransactionV1Header{
				BodyHash:  deploy.Header.BodyHash,
				ChainName: deploy.Header.ChainName,
				Timestamp: deploy.Header.Timestamp,
				TTL:       deploy.Header.TTL,
				InitiatorAddr: InitiatorAddr{
					PublicKey: &deploy.Header.Account,
				},
				PricingMode: PricingMode{
					Classic: &ClassicMode{
						GasPriceTolerance: 1,
						PaymentAmount:     paymentAmount,
						StandardPayment:   standardPayment,
					},
				},
			},
			TransactionV1Body: TransactionV1Body{
				Args:                  deploy.Session.Args(),
				Target:                NewTransactionTargetFromSession(deploy.Session),
				TransactionEntryPoint: transactionEntryPoint,
				TransactionScheduling: TransactionScheduling{
					Standard: &struct{}{},
				},
			},
			Approvals: deploy.Approvals,
		},
		OriginDeployV1: &deploy,
	}
}

type TransactionWrapper struct {
	Deploy        *Deploy        `json:"Deploy,omitempty"`
	TransactionV1 *TransactionV1 `json:"Version1,omitempty"`
}

type TransactionV1 struct {
	// Hex-encoded TransactionV1 hash
	TransactionV1Hash *key.Hash `json:"hash"`
	// The header portion of a TransactionV1
	TransactionV1Header TransactionV1Header `json:"header"`
	// Body of a `TransactionV1`
	TransactionV1Body TransactionV1Body `json:"body"`
	// List of signers and signatures for this `deploy`
	Approvals []Approval `json:"approvals"`
}

type TransactionV1Header struct {
	// `Hash` of the body part of this `Deploy`.
	BodyHash  key.Hash `json:"body_hash"`
	ChainName string   `json:"chain_name"`
	// `Timestamp` formatted as per RFC 3339
	Timestamp Timestamp `json:"timestamp"`
	// Duration of the `Deploy` in milliseconds (from timestamp).
	TTL Duration `json:"ttl"`
	// The address of the initiator of a TransactionV1.
	InitiatorAddr InitiatorAddr `json:"initiator_addr"`
	// Pricing mode of a Transaction.
	PricingMode PricingMode `json:"pricing_mode"`
}

type TransactionEntryPoint struct {
	Custom *struct {
		Type string
	}
	Transfer    *struct{}
	AddBid      *struct{}
	WithdrawBid *struct{}
	Delegate    *struct{}
	Undelegate  *struct{}
	Redelegate  *struct{}
	ActivateBid *struct{}
}

func (t *TransactionEntryPoint) UnmarshalJSON(data []byte) error {
	var custom struct {
		Custom string `json:"Custom"`
	}
	if err := json.Unmarshal(data, &custom); err == nil {
		*t = TransactionEntryPoint{
			Custom: &struct{ Type string }{Type: custom.Custom},
		}
		return nil
	}

	var key string
	if err := json.Unmarshal(data, &key); err != nil {
		return err
	}

	var entryPoint TransactionEntryPoint
	switch key {
	case TransactionEntryPointTransfer:
		entryPoint.Transfer = &struct{}{}
	case TransactionEntryPointAddBid:
		entryPoint.AddBid = &struct{}{}
	case TransactionEntryPointWithdrawBid:
		entryPoint.WithdrawBid = &struct{}{}
	case TransactionEntryPointDelegate:
		entryPoint.Delegate = &struct{}{}
	case TransactionEntryPointUndelegate:
		entryPoint.Undelegate = &struct{}{}
	case TransactionEntryPointRedelegate:
		entryPoint.Redelegate = &struct{}{}
	case TransactionEntryPointActivateBid:
		entryPoint.ActivateBid = &struct{}{}
	}

	*t = entryPoint
	return nil
}

const (
	TransactionEntryPointCustom      = "Custom"
	TransactionEntryPointTransfer    = "Transfer"
	TransactionEntryPointAddBid      = "AddBid"
	TransactionEntryPointWithdrawBid = "WithdrawBid"
	TransactionEntryPointDelegate    = "Delegate"
	TransactionEntryPointUndelegate  = "Undelegate"
	TransactionEntryPointRedelegate  = "Redelegate"
	TransactionEntryPointActivateBid = "ActivateBid"
)

type TransactionV1Body struct {
	Args *Args `json:"args,omitempty"`
	// Execution target of a Transaction.
	Target TransactionTarget `json:"target"`
	// Entry point of a Transaction.
	TransactionEntryPoint TransactionEntryPoint `json:"entry_point"`
	// Scheduling mode of a Transaction.
	TransactionScheduling TransactionScheduling `json:"scheduling"`
}

type TransactionScheduling struct {
	// No special scheduling applied.
	Standard *struct{}
	// Execution should be scheduled for the specified era.
	FutureEra *FutureEraScheduling
	// Execution should be scheduled for the specified timestamp or later.
	FutureTimestamp *FutureTimestampScheduling
}

func (t *TransactionScheduling) UnmarshalJSON(data []byte) error {
	var futureKey struct {
		EraID           *uint64 `json:"FutureEra"`
		FutureTimestamp *string `json:"FutureTimestamp"`
	}
	if err := json.Unmarshal(data, &futureKey); err == nil {
		if futureKey.FutureTimestamp != nil {
			*t = TransactionScheduling{
				FutureTimestamp: &FutureTimestampScheduling{
					TimeStamp: *futureKey.FutureTimestamp,
				},
			}
		}

		if futureKey.EraID != nil {
			*t = TransactionScheduling{
				FutureEra: &FutureEraScheduling{
					EraID: *futureKey.EraID,
				},
			}
		}
		return nil
	}

	var key string
	if err := json.Unmarshal(data, &key); err == nil && key == "Standard" {
		*t = TransactionScheduling{
			Standard: &struct{}{},
		}
		return nil
	}

	return nil
}

type FutureEraScheduling struct {
	EraID uint64
}
type FutureTimestampScheduling struct {
	TimeStamp string `json:"FutureTimestamp"`
}

type TransactionTarget struct {
	// The execution target is a native operation (e.g. a transfer).
	Native *struct{}
	// The execution target is a stored entity or package.
	Stored *StoredTarget
	// The execution target is the included module bytes, i.e. compiled Wasm.
	Session *SessionTarget
}

// NewTransactionTargetFromSession create new TransactionTarget from ExecutableDeployItem
func NewTransactionTargetFromSession(session ExecutableDeployItem) TransactionTarget {
	if session.Transfer != nil {
		return TransactionTarget{
			Native: &struct{}{},
		}
	}

	if session.ModuleBytes != nil {
		return TransactionTarget{
			Session: &SessionTarget{
				ModuleBytes: session.ModuleBytes.ModuleBytes,
				Runtime:     "VmCasperV1",
			},
		}
	}

	if session.StoredContractByHash != nil {
		hash := session.StoredContractByHash.Hash.Hash.ToHex()
		return TransactionTarget{
			Stored: &StoredTarget{
				ID: TransactionInvocationTarget{
					ByHash: &hash,
				},
				Runtime: "VmCasperV1",
			},
		}
	}

	if session.StoredContractByName != nil {
		return TransactionTarget{
			Stored: &StoredTarget{
				ID: TransactionInvocationTarget{
					ByName: &session.StoredContractByName.Name,
				},
				Runtime: "VmCasperV1",
			},
		}
	}

	if session.StoredVersionedContractByHash != nil {
		var version *uint32
		if storedVersion := session.StoredVersionedContractByHash.Version; storedVersion != nil {
			versionNum, err := storedVersion.Int64()
			if err == nil {
				temp := uint32(versionNum)
				version = &temp
			}
		}
		byHashTarget := ByPackageHashInvocationTarget{
			Addr:    session.StoredVersionedContractByHash.Hash.Hash,
			Version: version,
		}
		return TransactionTarget{
			Stored: &StoredTarget{
				ID: TransactionInvocationTarget{
					ByPackageHash: &byHashTarget,
				},
				Runtime: "VmCasperV1",
			},
		}
	}

	if session.StoredVersionedContractByName != nil {
		var version *uint32
		if storedVersion := session.StoredVersionedContractByName.Version; storedVersion != nil {
			versionNum, err := storedVersion.Int64()
			if err == nil {
				temp := uint32(versionNum)
				version = &temp
			}
		}
		byNameTarget := ByPackageNameInvocationTarget{
			Name:    session.StoredContractByName.Name,
			Version: version,
		}
		return TransactionTarget{
			Stored: &StoredTarget{
				ID: TransactionInvocationTarget{
					ByPackageName: &byNameTarget,
				},
				Runtime: "VmCasperV1",
			},
		}
	}

	return TransactionTarget{}
}

type StoredTarget struct {
	// Identifier of a `Stored` transaction target.
	ID      TransactionInvocationTarget `json:"id"`
	Runtime string                      `json:"runtime"`
}

type TransactionInvocationTarget struct {
	// Hex-encoded entity address identifying the invocable entity.
	ByHash *string `json:"ByHash"`
	// The alias identifying the invocable entity.
	ByName *string `json:"ByName"`
	// The address and optional version identifying the package.
	ByPackageHash *ByPackageHashInvocationTarget `json:"ByPackageHash"`
	// The alias and optional version identifying the package.
	ByPackageName *ByPackageNameInvocationTarget `json:"ByPackageName"`
}

// ByPackageHashInvocationTarget The address and optional version identifying the package.
type ByPackageHashInvocationTarget struct {
	Addr    key.Hash `json:"addr"`
	Version *uint32  `json:"version"`
}

// ByPackageNameInvocationTarget The alias and optional version identifying the package.
type ByPackageNameInvocationTarget struct {
	Name    string  `json:"name"`
	Version *uint32 `json:"version"`
}

type SessionTarget struct {
	Kind        string `json:"string"`
	ModuleBytes string `json:"module_bytes"`
	Runtime     string `json:"runtime"`
}

type PricingMode struct {
	// The original payment model, where the creator of the transaction specifies how much they will pay, at what gas price.
	Classic *ClassicMode `json:"Classic"`
	// The cost of the transaction is determined by the cost table, per the transaction kind.
	Fixed *FixedMode `json:"Fixed"`

	Reserved *ReservedMode `json:"reserved"`
}

type ClassicMode struct {
	// User-specified gas_price tolerance (minimum 1). This is interpreted to mean "do not include this transaction in a block if the current gas price is greater than this number"
	GasPriceTolerance uint8 `json:"gas_price_tolerance"`
	// User-specified payment amount.
	PaymentAmount uint64 `json:"payment_amount"`
	// Standard payment.
	StandardPayment bool `json:"standard_payment"`
}

type FixedMode struct {
	// 	// User-specified gas_price tolerance (minimum 1). This is interpreted to mean "do not include this transaction in a block if the current gas price is greater than this number"
	GasPriceTolerance uint8 `json:"gas_price_tolerance"`
}

type ReservedMode struct {
	// Pre-paid receipt
	Receipt key.Hash `json:"receipt"`
	// Price paid in the past to reserve space in a future block.
	PaidAmount uint64 `json:"paid_amount"`
	// The gas price at the time of reservation.
	StrikePrice uint `json:"strike_price"`
}

// TransactionHash A versioned wrapper for a transaction hash or deploy hash
type TransactionHash struct {
	Deploy      *key.Hash `json:"Deploy,omitempty"`
	Transaction *key.Hash `json:"Version1,omitempty"`
}
