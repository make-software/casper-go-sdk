package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types"
)

func Test_ExecutableItem_MarshalUnmarshal_ShouldBeSameResult(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			"item with module bytes",
			`{"ModuleBytes": {"module_bytes": "[72 bytes]","args": [["testName", "testVal"]]}}`,
		},
		{
			"item with stored contract by hash",
			`{"StoredContractByHash": {"hash": "c4c411864f7b717c27839e56f6f1ebe5da3f35ec0043f437324325d65a22afa4","entry_point": "pclphXwfYmCmdITj8hnh","args": [["testName", "testVal"]]}}`,
		},
		{
			"item with stored contract by name",
			`{"StoredContractByName": {"name": "U5A74bSZH8abT8HqVaK9","entry_point": "gIetSxltnRDvMhWdxTqQ","args": [["testName", "testVal"]]}}`,
		},
		{
			"item with stored versioned contract by hash",
			`{"StoredContractByHash": {"hash": "c4c411864f7b717c27839e56f6f1ebe5da3f35ec0043f437324325d65a22afa4","entry_point": "pclphXwfYmCmdITj8hnh","args": [["testName", "testVal"]]}}`,
		},
		{
			"item with stored versioned contract by name",
			`{"StoredVersionedContractByName": {"name": "lWJWKdZUEudSakJzw1tn","version": 1632552656, "entry_point": "S1cXRT3E1jyFlWBAIVQ8","args": [["testName", "testVal"]]}}`,
		},
		{
			"item with stored transfer",
			`{"Transfer": {"args": [["testName", "testVal"]]}}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var deployItem types.ExecutableDeployItem
			err := json.Unmarshal([]byte(test.data), &deployItem)
			require.NoError(t, err)

			result, err := json.Marshal(deployItem)
			require.NoError(t, err)
			assert.JSONEq(t, test.data, string(result))
		})
	}
}
