package cl_value

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
	"github.com/make-software/casper-go-sdk/types/key"
)

func Test_NewURefFromBuffer_toString(t *testing.T) {
	source := "000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f07"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := key.NewURefFromBytes(inBytes)
	require.NoError(t, err)
	assert.Equal(t, "uref-000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f-007", res.ToPrefixedString())
}

func Test_NewURefFromString_toBytes(t *testing.T) {
	source := "uref-000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f-007"
	res, err := key.NewURef(source)
	require.NoError(t, err)
	assert.Equal(t, "000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f07", hex.EncodeToString(res.Bytes()))
}

func Test_FromBytesByType_URef(t *testing.T) {
	source := "000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f07"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.FromBytesByType(inBytes, cltype.Uref)
	require.NoError(t, err)
	assert.Equal(t, "000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f07", hex.EncodeToString(res.Uref.Bytes()))
}
