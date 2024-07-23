package types

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types"
)

func Test_BlockProposer_Scan_System(t *testing.T) {
	res := &types.Proposer{}
	hexData, err := hex.DecodeString("00")
	require.NoError(t, err)
	err = res.Scan(hexData)
	require.NoError(t, err)
	assert.True(t, res.IsSystem())
}

func Test_BlockProposer_Scan_PublicKey(t *testing.T) {
	res := &types.Proposer{}
	hexData, err := hex.DecodeString("015a372b0e230bf9393e2df0b3de857bb0e17370884bb881f840cb1482bb2922cf")
	require.NoError(t, err)
	err = res.Scan(hexData)
	require.NoError(t, err)
	assert.False(t, res.IsSystem())
	publicKey, err := res.PublicKey()
	require.NoError(t, err)
	assert.Equal(t, "015a372b0e230bf9393e2df0b3de857bb0e17370884bb881f840cb1482bb2922cf", publicKey.ToHex())
}

func Test_BlockProposer_InvalidData(t *testing.T) {
	res := &types.Proposer{}
	err := res.Scan("0")
	assert.Error(t, err)
}

func Test_BlockProposer_Value_System(t *testing.T) {
	res, err := types.NewProposer("00")
	require.NoError(t, err)
	assert.True(t, res.IsSystem())
}

func Test_BlockProposer_Value_PublicKey(t *testing.T) {
	res, err := types.NewProposer("015a372b0e230bf9393e2df0b3de857bb0e17370884bb881f840cb1482bb2922cf")
	require.NoError(t, err)
	pubKey, err := res.PublicKey()
	require.NoError(t, err)
	assert.Equal(t, "015a372b0e230bf9393e2df0b3de857bb0e17370884bb881f840cb1482bb2922cf", pubKey.ToHex())
}
