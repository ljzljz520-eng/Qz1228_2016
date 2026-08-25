package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const Gesture15 = 15

var (
	ErrInvalidGesture    = errors.New("gesture must be 15")
	ErrEmptyLabel        = errors.New("label cannot be empty")
	ErrInvalidStatus     = errors.New("invalid record status")
	ErrInvalidTransition = errors.New("invalid record transition")
	ErrStaleVersion      = errors.New("record version is stale")
)

type Status string

const (
	StatusReceived   Status = "received"
	StatusReviewed   Status = "reviewed"
	StatusProcessing Status = "processing"
	StatusArchived   Status = "archived"
	StatusCancelled  Status = "cancelled"
)

type Record struct {
	ID        string    `json:"id"`
	Gesture   int       `json:"gesture"`
	Label     string    `json:"label"`
	Status    Status    `json:"status"`
	Version   int       `json:"version"`
	OwnerID   string    `json:"owner_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Active bool   `json:"active"`
}

type Event struct {
	ID        string    `json:"id"`
	RecordID  string    `json:"record_id"`
	Kind      string    `json:"kind"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type Audit struct {
	ID        string    `json:"id"`
	RecordID  string    `json:"record_id"`
	Action    string    `json:"action"`
	Outcome   string    `json:"outcome"`
	CreatedAt time.Time `json:"created_at"`
}

func NewRecord(id, owner, label string, now time.Time) Record {
	return Record{ID: strings.TrimSpace(id), Gesture: Gesture15, Label: strings.TrimSpace(label), Status: StatusReceived, Version: 1, OwnerID: strings.TrimSpace(owner), CreatedAt: now.UTC(), UpdatedAt: now.UTC()}
}

func (r Record) Validate() error {
	if r.Gesture != Gesture15 {
		return ErrInvalidGesture
	}
	if strings.TrimSpace(r.Label) == "" {
		return ErrEmptyLabel
	}
	if r.ID == "" || r.OwnerID == "" {
		return errors.New("record id and owner are required")
	}
	if r.Version < 1 {
		return ErrStaleVersion
	}
	if !validStatus(r.Status) {
		return fmt.Errorf("%w: %s", ErrInvalidStatus, r.Status)
	}
	return nil
}

func validStatus(s Status) bool {
	switch s {
	case StatusReceived, StatusReviewed, StatusProcessing, StatusArchived, StatusCancelled:
		return true
	default:
		return false
	}
}

func (r *Record) Touch(now time.Time) { r.UpdatedAt = now.UTC(); r.Version++ }

func (r Record) Clone() Record { return r }
