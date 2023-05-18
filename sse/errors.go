package sse

import (
	"errors"
	"fmt"
	"math"
)

type ErrUnknownEventType struct {
	error
	RawData []byte
}

func (e ErrUnknownEventType) Error() string {
	showLen := math.Min(float64(len(e.RawData)), 32)
	return fmt.Sprintf("%s, raw data: %s...", e.error.Error(), e.RawData[0:int(showLen)])
}

func NewErrUnknownEventType(data []byte) ErrUnknownEventType {
	return ErrUnknownEventType{
		error:   errors.New("event type has not registered"),
		RawData: data,
	}
}
