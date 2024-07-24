package cl_value

import (
	"encoding/hex"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
)

func Test_NewInt32FromBuffer_maxValue(t *testing.T) {
	maxInBytes := clvalue.NewCLInt32(math.MaxInt32).Bytes()
	res, err := clvalue.NewInt32FromBytes(maxInBytes)
	require.NoError(t, err)
	assert.Equal(t, int32(math.MaxInt32), res.Value())
}

func Test_NewInt32FromBuffer_negative(t *testing.T) {
	maxInBytes := clvalue.NewCLUInt32(math.MaxUint32).Bytes()
	res, err := clvalue.NewInt32FromBytes(maxInBytes)
	require.NoError(t, err)
	assert.Equal(t, int32(-1), res.Value())
}

func Test_NewUInt32IncompleteFormat_ShouldRaiseError(t *testing.T) {
	src, err := hex.DecodeString("0700")
	require.NoError(t, err)
	_, err = clvalue.NewInt32FromBytes(src)
	assert.Error(t, err)
}
