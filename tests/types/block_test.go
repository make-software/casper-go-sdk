package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types"
)

func Test_Block_V1_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/block/block_v1_example.json")
	assert.NoError(t, err)

	var block types.BlockV1
	err = json.Unmarshal(fixture, &block)
	assert.NoError(t, err)

	result, err := json.Marshal(block)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}

func Test_Block_V2_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/block/block_v2_example.json")
	assert.NoError(t, err)

	var block types.BlockV2
	err = json.Unmarshal(fixture, &block)
	assert.NoError(t, err)

	result, err := json.Marshal(block)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}

func Test_BlockSwitch_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/block/block_switch_example.json")
	assert.NoError(t, err)

	var block types.BlockV1
	err = json.Unmarshal(fixture, &block)
	assert.NoError(t, err)

	result, err := json.Marshal(block)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}

func Test_BlockSwitch_WithSystemProposal_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/block/block_switch_system_proposer.json")
	assert.NoError(t, err)

	var block types.BlockV1
	err = json.Unmarshal(fixture, &block)
	assert.NoError(t, err)

	result, err := json.Marshal(block)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}

func Test_BlockSwitch_WithSystemProposal_IsSystem_ShouldReturnTrue(t *testing.T) {
	fixture, err := os.ReadFile("../data/block/block_switch_system_proposer.json")
	assert.NoError(t, err)

	var block types.BlockV1
	err = json.Unmarshal(fixture, &block)
	require.NoError(t, err)
	assert.True(t, block.Body.Proposer.IsSystem())
	_, err = block.Body.Proposer.PublicKey()
	assert.Error(t, err)
	pubKey := block.Body.Proposer.PublicKeyOptional()
	assert.Nil(t, pubKey)
}

func Test_BlockProposal_PublicKey_ShouldWorkForNormalBlock(t *testing.T) {
	fixture, err := os.ReadFile("../data/block/block_v1_example.json")
	assert.NoError(t, err)

	var block types.BlockV1
	err = json.Unmarshal(fixture, &block)
	require.NoError(t, err)
	assert.False(t, block.Body.Proposer.IsSystem())
	result, err := block.Body.Proposer.PublicKey()
	assert.NoError(t, err)
	assert.Equal(t, "019e7b8bdec03ba83be4f5443d9f7f9111c77fec984ce9bb5bb7eb3da1e689c02d", result.String())
	pubKey := block.Body.Proposer.PublicKeyOptional()
	assert.Equal(t, "019e7b8bdec03ba83be4f5443d9f7f9111c77fec984ce9bb5bb7eb3da1e689c02d", pubKey.String())
}
