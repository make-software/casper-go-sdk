package cl_value

import (
	"encoding/hex"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_UInt64_ToString(t *testing.T) {
	assert.Equal(t, "18446744073709551615", clvalue.NewCLUInt64(math.MaxUint64).String())
}

func Test_NewUInt64FromBuffer_maxValue(t *testing.T) {
	maxInBytes := clvalue.NewCLUInt64(math.MaxUint64).Bytes()
	res, err := clvalue.NewUint64FromBytes(maxInBytes)
	require.NoError(t, err)
	assert.Equal(t, uint64(math.MaxUint64), res.Value())
}

func Test_NewUInt64IncompleteFormat_ShouldRaiseError(t *testing.T) {
	src, err := hex.DecodeString("07000000")
	require.NoError(t, err)
	_, err = clvalue.NewUint64FromBytes(src)
	assert.Error(t, err)
}

func Test_FromBytesByType_UInt64(t *testing.T) {
	maxInBytes := clvalue.NewCLUInt64(math.MaxUint64).Bytes()
	res, err := clvalue.FromBytesByType(maxInBytes, cltype.UInt64)
	require.NoError(t, err)
	assert.Equal(t, uint64(math.MaxUint64), res.UI64.Value())
}

func Test_UInt64_ToBytes(t *testing.T) {
	assert.Equal(t, "0004000000000000", hex.EncodeToString(clvalue.NewCLUInt64(1024).Bytes()))
}
