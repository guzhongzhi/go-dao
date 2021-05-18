package client

import (
	"bytes"
	"github.com/guzhongzhi/gmicro/logger"
	"io"
	"net/http"
)

type HTTPClient interface {
	GET(u string, in io.Reader, contentType string, headers map[string]string) (*http.Response, error)
	POST(u string, in io.Reader, contentType string, headers map[string]string) (*http.Response, error)
	PUT(u string, in io.Reader, contentType string, headers map[string]string) (*http.Response, error)
	PATCH(u string, in io.Reader, contentType string, headers map[string]string) (*http.Response, error)
	DELETE(u string, in io.Reader, contentType string, headers map[string]string) (*http.Response, error)
}

func NewHTTPClient(l logger.SuperLogger) HTTPClient {
	if l == nil {
		l = logger.Default()
	}
	return &httpClient{
		logger: l,
	}
}

type httpClient struct {
	logger logger.SuperLogger
}

func (s *httpClient) do(method string, u string, body io.Reader, contentType string, headers map[string]string) (*http.Response, error) {
	if contentType == "" {
		contentType = "application/x-www-form-urlencoded"
	}
	if body == nil {
		body = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	s.logger.Debugf("start to send http request to '%s', method=%s, headers=%v,content-type=%v", u, method, headers, contentType)
	req.Header.Set("Content-Type", contentType)
	s.logger.Debugf("end send http request to '%s'", u)

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	client := http.DefaultClient
	return client.Do(req)
}

func (s *httpClient) GET(u string, in io.Reader, contentType string, headers map[string]string) (*http.Response, error) {
	return s.do("GET", u, in, contentType, headers)
}

func (s *httpClient) POST(u string, in io.Reader, contentType string, headers map[string]string) (*http.Response, error) {
	return s.do("POST", u, in, contentType, headers)
}

func (s *httpClient) PUT(u string, in io.Reader, contentType string, headers map[string]string) (*http.Response, error) {
	return s.do("PUT", u, in, contentType, headers)
}

func (s *httpClient) PATCH(u string, in io.Reader, contentType string, headers map[string]string) (*http.Response, error) {
	return s.do("PATCH", u, in, contentType, headers)
}

func (s *httpClient) DELETE(u string, in io.Reader, contentType string, headers map[string]string) (*http.Response, error) {
	return s.do("DELETE", u, in, contentType, headers)
}
