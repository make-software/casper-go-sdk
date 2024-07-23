package helper

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/make-software/casper-go-sdk/v2/rpc"
)

var delayedTestNetHandlerState rpc.Handler

type delayedTestNetHandler struct {
	handler     rpc.Handler
	lastRpcCall time.Time
	delay       time.Duration
	mu          sync.Mutex
}

func (h *delayedTestNetHandler) ProcessCall(ctx context.Context, params rpc.RpcRequest) (rpc.RpcResponse, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	sub := time.Since(h.lastRpcCall)
	time.Sleep(h.delay - sub)
	h.lastRpcCall = time.Now()
	fmt.Println(h.lastRpcCall)
	return h.handler.ProcessCall(ctx, params)
}

func TestRpcClient(endpoint string) rpc.Client {
	return rpc.NewClient(TestRpcHttpHandler(endpoint, http.DefaultClient))
}

// TestRpcHttpHandler is a default handler that uses http.DefaultClient
func TestRpcHttpHandler(endpoint string, client *http.Client) rpc.Handler {
	if delayedTestNetHandlerState == nil {
		delayedTestNetHandlerState = &delayedTestNetHandler{
			handler:     rpc.NewHttpHandler(endpoint, client),
			lastRpcCall: time.Now(),
			delay:       time.Second * 3,
		}
	}
	return delayedTestNetHandlerState
}

func GetTestRpcClient(t *testing.T) rpc.Client {
	nodeUrl, exist := os.LookupEnv("TEST_NODE_RPC_API_URL")
	if !exist {
		t.Error("env TEST_NODE_RPC_API_URL not found")
	}

	return TestRpcClient(nodeUrl)
}

func GetTestLoggerClient(t *testing.T) rpc.Client {
	nodeUrl, exist := os.LookupEnv("TEST_NODE_RPC_API_URL")
	if !exist {
		t.Error("env TEST_NODE_RPC_API_URL not found")
	}

	httpClient := &http.Client{Transport: &LogTestTransport{}}
	handler := rpc.NewHttpHandler(nodeUrl, httpClient)
	loggerDecorator := NewTestLoggerDecorator(handler)
	return rpc.NewClient(loggerDecorator)
}
