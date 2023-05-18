package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

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
