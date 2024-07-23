package cl_type

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

func Test_ByteArray32_ToString(t *testing.T) {
	assert.Equal(t, "ByteArray: 32", cltype.NewByteArray(32).String())
}

func Test_ByteArray32_FromJson(t *testing.T) {
	res, err := cltype.FromRawJson(json.RawMessage(`{"ByteArray": 32}`))
	require.NoError(t, err)
	assert.Equal(t, "ByteArray: 32", res.String())
}

func Test_ByteArray32_ToBytes(t *testing.T) {
	bytes := cltype.NewByteArray(32).Bytes()
	assert.Equal(t, "0f20000000", hex.EncodeToString(bytes))
}

func Test_ByteArray32_FromBytes(t *testing.T) {
	inBytes, err := hex.DecodeString("0f20000000")
	require.NoError(t, err)
	res, err := cltype.FromBytes(inBytes)
	require.NoError(t, err)
	assert.Equal(t, cltype.NewByteArray(32), res)
}

func Test_ByteArray32_FromBytes_MalformedValue_ShouldBeErr(t *testing.T) {
	inBytes, err := hex.DecodeString("0f2000")
	require.NoError(t, err)
	_, err = cltype.FromBytes(inBytes)
	assert.Error(t, err)
}

func Test_ByteArray32_FromBuffer_StreamValue_ShouldTakeOnly4Bytes(t *testing.T) {
	inBytes, err := hex.DecodeString("0f200000001010")
	require.NoError(t, err)
	buf := bytes.NewBuffer(inBytes)
	_, err = cltype.FromBuffer(buf)
	assert.NoError(t, err)
	assert.Equal(t, "1010", hex.EncodeToString(buf.Bytes()))
}
