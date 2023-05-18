package cl_type

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_Uref_FromString(t *testing.T) {
	res, err := cltype.FromRawJson([]byte(cltype.TypeNameURef))
	require.NoError(t, err)
	assert.Equal(t, cltype.Uref, res)
}

func Test_Uref_GetName(t *testing.T) {
	assert.Equal(t, cltype.TypeNameURef, cltype.Uref.Name())
}

func Test_Uref_FromBytes(t *testing.T) {
	res, err := cltype.FromBytes(cltype.Uref.Bytes())
	require.NoError(t, err)
	assert.Equal(t, cltype.Uref, res)
}
