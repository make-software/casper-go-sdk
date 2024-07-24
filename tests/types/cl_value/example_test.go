package cl_value

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

func Test_DecodeTypeAndValue(t *testing.T) {
	//united representation of type and value in one row
	sourceString := `0f00000001000000030000004142430a000000110a01`
	hexData, err := hex.DecodeString(sourceString)
	require.NoError(t, err)
	data, err := clvalue.FromBytes(hexData)
	require.NoError(t, err)
	assert.Equal(t, int32(10), data.Map.Get("ABC").I32.Value())
}

func Test_EncodeTypeAndValue(t *testing.T) {
	//united representation of type and value in one row
	expectedString := `0f00000001000000030000004142430a000000110a01`
	expectedResult, err := hex.DecodeString(expectedString)
	require.NoError(t, err)
	dest := clvalue.NewCLMap(cltype.String, cltype.Int32)
	require.NoError(t, dest.Map.Append(*clvalue.NewCLString("ABC"), clvalue.NewCLInt32(10)))
	data, err := clvalue.ToBytesWithType(dest)
	require.NoError(t, err)
	assert.Equal(t, expectedResult, data)
}

func Test_EncodeMapValue(t *testing.T) {
	//representation of map value
	expectedValue := "01000000030000004142430a000000"
	dest := clvalue.NewCLMap(cltype.String, cltype.Int32)
	require.NoError(t, dest.Map.Append(*clvalue.NewCLString("ABC"), clvalue.NewCLInt32(10)))
	assert.Equal(t, expectedValue, hex.EncodeToString(dest.Bytes()))
}

func Test_DecodeMapValue(t *testing.T) {
	//representation of map value
	sourceStr := `01000000030000004142430a000000`
	sourceHex, err := hex.DecodeString(sourceStr)
	require.NoError(t, err)
	mapType := cltype.NewMap(cltype.String, cltype.Int32)
	expectedValue := clvalue.NewCLMap(cltype.String, cltype.Int32)
	require.NoError(t, expectedValue.Map.Append(*clvalue.NewCLString("ABC"), clvalue.NewCLInt32(10)))
	result, err := clvalue.FromBytesByType(sourceHex, mapType)
	require.NoError(t, err)
	assert.Equal(t, expectedValue, result)
}

func Test_EncodeMapType(t *testing.T) {
	mapType := cltype.Map{Key: cltype.String, Val: cltype.Int32}
	assert.Equal(t, "110a01", hex.EncodeToString(mapType.Bytes()))
}

func Test_DecodeMapType(t *testing.T) {
	sourceStr := "110a01"
	sourceHex, err := hex.DecodeString(sourceStr)
	require.NoError(t, err)
	expectedResult := &cltype.Map{Key: cltype.String, Val: cltype.Int32}
	res, err := cltype.FromBytes(sourceHex)
	require.NoError(t, err)
	assert.Equal(t, expectedResult, res)
}

func Test_DecodeArgs_FromJson(t *testing.T) {
	sourceData := `{"bytes":"01000000030000004142430a000000","cl_type":{"Map":{"key":"String","value":"I32"}}}`
	data, err := types.ArgsFromRawJson(json.RawMessage(sourceData))
	require.NoError(t, err)
	assert.Equal(t, int32(10), data.Map.Get("ABC").I32.Value())
}
