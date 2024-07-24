package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types"
)

func Test_DeployExecutionResult_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/deploy/deploy_execution_result_example.json")
	require.NoError(t, err)

	var executionResult types.DeployExecutionResult
	err = json.Unmarshal(fixture, &executionResult)
	require.NoError(t, err)

	result, err := json.Marshal(executionResult)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}
