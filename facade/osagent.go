package facade

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Bo0mer/executor"
	execModel "github.com/Bo0mer/executor/model"

	. "github.com/Bo0mer/os-agent/jobstore"
	. "github.com/Bo0mer/os-agent/model"
	"github.com/Bo0mer/os-agent/server"
)

type OSAgentFacade interface {
	CreateJob(server.Request, server.Response)
	GetJob(server.Request, server.Response)
}

type osAgentFacade struct {
	executor executor.Executor
	jobs     JobStore
}

func NewOSAgentFacade(e executor.Executor, jobStore JobStore) OSAgentFacade {
	return &osAgentFacade{
		executor: e,
		jobs:     jobStore,
	}
}

func (f *osAgentFacade) CreateJob(req server.Request, resp server.Response) {
	jobRequest, err := f.createJobRequest(req.Body())
	if err != nil {
		resp.SetStatusCode(http.StatusBadRequest)
		return
	}

	execCommand := f.createExecutorCommand(jobRequest.Command)

	if jobRequest.Async {
		f.executeCommandAsync(execCommand, resp)
		return
	} else {
		f.executeCommand(execCommand, resp)
	}
}

func (f *osAgentFacade) GetJob(req server.Request, resp server.Response) {
	jobIdValues, found := req.ParamValues("id")
	if !found || len(jobIdValues) > 1 {
		resp.SetStatusCode(http.StatusBadRequest)
		return
	}

	jobId := jobIdValues[0]
	job, found := f.jobs.Get(jobId)
	if !found {
		resp.SetStatusCode(http.StatusNotFound)
		return
	}

	f.setResponseBody(job, resp)
}

func (f *osAgentFacade) executeCommandAsync(cmd *execModel.Command, resp server.Response) {
	job := f.createJob()
	resultChan := f.executor.ExecuteAsync(*cmd)
	go f.waitForCommandResult(resultChan, job)

	f.setResponseBody(job, resp)
}

func (f *osAgentFacade) executeCommand(cmd *execModel.Command, resp server.Response) {
	job := f.createJob()
	resultChan := f.executor.ExecuteAsync(*cmd)
	f.waitForCommandResult(resultChan, job)

	job, _ = f.jobs.Get(job.Id)
	f.setResponseBody(job, resp)
}

func (f *osAgentFacade) waitForCommandResult(resultChan <-chan execModel.CommandResult, job Job) {
	result := <-resultChan

	job.Status = JOB_COMPLETED
	job.Result.Stdout, job.Result.Stderr = result.Stdout, result.Stderr
	job.Result.ExitCode = result.ExitCode

	var errorString string
	if result.Error != nil {
		errorString = result.Error.Error()
	}

	job.Result.Error = errorString
	f.jobs.Set(job)
}

func (f *osAgentFacade) createJobRequest(data []byte) (*JobRequest, error) {
	jobRequest := &JobRequest{}
	err := json.Unmarshal(data, jobRequest)
	return jobRequest, err
}

func (f *osAgentFacade) createJob() Job {
	job := Job{
		Id:     f.generateUniqueId(),
		Status: JOB_IN_PROCESS,
	}

	f.jobs.Set(job)

	return job
}

func (f *osAgentFacade) createExecutorCommand(c CommandRequest) *execModel.Command {
	return &execModel.Command{
		Name:           c.Name,
		Args:           c.Args,
		Env:            c.Env,
		UseIsolatedEnv: c.UseIsolatedEnv,
		WorkingDir:     c.WorkingDir,
		Stdin:          strings.NewReader(c.Input),
	}
}

func (f *osAgentFacade) setResponseBody(job Job, resp server.Response) {
	responseBody, _ := json.Marshal(job)
	resp.SetBody(responseBody)
	resp.SetStatusCode(http.StatusOK)
}

func (f *osAgentFacade) generateUniqueId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
