package cl_type

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_PublicKey_FromString(t *testing.T) {
	res, err := cltype.FromRawJson([]byte(cltype.TypeNamePublicKey))
	require.NoError(t, err)
	assert.Equal(t, cltype.PublicKey, res)
}

func Test_PublicKey_FromBytes(t *testing.T) {
	res, err := cltype.FromBytes(cltype.PublicKey.Bytes())
	require.NoError(t, err)
	assert.Equal(t, cltype.PublicKey, res)
}
