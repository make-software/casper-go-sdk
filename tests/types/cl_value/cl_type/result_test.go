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
	assert.Equal(t, "(Result: Bool)", cltype.NewResultType(cltype.Bool).String())
}

func Test_ResultBool_FromJson(t *testing.T) {
	res, err := cltype.FromRawJson(json.RawMessage(`{"Result": "Bool"}`))
	require.NoError(t, err)
	assert.Equal(t, "(Result: Bool)", res.String())
}

func Test_ResultBool_ToBytes(t *testing.T) {
	assert.Equal(t, "1000", hex.EncodeToString(cltype.NewResultType(cltype.Bool).Bytes()))
}

func Test_ResultBool_FromBytes(t *testing.T) {
	inBytes, err := hex.DecodeString("1000")
	require.NoError(t, err)
	res, err := cltype.FromBytes(inBytes)
	require.NoError(t, err)
	assert.Equal(t, cltype.NewResultType(cltype.Bool), res)
}
