package types

import "errors"

// TransactionRuntime SmartContract transaction types.
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

type TransactionRuntimeFromBytesDecoder struct{}

func (addr *TransactionRuntimeFromBytesDecoder) FromBytes(data []byte) (TransactionRuntime, []byte, error) {
	if len(data) < 1 {
		return "", nil, errors.New("insufficient bytes to decode TransactionRuntime")
	}

	tag := data[0]
	remainder := data[1:]

	switch tag {
	case uint8(TransactionRuntimeTagVmCasperV1):
		return TransactionRuntimeVmCasperV1, remainder, nil
	case uint8(TransactionRuntimeTagVmCasperV2):
		return TransactionRuntimeVmCasperV2, remainder, nil
	default:
		return "", nil, errors.New("unknown TransactionRuntime variant")
	}
}
