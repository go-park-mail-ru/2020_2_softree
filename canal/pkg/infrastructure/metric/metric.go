package metric

import "github.com/prometheus/client_golang/prometheus"

var Metric = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "hits"}, []string{"status"})
