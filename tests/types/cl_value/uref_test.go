package cl_value

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
	"github.com/make-software/casper-go-sdk/v2/types/key"
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

func Test_NewURefFromString_IncorrectFormat(t *testing.T) {
	source := "01506f4df2ac64a2233e787c430dc91dad5cee8eabd7d64555f64bdc1a4b4044d7"
	_, err := key.NewURef(source)
	require.Error(t, err)
	assert.Equal(t, err, key.ErrIncorrectUrefFormat)
}

func Test_FromBytesByType_URef(t *testing.T) {
	source := "000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f07"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.FromBytesByType(inBytes, cltype.Uref)
	require.NoError(t, err)
	assert.Equal(t, "000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f07", hex.EncodeToString(res.Uref.Bytes()))
}
