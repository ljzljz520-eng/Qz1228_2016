package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type RequestLog struct {
	Method   string
	Path     string
	Status   int
	Duration time.Duration
}
type Logger func(RequestLog)

func WithLogger(logger Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		wrapped := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		if logger != nil {
			logger(RequestLog{Method: r.Method, Path: r.URL.Path, Status: wrapped.status, Duration: time.Since(started)})
		}
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int)           { w.status = code; w.ResponseWriter.WriteHeader(code) }
func (w *statusWriter) Write(data []byte) (int, error) { return w.ResponseWriter.Write(data) }

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.Header().Set("Allow", "GET,POST,PATCH,OPTIONS")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func JSONOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPatch {
			content := r.Header.Get("Content-Type")
			if content != "" && !strings.Contains(content, "application/json") {
				writeJSON(w, http.StatusUnsupportedMediaType, map[string]string{"error": "application/json required"})
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func ErrorPayload(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func DecodeJSON(r *http.Request, value any) error { return json.NewDecoder(r.Body).Decode(value) }
