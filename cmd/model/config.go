package model

//swagger:model
type Config struct {
	// required: true
	Authentication Authentication `json:"authentication,omitempty"`
	// required: true
	Functions []Function `json:"functions,omitempty"`
}

type Authentication struct {
	// required: true
	AuthUrl string `json:"auth_url,omitempty"`
	// required: true
	Realm string `json:"realm,omitempty"`
	// required: true
	ClientId string `json:"client_id,omitempty"`
	// required: true
	ClientSecret string `json:"client_secret,omitempty"`
}

type Function struct {
	// required: true
	Id string `json:"id,omitempty"`
	// required: true
	Description string `json:"desc,omitempty"`
	// required: true
	Type string `json:"type,omitempty"`
	// required: true
	Urls []string `json:"urls,omitempty"`
	// required: true
	Body string `json:"body,omitempty"`
	// required: true
	AuthEnabled bool `json:"auth_enabled,omitempty"`
	// required: true
	Timeout int64 `json:"timeout,omitempty"`
}
