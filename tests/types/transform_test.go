package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types"
)

func Test_Transform_AddUInt512(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/AddUInt512.json")
	require.NoError(t, err)
	var transformKey types.TransformKey

	err = json.Unmarshal(fixture, &transformKey)
	require.NoError(t, err)

	val, err := transformKey.Transform.ParseAsUInt512()
	require.NoError(t, err)

	assert.True(t, transformKey.Transform.IsAddUint512())
	assert.EqualValues(t, 100000000, val.Value().Int64())
}

func Test_Transform_WriteDeployInfo(t *testing.T) {
	fixute, err := os.ReadFile("../data/transform/WriteDeployInfo.json")
	require.NoError(t, err)
	var transformKey types.TransformKey

	err = json.Unmarshal(fixute, &transformKey)
	require.NoError(t, err)

	val, err := transformKey.Transform.ParseAsWriteDeployInfo()
	require.NoError(t, err)

	assert.True(t, transformKey.Transform.IsWriteDeployInfo())
	assert.EqualValues(t, 1, len(val.Transfers))
}
