package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	ErrParamsUnmarshalHandler  = errors.New("failed to marshal rpc request's params")
	ErrBuildHttpRequestHandler = errors.New("failed to build http request")
	ErrProcessHttpRequest      = errors.New("failed to sent http request")
	ErrReadHttpResponseBody    = errors.New("failed to read http response body")
	ErrRpcResponseUnmarshal    = errors.New("failed to unmarshal rpc response")
)

// HttpHandler implements Handler interface using the HTTP protocol under the implementation
type HttpHandler struct {
	httpClient    *http.Client
	endpoint      string
	CustomHeaders map[string]string
}

// NewHttpHandler is a constructor for HttpHandler that suppose to configure http.Client
// examples of usage can be found here [Test_ConfigurableClient_GetDeploy]
func NewHttpHandler(endpoint string, client *http.Client) *HttpHandler {
	return &HttpHandler{
		httpClient: client,
		endpoint:   endpoint,
	}
}

// ProcessCall operates with an external RPC server through HTTP. It builds and processes the request,
// reads a response and handles errors. All logic with HTTP interaction is isolated here and can be replaced with
// other (more efficient) protocols.
func (c *HttpHandler) ProcessCall(ctx context.Context, params RpcRequest) (RpcResponse, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return RpcResponse{}, fmt.Errorf("%w, details: %s", ErrParamsUnmarshalHandler, err.Error())
	}

	request, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return RpcResponse{}, fmt.Errorf("%w, details: %s", ErrBuildHttpRequestHandler, err.Error())
	}
	request.Header.Add("Content-Type", "application/json")
	for name, val := range c.CustomHeaders {
		request.Header.Add(name, val)
	}
	request = request.WithContext(ctx)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return RpcResponse{}, fmt.Errorf("%w, details: %s", ErrProcessHttpRequest, err.Error())
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return RpcResponse{}, fmt.Errorf("http error from rpc, %w", &HttpError{
			SourceErr:  errors.New(resp.Status),
			StatusCode: resp.StatusCode,
		})
	}

	var rpcResponse RpcResponse
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return RpcResponse{}, fmt.Errorf("%w, details: %s", ErrReadHttpResponseBody, err.Error())
	}

	os.WriteFile("test.json", b, 0666)

	err = json.Unmarshal(b, &rpcResponse)
	if err != nil {
		return RpcResponse{}, fmt.Errorf("%w, details: %s", ErrRpcResponseUnmarshal, err.Error())
	}

	return rpcResponse, nil
}
