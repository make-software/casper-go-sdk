package cl_type

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

func Test_Int64_FromString(t *testing.T) {
	res, err := cltype.FromRawJson([]byte(cltype.TypeNameI64))
	require.NoError(t, err)
	assert.Equal(t, cltype.Int64, res)
}

func Test_Int64_GetName(t *testing.T) {
	assert.Equal(t, cltype.TypeNameI64, cltype.Int64.Name())
}

func Test_Int64_FromBytes(t *testing.T) {
	res, err := cltype.FromBytes([]byte{cltype.TypeIDI64})
	require.NoError(t, err)
	assert.Equal(t, cltype.Int64, res)
}
