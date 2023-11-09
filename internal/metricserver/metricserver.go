package metricserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Mr-Punder/go-alerting-service/internal/interfaces"
	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
)

type middlewareFunc func(next http.Handler) http.Handler

type MetrciServer struct {
	Log        interfaces.Logger
	middlwares []middlewareFunc
	mux        http.Handler
	address    string
	server     *http.Server
}

func NewMetricServer(adress string, mux http.Handler, Log interfaces.Logger) *MetrciServer {

	return &MetrciServer{
		address: adress,
		mux:     mux,
		Log:     Log,
	}
}

func (ms *MetrciServer) AddMidleware(funcs ...middlewareFunc) {
	ms.middlwares = append(ms.middlwares, funcs...)
}

func (ms *MetrciServer) RunServer() {
	handler := ms.mux

	for _, f := range ms.middlwares {
		handler = f(handler)
	}

	ms.server = &http.Server{
		Addr:    ms.address,
		Handler: handler,
	}
	ms.Log.Infof("Starting server on %s", ms.address)
	if err := ms.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		ms.Log.Errorf("starting server on %s error: %s", ms.address, err)
	}

}

func RestoreMetric(path string, met *map[string]metrics.Metrics, Log interfaces.Logger) error {
	data, err := os.ReadFile(path)
	if err != nil {
		Log.Errorf("Cann't open file %s", path)
		return err
	}

	if len(data) == 0 {
		Log.Errorf("file %s is empty", path)
		return err
	}

	err = json.Unmarshal(data, met)
	if err != nil {
		Log.Error("json decoding error")
		return err
	}
	Log.Infof("Metrics restored from file %s", path)
	return nil
}

func (ms *MetrciServer) Shutdown(ctx context.Context) error {
	return ms.server.Shutdown(ctx)
}

func SaveMetrics(stor interfaces.MetricSaver, saveInt int64, Log interfaces.Logger) {
	for range time.Tick(time.Duration(saveInt) * time.Second) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err := stor.Save(ctx)
		if err != nil {
			log.Printf("Ошибка сохранения метрик на диск: %v", err)
		}
		cancel()
	}

}
