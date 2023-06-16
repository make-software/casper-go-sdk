package cl_type

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_ResultBool_ToString(t *testing.T) {
	assert.Equal(t, "(Result: Ok(String), Err(String)", cltype.NewResultType(cltype.String, cltype.String).String())
}

func Test_ResultBool_FromJson(t *testing.T) {
	res, err := cltype.FromRawJson(json.RawMessage(`{"Result":{"ok":"String","err":"String"}}`))
	require.NoError(t, err)
	assert.Equal(t, "(Result: Ok(String), Err(String)", res.String())
}

func Test_ResultBool_ToBytes(t *testing.T) {
	assert.Equal(t, "100a0a", hex.EncodeToString(cltype.NewResultType(cltype.String, cltype.String).Bytes()))
}

func Test_ResultBool_FromBytes(t *testing.T) {
	inBytes, err := hex.DecodeString("100a0a")
	require.NoError(t, err)
	res, err := cltype.FromBytes(inBytes)
	require.NoError(t, err)
	assert.Equal(t, cltype.NewResultType(cltype.String, cltype.String), res)
}

func Test_ResultFromRawJson_InvalidJsonFormat_ExceptError(t *testing.T) {
	_, err := cltype.FromRawJson(json.RawMessage(`{"Result": "String"}`))
	assert.ErrorIs(t, cltype.ErrInvalidResultJsonFormat, err)
	_, err = cltype.FromRawJson(json.RawMessage(`{"Result":{"err":"String"}}`))
	assert.ErrorIs(t, cltype.ErrInvalidResultJsonFormat, err)
	_, err = cltype.FromRawJson(json.RawMessage(`{"Result":{"ok":"String"}}`))
	assert.ErrorIs(t, cltype.ErrInvalidResultJsonFormat, err)
}
