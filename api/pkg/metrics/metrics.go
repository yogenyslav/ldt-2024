package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics содержит описание всех метрик.
type Metrics struct {
	LoginCount prometheus.Counter
}

// New создает новый экземпляр метрик.
func New() *Metrics {
	return &Metrics{
		LoginCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "auth",
			Name:      "login_count",
			Help:      "Количество логинов",
		}),
	}
}

// Collect регистрирует сборку метрик.
func (m *Metrics) Collect() {
	prometheus.MustRegister(m.LoginCount)
}
