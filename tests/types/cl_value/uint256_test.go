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

func Test_UInt256_ToString(t *testing.T) {
	newInt, _ := new(big.Int).SetString("2767011611056432742164327421", 10)
	value := clvalue.NewCLUInt256(newInt).UI256.Value()
	compound := big.NewInt(10)
	res := new(big.Int).Add(value, compound)
	assert.Equal(t, "2767011611056432742164327431", res.String())
}

func Test_NewUInt256FromBuffer_randomValue(t *testing.T) {
	str := "1457ff1ada959f4eb565465457ff1ada959f4eb106"
	byteStr, err := hex.DecodeString(str)
	require.NoError(t, err)
	res, err := clvalue.NewUint256FromBytes(byteStr)
	require.NoError(t, err)
	s := res.Value().String()
	assert.Equal(t, "38208025587469414829876084373438588806885998423", s)
}

func Test_FromBytesByType_UInt256(t *testing.T) {
	newInt, _ := new(big.Int).SetString("2767011611056432742164327421", 10)
	bytesData := clvalue.NewCLUInt256(newInt).Bytes()
	res, err := clvalue.FromBytesByType(bytesData, cltype.UInt256)
	require.NoError(t, err)
	assert.Equal(t, "2767011611056432742164327421", res.UI256.String())
}
