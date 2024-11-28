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
