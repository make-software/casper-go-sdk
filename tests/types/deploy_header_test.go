package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/make-software/casper-go-sdk/types"
)

func Test_DeployHeader_MarshalUnmarshal_ShouldBeSameResult(t *testing.T) {
	tests := []struct {
		name        string
		fixturePath string
	}{
		{
			"deploy with StoredContractByName",
			"../data/deploy/deploy_header_with_deps.json",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := os.ReadFile(test.fixturePath)
			assert.NoError(t, err)

			var deploy types.DeployHeader
			err = json.Unmarshal(data, &deploy)
			assert.NoError(t, err)

			result, err := json.Marshal(deploy)
			assert.NoError(t, err)
			assert.JSONEq(t, string(data), string(result), test.name)
		})
	}
}
