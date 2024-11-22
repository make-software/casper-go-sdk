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

type TransactionVersion uint

const (
	TransactionVersionV1 TransactionVersion = iota
	TransactionDeploy
)

type Transaction struct {
	// Hex-encoded Transaction hash
	Hash key.Hash `json:"hash"`
	// The header portion of a Transaction
	Payload TransactionPayload `json:"payload"`
	// List of signers and signatures for this Transaction
	Approvals []Approval `json:"approvals"`

	// source DeployV1, nil if constructed from TransactionV1
	originDeployV1      *Deploy
	originTransactionV1 *TransactionV1
}

type TransactionPayload struct {
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
	// ============== Transaction Body =====================
	Args *Args `json:"args,omitempty"`
	// Execution target of a Transaction.
	Target TransactionTarget `json:"target"`
	// Entry point of a Transaction.
	EntryPoint TransactionEntryPoint `json:"entry_point"`
	// Scheduling mode of a Transaction.
	Scheduling TransactionScheduling `json:"scheduling"`
	// Transaction category
	Category uint8 `json:"transaction_category"`
}

type TransactionHeader struct {
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
}

func (t *Transaction) GetDeploy() *Deploy {
	return t.originDeployV1
}

func (t *Transaction) GetTransactionV1() *TransactionV1 {
	return t.originTransactionV1
}

func NewTransactionFromTransactionV1(v1 TransactionV1) (Transaction, error) {
	var (
		args                  *Args
		transactionTarget     *TransactionTarget
		transactionEntryPoint *TransactionEntryPoint
		transactionScheduling *TransactionScheduling
	)

	for key, rawData := range v1.Payload.Fields {
		var err error
		switch key {
		case 0:
			decoder := &ArgsFromBytesDecoder{}
			args, _, err = decoder.FromBytes(rawData)
		case 1:
			decoder := TransactionTargetFromBytesDecoder{}
			transactionTarget, _, err = decoder.FromBytes(rawData)
		case 2:
			decoder := TransactionEntryPointFromBytesDecoder{}
			transactionEntryPoint, _, err = decoder.FromBytes(rawData)
		case 3:
			decoder := TransactionSchedulingFromBytesDecoder{}
			transactionScheduling, _, err = decoder.FromBytes(rawData)
		default:
			return Transaction{}, errors.New("unsupported field key")
		}

		if err != nil {
			return Transaction{}, err
		}
	}

	return Transaction{
		Hash: v1.Hash,
		Payload: TransactionPayload{
			ChainName:     v1.Payload.ChainName,
			Timestamp:     v1.Payload.Timestamp,
			TTL:           v1.Payload.TTL,
			InitiatorAddr: v1.Payload.InitiatorAddr,
			PricingMode:   v1.Payload.PricingMode,
			Args:          args,
			Target:        *transactionTarget,
			EntryPoint:    *transactionEntryPoint,
			Scheduling:    *transactionScheduling,
		},
		Approvals:           v1.Approvals,
		originTransactionV1: &v1,
	}, nil
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

	// Use StandardPayment as true only for payments without explicit `payment amount`
	var standardPayment = paymentAmount == 0 && deploy.Payment.ModuleBytes == nil
	return Transaction{
		Hash: deploy.Hash,
		Payload: TransactionPayload{
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
			Category: uint8(transactionCategory),
		},
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
	// Body of a `TransactionV1`
	Body TransactionV1Body `json:"body"`
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

type TransactionV1Body struct {
	Args *Args `json:"args,omitempty"`
	// Execution target of a Transaction.
	Target TransactionTarget `json:"target"`
	// Entry point of a Transaction.
	TransactionEntryPoint TransactionEntryPoint `json:"entry_point"`
	// Scheduling mode of a Transaction.
	TransactionScheduling TransactionScheduling `json:"scheduling"`
	// Transaction category
	TransactionCategory uint8 `json:"transaction_category"`
}

// TransactionHash A versioned wrapper for a transaction hash or deploy hash
type TransactionHash struct {
	Deploy        *key.Hash `json:"Deploy,omitempty"`
	TransactionV1 *key.Hash `json:"Version1,omitempty"`
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
	//bodyBytes, err := t.Body.Bytes()
	//if err != nil {
	//	return err
	//}
	//
	//if t.Header.BodyHash != blake2b.Sum256(bodyBytes) {
	//	return ErrInvalidBodyHash
	//}
	//
	//if t.Hash != blake2b.Sum256(t.Header.Bytes()) {
	//	return ErrInvalidTransactionHash
	//}

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
