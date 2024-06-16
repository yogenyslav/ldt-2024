package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics содержит описание всех метрик.
type Metrics struct {
	NumberOfActivatedCompanies prometheus.Counter
	AdminUserCount             prometheus.Counter
	AnalystUserCount           prometheus.Counter
	BuyerUserCount             prometheus.Counter
}

// New создает новый экземпляр метрик.
func New() *Metrics {
	return &Metrics{
		NumberOfActivatedCompanies: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "admin",
			Name:      "number_of_companies_activated",
			Help:      "Количество активированных компаний",
		}),
		AdminUserCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "admin",
			Name:      "admin_user_count",
			Help:      "Количество администраторов",
		}),
		AnalystUserCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "admin",
			Name:      "analyst_user_count",
			Help:      "Количество аналитиков",
		}),
		BuyerUserCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "admin",
			Name:      "buyer_user_count",
			Help:      "Количество закупщиков",
		}),
	}
}

// Collect регистрирует сборку метрик.
func (m *Metrics) Collect() {
	prometheus.MustRegister(m.NumberOfActivatedCompanies)
	prometheus.MustRegister(m.AdminUserCount)
	prometheus.MustRegister(m.AnalystUserCount)
	prometheus.MustRegister(m.BuyerUserCount)
}
