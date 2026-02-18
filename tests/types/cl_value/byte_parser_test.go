package cl_value

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

func Test_ByteParser_ParseValueWithType_FromByteToMap(t *testing.T) {
	sourceData := `0f00000001000000030000004142430a000000110a01`
	decoded, err := hex.DecodeString(sourceData)
	require.NoError(t, err)
	data, _, err := clvalue.FromBytes(decoded)
	require.NoError(t, err)
	assert.Equal(t, int32(10), data.Map.Get("ABC").I32.Value())
}

func Test_ByteParser_ParseValueWithType_InvalidList(t *testing.T) {
	sourceData := `3e00000003080080de38cc2e8e110001000000000000000006f17367000000000000000000000000000500e87648170008009867f0b42e8e1108009867f0b42e8e110e0320000000a68dfbe1f00d679c5bdd10abe447c5009246d08d3d627a85cacab165180e482f090000004c6971756964697479`
	decoded, err := hex.DecodeString(sourceData)
	require.NoError(t, err)
	_, _, err = clvalue.FromBytes(decoded)
	require.Error(t, err)
	assert.Equal(t, err.Error(), "list size exceeds buffer size")
}

func Test_ByteParser_ToBytesWithType_FromMapToByte(t *testing.T) {
	exceptedString := `0f00000001000000030000004142430a000000110a01`
	exceptedResult, err := hex.DecodeString(exceptedString)
	require.NoError(t, err)
	dest := clvalue.NewCLMap(cltype.String, cltype.Int32)
	require.NoError(t, dest.Map.Append(*clvalue.NewCLString("ABC"), clvalue.NewCLInt32(10)))
	data, err := clvalue.ToBytesWithType(dest)
	require.NoError(t, err)
	assert.Equal(t, exceptedResult, data)
}
