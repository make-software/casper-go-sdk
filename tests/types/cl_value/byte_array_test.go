package cl_value

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_NewByteArrayFromBuffer_ToString(t *testing.T) {
	source := "989ca079a5e446071866331468ab949483162588d57ec13ba6bb051f1e15f8b7"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res := clvalue.NewByteArrayFromBytes(inBytes, cltype.NewByteArray(32))
	require.NoError(t, err)
	assert.Equal(t, "989ca079a5e446071866331468ab949483162588d57ec13ba6bb051f1e15f8b7", res.String())
}

func Test_FromBytesByType_ByteArrayToString(t *testing.T) {
	source := "989ca079a5e446071866331468ab949483162588d57ec13ba6bb051f1e15f8b7"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.FromBytesByType(inBytes, cltype.NewByteArray(32))
	require.NoError(t, err)
	assert.Equal(t, "989ca079a5e446071866331468ab949483162588d57ec13ba6bb051f1e15f8b7", res.String())
}
