package cl_value

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

func Test_UInt128_ToString(t *testing.T) {
	newInt, _ := new(big.Int).SetString("27670116110564327421", 10)
	value := clvalue.NewCLUInt128(newInt).UI128.Value()
	compound := big.NewInt(10)
	res := new(big.Int).Add(value, compound)
	assert.Equal(t, "27670116110564327431", res.String())
}

func Test_NewUInt128FromBuffer_randomValue(t *testing.T) {
	str := "0957ff1ada959f4eb106"
	byteStr, err := hex.DecodeString(str)
	require.NoError(t, err)
	res, err := clvalue.NewUint128FromBytes(byteStr)
	require.NoError(t, err)
	assert.Equal(t, "123456789101112131415", res.Value().String())
}

func Test_FromBytesByType_UInt128(t *testing.T) {
	newInt, _ := new(big.Int).SetString("27670116110564327421", 10)
	bytesData := clvalue.NewCLUInt128(newInt).Bytes()
	res, err := clvalue.FromBytesByType(bytesData, cltype.UInt128)
	require.NoError(t, err)
	assert.Equal(t, "27670116110564327421", res.UI128.String())
}
