package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/make-software/casper-go-sdk/types"
)

func Test_Contract_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/contract/contract_example.json")
	assert.NoError(t, err)

	var contract types.Contract
	err = json.Unmarshal(fixture, &contract)
	assert.NoError(t, err)

	result, err := json.Marshal(contract)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}
