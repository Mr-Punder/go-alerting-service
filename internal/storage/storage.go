package storage

import (
	"errors"

	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
)

// MemStorage is simple implementation of storage metrics storage with map
type MemStorage struct {
	Storage map[string]metrics.Metrics
}

// GetAll returns map with all metrics
func (stor *MemStorage) GetAll() map[string]metrics.Metrics {
	if stor.Storage == nil {
		stor.Storage = make(map[string]metrics.Metrics)
	}
	return stor.Storage
}

// Set stores metric
func (stor *MemStorage) Set(metric metrics.Metrics) error {
	if stor.Storage == nil {
		stor.Storage = make(map[string]metrics.Metrics)
	}
	if metric.MType == "gauge" {
		stor.Storage[metric.ID] = metric
		return nil
	}
	if metric.MType == "counter" {
		if st, ok := stor.Storage[metric.ID]; ok {
			st.Delta += metric.Delta
			stor.Storage[metric.ID] = st
		} else {
			stor.Storage[metric.ID] = metric
		}
		return nil
	}
	return errors.New("wrong type")
}

// Get returns one metric  and it's existence
// returns metrics.Metrics{}, false if metric is not found
func (stor *MemStorage) Get(metric metrics.Metrics) (metrics.Metrics, bool) {
	if stor.Storage == nil {
		stor.Storage = make(map[string]metrics.Metrics)
	}
	if metric.MType != "gauge" && metric.MType != "counter" {
		return metrics.Metrics{}, false
	}
	m, ok := stor.Storage[metric.ID]
	return m, ok
}

// Delete deletes one gauge by name and do nothibg if the metric does not exist
func (stor *MemStorage) DeleteGouge(metric metrics.Metrics) {
	if stor.Storage == nil {
		stor.Storage = make(map[string]metrics.Metrics)
	}
	delete(stor.Storage, metric.ID)
}
