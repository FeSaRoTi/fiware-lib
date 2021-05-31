package iotagent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/FeSaRoTi/fiware-lib/pkg/fiware"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

type ioTAgentTestServer struct {
	obj interface{}
}

func (i *ioTAgentTestServer) saveObject(obj interface{}) {
	i.obj = obj
}

// createTestServer creates a mocked iot-agent which behaves a bit like a real iot-agent
func createTestServer(client *http.Client, host string) *ioTAgentTestServer {
	i := &ioTAgentTestServer{}
	httpmock.ActivateNonDefault(client)
	// Get response mock
	httpmock.RegisterResponder("GET", fmt.Sprintf("=~^%s/iot/(services|devices)", host), func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("fiware-service") == "" || req.Header.Get("fiware-servicepath") == "" {
			return httpmock.NewJsonResponse(400, IoTAgentError{Name: "MISSING_HEADERS", Message: "Some headers were missing from the request: [\"fiware-service\",\"fiware-servicepath\"]"})
		}
		if req.Header.Get("fiware-service") == "fail" {
			return httpmock.NewJsonResponse(400, IoTAgentError{
				Name:    "WRONG_SYNTAX",
				Message: "Failed because you want to.",
			})
		}
		resp, err := httpmock.NewJsonResponse(200, i.obj)
		if err != nil {
			return httpmock.NewStringResponse(500, ""), nil
		}
		return resp, nil
	})

	// Create Service response mock
	httpmock.RegisterResponder("POST", fmt.Sprintf("=~^%s/iot/(services|devices)", host), func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("fiware-service") == "" || req.Header.Get("fiware-servicepath") == "" {
			return httpmock.NewJsonResponse(400, IoTAgentError{Name: "MISSING_HEADERS", Message: "Some headers were missing from the request: [\"fiware-service\",\"fiware-servicepath\"]"})
		}
		if req.Header.Get("fiware-service") == "fail" {
			return httpmock.NewJsonResponse(400, IoTAgentError{
				Name:    "WRONG_SYNTAX",
				Message: "Failed because you want to.",
			})
		}
		obj := reflect.TypeOf(i.obj)
		if err := json.NewDecoder(req.Body).Decode(obj); err != nil {
			return httpmock.NewJsonResponse(400, IoTAgentError{Name: "WRONG_SYNTAX", Message: "Wrong syntax in request: Errors found validating request."})
		}
		return httpmock.NewStringResponse(200, ""), nil
	})

	// Update Service response mock
	httpmock.RegisterResponder("PUT", fmt.Sprintf("=~^%s/iot/(services|devices)", host), func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("fiware-service") == "" || req.Header.Get("fiware-servicepath") == "" {
			return httpmock.NewJsonResponse(400, IoTAgentError{
				Name:    "MISSING_HEADERS",
				Message: "Some headers were missing from the request: [\"fiware-service\",\"fiware-servicepath\"]"})
		}
		if req.Header.Get("fiware-service") == "fail" {
			return httpmock.NewJsonResponse(400, IoTAgentError{
				Name:    "WRONG_SYNTAX",
				Message: "Failed because you want to.",
			})
		}
		obj := reflect.TypeOf(i.obj)
		if err := json.NewDecoder(req.Body).Decode(obj); err != nil {
			return httpmock.NewJsonResponse(400, IoTAgentError{Name: "WRONG_SYNTAX", Message: "Wrong syntax in request: Errors found validating request."})
		}
		return httpmock.NewStringResponse(200, ""), nil
	})

	// Delete Service response mock
	httpmock.RegisterResponder("DELETE", fmt.Sprintf("=~^%s/iot/(services|devices)", host), func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("fiware-service") == "" || req.Header.Get("fiware-servicepath") == "" {
			return httpmock.NewJsonResponse(400, IoTAgentError{Name: "MISSING_HEADERS", Message: "Some headers were missing from the request: [\"fiware-service\",\"fiware-servicepath\"]"})
		}
		if req.Header.Get("fiware-service") == "fail" {
			return httpmock.NewJsonResponse(400, IoTAgentError{
				Name:    "WRONG_SYNTAX",
				Message: "Failed because you want to.",
			})
		}
		if req.URL.Query().Get("apikey") == "" || req.URL.Query().Get("resource") == "" {
			return httpmock.NewJsonResponse(400, IoTAgentError{
				Name:    "MISSING_HEADERS",
				Message: "Some headers were missing from the request: [\"apikey\"]",
			})
		}
		if req.URL.Query().Get("apikey") == "" || req.URL.Query().Get("resource") == "" {
			return httpmock.NewJsonResponse(400, IoTAgentError{Name: "MISSING_HEADERS", Message: "Some headers were missing from the request: [\"resource\",\"apikey\"]"})
		}
		return httpmock.NewStringResponse(200, ""), nil
	})
	return i
}

