package store

import (
	"fmt"
	"gesture-nebula-console/domain"
	"go.etcd.io/bbolt"
)

func (s *Store) SaveEvent(ev domain.Event) error {
	if ev.ID == "" || ev.RecordID == "" || ev.Kind == "" {
		return fmt.Errorf("event id, record id, and kind are required")
	}
	data, err := encode(ev)
	if err != nil {
		return err
	}
	return s.withUpdate(func(db *bbolt.DB) error {
		return db.Update(func(tx *bbolt.Tx) error {
			b, e := bucket(tx, []byte("events"))
			if e != nil {
				return e
			}
			return b.Put([]byte(ev.ID), data)
		})
	})
}

func (s *Store) ListEvents(recordID string) ([]domain.Event, error) {
	events := make([]domain.Event, 0)
	err := s.withView(func(db *bbolt.DB) error {
		return db.View(func(tx *bbolt.Tx) error {
			b, e := bucket(tx, []byte("events"))
			if e != nil {
				return e
			}
			return b.ForEach(func(_, value []byte) error {
				var ev domain.Event
				if err := decode(value, &ev); err != nil {
					return err
				}
				if recordID == "" || ev.RecordID == recordID {
					events = append(events, ev)
				}
				return nil
			})
		})
	})
	return events, err
}
