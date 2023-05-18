package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/make-software/casper-go-sdk/types"
)

func Test_EraSummary_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/era/era_summary_example.json")
	assert.NoError(t, err)

	var era types.EraSummary
	err = json.Unmarshal(fixture, &era)
	assert.NoError(t, err)

	result, err := json.Marshal(era)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}
