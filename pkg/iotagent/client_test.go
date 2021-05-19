package iotagent

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/FeSaRoTi/fiware-lib/pkg/fiware"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

func TestCreateService(t *testing.T) {
	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()
	obj := fiware.ServiceGroup{
		Resource:   "/iot/d",
		Apikey:     "apiKey",
		Type:       "ul",
		Trust:      "asdfj1123",
		Cbhost:     "http://orion:1026",
		Protocol:   "ul",
		Commands:   []fiware.Command{},
		Attributes: []fiware.Attribute{},
		Lazy:       []fiware.Lazy{}}
	obj2 := fiware.IoTAgentGetServicesResponse{
		Services: []fiware.ServiceGroup{obj},
	}

	httpmock.RegisterResponder("GET", "http://iot-agent.de/iot/services", func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(200, obj2)
		if err != nil {
			return httpmock.NewStringResponse(500, ""), nil
		}
		return resp, nil
	})
	fiwareClient := fiware.NewClient(fiware.WithHTTPClient(client))
	iotClient := NewClient(Host("http://iot-agent.de"), HTTPClient(fiwareClient))
	resp, err := iotClient.GetServiceGroups()
	if err != nil {
		t.Error(err)
	}
	if len(resp.Services) != 1 {
		t.Errorf("There should be only 1 servicegroup but there are %d", len(resp.Services))
	}
	if resp.Services[0].Resource != "/iot/d" {
		t.Error("Wrong result")
	}
	if resp.Services[0].Apikey != obj.Apikey {
		t.Errorf("%v not equal to %v", resp.Services[0].Apikey, obj.Apikey)
	}
}

func TestClient_GetServiceGroups(t *testing.T) {
	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())

	obj := fiware.IoTAgentGetServicesResponse{
		Services: []fiware.ServiceGroup{
			{
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
				}},
		},
	}
	httpmock.RegisterResponder("GET", "http://iot-agent.de/iot/services", func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(200, obj)
		if err != nil {
			return httpmock.NewStringResponse(500, ""), nil
		}
		return resp, nil
	})
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
			name: "Get correct response",
			fields: fields{
				httpClient:   httpClient,
				Host:         "http://iot-agent.de",
				FiwareConfig: fiware.NewConfig(fiware.Service("berlin"), fiware.ServicePath("/")),
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
			got, err := c.GetServiceGroups()
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
