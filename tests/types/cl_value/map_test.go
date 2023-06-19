package cl_value

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
	"github.com/make-software/casper-go-sdk/types/key"
)

func Test_MapType_Bool_To_String(t *testing.T) {
	newtype := cltype.Map{Key: cltype.Bool, Val: &cltype.Map{
		Key: cltype.Bool,
		Val: cltype.Bool,
	}}
	res := newtype.String()
	assert.Equal(t, "Map (Bool: Map (Bool: Bool))", res)
}

func Test_MapType_String_To_String_1(t *testing.T) {
	newtype := cltype.Map{Key: cltype.String, Val: cltype.Bool}
	assert.Equal(t, "Map (String: Bool)", newtype.String())
}

func Test_Map_ToMap(t *testing.T) {
	key1, err := key.NewKey("account-hash-bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59")
	require.NoError(t, err)
	key2, err := key.NewKey("uref-7b12008bb757ee32caefb3f7a1f77d9f659ee7a4e21ad4950c4e0294000492eb-007")
	require.NoError(t, err)
	sourceMap := clvalue.NewCLMap(cltype.Key, cltype.UInt512)
	err = sourceMap.Map.Append(clvalue.NewCLKey(key1), *clvalue.NewCLUInt512(big.NewInt(123)))
	require.NoError(t, err)
	err = sourceMap.Map.Append(clvalue.NewCLKey(key2), *clvalue.NewCLUInt512(big.NewInt(124)))
	require.NoError(t, err)
	result := sourceMap.Map.Map()
	assert.Equal(t,
		result["uref-7b12008bb757ee32caefb3f7a1f77d9f659ee7a4e21ad4950c4e0294000492eb-007"].UI512.Value().Uint64(),
		uint64(124),
	)
}

func Test_Map_ToData(t *testing.T) {
	key1, err := key.NewKey("account-hash-bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59")
	require.NoError(t, err)
	key2, err := key.NewKey("uref-7b12008bb757ee32caefb3f7a1f77d9f659ee7a4e21ad4950c4e0294000492eb-007")
	require.NoError(t, err)
	sourceMap := clvalue.NewCLMap(cltype.Key, cltype.UInt512)
	err = sourceMap.Map.Append(clvalue.NewCLKey(key1), *clvalue.NewCLUInt512(big.NewInt(123)))
	require.NoError(t, err)
	err = sourceMap.Map.Append(clvalue.NewCLKey(key2), *clvalue.NewCLUInt512(big.NewInt(124)))
	require.NoError(t, err)
	result := sourceMap.Map.Data()
	assert.Equal(t,
		"account-hash-bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59",
		result[0].Inner1.Key.Account.ToPrefixedString(),
	)
	assert.Equal(t,
		"uref-7b12008bb757ee32caefb3f7a1f77d9f659ee7a4e21ad4950c4e0294000492eb-007",
		result[1].Inner1.Key.URef.String(),
	)
}

func Test_NewMapFromBuffer_IncompleteFormat_ShouldRaiseError(t *testing.T) {
	source := "0000"
	inBytes, err := hex.DecodeString(source)
	require.NoError(t, err)
	_, err = clvalue.NewMapFromBuffer(bytes.NewBuffer(inBytes), cltype.NewMap(cltype.String, cltype.String))
	assert.Error(t, err)
}
