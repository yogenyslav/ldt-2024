package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics содержит описание всех метрик.
type Metrics struct {
	LoginCount              prometheus.Counter
	NumberOfReports         prometheus.Counter
	NumberOfPredictRequests prometheus.Counter
	NumberOfStockRequests   prometheus.Counter
	PredictionErrors        prometheus.Counter
}

// New создает новый экземпляр метрик.
func New() *Metrics {
	return &Metrics{
		LoginCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "auth",
			Name:      "login_count",
			Help:      "Количество логинов",
		}),
		NumberOfReports: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "predictor",
			Name:      "number_of_reports",
			Help:      "Количество отчетов, присланных пользователям (всего)",
		}),
		NumberOfPredictRequests: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "predictor",
			Name:      "number_of_predict_requests",
			Help:      "Пользователь попросил реализовать прогноз на период",
		}),
		NumberOfStockRequests: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "predictor",
			Name:      "number_of_stock_requests",
			Help:      "Пользователь попросил отчет об остатках ",
		}),
		PredictionErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "predictor",
			Name:      "prediction_errors",
			Help:      "Количество неудачных запросов на построение отчета",
		}),
	}
}

// Collect регистрирует сборку метрик.
func (m *Metrics) Collect() {
	prometheus.MustRegister(m.LoginCount)
	prometheus.MustRegister(m.NumberOfReports)
	prometheus.MustRegister(m.NumberOfPredictRequests)
	prometheus.MustRegister(m.NumberOfStockRequests)
	prometheus.MustRegister(m.PredictionErrors)
}
