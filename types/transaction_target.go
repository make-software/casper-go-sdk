package types

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/make-software/casper-go-sdk/v2/types/key"
	"github.com/make-software/casper-go-sdk/v2/types/serialization"
	"github.com/make-software/casper-go-sdk/v2/types/serialization/encoding"
)

const (
	SessionIsInstallIndex uint16 = iota + 1
	SessionRuntimeIndex
	SessionModuleBytesIndex
	SessionTransferredValueIndex
	SessionSeedIndex
)

const (
	StoredIdIndex uint16 = iota + 1
	StoredRuntimeIndex
)

const (
	TransactionTargetTypeNative = iota
	TransactionTargetTypeStored
	TransactionTargetTypeSession
)

type TransactionTarget struct {
	// The execution target is a native operation (e.g. a transfer).
	Native *struct{}
	// The execution target is a stored entity or package.
	Stored *StoredTarget `json:"Stored"`
	// The execution target is the included module bytes, i.e. compiled Wasm.
	Session *SessionTarget `json:"Session"`
}

func (t *TransactionTarget) SerializedLength() int {
	envelope := serialization.CallTableSerializationEnvelope{}
	return envelope.EstimateSize(t.serializedFieldLengths())
}

type StoredTarget struct {
	ID      TransactionInvocationTarget `json:"id"`
	Runtime TransactionRuntime          `json:"runtime"`
}

type SessionTarget struct {
	ModuleBytes      []byte             `json:"module_bytes"`
	Runtime          TransactionRuntime `json:"runtime"`
	IsInstallUpgrade bool               `json:"is_install_upgrade"`
}

func (t *TransactionTarget) Bytes() ([]byte, error) {
	builder, err := serialization.NewCallTableSerializationEnvelopeBuilder(t.serializedFieldLengths())
	if err != nil {
		return nil, err
	}

	switch {
	case t.Native != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionTargetTypeNative}); err != nil {
			return nil, err
		}
	case t.Stored != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionTargetTypeStored}); err != nil {
			return nil, err
		}

		storedIDBytes, err := t.Stored.ID.Bytes()
		if err != nil {
			return nil, err
		}

		if err = builder.AddField(StoredIdIndex, storedIDBytes); err != nil {
			return nil, err
		}

		runtimeBytes, _ := encoding.NewStringToBytesEncoder(string(t.Stored.Runtime)).Bytes()
		if err = builder.AddField(StoredRuntimeIndex, runtimeBytes); err != nil {
			return nil, err
		}
	case t.Session != nil:
		if err = builder.AddField(TagFieldIndex, []byte{TransactionTargetTypeSession}); err != nil {
			return nil, err
		}

		IsInstallUpgradeBytes, _ := encoding.NewBoolToBytesEncoder(t.Session.IsInstallUpgrade).Bytes()
		if err = builder.AddField(SessionIsInstallIndex, IsInstallUpgradeBytes); err != nil {
			return nil, err
		}

		if err = builder.AddField(SessionRuntimeIndex, []byte{t.Session.Runtime.RuntimeTag()}); err != nil {
			return nil, err
		}

		moduleBytes, _ := encoding.NewBytesToBytesEncoder(t.Session.ModuleBytes).Bytes()
		if err = builder.AddField(SessionModuleBytesIndex, moduleBytes); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid TransactionTarget")
	}

	return builder.BinaryPayloadBytes()
}

func (t TransactionTarget) serializedFieldLengths() []int {
	switch {
	case t.Native != nil:
		return []int{
			encoding.U8SerializedLength,
		}
	case t.Stored != nil:
		return []int{
			encoding.U8SerializedLength,
			t.Stored.ID.SerializedLength(),
			encoding.U8SerializedLength,
			encoding.U64SerializedLength,
		}
	case t.Session != nil:
		return []int{
			encoding.U8SerializedLength,
			encoding.BoolSerializedLength,
			encoding.U8SerializedLength,
			encoding.BytesSerializedLength(t.Session.ModuleBytes),
		}
	default:
		return []int{}
	}
}

