package nebula

import (
	"context"
	"gesture-nebula-console/domain"
	"gesture-nebula-console/service"
	"testing"
	"time"
)

func TestWorkflowOne(t *testing.T) {
	_, svc := openFixture(t)
	r, err := svc.RegisterRecord(context.Background(), "w1", "u1", "amber")
	if err != nil {
		t.Fatal(err)
	}
	view, err := svc.QueryRecord(context.Background(), r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if view.Record.Status != domain.StatusReceived {
		t.Fatalf("status %s", view.Record.Status)
	}
	if len(view.Audits) != 1 {
		t.Fatalf("audits %d", len(view.Audits))
	}
}

func TestWorkflowTwo(t *testing.T) {
	_, svc := openFixture(t)
	r, err := svc.RegisterRecord(context.Background(), "w2", "u1", "violet")
	if err != nil {
		t.Fatal(err)
	}
	r, err = svc.ReviewRecord(context.Background(), r.ID, r.Version)
	if err != nil {
		t.Fatal(err)
	}
	r, err = svc.ArchiveRecord(context.Background(), r.ID, r.Version)
	if err != nil {
		t.Fatal(err)
	}
	if r.Status != domain.StatusArchived {
		t.Fatalf("status %s", r.Status)
	}
	view, err := svc.QueryRecord(context.Background(), r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(view.Audits) != 3 {
		t.Fatalf("audits %d", len(view.Audits))
	}
}

func TestWorkflowThree(t *testing.T) {
	_, svc := openFixture(t)
	r, err := svc.RegisterRecord(context.Background(), "w3", "u1", "old-label")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(time.Millisecond, cancel)
	defer timer.Stop()
	updated, err := svc.ProcessRecord(ctx, service.ProcessRequest{RecordID: r.ID, ExpectedVersion: r.Version, NewLabel: "new-label", Delay: 10 * time.Millisecond})
	view, queryErr := svc.QueryRecord(context.Background(), r.ID)
	if queryErr != nil {
		t.Fatal(queryErr)
	}
	if view.Record.Label != "old-label" {
		t.Fatalf("expected old label after cancellation, got %q", view.Record.Label)
	}
	if err == nil {
		t.Fatalf("expected cancellation, got updated record %+v", updated)
	}
}
