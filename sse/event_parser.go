package sse

import (
	"bytes"
	"fmt"
	"strconv"
)

var (
	headerID   = []byte("id:")
	headerData = []byte("data:")
)

type EventParser struct {
	eventsToParse map[EventType]string
}

func NewEventParser() *EventParser {
	return &EventParser{eventsToParse: make(map[EventType]string)}
}

func (e *EventParser) RegisterEvent(eventType EventType) {
	e.eventsToParse[eventType] = AllEventsNames[eventType]
}

func (e *EventParser) ParseRawEvent(data []byte) (RawEvent, error) {
	var eventID []byte
	var eventData []byte
	for _, line := range bytes.FieldsFunc(data, func(r rune) bool { return r == '\n' || r == '\r' }) {
		switch {
		case bytes.HasPrefix(line, headerID):
			eventID = append([]byte(nil), trimPrefix(len(headerID), line)...)
		case bytes.HasPrefix(line, headerData):
			// The spec allows for multiple data fields per event, concatenated them with "\n".
			eventData = append([]byte(nil), append(trimPrefix(len(headerData), line), byte('\n'))...)
		default:
			// Ignore any garbage that doesn't match what we're looking for.
		}
	}
	return e.parseEvenType(eventID, eventData)
}

func (e *EventParser) parseEvenType(eventID, eventData []byte) (RawEvent, error) {
	var eventType EventType
	for eType, typeName := range e.eventsToParse {
		trimmedData := trimPrefix(len([]byte("{\"")), eventData)
		if bytes.HasPrefix(trimmedData, []byte(typeName)) {
			eventType = eType
			break
		}
	}

	if eventType == 0 {
		return RawEvent{}, NewErrUnknownEventType(eventData)
	}

	if eventType == APIVersionEventType {
		return RawEvent{
			EventType: eventType,
			Data:      eventData,
		}, nil
	}

	parsedID, err := strconv.ParseUint(string(eventID), 10, 0)
	if err != nil {
		return RawEvent{}, fmt.Errorf("error parsing event id, %w", err)
	}

	return RawEvent{
		EventType: eventType,
		EventID:   parsedID,
		Data:      eventData,
	}, nil
}

func trimPrefix(size int, data []byte) []byte {
	if data == nil || len(data) < size {
		return data
	}

	data = data[size:]
	// Remove optional leading whitespace
	if len(data) > 0 && data[0] == 32 {
		data = data[1:]
	}
	// Remove trailing new line
	if len(data) > 0 && data[len(data)-1] == 10 {
		data = data[:len(data)-1]
	}
	return data
}
