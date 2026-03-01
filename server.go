package main

import (
   "fmt"
   "net/http"
   "time"

   "github.com/prometheus/client_golang/prometheus"
   "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
   pingCounter = prometheus.NewCounter(
       prometheus.CounterOpts{
           Name: "ping_request_count",
           Help: "No of request handled by Ping handler",
       },
   )
   activeRequests = prometheus.NewGauge(
       prometheus.GaugeOpts{
           Name: "active_requests",
           Help: "Current number of active requests",
       },
   )
   requestDuration = prometheus.NewHistogram(
       prometheus.HistogramOpts{
           Name:    "request_duration_seconds",
           Help:    "Histogram of request durations in seconds",
           Buckets: []float64{0.1, 0.5, 1, 2, 5},
       },
   )
)

func ping(w http.ResponseWriter, req *http.Request) {
   activeRequests.Inc()
   defer activeRequests.Dec()
   
   timer := prometheus.NewTimer(requestDuration)
   defer timer.ObserveDuration()

   pingCounter.Inc()
   time.Sleep(100 * time.Millisecond) // Simulate some work
   fmt.Fprintf(w, "pong")
}

func main() {
   prometheus.MustRegister(pingCounter, activeRequests, requestDuration)

   http.HandleFunc("/ping", ping)
   http.Handle("/metrics", promhttp.Handler())
   http.ListenAndServe(":8090", nil)
}