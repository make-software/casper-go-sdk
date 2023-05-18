package sse

import (
	"context"
	"errors"
	"net/http"
	"strconv"
)

// HttpConnection is responsible to establish connection with SSE server.
// Create Request, handle http error and provide a response.
type HttpConnection struct {
	httpClient *http.Client
	Headers    map[string]string
	URL        string
}

func NewHttpConnection(httpClient *http.Client, sourceUrl string) *HttpConnection {
	return &HttpConnection{
		httpClient: httpClient,
		URL:        sourceUrl,
		Headers:    make(map[string]string),
	}
}

func (c *HttpConnection) Request(ctx context.Context, lastEventID int) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, c.URL, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	if lastEventID != 0 {
		query := req.URL.Query()
		query.Add("start_from", strconv.Itoa(lastEventID))
		req.URL.RawQuery = query.Encode()
	}

	// Add user specified headers
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error invalid connect response code")
	}

	return resp, err
}
