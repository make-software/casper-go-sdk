package types

import (
	"encoding/json"
	"errors"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/key"
)

var ErrNamedKeyNotFound = errors.New("NamedKey not found")

// NamedKey is a key in an Account or Contract.
type NamedKey struct {
	// The name of the entry.
	Name string `json:"name"`
	// The value of the entry: a casper `key.Key` type.
	Key key.Key `json:"key"`
}

// NamedKeyValue A NamedKey value.
type NamedKeyValue struct {
	// The name of the `Key` encoded as a CLValue.
	Name clvalue.CLValue `json:"name"`
	// The actual `Key` encoded as a CLValue.
	NamedKey clvalue.CLValue `json:"named_key"`
}

func (t *NamedKeyValue) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}

	raw := struct {
		Name     Argument `json:"name"`
		NamedKey Argument `json:"named_key"`
	}{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	valueName, err := raw.Name.Value()
	if err != nil {
		return err
	}

	valueKey, err := raw.NamedKey.Value()
	if err != nil {
		return err
	}

	*t = NamedKeyValue{
		Name:     valueName,
		NamedKey: valueKey,
	}
	return nil
}

type NamedKeys []NamedKey

func (k NamedKeys) ToMap() map[string]string {
	result := make(map[string]string, len(k))
	for _, nk := range k {
		result[nk.Name] = nk.Key.String()
	}
	return result
}

func (k NamedKeys) Find(target string) (key.Key, error) {
	for _, one := range k {
		if one.Name == target {
			return one.Key, nil
		}
	}

	return key.Key{}, ErrNamedKeyNotFound
}
