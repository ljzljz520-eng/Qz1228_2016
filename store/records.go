package store

import (
	"fmt"
	"gesture-nebula-console/domain"
	"go.etcd.io/bbolt"
)

func (s *Store) SaveRecord(r domain.Record) error {
	if err := r.Validate(); err != nil {
		return err
	}
	data, err := encode(r)
	if err != nil {
		return err
	}
	return s.withUpdate(func(db *bbolt.DB) error {
		return db.Update(func(tx *bbolt.Tx) error {
			b, e := bucket(tx, []byte("records"))
			if e != nil {
				return e
			}
			return b.Put([]byte(r.ID), data)
		})
	})
}

func (s *Store) LoadRecord(id string) (domain.Record, error) {
	var r domain.Record
	err := s.withView(func(db *bbolt.DB) error {
		return db.View(func(tx *bbolt.Tx) error {
			b, e := bucket(tx, []byte("records"))
			if e != nil {
				return e
			}
			return decode(b.Get([]byte(id)), &r)
		})
	})
	if err != nil {
		return domain.Record{}, fmt.Errorf("load record %s: %w", id, err)
	}
	return r, nil
}

func (s *Store) DeleteRecord(id string) error {
	return s.withUpdate(func(db *bbolt.DB) error {
		return db.Update(func(tx *bbolt.Tx) error {
			b, e := bucket(tx, []byte("records"))
			if e != nil {
				return e
			}
			return b.Delete([]byte(id))
		})
	})
}

func (s *Store) ListRecords() ([]domain.Record, error) {
	items := make([]domain.Record, 0)
	err := s.withView(func(db *bbolt.DB) error {
		return db.View(func(tx *bbolt.Tx) error {
			b, e := bucket(tx, []byte("records"))
			if e != nil {
				return e
			}
			return b.ForEach(func(_, value []byte) error {
				var r domain.Record
				if err := decode(value, &r); err != nil {
					return err
				}
				items = append(items, r)
				return nil
			})
		})
	})
	return items, err
}

func (s *Store) UpdateRecord(id string, expected int, mutate func(*domain.Record) error) (domain.Record, error) {
	var out domain.Record
	err := s.withUpdate(func(db *bbolt.DB) error {
		return db.Update(func(tx *bbolt.Tx) error {
			b, e := bucket(tx, []byte("records"))
			if e != nil {
				return e
			}
			var r domain.Record
			if e = decode(b.Get([]byte(id)), &r); e != nil {
				return e
			}
			if e = domain.ValidateVersion(expected, r.Version); e != nil {
				return e
			}
			if e = mutate(&r); e != nil {
				return e
			}
			if e = r.Validate(); e != nil {
				return e
			}
			data, e := encode(r)
			if e != nil {
				return e
			}
			if e = b.Put([]byte(id), data); e == nil {
				out = r
			}
			return e
		})
	})
	return out, err
}
