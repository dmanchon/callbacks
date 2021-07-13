package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	apiDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redirect_durations_histogram_seconds",
			Help:    "Redirect endpoint latency distributions.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint", "error"},
	)
)

func health(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		apiDurationHistogram.WithLabelValues("health", "none").Observe(v)
	}))
	defer timer.ObserveDuration()
}

func handler(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		params := r.URL.Query()
		state := params.Get("state")
		if state == "" {
			http.Error(w, "'state' query param not present", http.StatusPreconditionFailed)
			apiDurationHistogram.WithLabelValues("redirect", "state_not_present").Observe(v)
			return
		}

		url, err := base64.StdEncoding.DecodeString(state)
		if err != nil {
			http.Error(w, "'state' is not valid base64 string", http.StatusPreconditionFailed)
			apiDurationHistogram.WithLabelValues("redirect", "invalid_base64").Observe(v)
			return
		}
		http.Redirect(w, r, string(url), http.StatusFound)
		apiDurationHistogram.WithLabelValues("redirect", "none").Observe(v)
	}))
	defer timer.ObserveDuration()
}

func GetEnv(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func main() {
	var port string
	var host string

	flag.StringVar(&host, "host", GetEnv("HOST", "localhost"), "Server host")
	flag.StringVar(&port, "port", GetEnv("PORT", "8080"), "Server port")
	flag.Parse()

	prometheus.MustRegister(apiDurationHistogram)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/redirect", handler)
	http.HandleFunc("/health", health)

	log.Printf("Serving app on %s:%s...", host, port)
	error := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil)
	if error != nil {
		panic(error)
	}
}
