package model

type JobStatus string

const (
	JOB_IN_PROCESS JobStatus = "IN_PROCESS"
	JOB_COMPLETED  JobStatus = "COMPLETED"
)

type JobRequest struct {
	Async   bool           `json:"async"`
	Command CommandRequest `json:"command"`
}

type Job struct {
	Id     string          `json:"id"`
	Status JobStatus       `json:"status"`
	Result CommandResponse `json:"result"`
}
