package cl_value

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

func Test_NewOptionFromBuffer_EmptyToString(t *testing.T) {
	source := "00"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.NewOptionFromBytes(inBytes, cltype.NewOptionType(cltype.Bool))
	require.NoError(t, err)
	assert.Equal(t, "", res.String())
}

func Test_NewOptionFromBuffer_EmptyToVal(t *testing.T) {
	source := "00"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.NewOptionFromBytes(inBytes, cltype.NewOptionType(cltype.Bool))
	require.NoError(t, err)
	assert.Nil(t, res.Value())
}

func Test_NewOptionFromBuffer_U32ToString(t *testing.T) {
	source := "010a000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.NewOptionFromBytes(inBytes, cltype.NewOptionType(cltype.UInt32))
	require.NoError(t, err)
	assert.Equal(t, "10", res.String())
}

func Test_NewOptionFromBuffer_U32ToVal(t *testing.T) {
	source := "010a000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.NewOptionFromBytes(inBytes, cltype.NewOptionType(cltype.UInt32))
	require.NoError(t, err)
	assert.Equal(t, uint32(10), res.Value().UI32.Value())
}

func Test_FromBytesByType_OptionU32ToVal(t *testing.T) {
	source := "010a000000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.FromBytesByType(inBytes, cltype.NewOptionType(cltype.UInt32))
	require.NoError(t, err)
	assert.Equal(t, uint32(10), res.Option.Value().UI32.Value())
}
