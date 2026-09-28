package rpc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/rpc"
	"github.com/make-software/casper-go-sdk/v2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetupSpecExecServer(t *testing.T, filePath string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		var requestParams rpc.RpcRequest
		err := decoder.Decode(&requestParams)
		require.NoError(t, err)
		fixture, err := os.ReadFile(filePath)
		require.NoError(t, err)
		_, err = rw.Write(fixture)
		require.NoError(t, err)
	}))
	return server
}

func Test_SpeculativeExecResult_UnmarshalJSON_V2(t *testing.T) {
	server := SetupSpecExecServer(t, "../data/rpc_response/specexec_cep18transfer_success_v200.json")
	defer server.Close()
	client := rpc.NewSpeculativeClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.SpeculativeExec(context.Background(), types.Deploy{}, nil)
	require.NoError(t, err)
	assert.Equal(t, "2.0.0", result.ApiVersion)
	assert.Equal(t, "2bd86fac7be3623f5a99cb43973eee8da93530c5b14947c4416778413a80a3fa", result.BlockHash.ToHex())
	assert.Equal(t, uint64(2500000000), result.ExecutionResult.Limit)
	assert.Equal(t, uint64(395091618), result.ExecutionResult.Consumed)
	assert.Equal(t, 2, len(result.ExecutionResult.Effects))
	assert.Equal(t, 1, len(result.Messages))
	assert.Nil(t, result.ExecutionResult.ErrorMessage)
}

func Test_SpeculativeExecResult_UnmarshalJSON_V2_WithError(t *testing.T) {
	server := SetupSpecExecServer(t, "../data/rpc_response/specexec_transfer_failure_v200.json")
	defer server.Close()
	client := rpc.NewSpeculativeClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.SpeculativeExec(context.Background(), types.Deploy{}, nil)
	require.NoError(t, err)
	require.NotNil(t, result.ExecutionResult.ErrorMessage)
	assert.Equal(t, "Mint(InsufficientFunds)", *result.ExecutionResult.ErrorMessage)
}

func Test_SpeculativeExecResult_UnmarshalJSON_V1(t *testing.T) {
	server := SetupSpecExecServer(t, "../data/rpc_response/specexec_transfer_success_v158.json")
	defer server.Close()
	client := rpc.NewSpeculativeClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.SpeculativeExec(context.Background(), types.Deploy{}, nil)
	require.NoError(t, err)
	require.Nil(t, result.ExecutionResult.ErrorMessage)
	assert.Equal(t, "1.5.0", result.ApiVersion)
	assert.Equal(t, "1dd95f1d19fb96ebcbaa377de98cad991b30fce4aa8316eb2af0d7fa1f0276a6", result.BlockHash.ToHex())
	assert.Equal(t, uint64(100000000), result.ExecutionResult.Cost)
	assert.Equal(t, uint64(100000000), result.ExecutionResult.Consumed)
	assert.Equal(t, 5, len(result.ExecutionResult.Effects))
	assert.Equal(t, 1, len(result.ExecutionResult.Transfers))
	assert.Nil(t, result.Messages)
}

func Test_SpeculativeExecResult_UnmarshalJSON_V1_WithError(t *testing.T) {
	server := SetupSpecExecServer(t, "../data/rpc_response/specexec_transfer_failure_v158.json")
	defer server.Close()
	client := rpc.NewSpeculativeClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.SpeculativeExec(context.Background(), types.Deploy{}, nil)
	require.NoError(t, err)
	require.NotNil(t, result.ExecutionResult.ErrorMessage)
	assert.Equal(t, "Insufficient payment", *result.ExecutionResult.ErrorMessage)
	assert.Equal(t, "1.5.0", result.ApiVersion)
	assert.Equal(t, "eda864c7da3f5765a7027e3aa234b05c4b6c33efae51adacd224dc5f3c1c7958", result.BlockHash.ToHex())
	assert.Equal(t, uint64(100000000), result.ExecutionResult.Consumed)
	assert.Equal(t, 2, len(result.ExecutionResult.Effects))
	assert.Equal(t, 0, len(result.ExecutionResult.Transfers))
	assert.Nil(t, result.Messages)
}
