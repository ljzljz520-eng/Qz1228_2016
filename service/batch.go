package service

import (
	"context"
	"fmt"
	"sync"

	"gesture-nebula-console/domain"
)

type BatchItem struct {
	ID      string
	OwnerID string
	Label   string
}
type BatchResult struct {
	Records []domain.Record
	Errors  map[string]error
}

func (s *Service) RegisterBatch(ctx context.Context, items []BatchItem) BatchResult {
	result := BatchResult{Records: make([]domain.Record, 0, len(items)), Errors: make(map[string]error)}
	for _, item := range items {
		if err := ctx.Err(); err != nil {
			result.Errors[item.ID] = err
			continue
		}
		record, err := s.RegisterRecord(ctx, item.ID, item.OwnerID, item.Label)
		if err != nil {
			result.Errors[item.ID] = err
			continue
		}
		result.Records = append(result.Records, record)
	}
	return result
}

func (s *Service) ProcessBatch(ctx context.Context, requests []ProcessRequest) BatchResult {
	result := BatchResult{Records: make([]domain.Record, 0, len(requests)), Errors: make(map[string]error)}
	for _, req := range requests {
		if err := ctx.Err(); err != nil {
			result.Errors[req.RecordID] = err
			continue
		}
		record, err := s.ProcessRecord(ctx, req)
		if err != nil {
			result.Errors[req.RecordID] = err
			continue
		}
		result.Records = append(result.Records, record)
	}
	return result
}

func (s *Service) ProcessConcurrent(ctx context.Context, requests []ProcessRequest, workers int) BatchResult {
	if workers < 1 {
		workers = 1
	}
	result := BatchResult{Records: make([]domain.Record, 0, len(requests)), Errors: make(map[string]error)}
	jobs := make(chan ProcessRequest)
	var mu sync.Mutex
	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		for req := range jobs {
			record, err := s.ProcessRecord(ctx, req)
			mu.Lock()
			if err != nil {
				result.Errors[req.RecordID] = err
			} else {
				result.Records = append(result.Records, record)
			}
			mu.Unlock()
		}
	}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker()
	}
	for _, req := range requests {
		if ctx.Err() != nil {
			result.Errors[req.RecordID] = ctx.Err()
			continue
		}
		jobs <- req
	}
	close(jobs)
	wg.Wait()
	return result
}

func (s *Service) RequireBatchSuccess(result BatchResult) error {
	if len(result.Errors) > 0 {
		return fmt.Errorf("batch completed with %d errors", len(result.Errors))
	}
	return nil
}

func (s *Service) ReviewAndArchive(ctx context.Context, id string, version int) (domain.Record, error) {
	reviewed, err := s.ReviewRecord(ctx, id, version)
	if err != nil {
		return domain.Record{}, err
	}
	return s.ArchiveRecord(ctx, reviewed.ID, reviewed.Version)
}

func (s *Service) ApplyPolicy(record domain.Record, policy domain.LabelPolicy) error {
	return policy.Check(record.Label)
}
