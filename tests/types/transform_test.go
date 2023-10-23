package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types"
)

func Test_Transform_AddUint512(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/AddUInt512.json")
	require.NoError(t, err)
	var transformKey types.TransformKey

	err = json.Unmarshal(fixture, &transformKey)
	require.NoError(t, err)

	val, err := transformKey.Transform.ParseAsUInt512()
	require.NoError(t, err)

	assert.EqualValues(t, 100000000, val.Value().Int64())
}