func (t *TransactionTarget) UnmarshalJSON(data []byte) error {
	var target struct {
		Stored  *StoredTarget `json:"Stored"`
		Session *struct {
			Runtime          TransactionRuntime `json:"runtime"`
			IsInstallUpgrade bool               `json:"is_install_upgrade"`
			Module           string             `json:"module_bytes"`
		} `json:"Session"`
	}
	if err := json.Unmarshal(data, &target); err == nil {
		if target.Session != nil {
			decodedBytes, err := hex.DecodeString(target.Session.Module)
			if err != nil {
				return err
			}

			*t = TransactionTarget{
				Session: &SessionTarget{
					ModuleBytes:      decodedBytes,
					Runtime:          target.Session.Runtime,
					IsInstallUpgrade: target.Session.IsInstallUpgrade,
				},
			}
		}

		if target.Stored != nil {
			*t = TransactionTarget{
				Stored: target.Stored,
			}
		}
		return nil
	}

	var key string
	if err := json.Unmarshal(data, &key); err == nil && key == "Native" {
		*t = TransactionTarget{
			Native: &struct{}{},
		}
		return nil
	}

	return nil
}

func (t TransactionTarget) MarshalJSON() ([]byte, error) {
	if t.Native != nil {
		return json.Marshal("Native")
	}

	if t.Stored != nil {
		return json.Marshal(struct {
			Stored *StoredTarget `json:"Stored"`
		}{
			Stored: t.Stored,
		})
	}

	if t.Session != nil {
		type sessionTarget struct {
			Runtime          TransactionRuntime `json:"runtime"`
			TransferredValue uint64             `json:"transferred_value"`
			IsInstallUpgrade bool               `json:"is_install_upgrade"`
			Seed             *key.Hash          `json:"seed,omitempty"`
			ModuleBytes      string             `json:"module_bytes"`
		}

		return json.Marshal(struct {
			Session sessionTarget `json:"Session"`
		}{
			Session: sessionTarget{
				Runtime:          t.Session.Runtime,
				IsInstallUpgrade: t.Session.IsInstallUpgrade,
				ModuleBytes:      hex.EncodeToString(t.Session.ModuleBytes),
			},
		})
	}

	return nil, errors.New("unknown target type")
}

// NewTransactionTargetFromSession create new TransactionTarget from ExecutableDeployItem
func NewTransactionTargetFromSession(session ExecutableDeployItem) TransactionTarget {
	if session.Transfer != nil {
		return TransactionTarget{
			Native: &struct{}{},
		}
	}

	if session.ModuleBytes != nil {
		decodedBytes, _ := hex.DecodeString(session.ModuleBytes.ModuleBytes)
		return TransactionTarget{
			Session: &SessionTarget{
				ModuleBytes: decodedBytes,
				Runtime:     "VmCasperV1",
			},
		}
	}

	if session.StoredContractByHash != nil {
		hash := session.StoredContractByHash.Hash.Hash
		return TransactionTarget{
			Stored: &StoredTarget{
				ID: TransactionInvocationTarget{
					ByHash: &hash,
				},
				Runtime: "VmCasperV1",
			},
		}
	}

	if session.StoredContractByName != nil {
		return TransactionTarget{
			Stored: &StoredTarget{
				ID: TransactionInvocationTarget{
					ByName: &session.StoredContractByName.Name,
				},
				Runtime: "VmCasperV1",
			},
		}
	}

	if session.StoredVersionedContractByHash != nil {
		var version *uint32
		if storedVersion := session.StoredVersionedContractByHash.Version; storedVersion != nil {
			versionNum, err := storedVersion.Int64()
			if err == nil {
				temp := uint32(versionNum)
				version = &temp
			}
		}
		byHashTarget := ByPackageHashInvocationTarget{
			Addr:    session.StoredVersionedContractByHash.Hash.Hash,
			Version: version,
		}
		return TransactionTarget{
			Stored: &StoredTarget{
				ID: TransactionInvocationTarget{
					ByPackageHash: &byHashTarget,
				},
				Runtime: "VmCasperV1",
			},
		}
	}

	if session.StoredVersionedContractByName != nil {
		var version *uint32
		if storedVersion := session.StoredVersionedContractByName.Version; storedVersion != nil {
			versionNum, err := storedVersion.Int64()
			if err == nil {
				temp := uint32(versionNum)
				version = &temp
			}
		}
		byNameTarget := ByPackageNameInvocationTarget{
			Version: version,
		}

		if session.StoredVersionedContractByName != nil {
			byNameTarget.Name = session.StoredVersionedContractByName.Name
		}

		return TransactionTarget{
			Stored: &StoredTarget{
				ID: TransactionInvocationTarget{
					ByPackageName: &byNameTarget,
				},
				Runtime: "VmCasperV1",
			},
		}
	}

	return TransactionTarget{}
}
