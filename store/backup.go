package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gesture-nebula-console/domain"
)

type Backup struct {
	CreatedAt time.Time       `json:"created_at"`
	Records   []domain.Record `json:"records"`
	Users     []domain.User   `json:"users"`
	Events    []domain.Event  `json:"events"`
	Audits    []domain.Audit  `json:"audits"`
}

func (s *Store) ExportBackup() (Backup, error) {
	records, err := s.ListRecords()
	if err != nil {
		return Backup{}, err
	}
	users, err := s.ListUsers()
	if err != nil {
		return Backup{}, err
	}
	events, err := s.ListEvents("")
	if err != nil {
		return Backup{}, err
	}
	audits, err := s.ListAudits("")
	if err != nil {
		return Backup{}, err
	}
	return Backup{CreatedAt: nowUTC(), Records: records, Users: users, Events: events, Audits: audits}, nil
}

func (s *Store) WriteBackup(path string) error {
	backup, err := s.ExportBackup()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

func ReadBackup(path string) (Backup, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Backup{}, err
	}
	var backup Backup
	if err := json.Unmarshal(data, &backup); err != nil {
		return Backup{}, fmt.Errorf("decode backup: %w", err)
	}
	return backup, nil
}

func (s *Store) RestoreBackup(backup Backup) error {
	for _, user := range backup.Users {
		if err := s.SaveUser(user); err != nil {
			return err
		}
	}
	for _, record := range backup.Records {
		if err := s.SaveRecord(record); err != nil {
			return err
		}
	}
	for _, event := range backup.Events {
		if err := s.SaveEvent(event); err != nil {
			return err
		}
	}
	for _, audit := range backup.Audits {
		if err := s.SaveAudit(audit); err != nil {
			return err
		}
	}
	return nil
}

func BackupSummary(backup Backup) map[string]int {
	return map[string]int{"records": len(backup.Records), "users": len(backup.Users), "events": len(backup.Events), "audits": len(backup.Audits)}
}
