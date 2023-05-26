package helper

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

type LogTestTransport struct {
	base http.Transport
}

func (l *LogTestTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	log.Printf("Request URL: %s", request.URL.String())
	log.Printf("Request Context: %v", request.Context())
	bodyBytes, err := io.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	bodyString := string(bodyBytes)
	log.Printf("Request Body: %s", bodyString)

	resp, err := l.base.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	log.Printf("Response URL: %s", resp.Request.URL.String())
	log.Printf("Response Status: %v", resp.StatusCode)
	bodyRespBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	bodyRespString := string(bodyRespBytes)
	log.Printf("Response Body: %s", bodyRespString)

	return resp, err
}
