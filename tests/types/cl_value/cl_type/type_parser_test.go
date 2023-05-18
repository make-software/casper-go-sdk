package cl_type

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_FromBuffer_Int32(t *testing.T) {
	res, err := cltype.FromBuffer(bytes.NewBuffer([]byte{1}))
	require.NoError(t, err)
	assert.Equal(t, cltype.Int32, res)
}
