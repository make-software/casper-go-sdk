package types

import (
	"encoding/json"
	"testing"

	sdk_types "github.com/make-software/casper-go-sdk/v2/types"
	"github.com/stretchr/testify/require"
)

func TestContractVersionKeyUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected sdk_types.ContractVersionKey
	}{
		{
			name:     "array",
			input:    `[2, 7]`,
			expected: sdk_types.ContractVersionKey{2, 7},
		},
		{
			name:     "object",
			input:    `{"contract_version": 7, "protocol_version_major": 2}`,
			expected: sdk_types.ContractVersionKey{2, 7},
		},
		{
			name:     "zero values",
			input:    `{"contract_version": 0, "protocol_version_major": 0}`,
			expected: sdk_types.ContractVersionKey{0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual sdk_types.ContractVersionKey
			require.NoError(t, json.Unmarshal([]byte(tt.input), &actual))
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestContractVersionKeyUnmarshalJSONRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "short array", input: `[2]`},
		{name: "long array", input: `[2, 7, 8]`},
		{name: "missing contract version", input: `{"protocol_version_major": 2}`},
		{name: "missing protocol version", input: `{"contract_version": 7}`},
		{name: "invalid JSON type", input: `"2,7"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var key sdk_types.ContractVersionKey
			require.Error(t, json.Unmarshal([]byte(tt.input), &key))
		})
	}
}
