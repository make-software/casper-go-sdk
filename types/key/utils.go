package key

import (
	"errors"
)

func NewHashFromString[T any](source string) (T, error) {
	var res any
	var k T
	var err error
	switch any(k).(type) {
	case Hash:
		res, err = NewHash(source)
	case URef:
		res, err = NewURef(source)
	case AccountHash:
		res, err = NewAccountHash(source)
	case TransferHash:
		res, err = NewTransferHash(source)
	case ContractHash:
		res, err = NewContract(source)
	case ContractPackageHash:
		res, err = NewContractPackage(source)
	default:
		err = errors.New("type is not found")
	}

	return res.(T), err
}

func StringsToHashList[T any](source []string) ([]T, error) {
	res := make([]T, 0, len(source))
	for _, one := range source {
		hash, err := NewHashFromString[T](one)
		if err != nil {
			return nil, err
		}
		res = append(res, hash)
	}

	return res, nil
}
