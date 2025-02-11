package types

import (
	"encoding/json"
	"errors"

	"github.com/make-software/casper-go-sdk/v2/types/key"
	"github.com/make-software/casper-go-sdk/v2/types/serialization"
	"github.com/make-software/casper-go-sdk/v2/types/serialization/encoding"
)

// TransactionRuntime Package transaction types.
const (
	TransactionRuntimeTagVmCasperV1 = iota
	TransactionRuntimeTagVmCasperV2
)

type TransactionRuntime struct {
	VmCasperV1 *struct{}
	// The execution target is a stored entity or package.
	VmCasperV2 *struct {
		Seed             *key.Hash `json:"seed,omitempty"`
		TransferredValue uint64    `json:"transferred_value"`
	} `json:"VmCasperV2,omitempty"`
}

func (t TransactionRuntime) IsVmCasperV1() bool {
	return t.VmCasperV1 != nil
}

func (t TransactionRuntime) IsVmCasperV2() bool {
	return t.VmCasperV2 != nil
}

func NewVmCasperV1TransactionRuntime() TransactionRuntime {
	return TransactionRuntime{
		VmCasperV1: &struct{}{},
	}
}

func NewVmCasperV2TransactionRuntime(transferredValue uint64, seed *key.Hash) TransactionRuntime {
	return TransactionRuntime{
		VmCasperV2: &struct {
			Seed             *key.Hash `json:"seed,omitempty"`
			TransferredValue uint64    `json:"transferred_value"`
		}{Seed: seed, TransferredValue: transferredValue},
	}
}

func (t TransactionRuntime) RuntimeTag() byte {
	if t.VmCasperV1 != nil {
		return TransactionRuntimeTagVmCasperV1
	} else if t.VmCasperV2 != nil {
		return TransactionRuntimeTagVmCasperV2
	}
	return 0
}

func (t *TransactionRuntime) Bytes() ([]byte, error) {
	builder, err := serialization.NewCallTableSerializationEnvelopeBuilder(t.serializedFieldLengths())
	if err != nil {
		return nil, err
	}
	if err = builder.AddField(TagFieldIndex, []byte{t.RuntimeTag()}); err != nil {
		return nil, err
	}
	return builder.BinaryPayloadBytes()
}

func (t *TransactionRuntime) SerializedLength() int {
	envelope := serialization.CallTableSerializationEnvelope{}
	return envelope.EstimateSize(t.serializedFieldLengths())
}

func (t TransactionRuntime) serializedFieldLengths() []int {
	if t.VmCasperV1 != nil {
		return []int{encoding.U8SerializedLength}
	} else if t.VmCasperV2 != nil {
		var seedSerializedLength int
		if t.VmCasperV2.Seed != nil {
			seedSerializedLength = encoding.StringSerializedLength(t.VmCasperV2.Seed.String())
		}

		return []int{
			encoding.U8SerializedLength,
			encoding.U64SerializedLength,
			encoding.U8SerializedLength + seedSerializedLength,
		}
	}
	return []int{}
}

func (t TransactionRuntime) MarshalJSON() ([]byte, error) {
	if t.VmCasperV1 != nil {
		return json.Marshal("VmCasperV1")
	}

	if t.VmCasperV2 != nil {
		return json.Marshal(t.VmCasperV2)
	}

	return nil, errors.New("unknown target runtime type")
}

func (t *TransactionRuntime) UnmarshalJSON(data []byte) error {
	var vmCasperV1 string
	if err := json.Unmarshal(data, &vmCasperV1); err == nil && vmCasperV1 == "VmCasperV1" {
		t.VmCasperV1 = &struct{}{}
		return nil
	}

	var runtime struct {
		VmCasperV2 *struct {
			Seed             *key.Hash `json:"seed,omitempty"`
			TransferredValue uint64    `json:"transferred_value"`
		} `json:"VmCasperV2"`
	}
	if err := json.Unmarshal(data, &runtime); err == nil && runtime.VmCasperV2 != nil {
		t.VmCasperV2 = runtime.VmCasperV2
		return nil
	}

	return errors.New("unknown target runtime type")
}
