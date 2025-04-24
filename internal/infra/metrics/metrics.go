package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	ApiRequestsTotal         *prometheus.CounterVec
	EventRegistrationsTotal  *prometheus.CounterVec
	UserStatusDecisionsTotal *prometheus.CounterVec
}

func New() *Metrics {
	m := &Metrics{
		ApiRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_requests_total",
				Help: "Total number of API requests",
			},
			[]string{"method", "status"},
		),
		EventRegistrationsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "event_registrations_total",
				Help: "Total number of event registrations",
			},
			[]string{"action"}, // "registered", "cancelled"
		),
		UserStatusDecisionsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "user_status_decisions_total",
				Help: "Total number of decisions on users status",
			},
			[]string{"decision"}, // "approved", "declined"
		),
	}

	prometheus.MustRegister(
		m.ApiRequestsTotal,
		m.EventRegistrationsTotal,
		m.UserStatusDecisionsTotal,
	)

	return m
}
