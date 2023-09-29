package storage

type MemStor interface {
	SetGouge(string, float64) error
	SetCounter(string, int64) error
}

type MemStorage struct {
	GaugeStorage   map[string]float64
	CounterStorage map[string]int64
}

func (stor *MemStorage) SetGouge(key string, val float64) error {
	stor.GaugeStorage[key] = val
	return nil
}

func (stor *MemStorage) SetCounter(key string, val int64) error {
	stor.CounterStorage[key] += val
	return nil
}
