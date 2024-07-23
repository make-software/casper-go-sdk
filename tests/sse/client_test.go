package sse

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/sse"
)

func Test_HttpConnection_request(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(json.RawMessage(`data: {"ApiVersion":"1.0.0"}`))
		require.NoError(t, err)
	}))

	client := sse.NewClient(server.URL)
	client.RegisterMiddleware(sse.Middleware(func(handler sse.HandlerFunc) sse.HandlerFunc {
		return func(ctx context.Context, event sse.RawEvent) error {
			log.Println("middleware registered")
			return handler(ctx, event)
		}
	}))
	client.RegisterHandler(sse.APIVersionEventType, func(ctx context.Context, event sse.RawEvent) error {
		data, err := event.ParseAsAPIVersionEvent()
		require.NoError(t, err)
		assert.Equal(t, "1.0.0", data.APIVersion)
		return nil
	})
	assert.Error(t, client.Start(context.Background(), 123))
}

func Test_withOneWorker_shouldProcessRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(json.RawMessage(`data: {"ApiVersion":"1.0.0"}`))
		require.NoError(t, err)
	}))

	client := sse.NewClient(server.URL)
	client.WorkersCount = 1
	var result bool
	client.RegisterHandler(sse.APIVersionEventType, func(ctx context.Context, event sse.RawEvent) error {
		result = true
		return nil
	})
	assert.Error(t, client.Start(context.Background(), -1))
	assert.True(t, result)
}
