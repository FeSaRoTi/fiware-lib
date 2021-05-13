package iotagent

// https://iotagent-node-lib.readthedocs.io/en/latest/api/index.html
import (
	"fmt"

	"github.com/FeSaRoTi/fiware-lib/pkg/fiware"
)

//Client defines the iot-agent client for requests
type Client struct {
	httpClient   *fiware.HTTPClient
	Host         string
	FiwareConfig *fiware.Config
	Resource     string
	Protocol     string
}

type ClientOpts func(*Client)

//Host sets the host of the iot-agent (https://iot-agent.de)
func Host(host string) ClientOpts {
	if host[len(host)-1:] == "/" {
		host = host[0 : len(host)-1]
	}
	return func(c *Client) {
		c.Host = host
	}
}

//Resource sets the resource of the iotagent
func Resource(res string) ClientOpts {
	return func(c *Client) {
		c.Resource = res
	}
}

// Protocol sets the protocol which is provided by the iot-agent
func Protocol(prot string) ClientOpts {
	return func(c *Client) {
		c.Protocol = prot
	}
}

// FiwareConfig sets the fiware configuration for this iot-agent
func FiwareConfig(conf *fiware.Config) ClientOpts {
	return func(c *Client) {
		c.FiwareConfig = conf
	}
}

// About returns the configuration from the iot-agent
func (c *Client) About() (*fiware.IoTAgentAboutResponse, error) {
	respObj := &fiware.IoTAgentAboutResponse{}
	resp, err := c.httpClient.R().SetResult(respObj).
		SetHeaders(c.FiwareConfig.GetHeader()).Get(fmt.Sprintf("%s/iot/about", c.Host))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("IoT-Agent response code: %d", resp.StatusCode())
	}
	return respObj, nil
}

// CreateService creates a service group in the iot-agent
func (c *Client) CreateService(service interface{}) error {
	resp, err := c.httpClient.R().SetBody(service).SetHeaders(c.FiwareConfig.GetHeader()).Post(fmt.Sprintf("%s/iot/services", c.Host))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("IoT-Agent respons with error code %d", resp.StatusCode())
	}
	return nil
}

// GetServiceGroups retuns a list of service groups
func (c *Client) GetServiceGroups() (*fiware.IoTAgentGetServicesResponse, error) {
	respObj := &fiware.IoTAgentGetServicesResponse{}
	resp, err := c.httpClient.R().SetResult(respObj).SetHeaders(c.FiwareConfig.GetHeader()).Get(fmt.Sprintf("%s/iot/services", c.Host))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("IoT-Agent respons with error code %d", resp.StatusCode())
	}
	return respObj, nil
}
