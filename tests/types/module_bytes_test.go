package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/make-software/casper-go-sdk/types"
)

func Test_MarshalUnmarshalModuleBytes_ShouldBeSameResult(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			"item with some args",
			`{"module_bytes": "[72 bytes]","args": [["testName", "testVal"]]}`,
		},
		{
			"item with empty args",
			`{"module_bytes": "[72 bytes]","args": []}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var deployItem types.ModuleBytes
			err := json.Unmarshal([]byte(test.data), &deployItem)
			assert.NoError(t, err)

			result, err := json.Marshal(deployItem)
			assert.NoError(t, err)
			assert.JSONEq(t, test.data, string(result))
		})
	}
}

func Test_MarshalUnmarshalOmittedModuleBytes_ShouldBeSameResult(t *testing.T) {
	data := `{"module_bytes": "[72 bytes]"}`
	var deployItem types.ModuleBytes
	err := json.Unmarshal([]byte(data), &deployItem)
	assert.NoError(t, err)

	result, err := json.Marshal(deployItem)
	assert.NoError(t, err)
	assert.JSONEq(t, data, string(result))
}
