package cl_type

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

func Test_OptionBool_ToString(t *testing.T) {
	assert.Equal(t, "(Option: Bool)", cltype.NewOptionType(cltype.Bool).String())
}

func Test_OptionBool_FromJson(t *testing.T) {
	res, err := cltype.FromRawJson(json.RawMessage(`{"Option": "Bool"}`))
	require.NoError(t, err)
	assert.Equal(t, "(Option: Bool)", res.String())
}

func Test_OptionBool_ToBytes(t *testing.T) {
	assert.Equal(t, "0d00", hex.EncodeToString(cltype.NewOptionType(cltype.Bool).Bytes()))
}

func Test_OptionBool_FromBytes(t *testing.T) {
	inBytes, err := hex.DecodeString("0d00")
	require.NoError(t, err)
	res, err := cltype.FromBytes(inBytes)
	require.NoError(t, err)
	assert.Equal(t, cltype.NewOptionType(cltype.Bool), res)
}
