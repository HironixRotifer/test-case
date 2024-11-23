package prometheus

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	cpuTemp    prometheus.Gauge
	hdFailures *prometheus.CounterVec

	allocated       prometheus.Gauge
	total_allocated prometheus.Gauge
	sys             prometheus.Gauge
	num_gc          prometheus.Gauge
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		cpuTemp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cpu_temperature_celsius",
			Help: "Current temperature of the CPU.",
		}),
		hdFailures: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "hd_errors_total",
				Help: "Number of hard-disk errors.",
			},
			[]string{"device"},
		),
	}

	// reg.MustRegister(m.cpuTemp)
	// m.hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()

	reg.MustRegister(m.hdFailures)
	reg.MustRegister(m.allocated)
	reg.MustRegister(m.total_allocated)
	reg.MustRegister(m.sys)
	reg.MustRegister(m.num_gc)

	return m
}

func (m *metrics) CollectMemUsageMetric(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	var memory runtime.MemStats

	for {
		select {
		case <-ctx.Done():
			return
		default:
			runtime.ReadMemStats(&memory)

			m.allocated.Set(float64(memory.Alloc))
			m.total_allocated.Set(float64(memory.TotalAlloc))
			m.sys.Set(float64(memory.Sys))
			m.num_gc.Set(float64(memory.NumGC))

			time.Sleep(10 * time.Second)
		}
	}
}
