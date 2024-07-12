package sse

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/casper-go-sdk/types/keypair"
)

const (
	APIVersionEventType EventType = iota + 1
	BlockAddedEventType
	DeployProcessedEventType
	DeployAcceptedEventType
	DeployExpiredEventType
	EventIDEventType
	FinalitySignatureType
	StepEventType
	FaultEventType
	ShutdownType
)

var AllEventsNames = map[EventType]string{
	APIVersionEventType:      "ApiVersion",
	BlockAddedEventType:      "BlockAdded",
	DeployProcessedEventType: "DeployProcessed",
	DeployAcceptedEventType:  "DeployAccepted",
	DeployExpiredEventType:   "DeployExpired",
	StepEventType:            "Step",
	FaultEventType:           "Fault",
	FinalitySignatureType:    "FinalitySignature",
	ShutdownType:             "Shutdown",
}

type (
	EventType = int
	EventData = json.RawMessage
)

type RawEvent struct {
	EventType EventType
	Data      EventData
	EventID   uint64
}

type APIVersionEvent struct {
	APIVersion string `json:"ApiVersion"`
}

// BlockAddedEvent definition
type (
	BlockAdded struct {
		BlockHash string      `json:"block_hash"`
		Block     types.Block `json:"block"`
	}

	BlockAddedEvent struct {
		BlockAdded BlockAdded `json:"BlockAdded"`
	}

	blockAddedV1 struct {
		BlockHash string        `json:"block_hash"`
		Block     types.BlockV1 `json:"block"`
	}

	blockAddedEventV1 struct {
		BlockAdded blockAddedV1 `json:"BlockAdded"`
	}

	blockAddedWrapper struct {
		BlockHash string             `json:"block_hash"`
		Block     types.BlockWrapper `json:"block"`
	}

	blockAddedEventWrapper struct {
		BlockAdded blockAddedWrapper `json:"BlockAdded"`
	}
)

func (t *BlockAddedEvent) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}

	wrapped := blockAddedEventWrapper{}
	if err := json.Unmarshal(data, &wrapped); err != nil {
		return err
	}

	if wrapped.BlockAdded.Block.BlockV1 != nil || wrapped.BlockAdded.Block.BlockV2 != nil {
		*t = BlockAddedEvent{
			BlockAdded: BlockAdded{
				BlockHash: wrapped.BlockAdded.BlockHash,
				Block:     types.NewBlockFromBlockWrapper(wrapped.BlockAdded.Block, nil),
			},
		}
		return nil
	}

	v1Event := blockAddedEventV1{}
	if err := json.Unmarshal(data, &v1Event); err != nil {
		return err
	}

	*t = BlockAddedEvent{
		BlockAdded: BlockAdded{
			BlockHash: wrapped.BlockAdded.BlockHash,
			Block:     types.NewBlockFromBlockV1(v1Event.BlockAdded.Block),
		},
	}
	return nil
}

type (
	DeployProcessedPayload struct {
		DeployHash      key.Hash                `json:"deploy_hash"`
		Account         keypair.PublicKey       `json:"account"`
		Timestamp       time.Time               `json:"timestamp"`
		TTL             string                  `json:"ttl"`
		BlockHash       key.Hash                `json:"block_hash"`
		ExecutionResult types.ExecutionResultV1 `json:"execution_result"`
	}
	DeployProcessedEvent struct {
		DeployProcessed DeployProcessedPayload `json:"DeployProcessed"`
	}
	DeployAcceptedEvent struct {
		DeployAccepted types.Deploy `json:"DeployAccepted"`
	}

	deployExpiredPayload struct {
		DeployHash key.Hash `json:"deploy_hash"`
	}

	deployExpiredEvent struct {
		DeployExpired deployExpiredPayload `json:"DeployExpired"`
	}

	TransactionAcceptedPayload struct {
		Transaction types.Transaction `json:"transaction"`
	}

	TransactionAcceptedEvent struct {
		TransactionAcceptedPayload TransactionAcceptedPayload `json:"TransactionAccepted"`
	}

	TransactionExpiredPayload struct {
		TransactionHash types.TransactionHash `json:"transaction_hash"`
	}

	TransactionExpiredEvent struct {
		TransactionExpiredPayload TransactionExpiredPayload `json:"TransactionExpired"`
	}

	TransactionProcessedPayload struct {
		BlockHash       key.Hash              `json:"block_hash"`
		TransactionHash types.TransactionHash `json:"transaction_hash"`
		InitiatorAddr   types.InitiatorAddr   `json:"initiator_addr"`
		Timestamp       time.Time             `json:"timestamp"`
		TTL             string                `json:"ttl"`
		ExecutionResult types.ExecutionResult `json:"execution_result"`
		Messages        []types.Message       `json:"messages"`
	}

	TransactionProcessedEvent struct {
		TransactionProcessedPayload TransactionProcessedPayload `json:"TransactionProcessed"`
	}
)

