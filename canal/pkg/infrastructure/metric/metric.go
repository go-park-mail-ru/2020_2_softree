package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

var metric = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "hits"}, []string{"status", "url"})

func RecordHitMetric(code int, url string) {
	metric.WithLabelValues(strconv.Itoa(code), url).Inc()
}

var timer = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "timerAction",
		Help:       "Timer running action",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"method"},
)

func RecordTimeMetric(start time.Time, src string) {
	timer.WithLabelValues(src).Observe(float64(time.Since(start).Seconds()))
}

func Initialize() {
	prometheus.MustRegister(metric)
	prometheus.MustRegister(timer)
}
