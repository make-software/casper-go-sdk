package cl_type

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_UInt8_FromString(t *testing.T) {
	res, err := cltype.FromRawJson([]byte(cltype.TypeNameU8))
	require.NoError(t, err)
	assert.Equal(t, cltype.UInt8, res)
}

func Test_UInt8_GetName(t *testing.T) {
	assert.Equal(t, cltype.TypeNameU8, cltype.UInt8.Name())
}

func Test_UInt8_FromBytes(t *testing.T) {
	res, err := cltype.FromBytes([]byte{cltype.TypeIDU8})
	require.NoError(t, err)
	assert.Equal(t, cltype.UInt8, res)
}
