package metrics

// Metric is a type of Go runtime parameter
type Metrics struct {
	ID    string  `json:"id,"`
	MType string  `json:"m_type,"`
	Delta int64   `json:"delta,omitempty"`
	Value float64 `json:"value,omitempty"`
}
