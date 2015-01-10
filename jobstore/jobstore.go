package jobstore

import (
	"sync"

	"github.com/Bo0mer/os-agent/model"
)

type JobStore interface {
	Set(model.Job)
	Get(string) (model.Job, bool)
}

type jobStore struct {
	m    sync.RWMutex
	jobs map[string]model.Job
}

func NewJobStore() JobStore {
	return &jobStore{
		jobs: make(map[string]model.Job),
	}
}

func (s *jobStore) Set(job model.Job) {
	s.m.Lock()
	s.jobs[job.Id] = job
	s.m.Unlock()
}

func (s *jobStore) Get(id string) (job model.Job, found bool) {
	s.m.RLock()
	job, found = s.jobs[id]
	s.m.RUnlock()
	return
}
