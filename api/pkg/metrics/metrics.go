package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	LoginCount prometheus.Counter
}

func New() *Metrics {
	return &Metrics{
		LoginCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "auth",
			Name:      "login_count",
			Help:      "Количество логинов",
		}),
	}
}

func (m *Metrics) Collect() {
	prometheus.MustRegister(m.LoginCount)
}
