package cl_type

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_UInt64_FromString(t *testing.T) {
	res, err := cltype.FromRawJson([]byte(cltype.TypeNameU64))
	require.NoError(t, err)
	assert.Equal(t, cltype.UInt64, res)
}

func Test_UInt64_GetName(t *testing.T) {
	assert.Equal(t, cltype.TypeNameU64, cltype.UInt64.Name())
}

func Test_UInt64_FromBytes(t *testing.T) {
	res, err := cltype.FromBytes([]byte{cltype.TypeIDU64})
	require.NoError(t, err)
	assert.Equal(t, cltype.UInt64, res)
}