func (t *TransactionAcceptedEvent) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}

	transactionEvent := struct {
		TransactionAcceptedPayload types.TransactionWrapper `json:"TransactionAccepted"`
	}{}
	if err := json.Unmarshal(data, &transactionEvent); err != nil {
		return err
	}

	if deploy := transactionEvent.TransactionAcceptedPayload.Deploy; deploy != nil {
		*t = TransactionAcceptedEvent{
			TransactionAcceptedPayload: TransactionAcceptedPayload{
				Transaction: types.NewTransactionFromDeploy(*deploy),
			},
		}
		return nil
	}

	if v1 := transactionEvent.TransactionAcceptedPayload.TransactionV1; v1 != nil {
		*t = TransactionAcceptedEvent{
			TransactionAcceptedPayload: TransactionAcceptedPayload{
				Transaction: types.NewTransactionFromTransactionV1(*v1),
			},
		}
		return nil
	}

	deployEvent := DeployAcceptedEvent{}
	if err := json.Unmarshal(data, &deployEvent); err != nil {
		return err
	}

	*t = TransactionAcceptedEvent{
		TransactionAcceptedPayload: TransactionAcceptedPayload{
			Transaction: types.NewTransactionFromDeploy(deployEvent.DeployAccepted),
		},
	}
	return nil
}

func (t *TransactionProcessedEvent) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}

	transactionEvent := struct {
		TransactionProcessedPayload TransactionProcessedPayload `json:"TransactionProcessed"`
	}{}
	if err := json.Unmarshal(data, &transactionEvent); err != nil {
		return err
	}

	if transactionEvent.TransactionProcessedPayload.TransactionHash.Transaction != nil ||
		transactionEvent.TransactionProcessedPayload.TransactionHash.Deploy != nil {
		*t = transactionEvent
		return nil
	}

	deployEvent := DeployProcessedEvent{}
	if err := json.Unmarshal(data, &deployEvent); err != nil {
		return err
	}

	*t = TransactionProcessedEvent{
		TransactionProcessedPayload: TransactionProcessedPayload{
			BlockHash: deployEvent.DeployProcessed.BlockHash,
			TransactionHash: types.TransactionHash{
				Deploy: &deployEvent.DeployProcessed.DeployHash,
			},
			InitiatorAddr: types.InitiatorAddr{
				PublicKey: &deployEvent.DeployProcessed.Account,
			},
			Timestamp:       deployEvent.DeployProcessed.Timestamp,
			TTL:             deployEvent.DeployProcessed.TTL,
			ExecutionResult: types.NewExecutionResultFromV1(deployEvent.DeployProcessed.ExecutionResult),
		},
	}
	return nil
}

func (t *TransactionExpiredEvent) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}

	transactionEvent := struct {
		TransactionExpiredPayload TransactionExpiredPayload `json:"TransactionExpired"`
	}{}
	if err := json.Unmarshal(data, &transactionEvent); err != nil {
		return err
	}

	if transactionEvent.TransactionExpiredPayload.TransactionHash.Transaction != nil ||
		transactionEvent.TransactionExpiredPayload.TransactionHash.Deploy != nil {
		*t = transactionEvent
		return nil
	}

	deployEvent := deployExpiredEvent{}
	if err := json.Unmarshal(data, &deployEvent); err != nil {
		return err
	}

	*t = TransactionExpiredEvent{
		TransactionExpiredPayload: TransactionExpiredPayload{
			TransactionHash: types.TransactionHash{
				Deploy: &deployEvent.DeployExpired.DeployHash,
			},
		},
	}
	return nil
}

type (
	FinalitySignatureV1 struct {
		BlockHash key.Hash          `json:"block_hash"`
		EraID     uint64            `json:"era_id"`
		Signature types.HexBytes    `json:"signature"`
		PublicKey keypair.PublicKey `json:"public_key"`
	}

	FinalitySignatureV2 struct {
		BlockHash     key.Hash          `json:"block_hash"`
		BlockHeight   *uint64           `json:"block_height"`
		ChainNameHash *key.Hash         `json:"chain_name_hash"`
		EraID         uint64            `json:"era_id"`
		Signature     types.HexBytes    `json:"signature"`
		PublicKey     keypair.PublicKey `json:"public_key"`
	}

	finalitySignatureWrapper struct {
		V1 *FinalitySignatureV1 `json:"V1"`
		V2 *FinalitySignatureV2 `json:"V2"`
	}

	finalitySignatureWrapperEvent struct {
		FinalitySignature finalitySignatureWrapper `json:"FinalitySignature"`
	}

	finalitySignatureV1Event struct {
		FinalitySignature FinalitySignatureV1 `json:"FinalitySignature"`
	}

	FinalitySignature struct {
		BlockHash     key.Hash          `json:"block_hash"`
		BlockHeight   *uint64           `json:"block_height"`
		ChainNameHash *key.Hash         `json:"chain_name_hash"`
		EraID         uint64            `json:"era_id"`
		Signature     types.HexBytes    `json:"signature"`
		PublicKey     keypair.PublicKey `json:"public_key"`

		OriginFinalitySignatureV1 *FinalitySignatureV1
	}

	FinalitySignatureEvent struct {
		FinalitySignature FinalitySignature `json:"FinalitySignature"`
	}
)

