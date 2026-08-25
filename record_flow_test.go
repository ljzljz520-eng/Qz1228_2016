package nebula

import (
	"context"
	"gesture-nebula-console/domain"
	"testing"
)

func TestRecordFlow15(t *testing.T) {
	_, svc := openFixture(t)
	r, err := svc.RegisterRecord(context.Background(), "r15", "u1", "  cobalt   arc ")
	if err != nil {
		t.Fatal(err)
	}
	if r.Gesture != domain.Gesture15 || r.Label != "cobalt arc" {
		t.Fatalf("unexpected record: %+v", r)
	}
	view, err := svc.QueryRecord(context.Background(), "r15")
	if err != nil {
		t.Fatal(err)
	}
	if view.Record.Label != "cobalt arc" || len(view.Audits) != 1 {
		t.Fatalf("unexpected view: %+v", view)
	}
}
