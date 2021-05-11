package iotmanager

import (
	"fmt"

	"github.com/FeSaRoTi/fiware-lib/pkg/fiware"
)

//Client defines the IoT-Manager client for requests
type Client struct {
	httpClient    *fiware.HTTPClient
	ProtocolsPath string
	Host          string
	FiwareConfig  *fiware.Config
}

//ClientOpts for functional options
type ClientOpts func(*Client)

// ProtocolPath sets the path for the iot-manager for reqeustinig all iot-devices. defaul is /iot/parotocols
func ProtocolPath(path string) ClientOpts {
	return func(c *Client) {
		c.ProtocolsPath = path
	}
}

// Host sets the host of the iot-manager (https://iot-manager.com)
func Host(host string) ClientOpts {
	return func(c *Client) {
		c.Host = host
	}
}

// FiwareClient set the httpclient for the iot-manager
func FiwareClient(client *fiware.HTTPClient) ClientOpts {
	return func(c *Client) {
		c.httpClient = client
	}
}

// FiwareConfig sets the fiware configuration for the iot manager
func FiwareConfig(conf *fiware.Config) ClientOpts {
	return func(c *Client) {
		c.FiwareConfig = conf
	}
}

//NewClient creates a new IoT-Manager Client
func NewClient(opts ...ClientOpts) *Client {
	c := &Client{
		httpClient:    fiware.NewClient(),
		ProtocolsPath: "/iot/protocols",
		FiwareConfig:  fiware.NewConfig(),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

//ListIoTAgents lists all subscribes IoT-Agents from the IoT-Manager
func (c *Client) ListIoTAgents() (*fiware.IoTManagerGetResponse, error) {
	respObj := &fiware.IoTManagerGetResponse{}
	resp, err := c.httpClient.R().SetResult(respObj).SetHeaders(c.FiwareConfig.GetHeader()).
		SetHeader("Accept", "application/json").Get(fmt.Sprintf("%s%s", c.Host, c.ProtocolsPath))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("IoT-Manager response code: %d", resp.StatusCode())
	}
	return respObj, nil
}
