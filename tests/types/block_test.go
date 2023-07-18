package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/types"
)

func Test_Block_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/block/block_example.json")
	assert.NoError(t, err)

	var block types.Block
	err = json.Unmarshal(fixture, &block)
	assert.NoError(t, err)

	result, err := json.Marshal(block)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}

func Test_BlockSwitch_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/block/block_switch_example.json")
	assert.NoError(t, err)

	var block types.Block
	err = json.Unmarshal(fixture, &block)
	assert.NoError(t, err)

	result, err := json.Marshal(block)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}

func Test_Block_Proposer_PublicKey_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/block/block_with_system_proposer.json")
	assert.NoError(t, err)

	var block casper.ChainGetBlockResult
	err = json.Unmarshal(fixture, &block)
	assert.NoError(t, err)

	// "00" in hex
	assert.True(t, len(block.Block.Body.Proposer) == 1)

	result, err := json.Marshal(block)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}
