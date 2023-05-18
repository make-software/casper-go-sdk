package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types"
)

func Test_AuctionState_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/auction/auction_state_example.json")
	require.NoError(t, err)

	var auction types.AuctionState
	err = json.Unmarshal(fixture, &auction)
	require.NoError(t, err)

	result, err := json.Marshal(auction)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}
