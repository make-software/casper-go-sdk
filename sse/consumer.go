package sse

import (
	"context"
	"errors"
	"fmt"
)

var ErrHandlerNotRegistered = errors.New("handler is not registered")

// HandlerFunc is the interface of function that should be implemented in each specific event handler.
type HandlerFunc func(context.Context, RawEvent) error

type Handler interface {
	Handle(context.Context, RawEvent) error
}

// Consumer is a service that registers event handlers and assigns events from the stream to specific handlers.
type Consumer struct {
	handlers map[EventType]HandlerFunc
}

func NewConsumer() *Consumer {
	return &Consumer{
		handlers: make(map[EventType]HandlerFunc),
	}
}

func (c *Consumer) RegisterHandler(eventType EventType, handler HandlerFunc) {
	c.handlers[eventType] = handler
}

func (c *Consumer) Run(ctx context.Context, events <-chan RawEvent, errCh chan<- error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case rawEvent, ok := <-events:
			if !ok {
				return errors.New("events stream was closed")
			}
			handler, ok := c.handlers[rawEvent.EventType]
			if !ok {
				return fmt.Errorf("%s, type: %s", ErrHandlerNotRegistered, AllEventsNames[rawEvent.EventType])
			}
			if err := handler(ctx, rawEvent); err != nil {
				errCh <- err
				continue
			}
		}
	}
}
