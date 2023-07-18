package rpc

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types"
)

func Test_Speculative_endpoint(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/speculative_exec.json")
	defer server.Close()
	var deployFixture types.Deploy
	fixture, err := os.ReadFile("../data/deploy/deploy_with_transfer.json")
	require.NoError(t, err)
	err = json.Unmarshal(fixture, &deployFixture)
	require.NoError(t, err)
	client := rpc.NewSpeculativeClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.SpeculativeExec(context.Background(), deployFixture, nil)
	require.NoError(t, err)
	assert.Equal(t, uint64(100000000), result.ExecutionResult.Success.Cost)
}
