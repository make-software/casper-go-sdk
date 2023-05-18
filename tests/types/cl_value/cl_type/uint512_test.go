package cl_type

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_UInt512_FromString(t *testing.T) {
	res, err := cltype.FromRawJson([]byte(cltype.TypeNameU512))
	require.NoError(t, err)
	assert.Equal(t, cltype.UInt512, res)
}

func Test_UInt512_GetName(t *testing.T) {
	assert.Equal(t, cltype.TypeNameU512, cltype.UInt512.Name())
}

func Test_UInt512_FromBytes(t *testing.T) {
	res, err := cltype.FromBytes(cltype.UInt512.Bytes())
	require.NoError(t, err)
	assert.Equal(t, cltype.UInt512, res)
}
