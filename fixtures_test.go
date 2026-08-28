package nebula

import (
	"path/filepath"
	"testing"
	"time"

	"gesture-nebula-console/domain"
	"gesture-nebula-console/service"
	"gesture-nebula-console/store"
)

func openFixture(t *testing.T) (*store.Store, *service.Service) {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "nebula.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	if err := s.SaveUser(domain.User{ID: "u1", Name: "Operator", Role: "operator", Active: true}); err != nil {
		t.Fatal(err)
	}
	return s, service.New(s)
}

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }
