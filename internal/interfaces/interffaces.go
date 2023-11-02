package interfaces

import "github.com/Mr-Punder/go-alerting-service/internal/metrics"

// Logger is a global logger interface
type Logger interface {
	Info(mes string)
	Errorf(str string, arg ...any)
	Error(mess string)
	Infof(str string, arg ...any)
	Debug(mess string)
}

type MetricsAllGetter interface {
	GetAll() map[string]metrics.Metrics
}

type MetricsGetter interface {
	Get(metric metrics.Metrics) (metrics.Metrics, bool)
}

type MetricsSetter interface {
	Set(metric metrics.Metrics) error
}

type MetricsDeleter interface {
	DeleteGouge(metric metrics.Metrics)
}

type MetricSaver interface {
	Save() error
}

// Memstorer is a general metrics storage interface
type MetricsStorer interface {
	MetricsDeleter
	MetricsGetter
	MetricsSetter
	MetricsAllGetter
}
