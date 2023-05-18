package cl_value

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue"
)

func Test_String_Decode_Example(t *testing.T) {
	src, err := hex.DecodeString("070000004142432d444546")
	require.NoError(t, err)
	result := clvalue.NewStringFromBytes(src)
	assert.Equal(t, "ABC-DEF", result.String())
}

func Test_String_Encode_Example(t *testing.T) {
	src, err := hex.DecodeString("070000004142432d444546")
	require.NoError(t, err)
	dest := clvalue.NewCLString("ABC-DEF")
	require.NoError(t, err)
	assert.Equal(t, src, dest.Bytes())
}
