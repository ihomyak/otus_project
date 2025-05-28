package middleware

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)
	httpDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10, 30, 60},
		},
		[]string{"path", "method", "status"},
	)
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		path := normalizePath(r.URL.Path)
		duration := time.Since(start).Seconds()
		httpRequestsTotal.WithLabelValues(
			path,
			r.Method,
			strconv.Itoa(rw.status),
		).Inc()
		httpDuration.WithLabelValues(
			path,
			r.Method,
			strconv.Itoa(rw.status),
		).Observe(duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func normalizePath(path string) string {
	path = strings.Split(path, "?")[0]

	uuidRegex := regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	path = uuidRegex.ReplaceAllString(path, "{uuid}")

	numRegex := regexp.MustCompile(`/\d+(/|$)`)
	path = numRegex.ReplaceAllString(path, "/{id}$1")

	idRegex := regexp.MustCompile(`(/bot/|/token/|/hook/)[^/]+(/|$)`)
	path = idRegex.ReplaceAllString(path, "$1{id}$2")

	return path
}
