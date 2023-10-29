package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mr-Punder/go-alerting-service/internal/logger"
	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	address := server.URL

	// zapLogger, err := logger.NewLogZap("info", "stdout", "stderr")
	logrus, err := logger.NewLogLogrus("info", "stdout")
	require.NoError(t, err)
	var simpleValue = 4.2
	var simpleDelta int64 = 2

	tests := []struct {
		name    string
		metric  []metrics.Metrics
		wantErr bool
	}{
		{
			name:    "empty",
			metric:  make([]metrics.Metrics, 0),
			wantErr: false,
		},
		{
			name: "simple metric",
			metric: []metrics.Metrics{
				{
					ID:    "metric 1",
					MType: "gauge",
					Value: &simpleValue,
				},
				{
					ID:    "metric 2",
					MType: "counter",
					Delta: &simpleDelta,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantErr {
				assert.NoError(t, sendMetrics(tt.metric, address, logrus))
			}
		})
	}
}
