package store

import (
	"gesture-nebula-console/domain"
	"sort"
	"strings"
)

type IndexEntry struct {
	ID        string
	Label     string
	Status    domain.Status
	UpdatedAt int64
}

func BuildIndex(records []domain.Record) map[string][]IndexEntry {
	index := make(map[string][]IndexEntry)
	for _, record := range records {
		key := strings.ToLower(strings.TrimSpace(record.Label))
		index[key] = append(index[key], IndexEntry{ID: record.ID, Label: record.Label, Status: record.Status, UpdatedAt: record.UpdatedAt.UnixNano()})
	}
	for key := range index {
		sort.Slice(index[key], func(i, j int) bool { return index[key][i].UpdatedAt < index[key][j].UpdatedAt })
	}
	return index
}

func LookupIndex(index map[string][]IndexEntry, label string) []IndexEntry {
	return append([]IndexEntry(nil), index[strings.ToLower(strings.TrimSpace(label))]...)
}

func MergeIndexes(left, right map[string][]IndexEntry) map[string][]IndexEntry {
	merged := make(map[string][]IndexEntry, len(left)+len(right))
	for key, values := range left {
		merged[key] = append([]IndexEntry(nil), values...)
	}
	for key, values := range right {
		merged[key] = append(merged[key], values...)
	}
	for key := range merged {
		sort.SliceStable(merged[key], func(i, j int) bool { return merged[key][i].ID < merged[key][j].ID })
	}
	return merged
}

func IndexLabels(index map[string][]IndexEntry) []string {
	labels := make([]string, 0, len(index))
	for key := range index {
		labels = append(labels, key)
	}
	sort.Strings(labels)
	return labels
}
