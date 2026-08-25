package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"gesture-nebula-console/domain"
	"go.etcd.io/bbolt"
)

var (
	ErrNotFound = errors.New("entity not found")
	ErrClosed   = errors.New("store is closed")
)

var bucketNames = [][]byte{[]byte("records"), []byte("users"), []byte("events"), []byte("audits")}

type Store struct {
	db   *bbolt.DB
	mu   sync.RWMutex
	path string
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	db, err := bbolt.Open(filepath.Clean(path), 0600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, fmt.Errorf("open bolt: %w", err)
	}
	s := &Store{db: db, path: path}
	if err := s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, e := tx.CreateBucketIfNotExists(name); e != nil {
				return e
			}
		}
		return nil
	}); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Path() string { s.mu.RLock(); defer s.mu.RUnlock(); return s.path }

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) withView(fn func(*bbolt.DB) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return ErrClosed
	}
	return fn(s.db)
}

func (s *Store) withUpdate(fn func(*bbolt.DB) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return ErrClosed
	}
	return fn(s.db)
}

func encode(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}
	return b, nil
}
func decode(data []byte, v any) error {
	if len(data) == 0 {
		return ErrNotFound
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	return nil
}

func bucket(tx *bbolt.Tx, name []byte) (*bbolt.Bucket, error) {
	b := tx.Bucket(name)
	if b == nil {
		return nil, fmt.Errorf("missing bucket %s", name)
	}
	return b, nil
}

func nowUTC() time.Time { return time.Now().UTC().Round(time.Microsecond) }

func ensureValidEntities(r domain.Record, u domain.User) error {
	if err := r.Validate(); err != nil {
		return err
	}
	return domain.ValidateUser(u)
}
