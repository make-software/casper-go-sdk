package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RpcError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (h *RpcError) Error() string {
	return fmt.Sprintf("key: %s, data: %s", h.Message, h.Data)
}

type HttpError struct {
	SourceErr  error
	StatusCode int
}

func (h *HttpError) Error() string {
	return fmt.Sprintf("Code: %d, err: %s", h.StatusCode, h.SourceErr.Error())
}

func (h *HttpError) Unwrap() error {
	return h.SourceErr
}

func (h *HttpError) IsNotFound() bool {
	return h.StatusCode == http.StatusNotFound
}
