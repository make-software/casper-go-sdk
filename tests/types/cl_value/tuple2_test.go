package cl_value

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_NewTuple2FromBuffer_U32ToString(t *testing.T) {
	source := "0a0000000b000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.NewTuple2FromBytes(inBytes, cltype.NewTuple2(cltype.UInt32, cltype.UInt32))
	require.NoError(t, err)
	assert.Equal(t, "(10, 11)", res.String())
}

func Test_NewTuple2FromBuffer_U32ToVal(t *testing.T) {
	source := "0a0000000b000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.NewTuple2FromBytes(inBytes, cltype.NewTuple2(cltype.UInt32, cltype.UInt32))
	require.NoError(t, err)
	assert.Equal(t, uint32(10), res.Value()[0].UI32.Value())
	assert.Equal(t, uint32(11), res.Value()[1].UI32.Value())
}

func Test_FromBytesByType_Tuple2U32ToVal(t *testing.T) {
	source := "0a0000000b000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.FromBytesByType(inBytes, cltype.NewTuple2(cltype.UInt32, cltype.UInt32))
	require.NoError(t, err)
	assert.Equal(t, "(10, 11)", res.String())
}
