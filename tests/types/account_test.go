package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types"
)

func Test_Account_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/account/account_example.json")
	assert.NoError(t, err)

	var account types.Account
	err = json.Unmarshal(fixture, &account)
	require.NoError(t, err)

	result, err := json.Marshal(account)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}
