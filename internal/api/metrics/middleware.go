package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	statusCodeLabel = "status_code"
	endpointLabel   = "endpoint"
)

var (
	httpResponseTime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{Name: "http_response_time"},
		[]string{statusCodeLabel, endpointLabel},
	)
	httpCalls = promauto.NewCounterVec(
		prometheus.CounterOpts{Name: "http_calls"},
		[]string{statusCodeLabel, endpointLabel},
	)
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		start := time.Now()
		defer func() {
			httpResponseTime.With(
				prometheus.Labels{
					statusCodeLabel: strconv.Itoa(ww.Status()),
					endpointLabel:   r.URL.Path,
				},
			).Set(float64(time.Since(start).Seconds()))
			httpCalls.With(
				prometheus.Labels{
					statusCodeLabel: strconv.Itoa(ww.Status()),
					endpointLabel:   r.URL.Path,
				},
			).Inc()
		}()

		next.ServeHTTP(ww, r)
	})
}
