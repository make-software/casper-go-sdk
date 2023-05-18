package sse

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"time"
)

var ErrFullStreamTimeoutError = errors.New("can't fill the stream, because it full")

// Streamer is a service that main responsibility is to fill the events' channel.
// The Connection management is isolated in this service. Service uses a HttpConnection to get HTTP response as a stream
// resource and provides it to the EventStreamReader, that supposes to parse bytes from the response's body.
// This design assumes to manage the connection and provide reconnection logic above of this service.
type Streamer struct {
	Connection   *HttpConnection
	eventParser  *EventParser
	StreamReader *EventStreamReader
	// This duration allows the stream's buffer to stay in fill up completely state, which could indicate
	// that the workers are working too slowly and have not received any messages.
	// If this period elapses without any messages being received, an ErrFullStreamTimeoutError will be thrown.
	BlockedStreamLimit time.Duration
}

// NewStreamer is the idiomatic way to create Streamer
func NewStreamer(
	client *HttpConnection,
	reader *EventStreamReader,
	blockedStreamLimit time.Duration,
) *Streamer {
	return &Streamer{
		Connection:         client,
		StreamReader:       reader,
		BlockedStreamLimit: blockedStreamLimit,
		eventParser:        NewEventParser(),
	}
}

// DefaultStreamer is a shortcut to fast start with Streamer
func DefaultStreamer(url string) *Streamer {
	return NewStreamer(
		NewHttpConnection(http.DefaultClient, url),
		&EventStreamReader{
			MaxBufferSize: 1024 * 1024 * 50, // 50 MB
		},
		30*time.Second,
	)
}

func (i *Streamer) RegisterEvent(eventType EventType) {
	i.eventParser.RegisterEvent(eventType)
}

func (i *Streamer) FillStream(ctx context.Context, lastEventID int, stream chan<- RawEvent, errorsCh chan<- error) error {
	response, err := i.Connection.Request(ctx, lastEventID)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	i.StreamReader.RegisterStream(response.Body)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			evenBytes, err := i.StreamReader.ReadEvent()
			if err != nil {
				return err
			}
			// Ignore empty events.
			if bytes.Equal(evenBytes, []byte(":")) {
				continue
			}
			eventData, err := i.eventParser.ParseRawEvent(evenBytes)
			if err != nil {
				errorsCh <- err
				continue
			}
			err = i.addData(ctx, stream, eventData)
			if err != nil {
				return err
			}
		}
	}
}

func (i *Streamer) addData(ctx context.Context, stream chan<- RawEvent, data RawEvent) error {
	stackTimer := time.NewTicker(i.BlockedStreamLimit)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case stream <- data:
		stackTimer.Reset(i.BlockedStreamLimit)
		return nil
	case <-stackTimer.C:
		return ErrFullStreamTimeoutError
	}
}
