package key

import (
	"fmt"
	"github.com/make-software/casper-go-sdk/v2/types/key"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewContract_Creation(t *testing.T) {
	prefixes := []string{
		"",
		"hash-",
		"contract-",
		"contract-wasm-",
		"entity-contract-",
	}

	rawKey := "eb2909ad8a38239540db975ecbbaf3ea2a553c91ba9b9ab6302229301382541d"

	for _, prefix := range prefixes {
		contractHash, err := key.NewContract(fmt.Sprintf("%s%s", prefix, rawKey))
		require.NoError(t, err)
		assert.Equal(t, rawKey, contractHash.Hash.ToHex())
	}
}