func (t *FinalitySignatureEvent) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}

	wrapped := finalitySignatureWrapperEvent{}
	if err := json.Unmarshal(data, &wrapped); err != nil {
		return err
	}

	if wrapped.FinalitySignature.V1 != nil {
		*t = FinalitySignatureEvent{
			FinalitySignature: FinalitySignature{
				BlockHash:                 wrapped.FinalitySignature.V1.BlockHash,
				EraID:                     wrapped.FinalitySignature.V1.EraID,
				Signature:                 wrapped.FinalitySignature.V1.Signature,
				PublicKey:                 wrapped.FinalitySignature.V1.PublicKey,
				OriginFinalitySignatureV1: wrapped.FinalitySignature.V1,
			},
		}
		return nil
	}

	if wrapped.FinalitySignature.V2 != nil {
		*t = FinalitySignatureEvent{
			FinalitySignature: FinalitySignature{
				BlockHash:     wrapped.FinalitySignature.V2.BlockHash,
				BlockHeight:   wrapped.FinalitySignature.V2.BlockHeight,
				ChainNameHash: wrapped.FinalitySignature.V2.ChainNameHash,
				EraID:         wrapped.FinalitySignature.V2.EraID,
				Signature:     wrapped.FinalitySignature.V2.Signature,
				PublicKey:     wrapped.FinalitySignature.V2.PublicKey,
			},
		}
		return nil
	}

	v1Event := finalitySignatureV1Event{}
	if err := json.Unmarshal(data, &v1Event); err != nil {
		return err
	}

	*t = FinalitySignatureEvent{
		FinalitySignature: FinalitySignature{
			BlockHash:                 v1Event.FinalitySignature.BlockHash,
			EraID:                     v1Event.FinalitySignature.EraID,
			Signature:                 v1Event.FinalitySignature.Signature,
			PublicKey:                 v1Event.FinalitySignature.PublicKey,
			OriginFinalitySignatureV1: &v1Event.FinalitySignature,
		},
	}
	return nil
}

type (
	FaultPayload struct {
		EraID     uint64            `json:"era_id"`
		PublicKey keypair.PublicKey `json:"public_key"`
		Timestamp types.Timestamp   `json:"timestamp"`
	}
	FaultEvent struct {
		Fault FaultPayload `json:"Fault"`
	}
)

type (
	StepPayload struct {
		EraID           uint64       `json:"era_id"`
		ExecutionEffect types.Effect `json:"execution_effect"`
	}
	StepEvent struct {
		Step StepPayload `json:"step"`
	}
)

func (d *RawEvent) ParseAsAPIVersionEvent() (APIVersionEvent, error) {
	return ParseEvent[APIVersionEvent](d.Data)
}

func (d *RawEvent) ParseAsDeployProcessedEvent() (DeployProcessedEvent, error) {
	return ParseEvent[DeployProcessedEvent](d.Data)
}

func (d *RawEvent) ParseAsBlockAddedEvent() (BlockAddedEvent, error) {
	return ParseEvent[BlockAddedEvent](d.Data)
}

func (d *RawEvent) ParseAsDeployAcceptedEvent() (DeployAcceptedEvent, error) {
	return ParseEvent[DeployAcceptedEvent](d.Data)
}

func (d *RawEvent) ParseAsFinalitySignatureEvent() (FinalitySignatureEvent, error) {
	return ParseEvent[FinalitySignatureEvent](d.Data)
}

func (d *RawEvent) ParseAsTransactionExpiredEvent() (TransactionExpiredEvent, error) {
	return ParseEvent[TransactionExpiredEvent](d.Data)
}

func (d *RawEvent) ParseAsTransactionProcessedEvent() (TransactionProcessedEvent, error) {
	return ParseEvent[TransactionProcessedEvent](d.Data)
}

func (d *RawEvent) ParseAsTransactionAcceptedEvent() (TransactionAcceptedEvent, error) {
	return ParseEvent[TransactionAcceptedEvent](d.Data)
}

func (d *RawEvent) ParseAsFaultEvent() (FaultEvent, error) {
	return ParseEvent[FaultEvent](d.Data)
}

func (d *RawEvent) ParseAsStepEvent() (StepEvent, error) {
	return ParseEvent[StepEvent](d.Data)
}

func ParseEvent[T interface{}](data []byte) (T, error) {
	var res T
	if err := json.Unmarshal(data, &res); err != nil {
		return res, err
	}
	return res, nil
}
