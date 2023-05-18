package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/make-software/casper-go-sdk/types"
)

func Test_ContractPackage_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/contract/contract_package_example.json")
	assert.NoError(t, err)

	var contractPackage types.ContractPackage
	err = json.Unmarshal(fixture, &contractPackage)
	assert.NoError(t, err)

	result, err := json.Marshal(contractPackage)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}
