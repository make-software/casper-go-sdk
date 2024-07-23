package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/key"
)

type TransformKind json.RawMessage

// UnmarshalJSON sets *m to a copy of data.
func (t *TransformKind) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*t = append((*t)[0:0], data...)
	return nil
}

// TransformV2 is an enumeration of transformation types used in the execution of a `deploy` for V2 version.
type TransformV2 struct {
	Key  key.Key       `json:"key"`
	Kind TransformKind `json:"kind"`
}

// TransformKey is an enumeration of transformation types used in the execution of a `deploy`.
type TransformKey struct {
	Key       key.Key   `json:"key"`
	Transform Transform `json:"transform"`
}

type Transform json.RawMessage

type WriteTransfer struct {
	ID         *uint64          `json:"id"`
	To         *key.AccountHash `json:"to"`
	DeployHash key.Hash         `json:"deploy_hash"`
	From       key.AccountHash  `json:"from"`
	Amount     clvalue.UInt512  `json:"amount"`
	Source     key.URef         `json:"source"`
	Target     key.URef         `json:"target"`
	Gas        uint             `json:"gas,string"`
}

// MarshalJSON returns m as the JSON encoding of m.
func (t Transform) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("null"), nil
	}
	return t, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (t *Transform) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*t = append((*t)[0:0], data...)
	return nil
}

func (t *Transform) IsWriteTransfer() bool {
	return strings.Contains(string(*t), "WriteTransfer")
}

func (t *Transform) ParseAsWriteTransfer() (*WriteTransfer, error) {
	type RawWriteTransferTransform struct {
		WriteTransfer `json:"WriteTransfer"`
	}

	jsonRes := RawWriteTransferTransform{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	return &jsonRes.WriteTransfer, nil
}

func (t *Transform) IsWriteAccount() bool {
	return strings.Contains(string(*t), "WriteAccount")
}

func (t *Transform) ParseAsWriteAccount() (key.AccountHash, error) {
	type RawWriteAccountTransform struct {
		WriteAccount key.AccountHash `json:"WriteAccount"`
	}

	jsonRes := RawWriteAccountTransform{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return key.AccountHash{}, err
	}

	return jsonRes.WriteAccount, nil
}

func (t *Transform) IsWriteContract() bool {
	return bytes.Equal(*t, []byte("\"WriteContract\""))
}

func (t *Transform) IsWriteWithdraw() bool {
	return strings.Contains(string(*t), "WriteWithdraw")
}

func (t *Transform) IsWriteUnbonding() bool {
	return strings.Contains(string(*t), "WriteUnbonding")
}

func (t *Transform) IsWriteCLValue() bool {
	return bytes.Contains(*t, []byte("\"WriteCLValue\""))
}

func (t *Transform) IsWriteBid() bool {
	return strings.Contains(string(*t), "WriteBid")
}

func (t *Transform) IsAddUint512() bool {
	return strings.Contains(string(*t), "AddUInt512")
}

func (t *Transform) IsWriteDeployInfo() bool {
	return strings.Contains(string(*t), "WriteDeployInfo")
}

func (t *Transform) ParseAsWriteWithdraws() ([]UnbondingPurse, error) {
	type RawWriteWithdrawals struct {
		UnbondingPurses []UnbondingPurse `json:"WriteWithdraw"`
	}

	jsonRes := RawWriteWithdrawals{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	return jsonRes.UnbondingPurses, nil
}

func (t *Transform) ParseAsWriteUnbondings() ([]UnbondingPurse, error) {
	type RawWriteUnbondings struct {
		UnbondingPurses []UnbondingPurse `json:"WriteUnbonding"`
	}

	jsonRes := RawWriteUnbondings{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	return jsonRes.UnbondingPurses, nil
}

func (t *Transform) ParseAsWriteCLValue() (*Argument, error) {
	type RawWriteCLValue struct {
		WriteCLValue Argument `json:"WriteCLValue"`
	}

	jsonRes := RawWriteCLValue{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	return &jsonRes.WriteCLValue, nil
}

func (t *Transform) ParseAsUInt512() (*clvalue.UInt512, error) {
	type RawUInt512 struct {
		UInt512 clvalue.UInt512 `json:"AddUInt512"`
	}

	jsonRes := RawUInt512{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	return &jsonRes.UInt512, nil
}

func (t *Transform) ParseAsWriteDeployInfo() (*DeployInfo, error) {
	type RawWriteDeployInfo struct {
		WriteDeployInfo DeployInfo `json:"WriteDeployInfo"`
	}

	jsonRes := RawWriteDeployInfo{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	return &jsonRes.WriteDeployInfo, nil
}
