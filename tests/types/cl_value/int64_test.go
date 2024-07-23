package cl_value

import (
	"encoding/hex"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

func Test_Int64_ToString(t *testing.T) {
	assert.Equal(t, "9223372036854775807", clvalue.NewCLInt64(math.MaxInt64).String())
}

func Test_NewInt64FromBuffer_maxValue(t *testing.T) {
	maxInBytes := clvalue.NewCLInt64(math.MaxInt64).Bytes()
	res, err := clvalue.NewInt64FromBytes(maxInBytes)
	require.NoError(t, err)
	assert.Equal(t, int64(math.MaxInt64), res.Value())
}

func Test_NewInt64FromBufferIncompleteFormat_ShouldRaiseError(t *testing.T) {
	src, err := hex.DecodeString("07000000")
	require.NoError(t, err)
	_, err = clvalue.NewInt64FromBytes(src)
	assert.Error(t, err)
}

func Test_FromBytesByType_Int64(t *testing.T) {
	maxInBytes := clvalue.NewCLInt64(math.MaxInt64).Bytes()
	res, err := clvalue.FromBytesByType(maxInBytes, cltype.Int64)
	require.NoError(t, err)
	assert.Equal(t, int64(math.MaxInt64), res.I64.Value())
}
