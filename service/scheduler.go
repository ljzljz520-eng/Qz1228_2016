package service

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"gesture-nebula-console/domain"
)

type Job struct {
	ID              string
	RecordID        string
	ExpectedVersion int
	Label           string
	DueAt           time.Time
}
type JobResult struct {
	JobID       string
	Record      domain.Record
	Error       error
	CompletedAt time.Time
}

type Scheduler struct {
	service *Service
	mu      sync.Mutex
	jobs    map[string]Job
	results map[string]JobResult
}

func NewScheduler(service *Service) *Scheduler {
	return &Scheduler{service: service, jobs: make(map[string]Job), results: make(map[string]JobResult)}
}

func (s *Scheduler) Enqueue(job Job) error {
	if s.service == nil || job.ID == "" || job.RecordID == "" || job.Label == "" {
		return errors.New("invalid job")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.jobs[job.ID]; exists {
		return errors.New("job already exists")
	}
	s.jobs[job.ID] = job
	return nil
}

func (s *Scheduler) Pending() []Job {
	s.mu.Lock()
	defer s.mu.Unlock()
	jobs := make([]Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}
	sort.Slice(jobs, func(i, j int) bool {
		if jobs[i].DueAt.Equal(jobs[j].DueAt) {
			return jobs[i].ID < jobs[j].ID
		}
		return jobs[i].DueAt.Before(jobs[j].DueAt)
	})
	return jobs
}

func (s *Scheduler) RunDue(ctx context.Context, now time.Time) []JobResult {
	jobs := s.Pending()
	results := make([]JobResult, 0)
	for _, job := range jobs {
		if job.DueAt.After(now) {
			continue
		}
		record, err := s.service.ProcessRecord(ctx, ProcessRequest{RecordID: job.RecordID, ExpectedVersion: job.ExpectedVersion, NewLabel: job.Label})
		result := JobResult{JobID: job.ID, Record: record, Error: err, CompletedAt: now.UTC()}
		s.mu.Lock()
		s.results[job.ID] = result
		delete(s.jobs, job.ID)
		s.mu.Unlock()
		results = append(results, result)
	}
	return results
}

func (s *Scheduler) Result(id string) (JobResult, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result, ok := s.results[id]
	return result, ok
}

func (s *Scheduler) Cancel(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[id]; !ok {
		return false
	}
	delete(s.jobs, id)
	return true
}

func (s *Scheduler) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs = make(map[string]Job)
	s.results = make(map[string]JobResult)
}

func (s *Scheduler) Health() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return map[string]int{"pending": len(s.jobs), "completed": len(s.results)}
}

func (s *Scheduler) BuildJob(record domain.Record, label string, due time.Time) Job {
	return Job{ID: s.service.nextID("job"), RecordID: record.ID, ExpectedVersion: record.Version, Label: label, DueAt: due.UTC()}
}
