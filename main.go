package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_count_total",
		Help: "Counter of HTTP requests made.",
	}, []string{"code", "method"})
	requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "A histogram of latencies for requests.",
		Buckets: append([]float64{0.000001, 0.001, 0.003}, prometheus.DefBuckets...),
	}, []string{"code", "method"})
	responseSize = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_response_size_bytes",
		Help:    "A histogram of response sizes for requests.",
		Buckets: []float64{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20},
	}, []string{"code", "method"})
)

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(responseSize)
}

func handler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	fmt.Fprintf(w, "OK\n")
}

func main() {
	wrapHandler := promhttp.InstrumentHandlerCounter(
		requestCount,
		promhttp.InstrumentHandlerDuration(
			requestDuration,
			promhttp.InstrumentHandlerResponseSize(responseSize, http.HandlerFunc(handler)),
		),
	)
	http.Handle("/", wrapHandler)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("start: https://localhost:8000")

	http.ListenAndServe(":8000", nil)
}
