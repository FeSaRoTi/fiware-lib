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
