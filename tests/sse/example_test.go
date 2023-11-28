//go:build example
// +build example

package sse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/sse"
)

func Test_Simple_Example(t *testing.T) {
	client := sse.NewClient("http://52.3.38.81:9999/events/main")
	defer client.Stop()
	client.RegisterHandler(sse.DeployProcessedEventType, func(ctx context.Context, rawEvent sse.RawEvent) error {
		log.Printf("eventID: %d, raw data: %s", rawEvent.EventID, rawEvent.Data)
		return nil
	})
	lastEventID := 1234
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()
	assert.Error(t, client.Start(ctx, lastEventID))
}

func Test_Example_ParseDeploy(t *testing.T) {
	client := sse.NewClient("http://52.3.38.81:9999/events/main")
	defer client.Stop()
	client.RegisterHandler(sse.DeployProcessedEventType, func(ctx context.Context, rawEvent sse.RawEvent) error {
		deploy, err := rawEvent.ParseAsDeployProcessedEvent()
		if err != nil {
			return err
		}
		log.Printf("Deploy hash: %s", deploy.DeployProcessed.DeployHash)
		return nil

	})
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()
	assert.Error(t, client.Start(ctx, 1234))
}

func Test_Example_WithErrorHandling(t *testing.T) {
	client := sse.NewClient("http://52.3.38.81:9999/events/main")
	defer client.Stop()
	client.ConsumerErrorHandler = func(errors <-chan error) {
		for one := range errors {
			fmt.Println(one.Error())
		}
	}
	client.RegisterHandler(sse.APIVersionEventType, TestHandler{}.ProcessEvent)
	client.RegisterHandler(sse.DeployProcessedEventType, TestHandler{}.ProcessEvent)
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()
	assert.Error(t, client.Start(ctx, 1234))
}

func Test_Example_WithCancelContext(t *testing.T) {
	client := sse.NewClient("http://52.3.38.81:9999/events/main")
	client.RegisterHandler(sse.DeployProcessedEventType, func(ctx context.Context, event sse.RawEvent) error {
		return errors.New("critical")
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client.ConsumerErrorHandler = func(errors <-chan error) {
		for one := range errors {
			fmt.Println(one)
			if one.Error() == "critical" {
				cancel()
			}
		}
	}
	assert.Error(t, client.Start(ctx, 1234))
}

func Test_Example_CloseSlowConsuming(t *testing.T) {
	client := sse.NewClient("http://52.3.38.81:9999/events/main")
	client.Streamer.BlockedStreamLimit = 900 * time.Millisecond
	client.RegisterHandler(sse.BlockAddedEventType, func(ctx context.Context, event sse.RawEvent) error {
		time.Sleep(time.Second)
		return nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	assert.Error(t, client.Start(ctx, 1234))
}

func Test_Example_WithMiddleware(t *testing.T) {
	client := sse.NewClient("http://52.3.38.81:9999/events/main")
	client.RegisterMiddleware(WithLogs)
	client.RegisterHandler(sse.BlockAddedEventType, TestHandler{}.ProcessEvent)
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()
	assert.Error(t, client.Start(ctx, 1234))
}

func Test_Example_WithContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	client := sse.NewClient("http://52.3.38.81:9999/events/main")
	client.RegisterHandler(sse.BlockAddedEventType, TestHandler{}.ProcessEvent)
	assert.Error(t, client.Start(ctx, 1234))
}

func Test_Example_With5Workers(t *testing.T) {
	client := sse.NewClient("http://52.3.38.81:9999/events/main")
	client.WorkersCount = 5
	client.RegisterHandler(sse.BlockAddedEventType, func(ctx context.Context, event sse.RawEvent) error {
		log.Printf("Worker: %d, eventID: %d", ctx.Value("workerID"), event.EventID)
		return nil
	})
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()
	assert.Error(t, client.Start(ctx, 1234))
}

func Test_SSE_WithConfigurations(t *testing.T) {
	// This is example how to create and configure services in idiomatic way
	streamer := sse.NewStreamer(
		sse.NewHttpConnection(
			&http.Client{
				Transport: &http.Transport{
					IdleConnTimeout:       20 * time.Minute,
					ResponseHeaderTimeout: time.Second * 30,
				},
				Timeout: 30 * time.Second,
			},
			"http://52.3.38.81:9999/events/main",
		),
		&sse.EventStreamReader{MaxBufferSize: 1024 * 1024 * 15}, // 15 MB
		2*time.Minute,
	)
	consumer := sse.NewConsumer()
	_ = streamer
	_ = consumer
}

func Test_Client_WithAuthorizationHeader(t *testing.T) {
	authToken := "1234567890"
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		auth := request.Header.Get("Authorization")
		if auth != authToken {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, err := writer.Write(json.RawMessage(`data: {"ApiVersion":"1.0.0"}`))
		require.NoError(t, err)
	}))

	client := sse.NewClient(server.URL)
	client.Streamer.Connection.Headers = map[string]string{"Authorization": authToken}
	ctx, cancel := context.WithCancel(context.Background())
	client.RegisterHandler(sse.APIVersionEventType, func(ctx context.Context, event sse.RawEvent) error {
		data, err := event.ParseAsAPIVersionEvent()
		require.NoError(t, err)
		assert.Equal(t, "1.0.0", data.APIVersion)
		cancel()
		return nil
	})
	err := client.Start(context.Background(), 123)
	if err != io.EOF {
		require.NoError(t, err)
	}
	<-ctx.Done()
}
