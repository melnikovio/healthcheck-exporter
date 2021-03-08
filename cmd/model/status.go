package model

//swagger:model
type Status struct {
	// required: true
	Task []Task `json:"tasks,omitempty"`
}

//swagger:model
type Task struct {
	// required: true
	Id string `json:"id,omitempty"`
	// required: true
	Status string `json:"status,omitempty"`
	// required: true
	SuccessChecks int `json:"success_checks,omitempty"`
	// required: true
	FailureChecks int `json:"failure_checks,omitempty"`
}
