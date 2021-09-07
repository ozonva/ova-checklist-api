package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics interface {
	CreateChecklistError()
	CreateChecklistSuccess()

	MultiCreateChecklistError()
	MultiCreateChecklistSuccess()

	RemoveChecklistError()
	RemoveChecklistSuccess()

	UpdateChecklistError()
	UpdateChecklistSuccess()
}

type metrics struct {
	createError   prometheus.Counter
	createSuccess prometheus.Counter

	multiCreateError   prometheus.Counter
	multiCreateSuccess prometheus.Counter

	removeError   prometheus.Counter
	removeSuccess prometheus.Counter

	updateError   prometheus.Counter
	updateSuccess prometheus.Counter
}

func (m *metrics) CreateChecklistError() {
	m.createError.Inc()
}

func (m *metrics) CreateChecklistSuccess() {
	m.createSuccess.Inc()
}

func (m *metrics) MultiCreateChecklistError() {
	m.multiCreateError.Inc()
}

func (m *metrics) MultiCreateChecklistSuccess() {
	m.multiCreateSuccess.Inc()
}

func (m *metrics) RemoveChecklistError() {
	m.removeError.Inc()
}

func (m *metrics) RemoveChecklistSuccess() {
	m.removeSuccess.Inc()
}

func (m *metrics) UpdateChecklistError() {
	m.updateError.Inc()
}

func (m *metrics) UpdateChecklistSuccess() {
	m.updateSuccess.Inc()
}

func registerGrpcApiMetrics(m *metrics) {
	m.createError = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "grpc_create_checklist_response_error",
		Subsystem: "ova_checklist_api",
	})
	m.createSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "grpc_create_checklist_response_success",
		Subsystem: "ova_checklist_api",
	})

	m.multiCreateError = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "grpc_multi_create_checklist_response_error",
		Subsystem: "ova_checklist_api",
	})
	m.multiCreateSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "grpc_multi_create_checklist_response_success",
		Subsystem: "ova_checklist_api",
	})

	m.removeError = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "grpc_remove_checklist_response_error",
		Subsystem: "ova_checklist_api",
	})
	m.removeSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "grpc_remove_checklist_response_success",
		Subsystem: "ova_checklist_api",
	})

	m.updateError = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "grpc_update_checklist_response_error",
		Subsystem: "ova_checklist_api",
	})
	m.updateSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "grpc_update_checklist_response_success",
		Subsystem: "ova_checklist_api",
	})

	prometheus.MustRegister(m.createError)
	prometheus.MustRegister(m.createSuccess)
	prometheus.MustRegister(m.multiCreateError)
	prometheus.MustRegister(m.multiCreateSuccess)
	prometheus.MustRegister(m.removeError)
	prometheus.MustRegister(m.removeSuccess)
	prometheus.MustRegister(m.updateError)
	prometheus.MustRegister(m.updateSuccess)
}

func NewMetrics() Metrics {
	m := &metrics{}
	registerGrpcApiMetrics(m)
	return m
}
