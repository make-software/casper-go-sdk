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
		DeployHash      key.Hash                    `json:"deploy_hash"`
		Account         string                      `json:"account"`
		Timestamp       time.Time                   `json:"timestamp"`
		TTL             string                      `json:"ttl"`
		BlockHash       key.Hash                    `json:"block_hash"`
		ExecutionResult types.ExecutionResultStatus `json:"execution_result"`
	}
	DeployProcessedEvent struct {
		DeployProcessed DeployProcessedPayload `json:"DeployProcessed"`
	}
	DeployAcceptedEvent struct {
		DeployAccepted types.Deploy `json:"DeployAccepted"`
	}
	DeployExpiredPayload struct {
		DeployHash key.Hash `json:"deploy_hash"`
	}
	DeployExpiredEvent struct {
		DeployExpired DeployExpiredPayload `json:"DeployExpired"`
	}
)

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
		FinalitySignatureV2

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
				FinalitySignatureV2: FinalitySignatureV2{
					BlockHash: wrapped.FinalitySignature.V1.BlockHash,
					EraID:     wrapped.FinalitySignature.V1.EraID,
					Signature: wrapped.FinalitySignature.V1.Signature,
					PublicKey: wrapped.FinalitySignature.V1.PublicKey,
				},
				OriginFinalitySignatureV1: wrapped.FinalitySignature.V1,
			},
		}
		return nil
	}

	if wrapped.FinalitySignature.V2 != nil {
		*t = FinalitySignatureEvent{
			FinalitySignature: FinalitySignature{
				FinalitySignatureV2: *wrapped.FinalitySignature.V2,
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
			FinalitySignatureV2: FinalitySignatureV2{
				BlockHash: v1Event.FinalitySignature.BlockHash,
				EraID:     v1Event.FinalitySignature.EraID,
				Signature: v1Event.FinalitySignature.Signature,
				PublicKey: v1Event.FinalitySignature.PublicKey,
			},
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
		// Todo: not sure, didn't found example to test
		Operations *[]types.Operation `json:"operations,omitempty"`
		// Todo: not sure, didn't found example to test
		Transform *types.TransformKey `json:"transform,omitempty"`
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

func (d *RawEvent) ParseAsDeployExpiredEvent() (DeployExpiredEvent, error) {
	return ParseEvent[DeployExpiredEvent](d.Data)
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
