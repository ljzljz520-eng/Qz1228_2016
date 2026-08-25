package store

import (
	"fmt"
	"gesture-nebula-console/domain"
	"go.etcd.io/bbolt"
)

type Transaction struct {
	Record domain.Record
	Events []domain.Event
	Audits []domain.Audit
}

func (s *Store) SaveTransaction(transaction Transaction) error {
	if err := transaction.Record.Validate(); err != nil {
		return err
	}
	return s.withUpdate(func(db *bbolt.DB) error {
		return db.Update(func(tx *bbolt.Tx) error {
			records, err := bucket(tx, []byte("records"))
			if err != nil {
				return err
			}
			recordData, err := encode(transaction.Record)
			if err != nil {
				return err
			}
			if err := records.Put([]byte(transaction.Record.ID), recordData); err != nil {
				return err
			}
			events, err := bucket(tx, []byte("events"))
			if err != nil {
				return err
			}
			for _, event := range transaction.Events {
				if event.RecordID != transaction.Record.ID {
					return fmt.Errorf("event %s belongs to another record", event.ID)
				}
				data, err := encode(event)
				if err != nil {
					return err
				}
				if err := events.Put([]byte(event.ID), data); err != nil {
					return err
				}
			}
			audits, err := bucket(tx, []byte("audits"))
			if err != nil {
				return err
			}
			for _, audit := range transaction.Audits {
				if audit.RecordID != transaction.Record.ID {
					return fmt.Errorf("audit %s belongs to another record", audit.ID)
				}
				data, err := encode(audit)
				if err != nil {
					return err
				}
				if err := audits.Put([]byte(audit.ID), data); err != nil {
					return err
				}
			}
			return nil
		})
	})
}

func (s *Store) LoadTransaction(id string) (Transaction, error) {
	record, err := s.LoadRecord(id)
	if err != nil {
		return Transaction{}, err
	}
	events, err := s.ListEvents(id)
	if err != nil {
		return Transaction{}, err
	}
	audits, err := s.ListAudits(id)
	if err != nil {
		return Transaction{}, err
	}
	return Transaction{Record: record, Events: events, Audits: audits}, nil
}

func (s *Store) ReplaceTransaction(transaction Transaction, expected int) error {
	current, err := s.LoadRecord(transaction.Record.ID)
	if err != nil {
		return err
	}
	if err := domain.ValidateVersion(expected, current.Version); err != nil {
		return err
	}
	return s.SaveTransaction(transaction)
}
