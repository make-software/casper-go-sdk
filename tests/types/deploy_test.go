package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types"
)

func Test_Deploy_MarshalUnmarshal_ShouldBeSameResult(t *testing.T) {
	tests := []struct {
		name        string
		fixturePath string
	}{
		{
			"deploy with StoredContractByName",
			"../data/deploy/deploy_with_stored_contract_by_name.json",
		},
		{
			"deploy with StoredContractByHash",
			"../data/deploy/deploy_with_stored_contract_by_hash.json",
		},
		{
			"deploy with StoredContractByHash with version",
			"../data/deploy/deploy_with_stored_contract_by_hash_with_version.json",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := os.ReadFile(test.fixturePath)
			require.NoError(t, err)

			var deploy types.Deploy
			err = json.Unmarshal(data, &deploy)
			require.NoError(t, err)

			result, err := json.Marshal(deploy)
			require.NoError(t, err)
			assert.JSONEq(t, string(data), string(result), test.name)
		})
	}
}
