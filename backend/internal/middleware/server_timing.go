package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type serverTimingContextKey struct{}

type serverTimingMetric struct {
	name     string
	duration time.Duration
}

type serverTimingRecorder struct {
	startedAt time.Time
	mu        sync.Mutex
	metrics   []serverTimingMetric
}

// NewServerTiming adds Server-Timing metrics to normal API responses.
func NewServerTiming() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := &serverTimingRecorder{startedAt: time.Now()}
			ctx := context.WithValue(r.Context(), serverTimingContextKey{}, recorder)
			tw := &serverTimingResponseWriter{
				ResponseWriter: w,
				recorder:       recorder,
			}

			next.ServeHTTP(tw, r.WithContext(ctx))
			tw.setHeader()
		})
	}
}

// AddServerTimingMetric records a duration for the current request if timing is enabled.
func AddServerTimingMetric(ctx context.Context, name string, duration time.Duration) {
	recorder, ok := ctx.Value(serverTimingContextKey{}).(*serverTimingRecorder)
	if !ok || recorder == nil {
		return
	}
	recorder.add(name, duration)
}

func addServerTimingSince(ctx context.Context, name string, startedAt time.Time) {
	AddServerTimingMetric(ctx, name, time.Since(startedAt))
}

func (r *serverTimingRecorder) add(name string, duration time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics = append(r.metrics, serverTimingMetric{
		name:     sanitizeServerTimingName(name),
		duration: duration,
	})
}

func (r *serverTimingRecorder) headerValue() string {
	r.mu.Lock()
	metrics := append([]serverTimingMetric(nil), r.metrics...)
	startedAt := r.startedAt
	r.mu.Unlock()

	metrics = append(metrics, serverTimingMetric{
		name:     "app",
		duration: time.Since(startedAt),
	})

	parts := make([]string, 0, len(metrics))
	for _, metric := range metrics {
		parts = append(parts, fmt.Sprintf("%s;dur=%.1f", metric.name, float64(metric.duration.Microseconds())/1000))
	}
	return strings.Join(parts, ", ")
}

type serverTimingResponseWriter struct {
	http.ResponseWriter
	recorder  *serverTimingRecorder
	headerSet bool
}

func (w *serverTimingResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *serverTimingResponseWriter) WriteHeader(statusCode int) {
	w.setHeader()
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *serverTimingResponseWriter) Write(b []byte) (int, error) {
	w.setHeader()
	return w.ResponseWriter.Write(b)
}

func (w *serverTimingResponseWriter) setHeader() {
	if w.headerSet {
		return
	}
	w.headerSet = true

	value := w.recorder.headerValue()
	if existing := w.Header().Get("Server-Timing"); existing != "" {
		value = existing + ", " + value
	}
	w.Header().Set("Server-Timing", value)
}

func sanitizeServerTimingName(name string) string {
	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '_' || r == '-' || r == '.':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	if b.Len() == 0 {
		return "metric"
	}
	return b.String()
}