func TestClient_GetService(t *testing.T) {
	client := resty.New()
	testServer := createTestServer(client.GetClient(), "http://iot-agent.de")

	obj := fiware.IoTAgentGetServicesResponse{
		Services: []fiware.ServiceGroup{getServiceGroup()}}
	testServer.saveObject(obj)
	httpClient := fiware.NewClient(fiware.WithHTTPClient(client))
	type fields struct {
		httpClient   *fiware.HTTPClient
		Host         string
		FiwareConfig *fiware.Config
		Resource     string
		Protocol     string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *fiware.IoTAgentGetServicesResponse
		wantErr bool
	}{
		{
			name: "Get Service Group",
			fields: fields{
				httpClient:   httpClient,
				Host:         "http://iot-agent.de",
				FiwareConfig: fiware.NewConfig(fiware.WithService("berlin"), fiware.WithServicePath("/")),
				Resource:     "/iot/d",
				Protocol:     "ul",
			},
			want:    &obj,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				httpClient:   tt.fields.httpClient,
				Host:         tt.fields.Host,
				FiwareConfig: tt.fields.FiwareConfig,
				Resource:     tt.fields.Resource,
				Protocol:     tt.fields.Protocol,
			}
			got, err := c.GetService()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetServiceGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetServiceGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}
func getDevices() fiware.Devices {
	return fiware.Devices{
		DeviceID:   "DeviceId",
		EntityName: "urn:test",
		EntityType: "Device001",
		Attributes: []fiware.Attributes{{
			ObjectID: "t",
			Name:     "temperature",
			Type:     "Number",
		}},
		Lazy: []fiware.Lazy{{
			ObjectID: "t",
			Name:     "temperature",
			Type:     "",
			Metadata: fiware.Metadata{
				Unitcode: fiware.Unitcode{
					Type:  "",
					Value: "",
				},
			},
		}},
		Commands: []fiware.Commands{{
			ObjectID: "t",
			Name:     "temperature",
			Type:     "Number",
		}},
		StaticAttributes: []fiware.StaticAttributes{{
			Name:  "SerialID",
			Value: "10238kl81",
			Type:  "Text",
		}},
	}
}

func getServiceGroup() fiware.ServiceGroup {
	return fiware.ServiceGroup{
		Resource: "/iot/d",
		Apikey:   "apiKey",
		Type:     "ul",
		Trust:    "asdfj1123",
		Cbhost:   "http://orion:1026",
		Protocol: "ul",
		Commands: []fiware.Command{
			{
				Name: "Test",
				Type: "int",
			},
		},
		Attributes: []fiware.Attribute{
			{
				Name: "Temperature",
				Type: "Number",
				Metadata: fiware.Metadata{
					Unitcode: fiware.Unitcode{
						Type:  "",
						Value: "",
					},
				},
			},
		},
		Lazy: []fiware.Lazy{
			{
				Name: "lazyAttribute",
				Type: "Text",
			},
		},
	}
}
func TestClient_CreateService(t *testing.T) {
	client := resty.New()
	testServer := createTestServer(client.GetClient(), "http://iot-agent.de")

	obj := fiware.IoTAgentCreateServiceGroupRequest{
		Services: []fiware.ServiceGroup{getServiceGroup()}}
	testServer.saveObject(obj)
	httpClient := fiware.NewClient(fiware.WithHTTPClient(client))
	type fields struct {
		httpClient   *fiware.HTTPClient
		Host         string
		FiwareConfig *fiware.Config
		Resource     string
		Protocol     string
	}
	type args struct {
		services *fiware.IoTAgentCreateServiceGroupRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create Service Group",
			fields: fields{
				httpClient:   httpClient,
				Host:         "http://iot-agent.de",
				FiwareConfig: fiware.NewConfig(fiware.WithService("berlin"), fiware.WithServicePath("/")),
				Resource:     "/iot/d",
				Protocol:     "ul",
			},
			wantErr: false,
			args: args{
				services: &obj,
			},
		},
		{
			name: "Handle error response",
			fields: fields{
				httpClient:   httpClient,
				Host:         "http://iot-agent.de",
				FiwareConfig: fiware.NewConfig(fiware.WithService("fail"), fiware.WithServicePath("/")),
				Resource:     "/iot/d",
				Protocol:     "ul",
			},
			wantErr: true,
			args: args{
				services: &obj,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				httpClient:   tt.fields.httpClient,
				Host:         tt.fields.Host,
				FiwareConfig: tt.fields.FiwareConfig,
				Resource:     tt.fields.Resource,
				Protocol:     tt.fields.Protocol,
			}
			if err := c.CreateService(tt.args.services); (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_UpdateService(t *testing.T) {
	client := resty.New()
	testServer := createTestServer(client.GetClient(), "http://iot-agent.de")

	obj := fiware.IoTAgentCreateServiceGroupRequest{
		Services: []fiware.ServiceGroup{getServiceGroup()}}
	testServer.saveObject(obj)
	httpClient := fiware.NewClient(fiware.WithHTTPClient(client))
	type fields struct {
		httpClient   *fiware.HTTPClient
		Host         string
		FiwareConfig *fiware.Config
		Resource     string
		Protocol     string
	}
	type args struct {
		services *fiware.IoTAgentCreateServiceGroupRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Update Service Group",
			fields: fields{
				httpClient:   httpClient,
				Host:         "http://iot-agent.de",
				FiwareConfig: fiware.NewConfig(fiware.WithService("berlin"), fiware.WithServicePath("/")),
				Resource:     "/iot/d",
				Protocol:     "ul",
			},
			wantErr: false,
			args: args{
				services: &obj,
			},
		},
		{
			name: "Handle error response",
			fields: fields{
				httpClient:   httpClient,
				Host:         "http://iot-agent.de",
				FiwareConfig: fiware.NewConfig(fiware.WithService("fail"), fiware.WithServicePath("/")),
				Resource:     "/iot/d",
				Protocol:     "ul",
			},
			wantErr: true,
			args: args{
				services: &obj,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				httpClient:   tt.fields.httpClient,
				Host:         tt.fields.Host,
				FiwareConfig: tt.fields.FiwareConfig,
				Resource:     tt.fields.Resource,
				Protocol:     tt.fields.Protocol,
			}
			if err := c.UpdateService(tt.args.services); (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdateService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_DeleteService(t *testing.T) {
	client := resty.New()
	createTestServer(client.GetClient(), "http://iot-agent.de")

	httpClient := fiware.NewClient(fiware.WithHTTPClient(client))
	type fields struct {
		httpClient   *fiware.HTTPClient
		Host         string
		FiwareConfig *fiware.Config
		Resource     string
		Protocol     string
	}
	type args struct {
		resource string
		apikey   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Delete Service Group",
			fields: fields{
				httpClient:   httpClient,
				Host:         "http://iot-agent.de",
				FiwareConfig: fiware.NewConfig(fiware.WithService("berlin"), fiware.WithServicePath("/")),
				Resource:     "/iot/d",
				Protocol:     "ul",
			},
			wantErr: false,
			args: args{
				resource: "/iot/d",
				apikey:   "jasdf9823",
			},
		},
		{
			name: "Get error when apikey is missing",
			fields: fields{
				httpClient:   httpClient,
				Host:         "http://iot-agent.de",
				FiwareConfig: fiware.NewConfig(fiware.WithService("fail"), fiware.WithServicePath("/")),
				Resource:     "/iot/d",
				Protocol:     "ul",
			},
			wantErr: true,
			args: args{
				resource: "/iot/d",
				apikey:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				httpClient:   tt.fields.httpClient,
				Host:         tt.fields.Host,
				FiwareConfig: tt.fields.FiwareConfig,
				Resource:     tt.fields.Resource,
				Protocol:     tt.fields.Protocol,
			}
			if err := c.DeleteService(tt.args.resource, tt.args.apikey); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CreateDevice(t *testing.T) {
	client := resty.New()
	testServer := createTestServer(client.GetClient(), "http://iot-agent.de")

	obj := fiware.IoTAgentCreateDeviceRequest{
		Devices: []fiware.Devices{getDevices()}}
	testServer.saveObject(obj)
	httpClient := fiware.NewClient(fiware.WithHTTPClient(client))
	type fields struct {
		httpClient   *fiware.HTTPClient
		Host         string
		FiwareConfig *fiware.Config
		Resource     string
		Protocol     string
	}
	type args struct {
		device *fiware.IoTAgentCreateDeviceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create Device",
			fields: fields{
				httpClient:   httpClient,
				Host:         "http://iot-agent.de",
				FiwareConfig: fiware.NewConfig(fiware.WithService("berlin"), fiware.WithServicePath("/")),
				Resource:     "/iot/d",
				Protocol:     "ul",
			},
			wantErr: false,
			args: args{
				device: &obj,
			},
		},
		{
			name: "Handle error response",
			fields: fields{
				httpClient:   httpClient,
				Host:         "http://iot-agent.de",
				FiwareConfig: fiware.NewConfig(fiware.WithService("fail"), fiware.WithServicePath("/")),
				Resource:     "/iot/d",
				Protocol:     "ul",
			},
			wantErr: true,
			args: args{
				device: &obj,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				httpClient:   tt.fields.httpClient,
				Host:         tt.fields.Host,
				FiwareConfig: tt.fields.FiwareConfig,
				Resource:     tt.fields.Resource,
				Protocol:     tt.fields.Protocol,
			}
			if err := c.CreateDevice(tt.args.device); (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateDevice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//func TestClient_GetDevices(t *testing.T) {
//	client := resty.New()
//	testServer := createTestServer(client.GetClient(), "http://iot-agent.de")
//
//	obj := fiware.IoTAgentGetDevicesResponse{
//		Devices: []fiware.Devices{getDevices()}}
//	testServer.saveObject(obj)
//	httpClient := fiware.NewClient(fiware.WithHTTPClient(client))
//	type fields struct {
//		httpClient   *fiware.HTTPClient
//		Host         string
//		FiwareConfig *fiware.Config
//		Resource     string
//		Protocol     string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		want    *fiware.IoTAgentGetDevicesResponse
//		wantErr bool
//	}{
//		{
//			name: "Create Device",
//			fields: fields{
//				httpClient:   httpClient,
//				Host:         "http://iot-agent.de",
//				FiwareConfig: fiware.NewConfig(fiware.WithService("berlin"), fiware.WithServicePath("/")),
//				Resource:     "/iot/d",
//				Protocol:     "ul",
//			},
//			want:    &obj,
//			wantErr: false,
//		},
//		{
//			name: "Handle error response",
//			fields: fields{
//				httpClient:   httpClient,
//				Host:         "http://iot-agent.de",
//				FiwareConfig: fiware.NewConfig(fiware.WithService("fail"), fiware.WithServicePath("/")),
//				Resource:     "/iot/d",
//				Protocol:     "ul",
//			},
//			want:    &obj,
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Client{
//				httpClient:   tt.fields.httpClient,
//				Host:         tt.fields.Host,
//				FiwareConfig: tt.fields.FiwareConfig,
//				Resource:     tt.fields.Resource,
//				Protocol:     tt.fields.Protocol,
//			}
//			got, err := c.GetDevices()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Client.GetDevices() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Client.GetDevices() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestClient_UpdateDevice(t *testing.T) {
	client := resty.New()
	testServer := createTestServer(client.GetClient(), "http://iot-agent.de")

	obj := fiware.IoTAgentCreateDeviceRequest{
		Devices: []fiware.Devices{getDevices()}}
	testServer.saveObject(obj)

	httpClient := fiware.NewClient(fiware.WithHTTPClient(client))
	type fields struct {
		httpClient   *fiware.HTTPClient
		Host         string
		FiwareConfig *fiware.Config
		Resource     string
		Protocol     string
	}
	type args struct {
		device *fiware.IoTAgentCreateDeviceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Update Device",
			args: args{
				device: &obj,
			},
			fields: fields{
				httpClient:   httpClient,
				Host:         "http://iot-agent.de",
				FiwareConfig: fiware.NewConfig(fiware.WithService("berlin"), fiware.WithServicePath("/")),
				Resource:     "/iot/d",
				Protocol:     "ul",
			},
			wantErr: false,
		},
		{
			name: "Handle error response",
			args: args{
				device: &obj,
			},
			fields: fields{
				httpClient:   httpClient,
				Host:         "http://iot-agent.de",
				FiwareConfig: fiware.NewConfig(fiware.WithService("fail"), fiware.WithServicePath("/")),
				Resource:     "/iot/d",
				Protocol:     "ul",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				httpClient:   tt.fields.httpClient,
				Host:         tt.fields.Host,
				FiwareConfig: tt.fields.FiwareConfig,
				Resource:     tt.fields.Resource,
				Protocol:     tt.fields.Protocol,
			}
			if err := c.UpdateDevice(tt.args.device); (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdateDevice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_DeleteDevice(t *testing.T) {
	type fields struct {
		httpClient   *fiware.HTTPClient
		Host         string
		FiwareConfig *fiware.Config
		Resource     string
		Protocol     string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				httpClient:   tt.fields.httpClient,
				Host:         tt.fields.Host,
				FiwareConfig: tt.fields.FiwareConfig,
				Resource:     tt.fields.Resource,
				Protocol:     tt.fields.Protocol,
			}
			if err := c.DeleteDevice(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteDevice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetDevices(t *testing.T) {
	client := resty.New()
	testServer := createTestServer(client.GetClient(), "http://iot-agent.de")

	obj := fiware.IoTAgentGetDeviceResponse{
		Devices: []fiware.Devices{getDevices()}}
	testServer.saveObject(obj)
	httpClient := fiware.NewClient(fiware.WithHTTPClient(client))
	type fields struct {
		httpClient   *fiware.HTTPClient
		Host         string
		FiwareConfig *fiware.Config
		Resource     string
		Protocol     string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *fiware.IoTAgentGetDeviceResponse
		wantErr bool
	}{
		{
			name: "Create Device",
			fields: fields{
				httpClient:   httpClient,
				Host:         "http://iot-agent.de",
				FiwareConfig: fiware.NewConfig(fiware.WithService("berlin"), fiware.WithServicePath("/")),
				Resource:     "/iot/d",
				Protocol:     "ul",
			},
			want:    &obj,
			wantErr: false,
		},
		{
			name: "Handle error response",
			fields: fields{
				httpClient:   httpClient,
				Host:         "http://iot-agent.de",
				FiwareConfig: fiware.NewConfig(fiware.WithService("fail"), fiware.WithServicePath("/")),
				Resource:     "/iot/d",
				Protocol:     "ul",
			},
			want:    &obj,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				httpClient:   tt.fields.httpClient,
				Host:         tt.fields.Host,
				FiwareConfig: tt.fields.FiwareConfig,
				Resource:     tt.fields.Resource,
				Protocol:     tt.fields.Protocol,
			}
			got, err := c.GetDevices()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetDevices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetDevices() = %v, want %v", got, tt.want)
			}
		})
	}
}
