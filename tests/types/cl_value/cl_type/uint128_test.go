package cl_type

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_UInt128_FromString(t *testing.T) {
	res, err := cltype.FromRawJson([]byte(cltype.TypeNameU128))
	require.NoError(t, err)
	assert.Equal(t, cltype.UInt128, res)
}

func Test_UInt128_GetName(t *testing.T) {
	assert.Equal(t, cltype.TypeNameU128, cltype.UInt128.Name())
}

func Test_UInt128_FromBytes(t *testing.T) {
	res, err := cltype.FromBytes(cltype.UInt128.Bytes())
	require.NoError(t, err)
	assert.Equal(t, cltype.UInt128, res)
}
