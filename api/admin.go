package api

import (
	"encoding/json"
	"net/http"
	"time"

	"gesture-nebula-console/domain"
	"gesture-nebula-console/reporting"
)

type statusResponse struct {
	Status      string    `json:"status"`
	GeneratedAt time.Time `json:"generated_at"`
}

func (s *Server) admin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path
	if path == "/admin/metrics" {
		metrics, err := s.Service.Metrics(r.Context())
		if err != nil {
			ErrorPayload(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, metrics)
		return
	}
	if path == "/admin/inspect" {
		report, err := s.Service.Inspect(r.Context())
		if err != nil {
			ErrorPayload(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, report)
		return
	}
	if path == "/admin/status" {
		writeJSON(w, http.StatusOK, statusResponse{Status: "ready", GeneratedAt: time.Now().UTC()})
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func encodeRecords(records []domain.Record) ([]byte, error) {
	return json.Marshal(reporting.BuildExport(records, time.Now().UTC().Format(time.RFC3339)))
}

func parseQueryStatus(value string) (domain.Status, bool) { return domain.ParseStatus(value) }
