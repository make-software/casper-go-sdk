package cl_type

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_ListBool_ToString(t *testing.T) {
	assert.Equal(t, "(List of Bool)", cltype.NewList(cltype.Bool).String())
}

func Test_ListBool_FromJson(t *testing.T) {
	res, err := cltype.FromRawJson(json.RawMessage(`{"List": "Bool"}`))
	require.NoError(t, err)
	assert.Equal(t, "(List of Bool)", res.String())
}

func Test_ListBool_ToBytes(t *testing.T) {
	assert.Equal(t, "0e00", hex.EncodeToString(cltype.NewList(cltype.Bool).Bytes()))
}

func Test_ListBool_FromBytes(t *testing.T) {
	inBytes, err := hex.DecodeString("0e00")
	require.NoError(t, err)
	res, err := cltype.FromBytes(inBytes)
	require.NoError(t, err)
	assert.Equal(t, cltype.NewList(cltype.Bool), res)
}
