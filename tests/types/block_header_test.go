package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/make-software/casper-go-sdk/types"
)

func Test_BlockHeader_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	//TODO: Timestamp format has rounded up in some cases not equal to specification
	t.Skip("Timestamp format is different")
	fixture, err := os.ReadFile("../data/block/block_header.json")
	assert.NoError(t, err)

	var block types.BlockHeader
	err = json.Unmarshal(fixture, &block)
	assert.NoError(t, err)

	result, err := json.Marshal(block)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}
