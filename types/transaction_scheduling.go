package types

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/serialization"
	"github.com/make-software/casper-go-sdk/v2/types/serialization/encoding"
)

const (
	TransactionSchedulingStandardTag = iota
	TransactionSchedulingFutureEraTag
	TransactionSchedulingFutureTimestampTag
)

const (
	FutureEraEraIDIndex           uint16 = 1
	FutureTimestampTimestampIndex uint16 = 1
)

type TransactionScheduling struct {
	// No special scheduling applied.
	Standard *struct{} `json:"Standard,omitempty"`
	// Execution should be scheduled for the specified era.
	FutureEra *FutureEraScheduling `json:"FutureTimestamp,omitempty"`
	// Execution should be scheduled for the specified timestamp or later.
	FutureTimestamp *FutureTimestampScheduling `json:"FutureEra,omitempty"`
}

func (t *TransactionScheduling) SerializedLength() int {
	envelope := serialization.CallTableSerializationEnvelope{}
	return envelope.EstimateSize(t.serializedFieldLengths())
}

func (t *TransactionScheduling) Tag() byte {
	switch {
	case t.Standard != nil:
		return TransactionSchedulingStandardTag
	case t.FutureEra != nil:
		return TransactionSchedulingFutureEraTag
	case t.FutureTimestamp != nil:
		return TransactionSchedulingFutureTimestampTag
	default:
		return 0
	}
}

func (t *TransactionScheduling) UnmarshalJSON(data []byte) error {
	var futureKey struct {
		EraID           *uint64    `json:"FutureEra"`
		FutureTimestamp *Timestamp `json:"FutureTimestamp"`
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

func (t TransactionScheduling) MarshalJSON() ([]byte, error) {
	if t.Standard != nil {
		return json.Marshal("Standard")
	}

	if t.FutureTimestamp != nil {
		return json.Marshal(struct {
			FutureTimestamp Timestamp `json:"FutureTimestamp"`
		}{
			FutureTimestamp: t.FutureTimestamp.TimeStamp,
		})
	}

	if t.FutureEra != nil {
		return json.Marshal(struct {
			FutureEra uint64 `json:"FutureEra"`
		}{
			FutureEra: t.FutureEra.EraID,
		})
	}

	return nil, errors.New("unknown scheduling type")
}

type FutureEraScheduling struct {
	EraID uint64
}

type FutureTimestampScheduling struct {
	TimeStamp Timestamp `json:"FutureTimestamp"`
}

func (t *TransactionScheduling) Bytes() ([]byte, error) {
	builder, err := serialization.NewCallTableSerializationEnvelopeBuilder(t.serializedFieldLengths())
	if err != nil {
		return nil, err
	}

	switch {
	case t.Standard != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionSchedulingStandardTag}); err != nil {
			return nil, err
		}
	case t.FutureEra != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionSchedulingFutureEraTag}); err != nil {
			return nil, err
		}

		eraIDBytes, _ := encoding.NewU64ToBytesEncoder(t.FutureEra.EraID).Bytes()
		if err = builder.AddField(FutureEraEraIDIndex, eraIDBytes); err != nil {
			return nil, err
		}
	case t.FutureTimestamp != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionSchedulingFutureTimestampTag}); err != nil {
			return nil, err
		}

		timestampBytes := clvalue.NewCLUInt64(uint64(time.Time(t.FutureTimestamp.TimeStamp).UnixMilli())).Bytes()
		if err = builder.AddField(FutureTimestampTimestampIndex, timestampBytes); err != nil {
			return nil, err
		}
	}
	return builder.BinaryPayloadBytes()
}

func (d TransactionScheduling) serializedFieldLengths() []int {
	switch {
	case d.Standard != nil:
		return []int{
			encoding.U8SerializedLength,
		}
	case d.FutureEra != nil:
		return []int{
			encoding.U8SerializedLength,
			encoding.U64SerializedLength,
		}
	case d.FutureTimestamp != nil:
		return []int{
			encoding.U8SerializedLength,
			encoding.U64SerializedLength,
		}
	default:
		return []int{}
	}
}
