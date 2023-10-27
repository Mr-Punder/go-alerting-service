package metrics

import (
	"encoding/json"
	"errors"
	"log"
)

// Metric is a type of Go runtime parameter
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m Metrics) MarshalJSON() ([]byte, error) {

	//log.Printf("started json marhsling %v", m)
	type MetricAlias Metrics
	var Delta int64
	var Value float64

	switch m.MType {
	case "gauge":
		if m.Value == nil {
			Value = 0
		} else {
			Value = *m.Value
		}
		aliasMetric := struct {
			*MetricAlias
			Value float64 `json:"value"`
		}{
			MetricAlias: (*MetricAlias)(&m),
			Value:       Value,
		}
		data, err := json.Marshal(aliasMetric)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		return data, nil
	case "counter":
		if m.Delta == nil {
			Delta = 0
		} else {
			Delta = *m.Delta
		}
		aliasMetric := struct {
			*MetricAlias

			Delta int64 `json:"delta"`
		}{
			MetricAlias: (*MetricAlias)(&m),
			Delta:       Delta,
		}
		data, err := json.Marshal(aliasMetric)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		return data, nil
	default:
		return []byte{}, errors.New("uncknown type")
	}

}

// func (m *Metrics) UnmarshalJSON(data []byte) (err error) {
// 	type MetricAlias Metrics

// 	aliasMetric := struct {
// 		*MetricAlias
// 		Value float64
// 		Delta int64
// 	}{
// 		MetricAlias: (*MetricAlias)(m),
// 	}

// 	if err := json.Unmarshal(data, &aliasMetric); err != nil {
// 		return err
// 	}
// 	m.Delta = &aliasMetric.Delta
// 	m.Value = &aliasMetric.Value
// 	return nil
// }
