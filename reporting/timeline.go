package reporting

import (
	"sort"
	"strings"

	"gesture-nebula-console/domain"
)

type TimelineEntry struct {
	At     string `json:"at"`
	Kind   string `json:"kind"`
	Detail string `json:"detail"`
}

func BuildTimeline(events []domain.Event, audits []domain.Audit) []TimelineEntry {
	entries := make([]TimelineEntry, 0, len(events)+len(audits))
	for _, ev := range events {
		entries = append(entries, TimelineEntry{At: ev.CreatedAt.Format("2006-01-02T15:04:05.000000Z"), Kind: ev.Kind, Detail: ev.Payload})
	}
	for _, a := range audits {
		entries = append(entries, TimelineEntry{At: a.CreatedAt.Format("2006-01-02T15:04:05.000000Z"), Kind: "audit:" + a.Action, Detail: a.Outcome})
	}
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].At < entries[j].At })
	return entries
}

func RenderTimeline(entries []TimelineEntry) string {
	lines := make([]string, 0, len(entries))
	for _, e := range entries {
		lines = append(lines, e.At+" "+e.Kind+" "+e.Detail)
	}
	return strings.Join(lines, "\n")
}

func CountKinds(entries []TimelineEntry) map[string]int {
	counts := make(map[string]int)
	for _, e := range entries {
		counts[e.Kind]++
	}
	return counts
}
