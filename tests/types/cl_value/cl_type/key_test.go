package cl_type

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_Key_FromString(t *testing.T) {
	res, err := cltype.FromRawJson([]byte(cltype.TypeNameKey))
	require.NoError(t, err)
	assert.Equal(t, cltype.Key, res)
}

func Test_Key_GetName(t *testing.T) {
	assert.Equal(t, cltype.TypeNameKey, cltype.Key.Name())
}

func Test_Key_FromBytes(t *testing.T) {
	res, err := cltype.FromBytes(cltype.Key.Bytes())
	require.NoError(t, err)
	assert.Equal(t, cltype.Key, res)
}
