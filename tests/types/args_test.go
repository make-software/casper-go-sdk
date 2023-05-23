package types

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
)

func Test_Args_Bytes(t *testing.T) {
	args := &types.Args{}
	args.AddArgument("test", clvalue.NewCLBool(true))
	res, err := args.Bytes()
	require.NoError(t, err)
	assert.Equal(t, "010000000400000074657374010000000100", hex.EncodeToString(res))
}

func Test_Args_Find(t *testing.T) {
	args := &types.Args{}
	args.AddArgument("test", clvalue.NewCLBool(true))
	res, err := args.Find("test")
	require.NoError(t, err)
	value, err := res.Value()
	require.NoError(t, err)
	assert.Equal(t, true, value.Bool.Value())
}

func Test_Args_Marshal(t *testing.T) {
	args := &types.Args{}
	args.AddArgument("test", clvalue.NewCLBool(true))
	res, err := json.Marshal(args)
	require.NoError(t, err)
	assert.JSONEq(t, `[["test", {"bytes": "01", "cl_type": "Bool"}]]`, string(res))
}

func Test_Args_GetBytesBool(t *testing.T) {
	args := &types.Args{}
	args.AddArgument("test", clvalue.NewCLBool(true))
	res, err := args.Find("test")
	require.NoError(t, err)
	value, err := res.Bytes()
	require.NoError(t, err)
	assert.Equal(t, "010000000100", value.String())
}
