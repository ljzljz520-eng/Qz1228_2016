package nebula

import (
	"gesture-nebula-console/domain"
	"gesture-nebula-console/store"
	"path/filepath"
	"testing"
	"time"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "reopen.db")
	s, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	r := domain.NewRecord("persist", "u1", "persistent", time.Unix(100, 0))
	if err := s.SaveUser(domain.User{ID: "u1", Name: "Persist User", Role: "operator", Active: true}); err != nil {
		t.Fatal(err)
	}
	if err := s.SaveRecord(r); err != nil {
		t.Fatal(err)
	}
	if err := s.SaveEvent(domain.Event{ID: "e1", RecordID: r.ID, Kind: "created", Payload: r.Label, CreatedAt: r.CreatedAt}); err != nil {
		t.Fatal(err)
	}
	if err := s.SaveAudit(domain.Audit{ID: "a1", RecordID: r.ID, Action: "register", Outcome: "accepted", CreatedAt: r.CreatedAt}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	loaded, err := s.LoadRecord(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Label != r.Label || loaded.Gesture != domain.Gesture15 {
		t.Fatalf("loaded %+v", loaded)
	}
	if events, _ := s.ListEvents(r.ID); len(events) != 1 {
		t.Fatalf("events %d", len(events))
	}
	if audits, _ := s.ListAudits(r.ID); len(audits) != 1 {
		t.Fatalf("audits %d", len(audits))
	}
}
