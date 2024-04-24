package cl_value

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_NewTupl32FromBuffer_U32ToString(t *testing.T) {
	source := "0a0000000b0000000c000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.NewTuple3FromBytes(inBytes, cltype.NewTuple3(cltype.UInt32, cltype.UInt32, cltype.UInt32))
	require.NoError(t, err)
	assert.Equal(t, "(10, 11, 12)", res.String())
}

func Test_NewTuple3FromBuffer_U32ToVal(t *testing.T) {
	source := "0a0000000b0000000c000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.NewTuple3FromBytes(inBytes, cltype.NewTuple3(cltype.UInt32, cltype.UInt32, cltype.UInt32))
	require.NoError(t, err)
	assert.Equal(t, uint32(10), res.Value()[0].UI32.Value())
	assert.Equal(t, uint32(11), res.Value()[1].UI32.Value())
	assert.Equal(t, uint32(12), res.Value()[2].UI32.Value())
}

func Test_FromBytesByType_Tuple3U32ToVal(t *testing.T) {
	source := "0a0000000b0000000c000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.FromBytesByType(inBytes, cltype.NewTuple3(cltype.UInt32, cltype.UInt32, cltype.UInt32))
	require.NoError(t, err)
	assert.Equal(t, "(10, 11, 12)", res.String())
}

func Test_NewCLTuple3(t *testing.T) {
	val1 := clvalue.NewCLUInt32(10)
	val2 := clvalue.NewCLUInt32(11)
	val3 := clvalue.NewCLUInt32(12)
	res := clvalue.NewCLTuple3(*val1, *val2, *val3)
	assert.Equal(t, "(10, 11, 12)", res.String())
}
