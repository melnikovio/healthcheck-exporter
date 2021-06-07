package model

type Job struct {
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
	// required: true
	WatchDog WatchDog `json:"watchdog,omitempty"`
}
