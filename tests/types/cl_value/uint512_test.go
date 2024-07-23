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

func Test_UInt512_ToString(t *testing.T) {
	newInt, _ := new(big.Int).SetString("2767011611056432742164327421432742164327421", 10)
	value := clvalue.NewCLUInt512(newInt).UI512.Value()
	compound := big.NewInt(10)
	res := new(big.Int).Add(value, compound)
	assert.Equal(t, "2767011611056432742164327421432742164327431", res.String())
}

func Test_NewUInt512FromBuffer_randomValue(t *testing.T) {
	str := "0957ff1ada959f4eb106"
	byteStr, err := hex.DecodeString(str)
	require.NoError(t, err)
	res, err := clvalue.NewUint512FromBytes(byteStr)
	require.NoError(t, err)
	assert.Equal(t, "123456789101112131415", res.Value().String())
}

func Test_FromBytesByType_UInt512(t *testing.T) {
	newInt, _ := new(big.Int).SetString("2767011611056432742164327421432742164327421", 10)
	bytesData := clvalue.NewCLUInt512(newInt).Bytes()
	res, err := clvalue.FromBytesByType(bytesData, cltype.UInt512)
	require.NoError(t, err)
	assert.Equal(t, "2767011611056432742164327421432742164327421", res.UI512.String())
}
