package rpc

import (
	"fmt"
	"net/http"
)

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h *RpcError) Error() string {
	return h.Message
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
