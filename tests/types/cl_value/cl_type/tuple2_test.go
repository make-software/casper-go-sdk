package cl_type

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_Tuple2Bool_ToString(t *testing.T) {
	assert.Equal(t, "Tuple2 (Bool, Bool)", cltype.NewTuple2(cltype.Bool, cltype.Bool).String())
}

func Test_Tuple2Bool_FromJson(t *testing.T) {
	res, err := cltype.FromRawJson(json.RawMessage(`{"Tuple2": ["Bool","Bool"]}`))
	require.NoError(t, err)
	assert.Equal(t, "Tuple2 (Bool, Bool)", res.String())
}

func Test_Tuple2Bool_ToBytes(t *testing.T) {
	assert.Equal(t, "130000", hex.EncodeToString(cltype.NewTuple2(cltype.Bool, cltype.Bool).Bytes()))
}

func Test_Tuple2Bool_FromBytes(t *testing.T) {
	inBytes, err := hex.DecodeString("130000")
	require.NoError(t, err)
	res, err := cltype.FromBytes(inBytes)
	require.NoError(t, err)
	assert.Equal(t, cltype.NewTuple2(cltype.Bool, cltype.Bool), res)
}
