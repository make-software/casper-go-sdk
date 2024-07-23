package cl_value

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

func Test_UInt8_ToString(t *testing.T) {
	assert.Equal(t, "255", clvalue.NewCLUint8(math.MaxUint8).String())
}

func Test_NewUInt8FromBuffer_maxValue(t *testing.T) {
	maxInBytes := clvalue.NewCLUint8(math.MaxUint8).Bytes()
	res, err := clvalue.NewUInt8FromBytes(maxInBytes)
	require.NoError(t, err)
	assert.Equal(t, uint8(math.MaxUint8), res.Value())
}

func Test_FromBytesByType_UInt8(t *testing.T) {
	maxInBytes := clvalue.NewCLUint8(math.MaxUint8).Bytes()
	res, err := clvalue.FromBytesByType(maxInBytes, cltype.UInt8)
	require.NoError(t, err)
	assert.Equal(t, uint8(math.MaxUint8), res.UI8.Value())
}
