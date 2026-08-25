package store

import (
	"bytes"
	"sort"
	"strings"

	"gesture-nebula-console/domain"
	"go.etcd.io/bbolt"
)

type RecordFilter struct {
	Status        domain.Status
	OwnerID       string
	LabelContains string
	Limit         int
}

func (s *Store) FindRecords(filter RecordFilter) ([]domain.Record, error) {
	items := make([]domain.Record, 0)
	err := s.withView(func(db *bbolt.DB) error {
		return db.View(func(tx *bbolt.Tx) error {
			b, err := bucket(tx, []byte("records"))
			if err != nil {
				return err
			}
			return b.ForEach(func(_, value []byte) error {
				var r domain.Record
				if err := decode(value, &r); err != nil {
					return err
				}
				if !matches(r, filter) {
					return nil
				}
				items = append(items, r)
				return nil
			})
		})
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return items[i].ID < items[j].ID
		}
		return items[i].UpdatedAt.Before(items[j].UpdatedAt)
	})
	if filter.Limit > 0 && len(items) > filter.Limit {
		items = items[:filter.Limit]
	}
	return items, nil
}

func matches(r domain.Record, filter RecordFilter) bool {
	if filter.Status != "" && r.Status != filter.Status {
		return false
	}
	if filter.OwnerID != "" && r.OwnerID != filter.OwnerID {
		return false
	}
	if filter.LabelContains != "" && !bytes.Contains(bytes.ToLower([]byte(r.Label)), bytes.ToLower([]byte(filter.LabelContains))) {
		return false
	}
	return true
}

func (s *Store) CountRecords(status domain.Status) (int, error) {
	records, err := s.FindRecords(RecordFilter{Status: status})
	return len(records), err
}

func (s *Store) SearchLabels(query string) ([]string, error) {
	records, err := s.FindRecords(RecordFilter{LabelContains: strings.TrimSpace(query)})
	if err != nil {
		return nil, err
	}
	labels := make([]string, 0, len(records))
	seen := make(map[string]bool)
	for _, r := range records {
		if !seen[r.Label] {
			seen[r.Label] = true
			labels = append(labels, r.Label)
		}
	}
	sort.Strings(labels)
	return labels, nil
}

func (s *Store) Snapshot() (map[string]int, error) {
	result := map[string]int{"records": 0, "users": 0, "events": 0, "audits": 0}
	err := s.withView(func(db *bbolt.DB) error {
		return db.View(func(tx *bbolt.Tx) error {
			for key := range result {
				b, err := bucket(tx, []byte(key))
				if err != nil {
					return err
				}
				count := 0
				if err := b.ForEach(func(_, _ []byte) error { count++; return nil }); err != nil {
					return err
				}
				result[key] = count
			}
			return nil
		})
	})
	return result, err
}
