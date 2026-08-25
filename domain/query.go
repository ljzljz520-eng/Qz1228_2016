package domain

import (
	"sort"
	"strings"
)

type Query struct {
	Text            string
	Status          Status
	OwnerID         string
	IncludeArchived bool
	Limit           int
}

func (q Query) Match(record Record) bool {
	if q.Status != "" && record.Status != q.Status {
		return false
	}
	if q.OwnerID != "" && record.OwnerID != q.OwnerID {
		return false
	}
	if !q.IncludeArchived && record.Status == StatusArchived {
		return false
	}
	if q.Text != "" {
		text := strings.ToLower(q.Text)
		if !strings.Contains(strings.ToLower(record.ID), text) && !strings.Contains(strings.ToLower(record.Label), text) {
			return false
		}
	}
	return true
}

func ApplyQuery(records []Record, q Query) []Record {
	matched := make([]Record, 0, len(records))
	for _, record := range records {
		if q.Match(record) {
			matched = append(matched, record)
		}
	}
	sort.SliceStable(matched, func(i, j int) bool {
		if matched[i].UpdatedAt.Equal(matched[j].UpdatedAt) {
			return matched[i].ID < matched[j].ID
		}
		return matched[i].UpdatedAt.After(matched[j].UpdatedAt)
	})
	if q.Limit > 0 && len(matched) > q.Limit {
		return matched[:q.Limit]
	}
	return matched
}

func Statuses() []Status {
	return []Status{StatusReceived, StatusReviewed, StatusProcessing, StatusArchived, StatusCancelled}
}

func ParseStatus(value string) (Status, bool) {
	candidate := Status(strings.ToLower(strings.TrimSpace(value)))
	for _, status := range Statuses() {
		if status == candidate {
			return status, true
		}
	}
	return "", false
}

func StatusLabel(status Status) string {
	switch status {
	case StatusReceived:
		return "Received"
	case StatusReviewed:
		return "Reviewed"
	case StatusProcessing:
		return "Processing"
	case StatusArchived:
		return "Archived"
	case StatusCancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

func CloneRecords(records []Record) []Record {
	cloned := make([]Record, len(records))
	copy(cloned, records)
	return cloned
}
