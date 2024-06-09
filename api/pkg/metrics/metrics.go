package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics is the struct that holds all metrics.
type Metrics struct {
	LoginCount prometheus.Counter
}

// New creates a new Metrics.
func New() *Metrics {
	return &Metrics{
		LoginCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "auth",
			Name:      "login_count",
			Help:      "Количество логинов",
		}),
	}
}

// Collect registers all metrics.
func (m *Metrics) Collect() {
	prometheus.MustRegister(m.LoginCount)
}
