package service

import (
	"context"
	"fmt"

	"gesture-nebula-console/domain"
	"gesture-nebula-console/store"
)

type RecordView struct {
	Record domain.Record
	Events []domain.Event
	Audits []domain.Audit
}

func (s *Service) QueryRecord(ctx context.Context, id string) (RecordView, error) {
	if err := s.requireStore(); err != nil {
		return RecordView{}, err
	}
	if err := s.contextErr(ctx); err != nil {
		return RecordView{}, err
	}
	r, err := s.Store.LoadRecord(id)
	if err != nil {
		return RecordView{}, err
	}
	events, err := s.Store.ListEvents(id)
	if err != nil {
		return RecordView{}, fmt.Errorf("events: %w", err)
	}
	audits, err := s.Store.ListAudits(id)
	if err != nil {
		return RecordView{}, fmt.Errorf("audits: %w", err)
	}
	return RecordView{Record: r, Events: events, Audits: audits}, nil
}

func (s *Service) ListRecords(ctx context.Context) ([]domain.Record, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	if err := s.contextErr(ctx); err != nil {
		return nil, err
	}
	return s.Store.FindRecords(store.RecordFilter{})
}

func (s *Service) History(ctx context.Context, id string) ([]domain.Audit, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	if err := s.contextErr(ctx); err != nil {
		return nil, err
	}
	return s.Store.ListAudits(id)
}

func (s *Service) UpdateLabel(ctx context.Context, id string, expected int, label string) (domain.Record, error) {
	return s.ProcessRecord(ctx, ProcessRequest{RecordID: id, ExpectedVersion: expected, NewLabel: label})
}

func (s *Service) Search(ctx context.Context, query domain.Query) ([]domain.Record, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	if err := s.contextErr(ctx); err != nil {
		return nil, err
	}
	records, err := s.Store.ListRecords()
	if err != nil {
		return nil, err
	}
	return domain.ApplyQuery(records, query), nil
}
