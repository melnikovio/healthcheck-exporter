package model

type WatchDog struct {
	// required: true
	Enabled bool `json:"enabled,omitempty"`
	// required: true
	Deployments []string `json:"deployments,omitempty"`
	// required: true
	Namespace string `json:"namespace,omitempty"`
	// required: true
	FailureThreshold int `json:"failureThreshold,omitempty"`
	// required: true
	AwaitAfterRestart int64 `json:"awaitAfterRestart,omitempty"`
}
