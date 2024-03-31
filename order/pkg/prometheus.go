package pkg

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrict struct {
	TotalMemory *prometheus.GaugeVec
	TotalCPU    *prometheus.GaugeVec
	Counter     *prometheus.CounterVec
	Histogram   *prometheus.HistogramVec
}

func NewPrometheus(reg prometheus.Registerer) *Metrict {
	metric := Metrict{
		TotalMemory: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "order",
			Name:      "total_memory",
			Help:      "Show total memory use in this service",
		}, []string{"version"}),
		TotalCPU: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "order",
			Name:      "total_cpu",
			Help:      "Show total cpu use in this service",
		}, []string{"version"}),
		Counter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "order",
			Name:      "counter_order",
			Help:      "Counter Order in all services",
		}, []string{"type", "service"}),
		Histogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "order",
			Name:      "histogram_order",
			Help:      "Show grafic in all order services",
		}, []string{"code", "method", "type", "services"}),
	}

	return &metric
}

func (m *Metrict) SetTotalCPU() {
	cpu := runtime.NumCPU()
	m.TotalCPU.WithLabelValues("v1").Set(float64(cpu))
}

func (m *Metrict) SetTotalMemory() {
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)

	m.TotalMemory.WithLabelValues("v1").Set(float64(memory.TotalAlloc))
}
