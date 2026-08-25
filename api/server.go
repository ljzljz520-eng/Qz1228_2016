package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"gesture-nebula-console/notifier"
	"gesture-nebula-console/service"
)

type Server struct {
	Service  *service.Service
	Notifier *notifier.Notifier
	Mux      *http.ServeMux
}

func NewServer(svc *service.Service, n *notifier.Notifier) *Server {
	s := &Server{Service: svc, Notifier: n, Mux: http.NewServeMux()}
	s.Mux.HandleFunc("/health", s.health)
	s.Mux.HandleFunc("/metrics", s.metrics)
	s.Mux.HandleFunc("/admin/", s.admin)
	s.Mux.HandleFunc("/records", s.records)
	s.Mux.HandleFunc("/records/", s.record)
	return s
}

func (s *Server) Handler() http.Handler { return CORS(JSONOnly(WithLogger(nil, logging(s.Mux)))) }

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Nebula-Service", "gesture-15")
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func pathID(path string) string { return strings.Trim(strings.TrimPrefix(path, "/records/"), "/") }
