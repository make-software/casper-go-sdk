package cl_type

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_Any_FromString(t *testing.T) {
	res, err := cltype.FromRawJson([]byte(cltype.TypeNameAny))
	require.NoError(t, err)
	assert.Equal(t, cltype.Any, res)
}

func Test_Any_GetName(t *testing.T) {
	assert.Equal(t, cltype.TypeNameAny, cltype.Any.Name())
}

func Test_Any_FromBytes(t *testing.T) {
	res, err := cltype.FromBytes(cltype.Any.Bytes())
	require.NoError(t, err)
	assert.Equal(t, cltype.Any, res)
}
