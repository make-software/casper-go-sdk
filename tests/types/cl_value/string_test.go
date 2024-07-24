package cl_value

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
)

func Test_String_Decode_Example(t *testing.T) {
	src, err := hex.DecodeString("070000004142432d444546")
	require.NoError(t, err)
	result, err := clvalue.NewStringFromBytes(src)
	require.NoError(t, err)
	assert.Equal(t, "ABC-DEF", result.String())
}

func Test_StringIncompleteFormat_ShouldBeError(t *testing.T) {
	src, err := hex.DecodeString("0700")
	require.NoError(t, err)
	_, err = clvalue.NewStringFromBytes(src)
	assert.Error(t, err)
}

func Test_String_Encode_Example(t *testing.T) {
	src, err := hex.DecodeString("070000004142432d444546")
	require.NoError(t, err)
	dest := clvalue.NewCLString("ABC-DEF")
	require.NoError(t, err)
	assert.Equal(t, src, dest.Bytes())
}
