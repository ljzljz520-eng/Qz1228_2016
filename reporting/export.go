package reporting

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"gesture-nebula-console/domain"
	"gesture-nebula-console/service"
)

type Export struct {
	GeneratedAt string                `json:"generated_at"`
	Total       int                   `json:"total"`
	ByStatus    map[domain.Status]int `json:"by_status"`
	Records     []Summary             `json:"records"`
}

func BuildExport(records []domain.Record, generatedAt string) Export {
	by := make(map[domain.Status]int)
	summaries := make([]Summary, 0, len(records))
	for _, r := range records {
		by[r.Status]++
		summaries = append(summaries, Summary{ID: r.ID, Label: r.Label, Status: r.Status, Version: r.Version})
	}
	sort.Slice(summaries, func(i, j int) bool { return summaries[i].ID < summaries[j].ID })
	return Export{GeneratedAt: generatedAt, Total: len(records), ByStatus: by, Records: summaries}
}

func EncodeExport(export Export) ([]byte, error) { return json.MarshalIndent(export, "", "  ") }

func RenderCSV(records []domain.Record) string {
	lines := []string{"id,gesture,label,status,version,owner_id"}
	for _, r := range records {
		lines = append(lines, fmt.Sprintf("%s,%d,%q,%s,%d,%s", r.ID, r.Gesture, r.Label, r.Status, r.Version, r.OwnerID))
	}
	return strings.Join(lines, "\n") + "\n"
}

func ExportView(view service.RecordView) Export {
	return BuildExport([]domain.Record{view.Record}, view.Record.UpdatedAt.Format("2006-01-02T15:04:05Z"))
}

func GroupByOwner(records []domain.Record) map[string][]domain.Record {
	grouped := make(map[string][]domain.Record)
	for _, r := range records {
		grouped[r.OwnerID] = append(grouped[r.OwnerID], r)
	}
	return grouped
}
