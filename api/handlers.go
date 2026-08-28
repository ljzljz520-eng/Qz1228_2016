package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"gesture-nebula-console/domain"
	"gesture-nebula-console/reporting"
	"gesture-nebula-console/service"
)

type recordRequest struct {
	ID              string `json:"id"`
	OwnerID         string `json:"owner_id"`
	Label           string `json:"label"`
	ExpectedVersion int    `json:"expected_version"`
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "gesture": "15"})
}

func (s *Server) metrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := s.Service.Metrics(r.Context())
	if err != nil {
		ErrorPayload(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, metrics)
}

func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		query := domain.Query{Text: r.URL.Query().Get("q"), OwnerID: r.URL.Query().Get("owner"), IncludeArchived: r.URL.Query().Get("archived") == "true"}
		if status, ok := domain.ParseStatus(r.URL.Query().Get("status")); ok {
			query.Status = status
		}
		items, err := s.Service.Search(r.Context(), query)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, items)
	case http.MethodPost:
		var req recordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		item, err := s.Service.RegisterRecord(r.Context(), req.ID, req.OwnerID, req.Label)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, item)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) record(w http.ResponseWriter, r *http.Request) {
	id := pathID(r.URL.Path)
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "record id required"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		view, err := s.Service.QueryRecord(r.Context(), id)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"view": view, "summary": reporting.BuildSummary(view), "timeline": reporting.BuildTimeline(view.Events, view.Audits)})
	case http.MethodPatch:
		var req recordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		item, err := s.Service.UpdateLabel(r.Context(), id, req.ExpectedVersion, req.Label)
		if err != nil {
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		if s.Notifier != nil {
			_ = s.Notifier.Dispatch(context.Background(), item)
		}
		writeJSON(w, http.StatusOK, item)
	case http.MethodPost:
		if strings.HasSuffix(r.URL.Path, "/review") {
			item, err := s.Service.ReviewRecord(r.Context(), id, 1)
			if err != nil {
				writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, item)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func DecodeRecord(data []byte) (domain.Record, error) {
	var r domain.Record
	err := json.Unmarshal(data, &r)
	return r, err
}

func ServiceForHandler(s *service.Service) http.Handler { return NewServer(s, nil).Handler() }
