package store

import (
	"fmt"
	"gesture-nebula-console/domain"
	"go.etcd.io/bbolt"
)

func (s *Store) SaveAudit(a domain.Audit) error {
	if a.ID == "" || a.RecordID == "" || a.Action == "" {
		return fmt.Errorf("audit id, record id, and action are required")
	}
	data, err := encode(a)
	if err != nil {
		return err
	}
	return s.withUpdate(func(db *bbolt.DB) error {
		return db.Update(func(tx *bbolt.Tx) error {
			b, e := bucket(tx, []byte("audits"))
			if e != nil {
				return e
			}
			return b.Put([]byte(a.ID), data)
		})
	})
}

func (s *Store) ListAudits(recordID string) ([]domain.Audit, error) {
	audits := make([]domain.Audit, 0)
	err := s.withView(func(db *bbolt.DB) error {
		return db.View(func(tx *bbolt.Tx) error {
			b, e := bucket(tx, []byte("audits"))
			if e != nil {
				return e
			}
			return b.ForEach(func(_, value []byte) error {
				var a domain.Audit
				if err := decode(value, &a); err != nil {
					return err
				}
				if recordID == "" || a.RecordID == recordID {
					audits = append(audits, a)
				}
				return nil
			})
		})
	})
	return audits, err
}
