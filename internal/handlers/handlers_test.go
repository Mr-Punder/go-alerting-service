package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mr-Punder/go-alerting-service/internal/logger"
	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
	"github.com/Mr-Punder/go-alerting-service/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, sentbody string, sentheaders map[string]string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer([]byte(sentbody)))
	require.NoError(t, err)

	for key, value := range sentheaders {
		req.Header.Set(key, value)
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	resp.Body.Close()
	return resp, string(body)
}

func TestMetricRouter(t *testing.T) {

	simpleValue := func() *float64 { var v float64 = 1.5; return &v }
	simpleDelta := func() *int64 { var d int64 = 2; return &d }

	tests := []struct {
		name        string
		method      string
		wantCode    int
		wantBody    string
		sentBody    string
		sentHeaders map[string]string
		wantHeaders map[string]string
		uri         string
		metrics     map[string]metrics.Metrics
	}{
		{
			name:     "new gauge",
			method:   http.MethodPost,
			wantCode: 200,
			uri:      "/update/gauge/anotherGaugeMetric/4.2",
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"counterMetric": {
					ID:    "counterMetric",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			wantHeaders: map[string]string{},
			sentHeaders: map[string]string{},
		},
		{
			name:     "update gauge",
			method:   http.MethodPost,
			wantCode: 200,
			uri:      "/update/gauge/gaugeMetric/4.2",
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"counterMetric": {
					ID:    "counterMetric",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			sentHeaders: map[string]string{},
			wantHeaders: map[string]string{},
		},
		{
			name:     "wrong gauge",
			method:   http.MethodPost,
			wantCode: http.StatusBadRequest,
			uri:      "/update/gauge/anotherCounterMtric/A5C",
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"counterMetric": {
					ID:    "counterMetric",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			sentHeaders: map[string]string{},
			wantHeaders: map[string]string{},
		},
		{
			name:     "noname gauge",
			method:   http.MethodPost,
			wantCode: http.StatusNotFound,
			uri:      "/update/gauge/4.2",
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"counterMetric": {
					ID:    "counterMetric",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			sentHeaders: map[string]string{},
			wantHeaders: map[string]string{},
		},
		{
			name:     "wrong method GET",
			method:   http.MethodGet,
			wantCode: http.StatusMethodNotAllowed,
			uri:      "/update/gauge/anotherGaugeMetric/4.2",
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"counterMetric": {
					ID:    "counterMetric",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			sentHeaders: map[string]string{},
			wantHeaders: map[string]string{},
		},
		{
			name:     "new counter",
			method:   http.MethodPost,
			wantCode: 200,
			uri:      "/update/counter/anothercounterMetric/42",
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"counterMetric": {
					ID:    "counterMetric",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			sentHeaders: map[string]string{},
			wantHeaders: map[string]string{},
		},
		{
			name:     "update counter",
			method:   http.MethodPost,
			wantCode: 200,
			uri:      "/update/counter/counterMetric/42",
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"counterMetric": {
					ID:    "counterMetric",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			sentHeaders: map[string]string{},
			wantHeaders: map[string]string{},
		},
		{
			name:     "wrong counter",
			method:   http.MethodPost,
			wantCode: http.StatusBadRequest,
			uri:      "/update/counter/anothercounterMtric/A5C",
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"counterMetric": {
					ID:    "counterMetric",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			sentHeaders: map[string]string{},
			wantHeaders: map[string]string{},
		},
		{
			name:     "noname counter",
			method:   http.MethodPost,
			wantCode: http.StatusNotFound,
			uri:      "/update/counter/42",
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"counterMetric": {
					ID:    "counterMetric",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			sentHeaders: map[string]string{},
			wantHeaders: map[string]string{},
		},
		{
			name:     "wrong method GET",
			method:   http.MethodGet,
			wantCode: http.StatusMethodNotAllowed,
			uri:      "/update/counter/anothercounterMetric/42",
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"counterMetric": {
					ID:    "counterMetric",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			sentHeaders: map[string]string{},
			wantHeaders: map[string]string{},
		},
		{
			name:     "get gauge",
			method:   http.MethodGet,
			wantCode: http.StatusOK,
			uri:      "/value/gauge/g",
			metrics: map[string]metrics.Metrics{
				"g": {
					ID:    "g",
					MType: "gauge",
					Value: simpleValue(),
				},
				"counterMetric": {
					ID:    "counterMetric",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			sentHeaders: map[string]string{},
			wantHeaders: map[string]string{"Content-Type": "text/plain"},
			wantBody:    "1.5",
		},
		{
			name:     "get unknown gauge",
			method:   http.MethodGet,
			wantCode: http.StatusNotFound,
			uri:      "/value/gauge/M",
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"counterMetric": {
					ID:    "counterMetric",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			sentHeaders: map[string]string{},
			wantHeaders: map[string]string{},
			wantBody:    "",
		},
		{
			name:     "get counter",
			method:   http.MethodGet,
			wantCode: http.StatusOK,
			uri:      "/value/counter/c",
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"c": {
					ID:    "c",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			sentHeaders: map[string]string{},
			wantHeaders: map[string]string{"Content-Type": "text/plain"},
			wantBody:    "2",
		},
		{
			name:     "get unknown counter",
			method:   http.MethodGet,
			wantCode: http.StatusNotFound,
			uri:      "/value/counter/M",
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"counterMetric": {
					ID:    "counterMetric",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			sentHeaders: map[string]string{},
			wantHeaders: map[string]string{},
			wantBody:    "",
		},
		{
			name:     "get html",
			method:   http.MethodGet,
			wantCode: http.StatusOK,
			uri:      "",
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"counterMetric": {
					ID:    "counterMetric",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			sentHeaders: map[string]string{},
			wantHeaders: map[string]string{"Content-Type": "text/html"},
			wantBody: "<html><body><h2>Gauge:</h2>" +
				fmt.Sprintf("<p>%s: %f</p>", "gaugeMetric", 1.5) +
				"<h2>Counter:</h2>" +
				fmt.Sprintf("<p>%s: %d</p>", "counterMetric", 2) +
				"</body></html>",
		},
		{
			name:     "wrong request post",
			method:   http.MethodPost,
			wantCode: http.StatusBadRequest,
			uri:      "/another",
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"counterMetric": {
					ID:    "counterMetric",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			sentHeaders: map[string]string{},
			wantHeaders: map[string]string{},
		},
		{
			name:     "wrong request get",
			method:   http.MethodGet,
			wantCode: http.StatusBadRequest,
			uri:      "/another",
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"counterMetric": {
					ID:    "counterMetric",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			sentHeaders: map[string]string{},
			wantHeaders: map[string]string{},
		},
		{
			name:        "new gauge with JSON",
			method:      http.MethodPost,
			wantCode:    200,
			uri:         "/update",
			sentBody:    `{"id":"g_new","type":"gauge","value":5.2}`,
			sentHeaders: map[string]string{"Content-Type": "application/json"},
			metrics: map[string]metrics.Metrics{
				"g": {
					ID:    "g",
					MType: "gauge",
					Value: simpleValue(),
				},
				"c": {
					ID:    "c",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			wantBody:    `{"id":"g_new","type":"gauge","value":5.2}`,
			wantHeaders: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:        "get gauge with JSON",
			method:      http.MethodPost,
			wantCode:    200,
			uri:         "/value",
			sentBody:    `{"id":"g","type":"gauge"}`,
			sentHeaders: map[string]string{"Content-Type": "application/json"},
			metrics: map[string]metrics.Metrics{
				"g": {
					ID:    "g",
					MType: "gauge",
					Value: simpleValue(),
				},
				"c": {
					ID:    "c",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			wantBody:    `{"id":"g","type":"gauge","value":1.5}`,
			wantHeaders: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:        "new counter with JSON",
			method:      http.MethodPost,
			wantCode:    200,
			uri:         "/update",
			sentBody:    `{"id":"c_new","type":"counter","delta":5}`,
			sentHeaders: map[string]string{"Content-Type": "application/json"},
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"c": {
					ID:    "c",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			wantBody:    `{"id":"c_new","type":"counter","delta":5}`,
			wantHeaders: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:        "new wrong type with JSON",
			method:      http.MethodPost,
			wantCode:    http.StatusBadRequest,
			uri:         "/update",
			sentBody:    `{"id":"c_new","type":"smth","delta":5}`,
			sentHeaders: map[string]string{"Content-Type": "application/json"},
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"c": {
					ID:    "c",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			wantBody:    "",
			wantHeaders: map[string]string{},
		},
		{
			name:        "update counter with JSON",
			method:      http.MethodPost,
			wantCode:    200,
			uri:         "/update",
			sentBody:    `{"id":"c","type":"counter","delta":5}`,
			sentHeaders: map[string]string{"Content-Type": "application/json"},
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"c": {
					ID:    "c",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			wantBody:    `{"id":"c","type":"counter","delta":7}`,
			wantHeaders: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:        "update gauge with JSON",
			method:      http.MethodPost,
			wantCode:    200,
			uri:         "/update",
			sentBody:    `{"id":"g","type":"gauge","value":5.2}`,
			sentHeaders: map[string]string{"Content-Type": "application/json"},
			metrics: map[string]metrics.Metrics{
				"g": {
					ID:    "g",
					MType: "gauge",
					Value: simpleValue(),
				},
				"c": {
					ID:    "c",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			wantBody:    `{"id":"g","type":"gauge","value":5.2}`,
			wantHeaders: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:        "get counter with JSON",
			method:      http.MethodPost,
			wantCode:    200,
			uri:         "/value",
			sentBody:    `{"id":"c","type":"counter"}`,
			sentHeaders: map[string]string{"Content-Type": "application/json"},
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"c": {
					ID:    "c",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			wantBody:    `{"id":"c","type":"counter","delta":2}`,
			wantHeaders: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:        "get unknown counter with JSON",
			method:      http.MethodPost,
			wantCode:    http.StatusNotFound,
			uri:         "/value",
			sentBody:    `{"id":"c_unk","type":"counter"}`,
			sentHeaders: map[string]string{"Content-Type": "application/json"},
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"c": {
					ID:    "c",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			wantBody:    "",
			wantHeaders: map[string]string{},
		},
		{
			name:        "get unknown gauge with JSON",
			method:      http.MethodPost,
			wantCode:    http.StatusNotFound,
			uri:         "/value",
			sentBody:    `{"id":"g_unk","type":"gauge"}`,
			sentHeaders: map[string]string{"Content-Type": "application/json"},
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"c": {
					ID:    "c",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			wantBody:    "",
			wantHeaders: map[string]string{},
		},
		{
			name:        "get wrong type with JSON",
			method:      http.MethodPost,
			wantCode:    http.StatusBadRequest,
			uri:         "/value",
			sentBody:    `{"id":"c_unk","type":"smth"}`,
			sentHeaders: map[string]string{"Content-Type": "application/json"},
			metrics: map[string]metrics.Metrics{
				"gaugeMetric": {
					ID:    "gaugeMetric",
					MType: "gauge",
					Value: simpleValue(),
				},
				"c": {
					ID:    "c",
					MType: "counter",
					Delta: simpleDelta(),
				},
			},
			wantBody:    "",
			wantHeaders: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stor := &storage.MemStorage{Storage: tt.metrics}
			Log, err := logger.NewLogZap("info", "./log.txt", "stderr")
			require.NoError(t, err)

			ts := httptest.NewServer(MetricRouter(stor, Log))
			defer ts.Close()
			resp, body := testRequest(t, ts, tt.method, tt.uri, tt.sentBody, tt.sentHeaders)
			defer resp.Body.Close()
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
