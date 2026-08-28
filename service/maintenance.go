package service

import (
	"context"
	"errors"
	"sort"
	"time"

	"gesture-nebula-console/domain"
	"gesture-nebula-console/store"
)

type MaintenanceReport struct {
	Checked       int
	Healthy       int
	MissingOwners []string
	InvalidLabels []string
	GeneratedAt   time.Time
}

func (s *Service) Inspect(ctx context.Context) (MaintenanceReport, error) {
	if err := s.requireStore(); err != nil {
		return MaintenanceReport{}, err
	}
	records, err := s.Store.ListRecords()
	if err != nil {
		return MaintenanceReport{}, err
	}
	report := MaintenanceReport{Checked: len(records), MissingOwners: make([]string, 0), InvalidLabels: make([]string, 0), GeneratedAt: s.timestamp()}
	for _, record := range records {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		healthy := record.Validate() == nil
		if _, err := s.Store.LoadUser(record.OwnerID); err != nil {
			report.MissingOwners = append(report.MissingOwners, record.ID)
			healthy = false
		}
		if err := domain.ValidateLabel(record.Label); err != nil {
			report.InvalidLabels = append(report.InvalidLabels, record.ID)
			healthy = false
		}
		if healthy {
			report.Healthy++
		}
	}
	sort.Strings(report.MissingOwners)
	sort.Strings(report.InvalidLabels)
	return report, nil
}

func (s *Service) PurgeCancelled(ctx context.Context) (int, error) {
	if err := s.requireStore(); err != nil {
		return 0, err
	}
	records, err := s.Store.FindRecords(store.RecordFilter{Status: domain.StatusCancelled})
	if err != nil {
		return 0, err
	}
	removed := 0
	for _, record := range records {
		if err := ctx.Err(); err != nil {
			return removed, err
		}
		if err := s.Store.DeleteRecord(record.ID); err != nil {
			return removed, err
		}
		removed++
	}
	return removed, nil
}

func (s *Service) VerifyRecord(ctx context.Context, id string) error {
	view, err := s.QueryRecord(ctx, id)
	if err != nil {
		return err
	}
	if err := view.Record.Validate(); err != nil {
		return err
	}
	if len(view.Audits) == 0 {
		return errors.New("record has no audit history")
	}
	return nil
}

func (s *Service) RebuildAudit(ctx context.Context, id string) error {
	record, err := s.Store.LoadRecord(id)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return s.recordAudit(record.ID, "reconcile", "verified")
}
