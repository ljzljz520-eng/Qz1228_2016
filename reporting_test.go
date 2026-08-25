package nebula

import (
	"gesture-nebula-console/domain"
	"gesture-nebula-console/reporting"
	"gesture-nebula-console/service"
	"testing"
	"time"
)

func TestReportingSummary(t *testing.T) {
	r := domain.NewRecord("report", "u", "label", time.Unix(1, 0))
	view := service.RecordView{Record: r, Events: []domain.Event{{ID: "e", RecordID: r.ID, Kind: "created"}}, Audits: []domain.Audit{{ID: "a", RecordID: r.ID, Action: "register", Outcome: "accepted", CreatedAt: r.CreatedAt}}}
	s := reporting.BuildSummary(view)
	if s.Label != "label" || s.LastAction != "register" || s.EventCount != 1 {
		t.Fatalf("summary %+v", s)
	}
	if reporting.FormatSummary(s) == "" {
		t.Fatal("empty summary")
	}
}
