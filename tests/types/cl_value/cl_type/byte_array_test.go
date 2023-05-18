package cl_type

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
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
