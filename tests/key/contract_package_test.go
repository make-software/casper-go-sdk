package key

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/make-software/casper-go-sdk/v2/casper"
)

func Test_ContractPackage_Constructor(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{"From hash", "ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"From universal prefixed hash", "hash-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"From prefixed hash", "contract-package-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"From legacy prefixed hash", "contract-package-wasmee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := casper.NewContractPackageHash(test.source)
			assert.NoError(t, err)
			assert.Equal(t, "ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", result.ToHex())
		})
	}
}

func Test_ContractPackage_ToPrefixedString(t *testing.T) {
	result, err := casper.NewContractPackageHash("ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760")
	assert.NoError(t, err)
	assert.Equal(t, "contract-package-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", result.ToPrefixedString())
}

func Test_ContractPackage_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		excepted string
	}{
		{"From hash", `"ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"`, `"ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"`},
		{"From universal prefixed hash", `"hash-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"`, `"hash-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"`},
		{"From prefixed hash", `"contract-package-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"`, `"contract-package-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"`},
		{"From legacy prefixed hash", `"contract-package-wasmee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"`, `"contract-package-wasmee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var data casper.ContractPackageHash
			err := json.Unmarshal([]byte(test.source), &data)
			assert.NoError(t, err)
			actual, err := json.Marshal(data)
			assert.NoError(t, err)
			assert.Equal(t, test.excepted, string(actual))
		})
	}
}
