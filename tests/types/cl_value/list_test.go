package cl_value

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_NewListFromBuffer_EmptyToString(t *testing.T) {
	source := "00000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.NewListFromBytes(inBytes, cltype.NewList(cltype.Bool))
	require.NoError(t, err)
	assert.Equal(t, "[]", res.String())
}

func Test_NewListFromBuffer_IncompleteFormat_ShouldRaiseError(t *testing.T) {
	source := "000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	_, err = clvalue.NewListFromBytes(inBytes, cltype.NewList(cltype.Bool))
	assert.Error(t, err)
}

func Test_NewListFromBuffer_EmptyToVal(t *testing.T) {
	source := "00000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.NewListFromBytes(inBytes, cltype.NewList(cltype.Bool))
	require.NoError(t, err)
	assert.Equal(t, []clvalue.CLValue{}, res.Elements)
}

func Test_NewListFromBuffer_U32ToString(t *testing.T) {
	source := "03000000010000000200000003000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.NewListFromBytes(inBytes, cltype.NewList(cltype.UInt32))
	require.NoError(t, err)
	assert.Equal(t, "[\"1\",\"2\",\"3\"]", res.String())
}

func Test_FromBytesByType_ListU32ToVal(t *testing.T) {
	source := "03000000010000000200000003000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.FromBytesByType(inBytes, cltype.NewList(cltype.UInt32))
	require.NoError(t, err)
	assert.Equal(t, "[\"1\",\"2\",\"3\"]", res.List.String())
}
