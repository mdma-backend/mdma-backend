package metrics

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	statusCodeLabel = "status_code"
	endpointLabel   = "endpoint"
	methodLabel     = "method"
)

var (
	uuidPattern = regexp.MustCompile(`.*\/([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})[\/.*]?`)
	idPattern   = regexp.MustCompile(`.*\/([0-9]+)[\/.*]?`)
)

var (
	httpResponseTime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{Name: "http_response_time"},
		[]string{statusCodeLabel, endpointLabel, methodLabel},
	)
	httpCalls = promauto.NewCounterVec(
		prometheus.CounterOpts{Name: "http_calls"},
		[]string{statusCodeLabel, endpointLabel, methodLabel},
	)
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		start := time.Now()
		defer func() {
			status := ww.Status()
			if status == 0 {
				status = 200
			}

			path := r.URL.Path
			if mm := uuidPattern.FindStringSubmatch(path); len(mm) == 2 {
				path = strings.Replace(path, mm[1], "{uuid}", -1)
			} else if mm := idPattern.FindStringSubmatch(path); len(mm) == 2 {
				path = strings.Replace(path, mm[1], "{id}", -1)
			}

			httpResponseTime.With(
				prometheus.Labels{
					statusCodeLabel: strconv.Itoa(status),
					endpointLabel:   path,
					methodLabel:     r.Method,
				},
			).Set(float64(time.Since(start).Seconds()))
			httpCalls.With(
				prometheus.Labels{
					statusCodeLabel: strconv.Itoa(status),
					endpointLabel:   path,
					methodLabel:     r.Method,
				},
			).Inc()
		}()

		next.ServeHTTP(ww, r)
	})
}
