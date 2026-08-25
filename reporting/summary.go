package reporting

import (
	"fmt"
	"sort"
	"strings"

	"gesture-nebula-console/domain"
	"gesture-nebula-console/service"
)

type Summary struct {
	ID         string        `json:"id"`
	Label      string        `json:"label"`
	Status     domain.Status `json:"status"`
	Version    int           `json:"version"`
	EventCount int           `json:"event_count"`
	AuditCount int           `json:"audit_count"`
	LastAction string        `json:"last_action"`
}

func BuildSummary(view service.RecordView) Summary {
	last := ""
	if len(view.Audits) > 0 {
		sort.Slice(view.Audits, func(i, j int) bool { return view.Audits[i].CreatedAt.Before(view.Audits[j].CreatedAt) })
		last = view.Audits[len(view.Audits)-1].Action
	}
	return Summary{ID: view.Record.ID, Label: view.Record.Label, Status: view.Record.Status, Version: view.Record.Version, EventCount: len(view.Events), AuditCount: len(view.Audits), LastAction: last}
}

func FormatSummary(summary Summary) string {
	return fmt.Sprintf("%s gesture=%d label=%q status=%s version=%d events=%d audits=%d last=%s", summary.ID, domain.Gesture15, summary.Label, summary.Status, summary.Version, summary.EventCount, summary.AuditCount, summary.LastAction)
}

func FilterByStatus(records []domain.Record, status domain.Status) []domain.Record {
	out := make([]domain.Record, 0, len(records))
	for _, r := range records {
		if status == "" || r.Status == status {
			out = append(out, r)
		}
	}
	return out
}

func Labels(records []domain.Record) string {
	values := make([]string, 0, len(records))
	for _, r := range records {
		values = append(values, strings.TrimSpace(r.Label))
	}
	sort.Strings(values)
	return strings.Join(values, ", ")
}
