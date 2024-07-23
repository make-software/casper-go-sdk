package cl_value

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

func Test_NewTuple1FromBuffer_U32ToString(t *testing.T) {
	source := "0a000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.NewTuple1FromBytes(inBytes, cltype.NewTuple1(cltype.UInt32))
	require.NoError(t, err)
	assert.Equal(t, "(10)", res.String())
}

func Test_NewTuple1FromBuffer_U32ToVal(t *testing.T) {
	source := "0a000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.NewTuple1FromBytes(inBytes, cltype.NewTuple1(cltype.UInt32))
	require.NoError(t, err)
	assert.Equal(t, uint32(10), res.Value().UI32.Value())
}

func Test_FromBytesByType_Tuple1U32ToVal(t *testing.T) {
	source := "0a000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.FromBytesByType(inBytes, cltype.NewTuple1(cltype.UInt32))
	require.NoError(t, err)
	assert.Equal(t, uint32(10), res.Tuple1.Value().UI32.Value())
}
