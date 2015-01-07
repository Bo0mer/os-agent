package facade_test

import (
	"encoding/json"
	"net/http"

	"github.com/Bo0mer/executor/fakes"

	. "github.com/Bo0mer/os-agent/model"
	. "github.com/Bo0mer/os-agent/osagent"
	. "github.com/Bo0mer/os-agent/server/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OSAgentFacade", func() {

	var osAgentFacade OSAgentFacade
	var executor *fakes.FakeExecutor
	var req *FakeRequest
	var resp *FakeResponse

	BeforeEach(func() {
		req = new(FakeRequest)
		resp = new(FakeResponse)
		executor = new(fakes.FakeExecutor)
		osAgentFacade = NewOSAgentFacade(executor)
	})

	JustBeforeEach(func() {
		osAgentFacade.ExecuteCommand(req, resp)
	})

	Context("When the request is invalid", func() {
		BeforeEach(func() {
			req.BodyReturns([]byte("invalid request"))
		})

		It("should return an status bad request", func() {
			Expect(resp.SetStatusCodeCallCount()).To(Equal(1))
			Expect(resp.SetStatusCodeArgsForCall(0)).To(Equal(http.StatusBadRequest))
		})

		It("should have not executed a command", func() {
			Expect(executor.ExecuteCallCount()).To(Equal(0))
		})

	})

	Context("when the request is valid", func() {
		var commandRequest CommandRequest

		BeforeEach(func() {
			commandRequest = CommandRequest{
				Name: "ls",
				Args: []string{"-la"},
			}
			requestBody, _ := json.Marshal(commandRequest)
			req.BodyReturns(requestBody)

		})

		It("should have called the executor with the proper command", func() {
			Expect(executor.ExecuteCallCount()).To(Equal(1))

			execCmd := executor.ExecuteArgsForCall(0)
			Expect(execCmd.Name).To(Equal("ls"))
			Expect(execCmd.Args).To(Equal([]string{"-la"}))
		})

		Context("and the command execution is finished", func() {
			BeforeEach(func() {
				executor.ExecuteReturns("stdout", "stderr", 0, nil)
			})

			It("should return status code 200 OK", func() {
				Expect(resp.SetStatusCodeArgsForCall(0)).To(Equal(http.StatusOK))
			})

			It("should return the proper result", func() {
				commandResponse := &CommandResponse{}
				json.Unmarshal(resp.SetBodyArgsForCall(0), commandResponse)

				Expect(commandResponse.Stdout).To(Equal("stdout"))
				Expect(commandResponse.Stderr).To(Equal("stderr"))
				Expect(commandResponse.ExitCode).To(Equal(0))
			})
		})

	})
})
