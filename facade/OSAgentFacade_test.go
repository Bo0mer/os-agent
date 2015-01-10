package facade_test

import (
	"encoding/json"
	"net/http"

	"github.com/Bo0mer/executor/fakes"
	execModel "github.com/Bo0mer/executor/model"

	. "github.com/Bo0mer/os-agent/facade"
	. "github.com/Bo0mer/os-agent/jobstore/fakes"
	. "github.com/Bo0mer/os-agent/model"
	. "github.com/Bo0mer/os-agent/server/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OSAgentFacade", func() {

	var osAgentFacade OSAgentFacade
	var executor *fakes.FakeExecutor
	var req *FakeRequest
	var resp *FakeResponse
	var fakeJobStore *FakeJobStore

	var itBehavesLikeBadRequest = func() {
		It("should return status Bad Request", func() {
			Expect(resp.SetStatusCodeCallCount()).To(Equal(1))
			Expect(resp.SetStatusCodeArgsForCall(0)).To(Equal(http.StatusBadRequest))
		})
	}

	var doInvalidRequest = func() {
		req.BodyReturns([]byte("invalid request"))
	}

	var doValidRequest = func(async bool) {
		commandRequest := CommandRequest{
			Name: "ls",
			Args: []string{"-la"},
		}

		jobRequest := JobRequest{
			Async:   async,
			Command: commandRequest,
		}

		requestBody, _ := json.Marshal(jobRequest)
		req.BodyReturns(requestBody)
	}

	BeforeEach(func() {
		req = new(FakeRequest)
		resp = new(FakeResponse)
		executor = new(fakes.FakeExecutor)
		fakeJobStore = new(FakeJobStore)
		osAgentFacade = NewOSAgentFacade(executor, fakeJobStore)
	})

	Describe("CreateJob", func() {

		JustBeforeEach(func() {
			osAgentFacade.CreateJob(req, resp)
		})

		Context("When the request is invalid", func() {
			BeforeEach(func() {
				doInvalidRequest()
			})

			itBehavesLikeBadRequest()

			It("should have not executed a command", func() {
				Expect(executor.ExecuteAsyncCallCount()).To(Equal(0))
			})

		})

		Context("when the request is valid", func() {
			var actualJob Job

			BeforeEach(func() {
				executor.ExecuteAsyncStub = func(c execModel.Command) <-chan execModel.CommandResult {
					resultChan := make(chan execModel.CommandResult, 1)
					resultChan <- execModel.CommandResult{
						Stdout:   "stdout",
						Stderr:   "stderr",
						ExitCode: 0,
						Error:    nil,
					}
					return resultChan
				}

				fakeJobStore.SetStub = func(job Job) {
					actualJob = job
				}

				fakeJobStore.GetStub = func(id string) (Job, bool) {
					return actualJob, true
				}

			})

			Context("and NOT async", func() {
				BeforeEach(func() {
					doValidRequest(false)
				})

				It("should have called the executor with the proper command", func() {
					Expect(executor.ExecuteAsyncCallCount()).To(Equal(1))

					execCmd := executor.ExecuteAsyncArgsForCall(0)
					Expect(execCmd.Name).To(Equal("ls"))
					Expect(execCmd.Args).To(Equal([]string{"-la"}))
				})

				It("should return status code 200 OK", func() {
					Expect(resp.SetStatusCodeArgsForCall(0)).To(Equal(http.StatusOK))
				})

				It("should return the proper result", func() {
					job := &Job{}
					json.Unmarshal(resp.SetBodyArgsForCall(0), job)

					Expect(*job).To(Equal(actualJob))
				})

			})

			Context("and async", func() {
				BeforeEach(func() {
					doValidRequest(true)
				})

				It("should have called the executor with the proper command", func() {
					Expect(executor.ExecuteAsyncCallCount()).To(Equal(1))

					execCmd := executor.ExecuteAsyncArgsForCall(0)
					Expect(execCmd.Name).To(Equal("ls"))
					Expect(execCmd.Args).To(Equal([]string{"-la"}))
				})

				It("should have returned the proper job", func() {
					returnedJob := &Job{}
					json.Unmarshal(resp.SetBodyArgsForCall(0), returnedJob)

					Expect(*returnedJob).To(Equal(actualJob))
				})

				It("should return status 200 OK", func() {
					Expect(resp.SetStatusCodeArgsForCall(0)).To(Equal(http.StatusOK))
				})
			})

		})

	})

	Describe("GetJob", func() {

		JustBeforeEach(func() {
			osAgentFacade.GetJob(req, resp)
		})

		Context("when the id parameter is missing", func() {
			BeforeEach(func() {
				req.ParamValuesReturns([]string{}, false)
			})

			It("should return status Bad Request", func() {
				status := resp.SetStatusCodeArgsForCall(0)
				Expect(status).To(Equal(http.StatusBadRequest))
			})

		})

		Context("when there are more than one id params", func() {
			BeforeEach(func() {
				req.ParamValuesReturns([]string{"one", "two"}, true)
			})

			It("should return status Bad Request", func() {
				status := resp.SetStatusCodeArgsForCall(0)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when the request is valid", func() {

			var jobId string

			BeforeEach(func() {
				jobId = "id-of-the-job"
				req.ParamValuesReturns([]string{jobId}, true)
			})

			Context("and the job is missing", func() {
				BeforeEach(func() {
					fakeJobStore.GetStub = func(id string) (Job, bool) {
						return Job{}, false
					}

				})

				It("should have looked for the right job", func() {
					searchedJobId := fakeJobStore.GetArgsForCall(0)
					Expect(searchedJobId).To(Equal(jobId))
				})

				It("should return status code 400 Not Found", func() {
					status := resp.SetStatusCodeArgsForCall(0)
					Expect(status).To(Equal(http.StatusNotFound))
				})

			})

			Context("and the job is present", func() {

				var job Job

				BeforeEach(func() {
					job = Job{
						Id: "id",
					}
					fakeJobStore.GetStub = func(id string) (Job, bool) {
						return job, true
					}
				})

				It("should return the job", func() {
					responseBody := resp.SetBodyArgsForCall(0)
					returnedJob := &Job{}
					json.Unmarshal(responseBody, returnedJob)

					Expect(*returnedJob).To(Equal(job))
				})

				It("should return status code 200 OK", func() {
					status := resp.SetStatusCodeArgsForCall(0)
					Expect(status).To(Equal(http.StatusOK))
				})
			})
		})

	})

})
