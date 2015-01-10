package jobstore_test

import (
	. "github.com/Bo0mer/os-agent/jobstore"
	. "github.com/Bo0mer/os-agent/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jobstore", func() {

	var jobStore JobStore

	BeforeEach(func() {
		jobStore = NewJobStore()
	})

	var addTestJob = func() (string, Job) {
		jobId := "id"
		job := Job{
			Id: jobId,
		}

		jobStore.Set(job)

		return jobId, job
	}

	Describe("Set", func() {

		var jobId string
		var job Job

		BeforeEach(func() {
			jobId, job = addTestJob()
		})

		It("should have added the job", func() {
			_, found := jobStore.Get(jobId)
			Expect(found).To(BeTrue())
		})

	})

	Describe("Get", func() {

		Context("when the job is missing", func() {
			It("should return not found", func() {
				_, found := jobStore.Get("missing")
				Expect(found).To(BeFalse())
			})
		})

		Context("when the job is present", func() {

			var jobId string
			var job Job

			BeforeEach(func() {
				jobId, job = addTestJob()
			})

			It("should return it", func() {
				actualJob, found := jobStore.Get(jobId)
				Expect(found).To(BeTrue())
				Expect(actualJob).To(Equal(job))
			})
		})
	})

})
