package model

//swagger:model
type Config struct {
	// required: true
	Authentication Authentication `json:"authentication,omitempty"`
	// required: true
	PushGateway PushGateway `json:"push_gateway,omitempty"`
	// required: true
	Jobs []Job `json:"jobs,omitempty"`
}

type Authentication struct {
	// required: true
	AuthUrl string `json:"auth_url,omitempty"`
	// required: true
	ClientId string `json:"client_id,omitempty"`
	// required: true
	ClientSecret string `json:"client_secret,omitempty"`
}

type PushGateway struct {
	// required: true
	Address string `json:"address,omitempty"`
}
