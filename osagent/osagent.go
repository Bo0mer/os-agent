package facade

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Bo0mer/executor"
	execModel "github.com/Bo0mer/executor/model"

	. "github.com/Bo0mer/os-agent/model"
	"github.com/Bo0mer/os-agent/server"
)

type OSAgentFacade interface {
	ExecuteCommand(server.Request, server.Response)
}

type osAgentFacade struct {
	executor executor.Executor
}

func NewOSAgentFacade(e executor.Executor) OSAgentFacade {
	return &osAgentFacade{
		executor: e,
	}
}

func (f *osAgentFacade) ExecuteCommand(req server.Request, resp server.Response) {
	execCommand, err := f.createExecutorCommand(req.Body())
	if err != nil {
		resp.SetStatusCode(http.StatusBadRequest)
		return
	}

	stdout, stderr, exitCode, _ := f.executor.Execute(*execCommand)
	commandResp := &CommandResponse{
		Stdout:   stdout,
		Stderr:   stderr,
		ExitCode: exitCode,
	}

	responseBody, err := json.Marshal(commandResp)
	if err != nil {
		resp.SetStatusCode(http.StatusInternalServerError)
		return
	}
	resp.SetBody(responseBody)
	resp.SetStatusCode(http.StatusOK)
}

func (f *osAgentFacade) createExecutorCommand(requestBody []byte) (*execModel.Command, error) {
	var commandReq = &CommandRequest{}
	err := json.Unmarshal(requestBody, commandReq)
	if err != nil {
		return nil, err
	}

	execCmd := &execModel.Command{
		Name:           commandReq.Name,
		Args:           commandReq.Args,
		Env:            commandReq.Env,
		UseIsolatedEnv: commandReq.UseIsolatedEnv,
		WorkingDir:     commandReq.WorkingDir,
		Stdin:          strings.NewReader(commandReq.Input),
	}

	return execCmd, nil
}
