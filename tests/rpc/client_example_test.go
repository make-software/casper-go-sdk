//go:build example
// +build example

package rpc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/rpc"
	"github.com/make-software/casper-go-sdk/v2/tests/helper"
	"github.com/make-software/casper-go-sdk/v2/types"
)

func Test_ConfigurableClient_GetDeploy(t *testing.T) {
	deployHash := "0009ea4441f4700325d9c38b0b6df415537596e1204abe4f6a94b6996aebf2f1"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fixture, err := os.ReadFile("../data/deploy/get_raw_rpc_deploy.json")
		require.NoError(t, err)
		rw.Write(fixture)
	}))
	defer server.Close()

	handler := rpc.NewHttpHandler(server.URL, http.DefaultClient)
	loggerDecorator := helper.NewTestLoggerDecorator(handler)
	client := rpc.NewClient(loggerDecorator)
	ctx := context.Background()
	ctx = rpc.WithRequestId(ctx, 123)
	result, err := client.GetDeploy(ctx, deployHash)
	require.NoError(t, err)

	assert.Equal(t, deployHash, result.Deploy.Hash.ToHex())
}

func Test_Client_WithRetries_GetDeploy(t *testing.T) {
	deployHash := "0009ea4441f4700325d9c38b0b6df415537596e1204abe4f6a94b6996aebf2f1"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	handler := rpc.NewHttpHandler(server.URL, http.DefaultClient)
	loggerDecorator := helper.NewTestLoggerDecorator(handler)
	retriesDecorator := helper.NewTestRetriesDecorator(loggerDecorator, 4, 2)
	client := rpc.NewClient(retriesDecorator)

	_, err := client.GetDeploy(context.Background(), deployHash)
	assert.Error(t, err)
}

func Test_SpeculativeExec(t *testing.T) {
	fixture, err := os.ReadFile("../data/deploy/deploy_with_transfer.json")
	require.NoError(t, err)
	var deployFixture types.Deploy
	err = json.Unmarshal(fixture, &deployFixture)
	require.NoError(t, err)
	client := rpc.NewSpeculativeClient(casper.NewRPCHandler("http://127.0.0.1:25102/rpc", http.DefaultClient))
	result, err := client.SpeculativeExec(context.Background(), deployFixture, nil)
	require.NoError(t, err)
	assert.Equal(t, uint64(100000000), result.ExecutionResult.Cost)
}

func Test_Client_RPCGetStatus_WithAuthorizationHeader(t *testing.T) {
	authToken := "1234567890"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		if auth != authToken {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		fixture, err := os.ReadFile("../data/rpc_response/get_status.json")
		require.NoError(t, err)
		rw.Write(fixture)
		rw.WriteHeader(http.StatusOK)
	}))
	handler := casper.NewRPCHandler(server.URL, http.DefaultClient)
	handler.CustomHeaders = map[string]string{"Authorization": authToken}
	client := casper.NewRPCClient(handler)

	status, err := client.GetStatus(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "2.0.0", status.APIVersion)
	assert.Equal(t, "5.4.3", status.ProtocolVersion)
	assert.NotEmpty(t, status.LatestSwitchBlockHash)
}
