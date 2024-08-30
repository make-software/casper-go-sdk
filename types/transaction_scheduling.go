package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
)

const (
	TransactionSchedulingNativeTag = iota
	TransactionSchedulingFutureEraTag
	TransactionSchedulingFutureTimestampTag
)

type TransactionScheduling struct {
	// No special scheduling applied.
	Standard *struct{} `json:"Standard,omitempty"`
	// Execution should be scheduled for the specified era.
	FutureEra *FutureEraScheduling `json:"FutureTimestamp,omitempty"`
	// Execution should be scheduled for the specified timestamp or later.
	FutureTimestamp *FutureTimestampScheduling `json:"FutureEra,omitempty"`
}

func (t *TransactionScheduling) Tag() byte {
	switch {
	case t.Standard != nil:
		return TransactionSchedulingNativeTag
	case t.FutureEra != nil:
		return TransactionSchedulingFutureEraTag
	case t.FutureTimestamp != nil:
		return TransactionSchedulingFutureTimestampTag
	default:
		return 0
	}
}

func (t *TransactionScheduling) Bytes() []byte {
	result := make([]byte, 0, 2)
	result = append(result, t.Tag())

	if t.FutureEra != nil {
		result = append(result, clvalue.NewCLUInt64(t.FutureEra.EraID).Bytes()...)
	} else if t.FutureTimestamp != nil {
		result = append(result, clvalue.NewCLUInt64(uint64(time.Time(t.FutureTimestamp.TimeStamp).UnixMilli())).Bytes()...)
	}
	return result
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
	if err := json.Unmarshal(data, &key); err != nil {
		return err
	}
	if key == "Standard" {
		*t = TransactionScheduling{
			Standard: &struct{}{},
		}
		return nil
	}

	return fmt.Errorf("unknown transaction scheduling type: %s", key)
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
