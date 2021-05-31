package fiware

//IoTManagerGetResponse representation of the response from the IoT-Manager
type IoTManagerGetResponse struct {
	Count     int `json:"count"`
	Protocols []struct {
		ID          string `json:"_id"`
		Iotagent    string `json:"iotagent"`
		Resource    string `json:"resource"`
		Protocol    string `json:"protocol"`
		Description string `json:"description"`
		V           int    `json:"__v"`
	} `json:"protocols"`
}

type IoTAgentAboutResponse struct {
	Libversion string `json:"libVersion"`
	Port       string `json:"port"`
	Baseroot   string `json:"baseRoot"`
	Version    string `json:"version"`
}

//type IoTAgentGetServicesResponse struct {
//	Services []struct {
//		Resource string `json:"resource"`
//		Apikey   string `json:"apikey"`
//		Type     string `json:"type"`
//		Trust    string `json:"trust"`
//		Cbhost   string `json:"cbHost"`
//		Protocol string `json:"protocol"`
//		Commands []struct {
//			Name string `json:"name"`
//			Type string `json:"type"`
//		} `json:"commands"`
//		Attributes []struct {
//			Name     string `json:"name"`
//			Type     string `json:"type"`
//			Metadata struct {
//				Unitcode struct {
//					Type  string `json:"type"`
//					Value string `json:"value"`
//				} `json:"unitCode"`
//			} `json:"metadata"`
//		} `json:"attributes"`
//		Lazy []struct {
//			Name string `json:"name"`
//			Type string `json:"type"`
//		} `json:"lazy"`
//	} `json:"services"`
//}

type IoTAgentGetServicesResponse struct {
	Services []ServiceGroup `json:"services"`
}
type IoTAgentCreateServiceGroupRequest struct {
	Services []ServiceGroup `json:"services"`
}
type Command struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
type Unitcode struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type Metadata struct {
	Unitcode Unitcode `json:"unitCode"`
}
type Attribute struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Metadata Metadata `json:"metadata"`
}
type ServiceGroup struct {
	Resource   string      `json:"resource"`
	Apikey     string      `json:"apikey"`
	Type       string      `json:"type"`
	Trust      string      `json:"trust"`
	Cbhost     string      `json:"cbHost"`
	Protocol   string      `json:"protocol"`
	Commands   []Command   `json:"commands"`
	Attributes []Attribute `json:"attributes"`
	Lazy       []Lazy      `json:"lazy"`
}
type IoTAgentCreateDeviceRequest struct {
	Devices []Devices `json:"devices"`
}
type IoTAgentGetDevicesResponse struct {
	Devices []Devices `json:"devices"`
}
type Attributes struct {
	ObjectID string `json:"object_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}
type Lazy struct {
	ObjectID string   `json:"object_id,omitempty"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Metadata Metadata `json:"metadata,omitempty"`
}
type Commands struct {
	ObjectID string `json:"object_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}
type StaticAttributes struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}
type Devices struct {
	DeviceID         string             `json:"device_id"`
	EntityName       string             `json:"entity_name"`
	EntityType       string             `json:"entity_type"`
	Attributes       []Attributes       `json:"attributes,omitempty"`
	Lazy             []Lazy             `json:"lazy,omitempty"`
	Commands         []Commands         `json:"commands,omitempty"`
	StaticAttributes []StaticAttributes `json:"static_attributes,omitempty"`
}
