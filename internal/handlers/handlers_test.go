package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mr-Punder/go-alerting-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	resp.Body.Close()
	return resp, string(body)
}

func TestMetricRouter(t *testing.T) {

	tests := []struct {
		name          string
		method        string
		wantCode      int
		wantBody      string
		wantHeaders   map[string]string
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
			wantHeaders:   map[string]string{},
		},
		{
			name:          "update gauge",
			method:        http.MethodPost,
			wantCode:      200,
			uri:           "/update/gauge/gaugeMetric/4.2",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
			wantHeaders:   map[string]string{},
		},
		{
			name:          "wrong gauge",
			method:        http.MethodPost,
			wantCode:      http.StatusBadRequest,
			uri:           "/update/gauge/anotherCounterMtric/A5C",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
			wantHeaders:   map[string]string{},
		},
		{
			name:          "noname gauge",
			method:        http.MethodPost,
			wantCode:      http.StatusNotFound,
			uri:           "/update/gauge/4.2",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
			wantHeaders:   map[string]string{},
		},
		{
			name:          "wrong method GET",
			method:        http.MethodGet,
			wantCode:      http.StatusMethodNotAllowed,
			uri:           "/update/gauge/anotherGaugeMetric/4.2",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
			wantHeaders:   map[string]string{},
		},
		{
			name:          "new counter",
			method:        http.MethodPost,
			wantCode:      200,
			uri:           "/update/counter/anothercounterMetric/42",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
			wantHeaders:   map[string]string{},
		},
		{
			name:          "update counter",
			method:        http.MethodPost,
			wantCode:      200,
			uri:           "/update/counter/counterMetric/42",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
			wantHeaders:   map[string]string{},
		},
		{
			name:          "wrong counter",
			method:        http.MethodPost,
			wantCode:      http.StatusBadRequest,
			uri:           "/update/counter/anothercounterMtric/A5C",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
			wantHeaders:   map[string]string{},
		},
		{
			name:          "noname counter",
			method:        http.MethodPost,
			wantCode:      http.StatusNotFound,
			uri:           "/update/counter/42",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
			wantHeaders:   map[string]string{},
		},
		{
			name:          "wrong method GET",
			method:        http.MethodGet,
			wantCode:      http.StatusMethodNotAllowed,
			uri:           "/update/counter/anothercounterMetric/42",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
			wantHeaders:   map[string]string{},
		},
		{
			name:          "get gauge",
			method:        http.MethodGet,
			wantCode:      http.StatusOK,
			uri:           "/value/gauge/g",
			gaugeMetric:   map[string]float64{"g": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
			wantHeaders:   map[string]string{"Content-Type": "text/plain"},
			wantBody:      "1.5",
		},
		{
			name:          "get unknown gauge",
			method:        http.MethodGet,
			wantCode:      http.StatusNotFound,
			uri:           "/value/gauge/M",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
			wantHeaders:   map[string]string{},
			wantBody:      "",
		},
		{
			name:          "get counter",
			method:        http.MethodGet,
			wantCode:      http.StatusOK,
			uri:           "/value/counter/c",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"c": 2},
			wantHeaders:   map[string]string{"Content-Type": "text/plain"},
			wantBody:      "2",
		},
		{
			name:          "get unknown counter",
			method:        http.MethodGet,
			wantCode:      http.StatusNotFound,
			uri:           "/value/counter/M",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
			wantHeaders:   map[string]string{},
			wantBody:      "",
		},
		{
			name:          "get html",
			method:        http.MethodGet,
			wantCode:      http.StatusOK,
			uri:           "",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
			wantHeaders:   map[string]string{"Content-Type": "text/html"},
			wantBody: "<html><body><h2>Gauge:</h2>" +
				fmt.Sprintf("<p>%s: %f</p>", "gaugeMetric", 1.5) +
				"<h2>Counter:</h2>" +
				fmt.Sprintf("<p>%s: %d</p>", "counterMetric", 2) +
				"</body></html>",
		},
		{
			name:          "wrong request post",
			method:        http.MethodPost,
			wantCode:      http.StatusBadRequest,
			uri:           "/another",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
			wantHeaders:   map[string]string{},
		},
		{
			name:          "wrong request get",
			method:        http.MethodGet,
			wantCode:      http.StatusBadRequest,
			uri:           "/another",
			gaugeMetric:   map[string]float64{"gaugeMetric": 1.5},
			counterMetric: map[string]int64{"counterMetric": 2},
			wantHeaders:   map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stor := &storage.MemStorage{GaugeStorage: tt.gaugeMetric, CounterStorage: tt.counterMetric}
			ts := httptest.NewServer(MetricRouter(stor))
			defer ts.Close()
			resp, body := testRequest(t, ts, tt.method, tt.uri)
			assert.Equal(t, tt.wantCode, resp.StatusCode)
			if tt.wantBody != "" {
				assert.Equal(t, tt.wantBody, body)
			}
			if len(tt.wantHeaders) > 0 {
				for key, value := range tt.wantHeaders {
					assert.Equal(t, value, resp.Header.Get(key))
				}
			}

		})
	}
}
