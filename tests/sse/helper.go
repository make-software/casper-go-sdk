package sse

import (
	"context"
	"errors"
	"log"

	"github.com/make-software/casper-go-sdk/sse"
)

type TestHandler struct{}

func (t TestHandler) ProcessEvent(ctx context.Context, event sse.RawEvent) error {
	log.Println(ctx.Value("EventID"))
	log.Printf("handler %s", event.Data)
	return errors.New("test handler err")
}

func WithLogs(handler sse.HandlerFunc) sse.HandlerFunc {
	return func(ctx context.Context, event sse.RawEvent) error {
		log.Println("Log middleware")
		log.Println(event.EventID)
		return handler(ctx, event)
	}
}
