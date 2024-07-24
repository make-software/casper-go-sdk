package types

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/key"
)

const (
	TransactionTargetTypeNative = iota
	TransactionTargetTypeStored
	TransactionTargetTypeSession
)

const (
	InvocationTargetTagByHash = iota
	InvocationTargetTagByName
	InvocationTargetTagByPackageHash
	InvocationTargetTagByPackageName
)

type TransactionTarget struct {
	// The execution target is a native operation (e.g. a transfer).
	Native *struct{}
	// The execution target is a stored entity or package.
	Stored *StoredTarget `json:"Stored"`
	// The execution target is the included module bytes, i.e. compiled Wasm.
	Session *SessionTarget `json:"Session"`
}

func (t *TransactionTarget) Bytes() ([]byte, error) {
	result := make([]byte, 0, 8)

	if t.Native != nil {
		result = append(result, TransactionTargetTypeNative)
	} else if t.Stored != nil {
		result = append(result, TransactionTargetTypeStored)
		if byHash := t.Stored.ID.ByHash; byHash != nil {
			result = append(result, InvocationTargetTagByHash)
			result = append(result, byHash.Bytes()...)
		} else if byName := t.Stored.ID.ByName; byName != nil {
			result = append(result, InvocationTargetTagByName)
			result = append(result, clvalue.NewCLString(*byName).Bytes()...)
		} else if byPackageHash := t.Stored.ID.ByPackageHash; byPackageHash != nil {
			result = append(result, InvocationTargetTagByPackageHash)
			result = append(result, byPackageHash.Addr.Bytes()...)
			if byPackageHash.Version != nil {
				result = append(result, 1)
				result = append(result, clvalue.NewCLUInt32(*byPackageHash.Version).Bytes()...)
			} else {
				result = append(result, 0)
			}
		} else if byPackageName := t.Stored.ID.ByPackageName; byPackageName != nil {
			result = append(result, InvocationTargetTagByPackageName)
			result = append(result, clvalue.NewCLString(byPackageName.Name).Bytes()...)
			if byPackageHash.Version != nil {
				result = append(result, 1)
				result = append(result, clvalue.NewCLUInt32(*byPackageName.Version).Bytes()...)
			} else {
				result = append(result, 0)
			}
		}
		result = append(result, t.Stored.Runtime.RuntimeTag())
	} else if t.Session != nil {
		result = append(result, TransactionTargetTypeSession)
		if len(t.Session.ModuleBytes) == 0 {
			result = append(result, clvalue.NewCLInt32(0).Bytes()...)
		} else {
			bytes, err := hex.DecodeString(t.Session.ModuleBytes)
			if err != nil {
				return nil, err
			}
			result = append(result, clvalue.NewCLUInt32(uint32(len(bytes))).Bytes()...)
			result = append(result, clvalue.NewCLByteArray(bytes).Bytes()...)
		}

		result = append(result, t.Session.Runtime.RuntimeTag())
	}

	return result, nil
}

func (t *TransactionTarget) UnmarshalJSON(data []byte) error {
	var target struct {
		Stored  *StoredTarget  `json:"Stored"`
		Session *SessionTarget `json:"Session"`
	}
	if err := json.Unmarshal(data, &target); err == nil {
		if target.Session != nil {
			*t = TransactionTarget{
				Session: target.Session,
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
		return json.Marshal(struct {
			Session *SessionTarget `json:"Session"`
		}{
			Session: t.Session,
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
		return TransactionTarget{
			Session: &SessionTarget{
				ModuleBytes: session.ModuleBytes.ModuleBytes,
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

type StoredTarget struct {
	// Identifier of a `Stored` transaction target.
	ID      TransactionInvocationTarget `json:"id"`
	Runtime TransactionRuntime          `json:"runtime"`
}

type TransactionInvocationTarget struct {
	// Hex-encoded entity address identifying the invocable entity.
	ByHash *key.Hash `json:"ByHash"`
	// The alias identifying the invocable entity.
	ByName *string `json:"ByName"`
	// The address and optional version identifying the package.
	ByPackageHash *ByPackageHashInvocationTarget `json:"ByPackageHash"`
	// The alias and optional version identifying the package.
	ByPackageName *ByPackageNameInvocationTarget `json:"ByPackageName"`
}

// ByPackageHashInvocationTarget The address and optional version identifying the package.
type ByPackageHashInvocationTarget struct {
	Addr    key.Hash `json:"addr"`
	Version *uint32  `json:"version"`
}

// ByPackageNameInvocationTarget The alias and optional version identifying the package.
type ByPackageNameInvocationTarget struct {
	Name    string  `json:"name"`
	Version *uint32 `json:"version"`
}

type SessionTarget struct {
	ModuleBytes string             `json:"module_bytes"`
	Runtime     TransactionRuntime `json:"runtime"`
}
