package fiware

import (
	resty "github.com/go-resty/resty/v2"
)

// ClientOption for the http client
type ClientOption func(*HTTPClient)

// ClientInf interface for all http clients
type ClientInf interface {
	NewClient() *HTTPClient
}

// HTTPClient Http HTTPClient wrapper
type HTTPClient struct {
	*resty.Client
	Host string
}

// Host sets the hostdomain for the server you want to communicate with (https://example.de)
func Host(host string) ClientOption {
	if host[len(host)-1:] == "/" {
		host = host[0 : len(host)-1]
	}
	return func(c *HTTPClient) {
		c.Host = host
	}
}

// NewClient creates a new http client
func NewClient(opts ...ClientOption) *HTTPClient {
	c := &HTTPClient{
		Client: resty.New(),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
