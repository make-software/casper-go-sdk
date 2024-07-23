package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types"
)

func Test_AddressableEntity_Account_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/addressable_entity/addressable_entity_example.json")
	assert.NoError(t, err)

	var account types.AddressableEntity
	err = json.Unmarshal(fixture, &account)
	require.NoError(t, err)

	result, err := json.Marshal(account)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}
