package rpc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/rpc"
)

func Test_UnmarshalRpcRequest_withIDAsStringVal(t *testing.T) {
	data := `{
  "jsonrpc": "2.0",
  "id": "0",
  "method": "state_get_auction_info",
  "params": {}
}`
	var sut rpc.RpcRequest
	err := json.Unmarshal([]byte(data), &sut)
	require.NoError(t, err)
	result, err := json.Marshal(sut)
	require.NoError(t, err)
	assert.JSONEq(t, data, string(result))
}

func Test_UnmarshalRpcRequest_withIDAsIntVal(t *testing.T) {
	data := `{
  "jsonrpc": "2.0",
  "id": 0,
  "method": "state_get_auction_info",
  "params": {}
}`
	t.Skip("Don't support int values yet")
	var sut rpc.RpcRequest
	err := json.Unmarshal([]byte(data), &sut)
	require.NoError(t, err)
	result, err := json.Marshal(sut)
	require.NoError(t, err)
	assert.JSONEq(t, data, string(result))
}

func Test_UnmarshalRpcRequest_withoutID(t *testing.T) {
	data := `{
  "jsonrpc": "2.0",
  "method": "state_get_auction_info",
  "params": {}
}`
	var sut rpc.RpcRequest
	err := json.Unmarshal([]byte(data), &sut)
	require.NoError(t, err)
	result, err := json.Marshal(sut)
	require.NoError(t, err)
	assert.JSONEq(t, data, string(result))
}

func Test_NewParamBlockByHeight_withZero(t *testing.T) {
	tmp := rpc.NewParamBlockByHeight(0)
	res, err := json.Marshal(tmp)
	require.NoError(t, err)
	assert.JSONEq(t, `{"block_identifier":{"Height":0}}`, string(res))
}
