package store

import (
	"fmt"
	"gesture-nebula-console/domain"
	"go.etcd.io/bbolt"
)

func (s *Store) SaveUser(u domain.User) error {
	if err := domain.ValidateUser(u); err != nil {
		return err
	}
	data, err := encode(u)
	if err != nil {
		return err
	}
	return s.withUpdate(func(db *bbolt.DB) error {
		return db.Update(func(tx *bbolt.Tx) error {
			b, e := bucket(tx, []byte("users"))
			if e != nil {
				return e
			}
			return b.Put([]byte(u.ID), data)
		})
	})
}

func (s *Store) LoadUser(id string) (domain.User, error) {
	var u domain.User
	err := s.withView(func(db *bbolt.DB) error {
		return db.View(func(tx *bbolt.Tx) error {
			b, e := bucket(tx, []byte("users"))
			if e != nil {
				return e
			}
			return decode(b.Get([]byte(id)), &u)
		})
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("load user %s: %w", id, err)
	}
	return u, nil
}

func (s *Store) ListUsers() ([]domain.User, error) {
	users := make([]domain.User, 0)
	err := s.withView(func(db *bbolt.DB) error {
		return db.View(func(tx *bbolt.Tx) error {
			b, e := bucket(tx, []byte("users"))
			if e != nil {
				return e
			}
			return b.ForEach(func(_, value []byte) error {
				var u domain.User
				if err := decode(value, &u); err != nil {
					return err
				}
				users = append(users, u)
				return nil
			})
		})
	})
	return users, err
}
