package service

import (
	"context"
	"fmt"

	"gesture-nebula-console/domain"
)

func (s *Service) RegisterRecord(ctx context.Context, id, owner, label string) (domain.Record, error) {
	if err := s.requireStore(); err != nil {
		return domain.Record{}, err
	}
	if err := s.contextErr(ctx); err != nil {
		return domain.Record{}, err
	}
	u, err := s.Store.LoadUser(owner)
	if err != nil {
		return domain.Record{}, fmt.Errorf("load owner: %w", err)
	}
	if err := domain.ValidateUser(u); err != nil {
		return domain.Record{}, err
	}
	clean := domain.NormalizeLabel(label)
	if err := domain.ValidateLabel(clean); err != nil {
		return domain.Record{}, err
	}
	r := domain.NewRecord(id, owner, clean, s.timestamp())
	if err := s.Store.SaveRecord(r); err != nil {
		return domain.Record{}, err
	}
	return r, s.recordAudit(r.ID, "register", "accepted")
}

func (s *Service) ReviewRecord(ctx context.Context, id string, expected int) (domain.Record, error) {
	if err := s.requireStore(); err != nil {
		return domain.Record{}, err
	}
	if err := s.contextErr(ctx); err != nil {
		return domain.Record{}, err
	}
	r, err := s.Store.UpdateRecord(id, expected, func(r *domain.Record) error { return r.ApplyTransition(domain.StatusReviewed, s.timestamp()) })
	if err != nil {
		return domain.Record{}, err
	}
	return r, s.recordAudit(id, "review", "approved")
}

func (s *Service) ArchiveRecord(ctx context.Context, id string, expected int) (domain.Record, error) {
	if err := s.requireStore(); err != nil {
		return domain.Record{}, err
	}
	if err := s.contextErr(ctx); err != nil {
		return domain.Record{}, err
	}
	r, err := s.Store.UpdateRecord(id, expected, func(r *domain.Record) error { return r.ApplyTransition(domain.StatusArchived, s.timestamp()) })
	if err != nil {
		return domain.Record{}, err
	}
	return r, s.recordAudit(id, "archive", "completed")
}

func (s *Service) recordAudit(recordID, action, outcome string) error {
	return s.Store.SaveAudit(domain.Audit{ID: s.nextID("audit"), RecordID: recordID, Action: action, Outcome: outcome, CreatedAt: s.timestamp()})
}
