package storage

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/Mr-Punder/go-alerting-service/internal/interfaces"
	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
)

// MemStorage is simple implementation of storage metrics storage with map
type MemStorage struct {
	syncSave bool
	log      interfaces.Logger
	file     *os.File
	encoder  *json.Encoder
	mu       sync.Mutex
	storage  map[string]metrics.Metrics
}

func NewMemStorage(metrics map[string]metrics.Metrics, ss bool, path string, log interfaces.Logger) (*MemStorage, error) {
	var file *os.File
	var err error

	if path != "" {
		file, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
		if err != nil {
			log.Infof("cann't open file %s", err)
			log.Infof("Trying to create dir %s", filepath.Dir(path))
			err = os.MkdirAll(filepath.Dir(path), 0777)
			if err != nil {
				log.Errorf("creating directory %s", err)
				return nil, err
			}
			file, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
			if err != nil {
				log.Errorf("cann't open file %s", err)
				return nil, err

			}
		}
		log.Info("File opened")
	}

	return &MemStorage{
		syncSave: ss,
		log:      log,
		storage:  metrics,
		file:     file,
		encoder:  json.NewEncoder(file),
	}, nil
}

func (stor *MemStorage) Close() error {
	err := stor.file.Close()
	if err != nil {
		stor.log.Errorf("Error closing file", err)
		return err
	}
	stor.log.Info("File closed")
	return nil
}

func (stor *MemStorage) Ping() error {
	return nil
}

// GetAll returns map with all metrics
func (stor *MemStorage) GetAll(ctx context.Context) map[string]metrics.Metrics {
	if stor.storage == nil {
		stor.storage = make(map[string]metrics.Metrics)
	}
	return stor.storage
}

// Set stores metric
func (stor *MemStorage) Set(ctx context.Context, metric metrics.Metrics) error {
	if stor.storage == nil {
		stor.storage = make(map[string]metrics.Metrics)
	}
	if metric.MType == "gauge" {
		stor.mu.Lock()
		stor.storage[metric.ID] = metric
		stor.mu.Unlock()

	} else if metric.MType == "counter" {
		stor.mu.Lock()

		if st, ok := stor.storage[metric.ID]; ok {
			*st.Delta += *metric.Delta
			stor.storage[metric.ID] = st
		} else {
			stor.storage[metric.ID] = metric
		}
		stor.mu.Unlock()

	} else {
		return errors.New("wrong type")

	}
	if stor.syncSave {
		err := stor.Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Get returns one metric  and it's existence
// returns metrics.Metrics{}, false if metric is not found
func (stor *MemStorage) Get(ctx context.Context, metric metrics.Metrics) (metrics.Metrics, bool) {
	if stor.storage == nil {
		stor.storage = make(map[string]metrics.Metrics)
	}
	if metric.MType != "gauge" && metric.MType != "counter" {
		return metrics.Metrics{}, false
	}
	m, ok := stor.storage[metric.ID]
	return m, ok
}

// Delete deletes one gauge by name and do nothibg if the metric does not exist
func (stor *MemStorage) Delete(ctx context.Context, metric metrics.Metrics) error {
	if stor.storage == nil {
		stor.storage = make(map[string]metrics.Metrics)
	}
	delete(stor.storage, metric.ID)
	return nil
}

func (stor *MemStorage) Save(ctx context.Context) error {
	if _, err := stor.file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	err := stor.file.Truncate(0)
	if err != nil {
		stor.log.Errorf("Truncate file %s", err)
		return err
	}

	stor.mu.Lock()
	err = stor.encoder.Encode(stor.storage)
	if err != nil {
		stor.log.Errorf("encode metrics %s", err)
		return err
	}
	stor.mu.Unlock()
	stor.log.Info("Metrics saved")
	return nil
}

func (stor *MemStorage) SetAll(ctx context.Context, metrics []metrics.Metrics) error {
	for _, metric := range metrics {
		err := stor.Set(ctx, metric)
		if err != nil {
			return err
		}
	}

	return nil
}
