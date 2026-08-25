package service

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"gesture-nebula-console/store"
)

var ErrCancelled = errors.New("operation cancelled")

type Clock interface{ Now() time.Time }

type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now().UTC().Round(time.Microsecond) }

type Service struct {
	Store    *store.Store
	Clock    Clock
	sequence atomic.Uint64
}

func New(s *store.Store) *Service { return &Service{Store: s, Clock: RealClock{}} }

func (s *Service) nextID(prefix string) string {
	n := s.sequence.Add(1)
	return fmt.Sprintf("%s-%06d", prefix, n)
}

func (s *Service) requireStore() error {
	if s == nil || s.Store == nil {
		return errors.New("service store is unavailable")
	}
	return nil
}

func (s *Service) contextErr(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (s *Service) timestamp() time.Time {
	if s.Clock == nil {
		s.Clock = RealClock{}
	}
	return s.Clock.Now().UTC()
}
