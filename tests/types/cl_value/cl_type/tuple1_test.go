package cl_type

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_TupleBool_ToString(t *testing.T) {
	assert.Equal(t, "Tuple1 (Bool)", cltype.NewTuple1(cltype.Bool).String())
}

func Test_TupleBool_FromJson(t *testing.T) {
	res, err := cltype.FromRawJson(json.RawMessage(`{"Tuple1": ["Bool"]}`))
	require.NoError(t, err)
	assert.Equal(t, "Tuple1 (Bool)", res.String())
}

func Test_TupleBool_ToBytes(t *testing.T) {
	assert.Equal(t, "1200", hex.EncodeToString(cltype.NewTuple1(cltype.Bool).Bytes()))
}

func Test_TupleBool_FromBytes(t *testing.T) {
	inBytes, err := hex.DecodeString("1200")
	require.NoError(t, err)
	res, err := cltype.FromBytes(inBytes)
	require.NoError(t, err)
	assert.Equal(t, cltype.NewTuple1(cltype.Bool), res)
}
