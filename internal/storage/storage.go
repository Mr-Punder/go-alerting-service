package storage

type MemStor interface {
	SetGouge(string, float64) error
	SetCounter(string, int64) error
	GetGouge(string) (float64, bool)
	GetCounter(string) (int64, bool)
	DeleteGouge(string)
	DeleteCouner(string)
}

type MemStorage struct {
	GaugeStorage   map[string]float64
	CounterStorage map[string]int64
}

func (stor *MemStorage) SetGouge(key string, val float64) error {
	if stor.GaugeStorage == nil {
		stor.GaugeStorage = make(map[string]float64)
	}
	stor.GaugeStorage[key] = val
	return nil
}

func (stor *MemStorage) SetCounter(key string, val int64) error {
	if stor.CounterStorage == nil {
		stor.CounterStorage = make(map[string]int64)
	}
	stor.CounterStorage[key] += val
	return nil
}

func (stor *MemStorage) GetGouge(name string) (float64, bool) {
	if stor.GaugeStorage == nil {
		stor.GaugeStorage = make(map[string]float64)
	}
	val, ok := stor.GaugeStorage[name]
	return val, ok
}

func (stor *MemStorage) GetCounter(name string) (int64, bool) {
	if stor.CounterStorage == nil {
		stor.CounterStorage = make(map[string]int64)
	}
	val, ok := stor.CounterStorage[name]
	return val, ok
}

func (stor *MemStorage) DeleteGouge(name string) {
	delete(stor.GaugeStorage, name)
}

func (stor *MemStorage) DeleteCouner(name string) {
	delete(stor.GaugeStorage, name)
}
