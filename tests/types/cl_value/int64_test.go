package cl_value

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_Int64_ToString(t *testing.T) {
	assert.Equal(t, "9223372036854775807", clvalue.NewCLInt64(math.MaxInt64).String())
}

func Test_NewInt64FromBuffer_maxValue(t *testing.T) {
	maxInBytes := clvalue.NewCLInt64(math.MaxInt64).Bytes()
	res := clvalue.NewInt64FromBytes(maxInBytes)
	assert.Equal(t, int64(math.MaxInt64), res.Value())
}

func Test_FromBytesByType_Int64(t *testing.T) {
	maxInBytes := clvalue.NewCLInt64(math.MaxInt64).Bytes()
	res, err := clvalue.FromBytesByType(maxInBytes, cltype.Int64)
	require.NoError(t, err)
	assert.Equal(t, int64(math.MaxInt64), res.I64.Value())
}
