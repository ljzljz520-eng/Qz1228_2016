package service

import (
	"context"
	"sort"
	"time"

	"gesture-nebula-console/domain"
)

type Metrics struct {
	Total      int
	Received   int
	Reviewed   int
	Processing int
	Archived   int
	Cancelled  int
	Labels     []string
	AsOf       time.Time
}

func (s *Service) Metrics(ctx context.Context) (Metrics, error) {
	records, err := s.ListRecords(ctx)
	if err != nil {
		return Metrics{}, err
	}
	m := Metrics{Total: len(records), Labels: make([]string, 0), AsOf: s.timestamp()}
	seen := make(map[string]bool)
	for _, record := range records {
		switch record.Status {
		case domain.StatusReceived:
			m.Received++
		case domain.StatusReviewed:
			m.Reviewed++
		case domain.StatusProcessing:
			m.Processing++
		case domain.StatusArchived:
			m.Archived++
		case domain.StatusCancelled:
			m.Cancelled++
		}
		if !seen[record.Label] {
			seen[record.Label] = true
			m.Labels = append(m.Labels, record.Label)
		}
	}
	sort.Strings(m.Labels)
	return m, nil
}

func (m Metrics) Active() int   { return m.Received + m.Reviewed + m.Processing }
func (m Metrics) Terminal() int { return m.Archived + m.Cancelled }
func (m Metrics) HasLabel(label string) bool {
	for _, value := range m.Labels {
		if value == label {
			return true
		}
	}
	return false
}

func (s *Service) Reconcile(ctx context.Context, records []domain.Record) ([]domain.Record, error) {
	result := make([]domain.Record, 0, len(records))
	for _, record := range records {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		stored, err := s.Store.LoadRecord(record.ID)
		if err != nil {
			return result, err
		}
		if stored.Version >= record.Version {
			result = append(result, stored)
		} else {
			result = append(result, record)
		}
	}
	return result, nil
}

func (s *Service) TouchRecord(ctx context.Context, id string, expected int) (domain.Record, error) {
	if err := s.contextErr(ctx); err != nil {
		return domain.Record{}, err
	}
	return s.Store.UpdateRecord(id, expected, func(record *domain.Record) error { record.Touch(s.timestamp()); return nil })
}
