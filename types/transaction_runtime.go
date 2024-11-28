package types

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
