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

type IoTAgentError struct {
	Name       string `json:"name"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code,omitempty"`
}

func (i *IoTAgentError) Error() string {
	return fmt.Sprintf("%s: %s - http code: %v", i.Name, i.Message, i.StatusCode)
}

type ClientOpts func(*Client)

func NewClient(opts ...ClientOpts) *Client {
	client := &Client{
		FiwareConfig: fiware.NewConfig(),
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// HTTPClient sets the http client
func HTTPClient(client *fiware.HTTPClient) ClientOpts {

	return func(c *Client) {
		c.httpClient = client
	}
}

//WithHost sets the host of the iot-agent (https://iot-agent.de)
func WithHost(host string) ClientOpts {
	if host[len(host)-1:] == "/" {
		host = host[0 : len(host)-1]
	}
	return func(c *Client) {
		c.Host = host
	}
}

//WithResource sets the resource of the iotagent
func WithResource(res string) ClientOpts {
	return func(c *Client) {
		c.Resource = res
	}
}

// WithProtocol sets the protocol which is provided by the iot-agent
func WithProtocol(prot string) ClientOpts {
	return func(c *Client) {
		c.Protocol = prot
	}
}

// WithFiwareConfig sets the fiware configuration for this iot-agent
func WithFiwareConfig(conf *fiware.Config) ClientOpts {
	return func(c *Client) {
		c.FiwareConfig = conf
	}
}

// About returns the configuration from the iot-agent
func (c *Client) About() (*fiware.IoTAgentAboutResponse, error) {
	respObj := &fiware.IoTAgentAboutResponse{}
	resp, err := c.httpClient.R().SetResult(respObj).SetError(&IoTAgentError{}).
		SetHeaders(c.FiwareConfig.GetHeader()).Get(fmt.Sprintf("%s/iot/about", c.Host))
	if err != nil {
		return nil, err
	}
	if statusCode := resp.StatusCode(); statusCode != 200 {
		httpErr := resp.Error().(*IoTAgentError)
		httpErr.StatusCode = statusCode
		return nil, httpErr
	}
	return respObj, nil
}

// CreateService creates a service group in the iot-agent
func (c *Client) CreateService(services *fiware.IoTAgentCreateServiceGroupRequest) error {
	resp, err := c.httpClient.R().SetBody(services).SetError(&IoTAgentError{}).
		SetHeaders(c.FiwareConfig.GetHeader()).Post(fmt.Sprintf("%s/iot/services", c.Host))
	if err != nil {
		return err
	}
	if statusCode := resp.StatusCode(); statusCode != 200 {
		httpErr := resp.Error().(*IoTAgentError)
		httpErr.StatusCode = statusCode
		return httpErr
	}
	return nil
}

// UpdateService updates a service group which is identified by api key and resource
func (c *Client) UpdateService(services interface{}) error {
	resp, err := c.httpClient.R().SetBody(services).SetError(&IoTAgentError{}).SetHeaders(c.FiwareConfig.GetHeader()).Put(fmt.Sprintf("%s/iot/services", c.Host))
	if err != nil {
		return err
	}
	if statusCode := resp.StatusCode(); statusCode != 200 {
		httpErr := resp.Error().(*IoTAgentError)
		httpErr.StatusCode = statusCode
		return httpErr
	}
	return nil
}

// GetService retuns a list of service groups
func (c *Client) GetService() (*fiware.IoTAgentGetServicesResponse, error) {
	resp, err := c.httpClient.R().
		SetResult(&fiware.IoTAgentGetServicesResponse{}).SetError(&IoTAgentError{}).
		SetHeaders(c.FiwareConfig.GetHeader()).Get(fmt.Sprintf("%s/iot/services", c.Host))
	if err != nil {
		return nil, err
	}
	if statusCode := resp.StatusCode(); statusCode != 200 {
		httpErr := resp.Error().(*IoTAgentError)
		httpErr.StatusCode = statusCode
		return nil, httpErr
	}
	return resp.Result().(*fiware.IoTAgentGetServicesResponse), nil
}

// DeleteService deletes a service group specified by the apikey and the resource
func (c *Client) DeleteService(resource string, apikey string) error {
	resp, err := c.httpClient.R().SetError(&IoTAgentError{}).
		SetHeaders(c.FiwareConfig.GetHeader()).
		SetQueryParam("apikey", apikey).SetQueryParam("resource", resource).Delete(fmt.Sprintf("%s/iot/services", c.Host))
	if err != nil {
		return err
	}
	if statusCode := resp.StatusCode(); statusCode != 200 {
		httpErr := resp.Error().(*IoTAgentError)
		httpErr.StatusCode = statusCode
		return httpErr
	}
	return nil
}

// CreateDevice creates a device in the iot-agent
func (c *Client) CreateDevice(device *fiware.IoTAgentCreateDeviceRequest) error {
	resp, err := c.httpClient.R().SetError(&IoTAgentError{}).SetBody(device).
		SetHeaders(c.FiwareConfig.GetHeader()).Post(fmt.Sprintf("%s/iot/devices", c.Host))
	if err != nil {
		return err
	}
	if statusCode := resp.StatusCode(); statusCode != 200 {
		httpErr := resp.Error().(*IoTAgentError)
		httpErr.StatusCode = statusCode
		return httpErr
	}
	return nil
}

// GetDevices gets all devices saved in the iot-agent
func (c *Client) GetDevices() (*fiware.IoTAgentGetDevicesResponse, error) {
	resp, err := c.httpClient.R().SetError(&IoTAgentError{}).SetBody(&fiware.IoTAgentGetDevicesResponse{}).
		SetHeaders(c.FiwareConfig.GetHeader()).Get(fmt.Sprintf("%s/iot/devices", c.Host))
	if err != nil {
		return nil, err
	}
	if statusCode := resp.StatusCode(); statusCode != 200 {
		httpErr := resp.Error().(*IoTAgentError)
		httpErr.StatusCode = statusCode
		return nil, httpErr
	}
	return resp.Result().(*fiware.IoTAgentGetDevicesResponse), nil
}

// UpdateDevice updates a device in the iot-agent
func (c *Client) UpdateDevice(device interface{}) error {
	resp, err := c.httpClient.R().SetError(&IoTAgentError{}).SetBody(device).
		SetHeaders(c.FiwareConfig.GetHeader()).Put(fmt.Sprintf("%s/iot/devices", c.Host))
	if err != nil {
		return err
	}
	if statusCode := resp.StatusCode(); statusCode != 200 {
		httpErr := resp.Error().(*IoTAgentError)
		httpErr.StatusCode = statusCode
		return httpErr
	}
	return nil
}

// DeleteDevice deletes a device by the given device identified
func (c *Client) DeleteDevice(id string) error {
	resp, err := c.httpClient.R().SetError(&IoTAgentError{}).SetHeaders(c.FiwareConfig.GetHeader()).
		Delete(fmt.Sprintf("%s/iot/devices/:%s", c.Host, id))
	if err != nil {
		return err
	}
	if statusCode := resp.StatusCode(); statusCode != 200 {
		httpErr := resp.Error().(*IoTAgentError)
		httpErr.StatusCode = statusCode
		return httpErr
	}
	return nil
}
