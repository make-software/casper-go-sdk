package types

import (
	"errors"

	"github.com/make-software/casper-go-sdk/types/key"
)

var ErrNamedKeyNotFound = errors.New("NamedKey not found")

// NamedKey is a key in an Account or Contract.
type NamedKey struct {
	// The name of the entry.
	Name string `json:"name"`
	// The value of the entry: a casper `key.Key` type.
	Key key.Key `json:"key"`
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
