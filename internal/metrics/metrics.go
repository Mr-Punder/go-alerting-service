package metrics

import (
	"encoding/json"
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

	if m.Delta == nil {
		Delta = 0
	} else {
		Delta = *m.Delta
	}
	if m.Value == nil {
		Value = 0
	} else {
		Value = *m.Value
	}

	aliasMetric := struct {
		*MetricAlias
		Value float64 `json:"value,omitempty"`
		Delta int64   `json:"delta,omitempty"`
	}{
		MetricAlias: (*MetricAlias)(&m),
		Delta:       Delta,
		Value:       Value,
	}
	// log.Printf("aliasMetric %v", aliasMetric)

	data, err := json.Marshal(aliasMetric)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	//log.Printf("finished json marhsling with %s", string(data))

	return data, nil
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
