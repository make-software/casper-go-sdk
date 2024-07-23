package cl_type

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

func Test_UInt256_FromString(t *testing.T) {
	res, err := cltype.FromRawJson([]byte(cltype.TypeNameU256))
	require.NoError(t, err)
	assert.Equal(t, cltype.UInt256, res)
}

func Test_UInt256_GetName(t *testing.T) {
	assert.Equal(t, cltype.TypeNameU256, cltype.UInt256.Name())
}

func Test_UInt256_FromBytes(t *testing.T) {
	res, err := cltype.FromBytes(cltype.UInt256.Bytes())
	require.NoError(t, err)
	assert.Equal(t, cltype.UInt256, res)
}
