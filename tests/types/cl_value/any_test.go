package cl_value

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

func Test_Any_Decode_Example(t *testing.T) {
	src, err := hex.DecodeString("4142432d444546")
	require.NoError(t, err)
	result := clvalue.NewAnyFromBytes(src)
	assert.Equal(t, "ABC-DEF", result.String())
}

func Test_Any_Encode_Example(t *testing.T) {
	src, err := hex.DecodeString("4142432d444546")
	require.NoError(t, err)
	dest, err := clvalue.FromBytesByType(src, cltype.Any)
	require.NoError(t, err)
	assert.Equal(t, src, dest.Bytes())
}
