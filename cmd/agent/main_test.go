package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
	"github.com/stretchr/testify/assert"
)

func TestSendMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	address := server.URL

	tests := []struct {
		name    string
		metric  []metrics.Metric
		wantErr bool
	}{
		{
			name:    "empty",
			metric:  make([]metrics.Metric, 0),
			wantErr: false,
		},
		{
			name: "simple metric",
			metric: []metrics.Metric{
				{
					Name: "metric 1",
					Type: "gauge",
					Val:  "4.2",
				},
				{
					Name: "metric 2",
					Type: "counter",
					Val:  "42",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				assert.NoError(t, sendMetrics(tt.metric, address))
			}
		})
	}
}
