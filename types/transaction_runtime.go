package types

import (
	"github.com/make-software/casper-go-sdk/v2/types/serialization"
	"github.com/make-software/casper-go-sdk/v2/types/serialization/encoding"
)

// TransactionRuntime Package transaction types.
type TransactionRuntime string

const (
	TransactionRuntimeTagVmCasperV1 = iota
	TransactionRuntimeTagVmCasperV2
)

const (
	TransactionRuntimeVmCasperV1 TransactionRuntime = "VmCasperV1"
	TransactionRuntimeVmCasperV2 TransactionRuntime = "VmCasperV2"
)

func (t TransactionRuntime) RuntimeTag() byte {
	if t == TransactionRuntimeVmCasperV1 {
		return TransactionRuntimeTagVmCasperV1
	} else if t == TransactionRuntimeVmCasperV2 {
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

func (t TransactionRuntime) serializedFieldLengths() []int {
	return []int{
		encoding.U8SerializedLength,
	}
}
