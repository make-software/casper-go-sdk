package cl_type

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_Unit_FromString(t *testing.T) {
	res, err := cltype.FromRawJson([]byte(cltype.TypeNameUnit))
	require.NoError(t, err)
	assert.Equal(t, cltype.Unit, res)
}

func Test_Unit_GetName(t *testing.T) {
	assert.Equal(t, cltype.TypeNameUnit, cltype.Unit.Name())
}

func Test_Unit_FromBytes(t *testing.T) {
	res, err := cltype.FromBytes(cltype.Unit.Bytes())
	require.NoError(t, err)
	assert.Equal(t, cltype.Unit, res)
}
