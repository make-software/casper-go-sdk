package types

import (
	"encoding/json"
	"errors"
	"strconv"

	"golang.org/x/crypto/blake2b"

	"github.com/make-software/casper-go-sdk/v2/types/key"
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

var (
	ErrInvalidBodyHash          = errors.New("invalid body hash")
	ErrInvalidTransactionHash   = errors.New("invalid transaction hash")
	ErrInvalidApprovalSignature = errors.New("invalid approval signature")
)

type TransactionCategory uint

const (
	TransactionCategoryMint TransactionCategory = iota
	TransactionCategoryAuction
	TransactionCategoryInstallUpgrade
	TransactionCategoryLarge
	TransactionCategoryMedium
	TransactionCategorySmall
)

type Transaction struct {
	// Hex-encoded Transaction hash
	Hash key.Hash `json:"hash"`
	// Transaction chain name
	ChainName string `json:"chain_name"`
	// `Timestamp` formatted as per RFC 3339
	Timestamp Timestamp `json:"timestamp"`
	// Duration of the `TransactionV1` in milliseconds (from timestamp).
	TTL Duration `json:"ttl"`
	// The address of the initiator of a Transaction.
	InitiatorAddr InitiatorAddr `json:"initiator_addr"`
	// Pricing mode of a Transaction.
	PricingMode PricingMode `json:"pricing_mode"`
	// Args transaction args
	Args *Args `json:"args,omitempty"`
	// Execution target of a Transaction.
	Target TransactionTarget `json:"target"`
	// Entry point of a Transaction.
	EntryPoint TransactionEntryPoint `json:"entry_point"`
	// Scheduling mode of a Transaction.
	Scheduling TransactionScheduling `json:"scheduling"`
	// List of signers and signatures for this Transaction
	Approvals []Approval `json:"approvals"`
	// Transaction category
	Category uint8 `json:"transaction_category"`
	// source DeployV1, nil if constructed from TransactionV1
	originDeployV1      *Deploy
	originTransactionV1 *TransactionV1
}

func (t *Transaction) GetDeploy() *Deploy {
	return t.originDeployV1
}

func (t *Transaction) GetTransactionV1() *TransactionV1 {
	return t.originTransactionV1
}

func NewTransactionFromTransactionV1(v1 TransactionV1) Transaction {
	return Transaction{
		Hash:                v1.Hash,
		ChainName:           v1.Payload.ChainName,
		Timestamp:           v1.Payload.Timestamp,
		TTL:                 v1.Payload.TTL,
		InitiatorAddr:       v1.Payload.InitiatorAddr,
		PricingMode:         v1.Payload.PricingMode,
		Args:                v1.Payload.Fields.NamedArgs.Args,
		Target:              v1.Payload.Fields.Target,
		EntryPoint:          v1.Payload.Fields.TransactionEntryPoint,
		Scheduling:          v1.Payload.Fields.TransactionScheduling,
		Approvals:           v1.Approvals,
		originTransactionV1: &v1,
	}
}

func NewTransactionFromDeploy(deploy Deploy) Transaction {
	var (
		paymentAmount         uint64
		transactionEntryPoint TransactionEntryPoint
		transactionCategory   = TransactionCategoryLarge
	)

	if deploy.Session.Transfer != nil {
		transactionCategory = TransactionCategoryMint
		transactionEntryPoint.Transfer = &struct{}{}
	} else if deploy.Session.ModuleBytes != nil {
		transactionEntryPoint.Call = &struct{}{}
	} else {
		var entrypoint string
		switch {
		case deploy.Session.StoredContractByHash != nil:
			entrypoint = deploy.Session.StoredContractByHash.EntryPoint
		case deploy.Session.StoredContractByName != nil:
			entrypoint = deploy.Session.StoredContractByName.EntryPoint
		case deploy.Session.StoredVersionedContractByHash != nil:
			entrypoint = deploy.Session.StoredVersionedContractByHash.EntryPoint
		case deploy.Session.StoredVersionedContractByName != nil:
			entrypoint = deploy.Session.StoredVersionedContractByName.EntryPoint
		}
		transactionEntryPoint.Custom = &entrypoint
	}

	if args := deploy.Payment.Args(); args != nil {
		argument, err := args.Find("amount")
		if err == nil {
			parsed, err := argument.Parsed()
			if err == nil {
				var amount string
				if err := json.Unmarshal(parsed, &amount); err == nil {
					paymentAmount, _ = strconv.ParseUint(amount, 10, 64)
				}
			}
		}
	}

	// Use StandardPayment as true only for payments with module_bytes: ""
	var standardPayment = deploy.Payment.ModuleBytes != nil && deploy.Payment.ModuleBytes.ModuleBytes == ""
	return Transaction{
		Hash:      deploy.Hash,
		ChainName: deploy.Header.ChainName,
		Timestamp: deploy.Header.Timestamp,
		TTL:       deploy.Header.TTL,
		InitiatorAddr: InitiatorAddr{
			PublicKey: &deploy.Header.Account,
		},
		PricingMode: PricingMode{
			Limited: &LimitedMode{
				GasPriceTolerance: 1,
				PaymentAmount:     paymentAmount,
				StandardPayment:   standardPayment,
			},
		},
		Args:       deploy.Session.Args(),
		Target:     NewTransactionTargetFromSession(deploy.Session),
		EntryPoint: transactionEntryPoint,
		Scheduling: TransactionScheduling{
			Standard: &struct{}{},
		},
		Category:       uint8(transactionCategory),
		Approvals:      deploy.Approvals,
		originDeployV1: &deploy,
	}
}

type TransactionWrapper struct {
	Deploy        *Deploy        `json:"Deploy,omitempty"`
	TransactionV1 *TransactionV1 `json:"Version1,omitempty"`
}

type TransactionV1 struct {
	// Hex-encoded TransactionV1 hash
	Hash key.Hash `json:"hash"`
	// Transaction payload
	Payload TransactionV1Payload `json:"payload"`
	// List of signers and signatures for this `deploy`
	Approvals []Approval `json:"approvals"`
}

// TransactionHash A versioned wrapper for a transaction hash or deploy hash
type TransactionHash struct {
	Deploy        *key.Hash `json:"Deploy,omitempty"`
	TransactionV1 *key.Hash `json:"Version1,omitempty"`
}

func (t *TransactionHash) String() string {
	if t.Deploy != nil {
		return t.Deploy.String()
	} else if t.TransactionV1 != nil {
		return t.TransactionV1.String()
	} else {
		return ""
	}
}

func (t *TransactionHash) ToHash() key.Hash {
	if t.Deploy != nil {
		return *t.Deploy
	} else if t.TransactionV1 != nil {
		return *t.TransactionV1
	} else {
		return key.Hash{}
	}
}

func (t *TransactionV1) Sign(keys keypair.PrivateKey) error {
	signature, err := keys.Sign(t.Hash.Bytes())
	if err != nil {
		return err
	}
	approval := Approval{
		Signer:    keys.PublicKey(),
		Signature: signature,
	}

	if t.Approvals == nil {
		t.Approvals = make([]Approval, 0, 1)
	}

	t.Approvals = append(t.Approvals, approval)
	return nil
}

func (t *TransactionV1) Validate() error {
	payloadBytes, err := t.Payload.Bytes()
	if err != nil {
		return err
	}

	if t.Hash != blake2b.Sum256(payloadBytes) {
		return ErrInvalidTransactionHash
	}

	for _, one := range t.Approvals {
		if one.Signer.VerifySignature(t.Hash.Bytes(), one.Signature) != nil {
			return ErrInvalidApprovalSignature
		}
	}

	return nil
}

func NewTransactionV1(hash key.Hash, payload TransactionV1Payload, approvals []Approval) *TransactionV1 {
	return &TransactionV1{
		Hash:      hash,
		Payload:   payload,
		Approvals: approvals,
	}
}

func MakeTransactionV1(transactionPayload TransactionV1Payload) (*TransactionV1, error) {
	payloadBytes, err := transactionPayload.Bytes()
	if err != nil {
		return nil, err
	}

	transactionHash := blake2b.Sum256(payloadBytes)
	return NewTransactionV1(transactionHash, transactionPayload, make([]Approval, 0)), nil
}
