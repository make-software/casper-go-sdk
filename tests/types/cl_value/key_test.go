package cl_value

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
	"github.com/make-software/casper-go-sdk/types/key"
)

func Test_KeyHash_ToString(t *testing.T) {
	source := "019d48b4264218bb732141b2975a480b049c131b92b711f9ff823962f745ec966a"
	decoded, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := key.NewKeyFromBytes(decoded)
	require.NoError(t, err)
	assert.Equal(t, "hash-9d48b4264218bb732141b2975a480b049c131b92b711f9ff823962f745ec966a", res.ToPrefixedString())
}

func Test_KeyHash_ToBytes(t *testing.T) {
	res, err := key.NewKey("hash-9d48b4264218bb732141b2975a480b049c131b92b711f9ff823962f745ec966a")
	require.NoError(t, err)
	assert.Equal(t, "019d48b4264218bb732141b2975a480b049c131b92b711f9ff823962f745ec966a", hex.EncodeToString(res.Bytes()))
}

func Test_KeyEraID_ToString(t *testing.T) {
	source := "050004000000000000"
	decoded, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := key.NewKeyFromBytes(decoded)
	require.NoError(t, err)
	assert.Equal(t, "era-1024", res.ToPrefixedString())
}

func Test_KeyEraID_ToBytes(t *testing.T) {
	res, err := key.NewKey("era-1024")
	require.NoError(t, err)
	assert.Equal(t, "050004000000000000", hex.EncodeToString(res.Bytes()))
}

func Test_KeyEraSummary_ToBytes(t *testing.T) {
	res, err := key.NewKey("era-summary-0000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)
	assert.Equal(t, "0b0000000000000000000000000000000000000000000000000000000000000000", hex.EncodeToString(res.Bytes()))
}

func Test_KeyEraID_FromBytesByType(t *testing.T) {
	source := "050004000000000000"
	decoded, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.FromBytesByType(decoded, cltype.Key)
	require.NoError(t, err)
	assert.Equal(t, "era-1024", res.Key.ToPrefixedString())
}
