package sse

import (
	"context"
	"log"

	"golang.org/x/sync/errgroup"
)

type CtxWorkerID string

const CtxWorkerIDKey CtxWorkerID = "workerID"

type MiddlewareHandler interface {
	Process(handler HandlerFunc) HandlerFunc
}

type Middleware func(handler HandlerFunc) HandlerFunc

// Client is a facade that provide convenient interface to process data from the stream,
// and unites Streamer and Consumer under implementation. Also, the Client allows to register global middleware
// that will be applied for all handlers.
type Client struct {
	Streamer *Streamer
	Consumer *Consumer

	EventStream          chan RawEvent
	streamErrors         chan error
	consumerErrors       chan error
	StreamErrorHandler   func(<-chan error)
	ConsumerErrorHandler func(<-chan error)
	middlewares          []Middleware
	WorkersCount         int
}

func NewClient(url string) *Client {
	return &Client{
		Streamer:             DefaultStreamer(url),
		Consumer:             NewConsumer(),
		EventStream:          make(chan RawEvent, 10),
		streamErrors:         make(chan error, 1),
		consumerErrors:       make(chan error, 1),
		StreamErrorHandler:   logErrors,
		ConsumerErrorHandler: logErrors,
		WorkersCount:         1,
	}
}

func (p *Client) Start(ctx context.Context, lastEventID int) error {
	groupErrs, ctx := errgroup.WithContext(ctx)
	groupErrs.Go(func() error {
		return p.Streamer.FillStream(ctx, lastEventID, p.EventStream, p.streamErrors)
	})
	for i := 0; i < p.WorkersCount; i++ {
		newCtx := context.WithValue(ctx, CtxWorkerIDKey, i)
		groupErrs.Go(func() error {
			return p.Consumer.Run(newCtx, p.EventStream, p.consumerErrors)
		})
	}
	go p.StreamErrorHandler(p.streamErrors)
	go p.ConsumerErrorHandler(p.consumerErrors)
	return groupErrs.Wait()
}

func (p *Client) RegisterMiddleware(one Middleware) {
	p.middlewares = append(p.middlewares, one)
}

func (p *Client) RegisterHandler(eventType EventType, handler HandlerFunc) {
	p.Streamer.RegisterEvent(eventType)
	// Loop backwards through the middleware invoking each one. Replace the
	// handler with the new wrapped handler. Looping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(p.middlewares) - 1; i >= 0; i-- {
		handler = p.middlewares[i](handler)
	}
	p.Consumer.RegisterHandler(eventType, handler)
}

func logErrors(source <-chan error) {
	for one := range source {
		log.Println(one)
	}
}

func (p *Client) Stop() {
	close(p.EventStream)
	close(p.streamErrors)
	close(p.consumerErrors)
}
