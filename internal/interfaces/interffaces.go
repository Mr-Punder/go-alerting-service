package interfaces

import (
	"context"

	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
)

// Logger is a global logger interface
type Logger interface {
	Info(mes string)
	Errorf(str string, arg ...any)
	Error(mess string)
	Infof(str string, arg ...any)
	Debug(mess string)
}

type MetricsGetter interface {
	GetAll(ctx context.Context) map[string]metrics.Metrics
	Get(ctx context.Context, metric metrics.Metrics) (metrics.Metrics, bool)
}

type MetricsSetter interface {
	Set(ctx context.Context, metric metrics.Metrics) error
	SetAll(ctx context.Context, metrics []metrics.Metrics) error
}

type MetricsDeleter interface {
	Delete(ctx context.Context, metric metrics.Metrics) error
}

type MetricSaver interface {
	Save(ctx context.Context) error
}

type MetricPinger interface {
	Ping() error
}

// Memstorer is a general metrics storage interface
type MetricsStorer interface {
	MetricsDeleter
	MetricsGetter
	MetricsSetter
	MetricPinger
}
