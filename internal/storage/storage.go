package storage

import (
	"context"
	"log"
	"time"

	"github.com/Mr-Punder/go-alerting-service/internal/logger"
	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
	"github.com/Mr-Punder/go-alerting-service/internal/metricserver"
	"github.com/Mr-Punder/go-alerting-service/internal/server/config"
)

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

func NewStorage(conf *config.Config, log logger.Logger) (MetricsStorer, func() error, error) {

	if conf.DBstring != "" {

		db, err := NewPostgreDB(conf.DBstring, log)
		if err != nil {
			log.Errorf("Error opening database", err)
			return nil, nil, err
		}
		err = db.Ping()
		if err != nil {
			log.Errorf("Db ping error %s", err)
		}
		log.Infof("Database is opened with dsn %s", conf.DBstring)

		return db, db.Close, nil
	}
	met := make(map[string]metrics.Metrics, 0)
	if conf.Restore && conf.FileStoragePath != "" {
		if err := metricserver.RestoreMetric(conf.FileStoragePath, &met, log); err != nil {
			log.Errorf("failed to restor metrics %s", err)
		}

	}

	syncSave := conf.StoreInterval <= 0 && conf.FileStoragePath != ""
	stor, err := NewMemStorage(met, syncSave, conf.FileStoragePath, log)
	if err != nil {
		log.Error("Cann't create storage")
		return nil, nil, err
	}
	log.Info("Storage created")

	if conf.StoreInterval > 0 && conf.FileStoragePath != "" {
		go SaveMetrics(stor, conf.StoreInterval, log)
		log.Info("Started metric saving goroutine")
	}

	return stor, stor.Close, nil
}

func SaveMetrics(stor MetricSaver, saveInt int64, Log logger.Logger) {
	for range time.Tick(time.Duration(saveInt) * time.Second) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err := stor.Save(ctx)
		if err != nil {
			log.Printf("Ошибка сохранения метрик на диск: %v", err)
		}
		cancel()
	}

}
