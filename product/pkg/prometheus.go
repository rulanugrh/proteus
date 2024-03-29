package pkg

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
)

type Metric struct {
	TotalMemory *prometheus.GaugeVec
	TotalCPU    *prometheus.GaugeVec
	Counter     *prometheus.CounterVec
	Histogram   *prometheus.HistogramVec
}

func NewPrometheus(reg prometheus.Registerer) *Metric {
	metric := Metric{
		TotalMemory: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "product",
			Name:      "total_memory",
			Help:      "Total memory used on product services",
		}, []string{"version"}),
		TotalCPU: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "product",
			Name:      "total_cpu",
			Help:      "Total cpu used on product services",
		}, []string{"version"}),
		Counter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "product",
			Name:      "counter_product",
			Help:      "Counter product on all services",
		}, []string{"type"}),
		Histogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "product",
			Name:      "histogram_product",
			Help:      "Histogram product on all services",
			Buckets:   []float64{0.1, 0.15, 0.2, 0.25, 0.3},
		}, []string{"code", "method", "type"}),
	}

	reg.MustRegister(metric.TotalCPU, metric.Histogram, metric.TotalMemory, metric.Counter)
	return &metric
}

func (m *Metric) SetTotalCPU() {
	numCPU := runtime.NumCPU()
	m.TotalCPU.WithLabelValues("v1").Set(float64(numCPU))
}

func (m *Metric) SetTotalMemory() {
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	m.TotalMemory.WithLabelValues("v1").Set(float64(memory.TotalAlloc))
}
