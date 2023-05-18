//go:build example
// +build example

package rpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/tests/helper"
)

func Test_DefaultClient_GetDeploy_Example(t *testing.T) {
	deployHash := "0009ea4441f4700325d9c38b0b6df415537596e1204abe4f6a94b6996aebf2f1"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fixture, err := os.ReadFile("../data/deploy/get_raw_rpc_deploy.json")
		require.NoError(t, err)
		rw.Write(fixture)
	}))
	defer server.Close()

	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	deployResult, err := client.GetDeploy(context.Background(), deployHash)
	require.NoError(t, err)

	assert.Equal(t, deployHash, deployResult.Deploy.Hash.ToHex())
}

func Test_ConfigurableClient_GetDeploy(t *testing.T) {
	deployHash := "0009ea4441f4700325d9c38b0b6df415537596e1204abe4f6a94b6996aebf2f1"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fixture, err := os.ReadFile("../data/deploy/get_raw_rpc_deploy.json")
		require.NoError(t, err)
		rw.Write(fixture)
	}))
	defer server.Close()

	httpClient := &http.Client{Transport: &helper.LogTestTransport{}}
	handler := rpc.NewHttpHandler(server.URL, httpClient)
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
