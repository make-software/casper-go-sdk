package cl_value

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

func Test_Unit_ToString(t *testing.T) {
	assert.Equal(t, "nil", clvalue.NewCLUnit().String())
}

func Test_Unit_ToVal(t *testing.T) {
	assert.Equal(t, nil, clvalue.NewCLUnit().Unit.Value().Type())
}

func Test_NewUnitFromBuffer_randomValue(t *testing.T) {
	res, err := clvalue.NewUnitFromBytes([]byte{})
	require.NoError(t, err)
	assert.Equal(t, "nil", res.Value().String())
}

func Test_FromBytesByType_Unit(t *testing.T) {
	res, err := clvalue.FromBytesByType([]byte{}, cltype.Unit)
	require.NoError(t, err)
	assert.Equal(t, "nil", res.Unit.String())
}
