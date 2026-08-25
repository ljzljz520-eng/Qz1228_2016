package service

import (
	"context"
	"fmt"
	"time"

	"gesture-nebula-console/domain"
)

type ProcessRequest struct {
	RecordID        string
	ExpectedVersion int
	NewLabel        string
	Delay           time.Duration
}

func (s *Service) ProcessRecord(ctx context.Context, req ProcessRequest) (domain.Record, error) {
	if err := s.requireStore(); err != nil {
		return domain.Record{}, err
	}
	if err := s.contextErr(ctx); err != nil {
		return domain.Record{}, err
	}
	if err := domain.ValidateLabel(req.NewLabel); err != nil {
		return domain.Record{}, err
	}
	if req.Delay > 0 {
		time.Sleep(req.Delay)
	}
	return s.processWithContext(context.Background(), req)
}

func (s *Service) processWithContext(ctx context.Context, req ProcessRequest) (domain.Record, error) {
	if err := s.contextErr(ctx); err != nil {
		return domain.Record{}, ErrCancelled
	}
	clean := domain.NormalizeLabel(req.NewLabel)
	r, err := s.Store.UpdateRecord(req.RecordID, req.ExpectedVersion, func(r *domain.Record) error {
		if err := s.contextErr(ctx); err != nil {
			return ErrCancelled
		}
		if domain.IsTerminal(r.Status) {
			return fmt.Errorf("record %s is terminal", r.ID)
		}
		r.Label = clean
		return r.ApplyTransition(domain.StatusProcessing, s.timestamp())
	})
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.Store.SaveEvent(domain.Event{ID: s.nextID("event"), RecordID: r.ID, Kind: "label-updated", Payload: clean, CreatedAt: s.timestamp()}); err != nil {
		return domain.Record{}, err
	}
	if err := s.recordAudit(r.ID, "process", "updated"); err != nil {
		return domain.Record{}, err
	}
	return r, nil
}

func (s *Service) CancelRecord(ctx context.Context, id string, expected int) (domain.Record, error) {
	if err := s.requireStore(); err != nil {
		return domain.Record{}, err
	}
	if err := s.contextErr(ctx); err != nil {
		return domain.Record{}, err
	}
	r, err := s.Store.UpdateRecord(id, expected, func(r *domain.Record) error { return r.ApplyTransition(domain.StatusCancelled, s.timestamp()) })
	if err != nil {
		return domain.Record{}, err
	}
	return r, s.recordAudit(id, "cancel", "cancelled")
}
