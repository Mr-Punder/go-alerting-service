package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mr-Punder/go-alerting-service/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestGaugeUpd(t *testing.T) {

	tests := []struct {
		name          string
		method        string
		wantCode      int
		uri           string
		gaugeMetric   map[string]float64
		counterMetric map[string]int64
	}{
		{
			name:          "new gauge",
			method:        http.MethodPost,
			wantCode:      200,
			uri:           "/update/gauge/anotherGaugeMetric/4.2",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
		},
		{
			name:          "update gauge",
			method:        http.MethodPost,
			wantCode:      200,
			uri:           "/update/gauge/gaugeMetric/4.2",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
		},
		{
			name:          "wrong gauge",
			method:        http.MethodPost,
			wantCode:      http.StatusBadRequest,
			uri:           "/update/gauge/anotherCounterMtric/A5C",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
		},
		{
			name:          "noname gauge",
			method:        http.MethodPost,
			wantCode:      http.StatusNotFound,
			uri:           "/update/gauge/4.2",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
		},
		{
			name:          "wrong method GET",
			method:        http.MethodGet,
			wantCode:      http.StatusMethodNotAllowed,
			uri:           "/update/gauge/anotherGaugeMetric/4.2",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.uri, nil)
			w := httptest.NewRecorder()
			stor := &storage.MemStorage{GaugeStorage: tt.gaugeMetric, CounterStorage: tt.counterMetric}
			GaugeUpd(stor)(w, r)
			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func TestCounterUpd(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		wantCode      int
		uri           string
		gaugeMetric   map[string]float64
		counterMetric map[string]int64
	}{
		{
			name:          "new counter",
			method:        http.MethodPost,
			wantCode:      200,
			uri:           "/update/counter/anothercounterMetric/42",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
		},
		{
			name:          "update counter",
			method:        http.MethodPost,
			wantCode:      200,
			uri:           "/update/counter/counterMetric/42",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
		},
		{
			name:          "wrong counter",
			method:        http.MethodPost,
			wantCode:      http.StatusBadRequest,
			uri:           "/update/counter/anothercounterMtric/A5C",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
		},
		{
			name:          "noname counter",
			method:        http.MethodPost,
			wantCode:      http.StatusNotFound,
			uri:           "/update/counter/42",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
		},
		{
			name:          "wrong method GET",
			method:        http.MethodGet,
			wantCode:      http.StatusMethodNotAllowed,
			uri:           "/update/counter/anothercounterMetric/42",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.uri, nil)
			w := httptest.NewRecorder()
			stor := &storage.MemStorage{GaugeStorage: tt.gaugeMetric, CounterStorage: tt.counterMetric}
			CounterUpd(stor)(w, r)
			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
