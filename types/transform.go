package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/key"
)

// Transform is an enumeration of transformation types used in the execution of a `transaction` for V2 version.
type Transform struct {
	Key  key.Key       `json:"key"`
	Kind TransformKind `json:"kind"`
}

// TransformKey is an enumeration of transformation types used in the execution of a `deploy`.
type TransformKey struct {
	Key       key.Key       `json:"key"`
	Transform TransformKind `json:"transform"`
}

type NamedKeyKind struct {
	NamedKey Argument `json:"named_key"`
	Name     Argument `json:"name"`
}

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

type TransformKind json.RawMessage

// UnmarshalJSON sets *m to a copy of data.
func (t *TransformKind) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*t = append((*t)[0:0], data...)
	return nil
}

// MarshalJSON returns m as the JSON encoding of m.
func (t TransformKind) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("null"), nil
	}
	return t, nil
}

func (t *TransformKind) IsWriteTransfer() bool {
	return strings.Contains(string(*t), "WriteTransfer")
}

func (t *TransformKind) ParseAsWriteTransfer() (*WriteTransfer, error) {
	type RawWriteTransferTransform struct {
		WriteTransfer *WriteTransfer `json:"WriteTransfer"`
	}

	jsonRes := RawWriteTransferTransform{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	if jsonRes.WriteTransfer == nil {
		return nil, errors.New("error: empty response")
	}

	return jsonRes.WriteTransfer, nil
}

func (t *TransformKind) IsWriteAccount() bool {
	return strings.Contains(string(*t), "WriteAccount")
}

func (t *TransformKind) ParseAsWriteAccount() (key.AccountHash, error) {
	type RawWriteAccountTransform struct {
		WriteAccount key.AccountHash `json:"WriteAccount"`
	}

	jsonRes := RawWriteAccountTransform{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return key.AccountHash{}, err
	}

	return jsonRes.WriteAccount, nil
}

func (t *TransformKind) IsWriteContract() bool {
	return bytes.Equal(*t, []byte("\"WriteContract\""))
}

func (t *TransformKind) IsWriteWithdraw() bool {
	return strings.Contains(string(*t), "WriteWithdraw")
}

func (t *TransformKind) IsWriteUnbonding() bool {
	return strings.Contains(string(*t), "WriteUnbonding")
}

func (t *TransformKind) IsWriteCLValue() bool {
	return bytes.Contains(*t, []byte("CLValue"))
}

func (t *TransformKind) IsWritePackage() bool {
	return bytes.Contains(*t, []byte("\"Package\""))
}

func (t *TransformKind) IsWriteAddressableEntity() bool {
	return bytes.Contains(*t, []byte("\"AddressableEntity\""))
}

func (t *TransformKind) IsWriteBidKind() bool {
	return bytes.Contains(*t, []byte("\"BidKind\""))
}

func (t *TransformKind) IsWriteNamedKey() bool {
	return bytes.Contains(*t, []byte("\"NamedKey\""))
}

func (t *TransformKind) IsWriteMessage() bool {
	return bytes.Contains(*t, []byte("\"Message\""))
}

func (t *TransformKind) IsWriteMessageTopic() bool {
	return bytes.Contains(*t, []byte("\"MessageTopic\""))
}

func (t *TransformKind) IsWriteBid() bool {
	return strings.Contains(string(*t), "WriteBid")
}

func (t *TransformKind) IsAddUint512() bool {
	return strings.Contains(string(*t), "AddUInt512")
}

func (t *TransformKind) IsWriteDeployInfo() bool {
	return strings.Contains(string(*t), "WriteDeployInfo")
}

func (t *TransformKind) ParseAsWriteWithdraws() ([]UnbondingPurse, error) {
	type RawWriteWithdrawals struct {
		UnbondingPurses []UnbondingPurse `json:"WriteWithdraw"`
	}

	jsonRes := RawWriteWithdrawals{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	return jsonRes.UnbondingPurses, nil
}

func (t *TransformKind) ParseAsWriteAddressableEntity() (*AddressableEntity, error) {
	type rawData struct {
		Write *struct {
			AddressableEntity *AddressableEntity `json:"AddressableEntity"`
		} `json:"Write"`
	}

	jsonRes := rawData{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	if jsonRes.Write == nil || jsonRes.Write.AddressableEntity == nil {
		return nil, errors.New("error: empty response")
	}

	return jsonRes.Write.AddressableEntity, nil
}

func (t *TransformKind) ParseAsWritePackage() (*Package, error) {
	type rawData struct {
		Write *struct {
			Package *Package `json:"Package"`
		} `json:"Write"`
	}

	jsonRes := rawData{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	if jsonRes.Write == nil || jsonRes.Write.Package == nil {
		return nil, errors.New("error: empty response")
	}

	return jsonRes.Write.Package, nil
}

func (t *TransformKind) ParseAsWriteBidKind() (*BidKind, error) {
	type rawData struct {
		Write *struct {
			BidKind *BidKind `json:"BidKind"`
		} `json:"Write"`
	}

	jsonRes := rawData{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	if jsonRes.Write == nil || jsonRes.Write.BidKind == nil {
		return nil, errors.New("error: empty response")
	}

	return jsonRes.Write.BidKind, nil
}

func (t *TransformKind) ParseAsWriteNamedKey() (*NamedKeyKind, error) {
	type rawData struct {
		Write *struct {
			NamedKey *NamedKeyKind `json:"NamedKey"`
		} `json:"Write"`
	}

	jsonRes := rawData{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	if jsonRes.Write == nil || jsonRes.Write.NamedKey == nil {
		return nil, errors.New("error: empty response")
	}

	return jsonRes.Write.NamedKey, nil
}

func (t *TransformKind) ParseAsWriteMessage() (*MessageChecksum, error) {
	type rawData struct {
		Write *struct {
			Message *MessageChecksum `json:"Message"`
		} `json:"Write"`
	}

	jsonRes := rawData{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	if jsonRes.Write == nil || jsonRes.Write.Message == nil {
		return nil, errors.New("error: empty response")
	}

	return jsonRes.Write.Message, nil
}

func (t *TransformKind) ParseAsWriteMessageTopic() (*MessageTopicSummary, error) {
	type rawData struct {
		Write *struct {
			MessageTopic *MessageTopicSummary `json:"MessageTopic"`
		} `json:"Write"`
	}

	jsonRes := rawData{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	if jsonRes.Write == nil || jsonRes.Write.MessageTopic == nil {
		return nil, errors.New("error: empty response")
	}

	return jsonRes.Write.MessageTopic, nil
}

func (t *TransformKind) ParseAsWriteUnbondings() ([]UnbondingPurse, error) {
	type RawWriteUnbondings struct {
		UnbondingPurses []UnbondingPurse `json:"WriteUnbonding"`
	}

	jsonRes := RawWriteUnbondings{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	return jsonRes.UnbondingPurses, nil
}

func (t *TransformKind) ParseAsUInt512() (*clvalue.UInt512, error) {
	type RawUInt512 struct {
		UInt512 clvalue.UInt512 `json:"AddUInt512"`
	}

	jsonRes := RawUInt512{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	return &jsonRes.UInt512, nil
}

func (t *TransformKind) ParseAsWriteDeployInfo() (*DeployInfo, error) {
	type RawWriteDeployInfo struct {
		WriteDeployInfo *DeployInfo `json:"WriteDeployInfo"`
	}

	jsonRes := RawWriteDeployInfo{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	if jsonRes.WriteDeployInfo == nil {
		return nil, errors.New("error: empty response")
	}

	return jsonRes.WriteDeployInfo, nil
}

func (t *TransformKind) ParseAsWriteCLValue() (*Argument, error) {
	type RawWriteCLValue struct {
		WriteCLValue *Argument `json:"WriteCLValue"`
	}

	jsonRes := RawWriteCLValue{}
	err := json.Unmarshal(*t, &jsonRes)
	if err == nil && jsonRes.WriteCLValue != nil {
		return jsonRes.WriteCLValue, nil
	}

	type RawWriteCLValueV2 struct {
		Write *struct {
			CLValue *Argument `json:"CLValue"`
		} `json:"Write"`
	}

	jsonResV2 := RawWriteCLValueV2{}
	err = json.Unmarshal(*t, &jsonResV2)
	if err == nil && jsonResV2.Write != nil {
		return jsonResV2.Write.CLValue, nil
	}

	return nil, err
}
