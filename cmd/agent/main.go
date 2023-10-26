package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"

	"time"

	"github.com/Mr-Punder/go-alerting-service/internal/agent/config"
	"github.com/Mr-Punder/go-alerting-service/internal/logger"
	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
	"github.com/go-resty/resty/v2"
)

func main() {
	config.ParseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func sendMetrics(metrics []metrics.Metrics, addres string) error {
	init := fmt.Sprintf("%s/update", addres)
	client := resty.New()

	for _, metric := range metrics {
		url := init
		body, err := json.Marshal(metric)
		if err != nil {
			panic(err)
		}

		resp, err := client.R().SetHeader("Content-Type", "application/json").SetBody(body).Post(url)

		if err != nil {
			return err
		}
		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}
	}
	return nil
}

func run() error {
	zapLogger, err := logger.NewLogZap("info", "stdout", "stderr")
	if err != nil {
		return err
	}
	pollTicker := time.NewTicker(config.PollInterval)
	reportTicker := time.NewTicker(config.ReportInterval)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	metric := Collect()
	zapLogger.Info("metrics collected")
	for {
		select {
		case <-pollTicker.C:
			metric = Collect()

		case <-reportTicker.C:
			address := "http://" + config.ServerAddress
			zapLogger.Info(fmt.Sprintf("metrics sent to %s", address))
			err := sendMetrics(metric, address)
			if err != nil {
				return err
			}
		}
	}

}

// Collect returns slice of current runtime metrics
func Collect() []metrics.Metrics {

	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	return []metrics.Metrics{
		{
			ID:    "Alloc",
			MType: "gauge",
			Value: float64(memStats.Alloc),
		},
		{
			ID:    "BuckHashSys",
			MType: "gauge",
			Value: float64(memStats.BuckHashSys),
		},
		{
			ID:    "Frees",
			MType: "gauge",
			Value: float64(memStats.Frees),
		},
		{
			ID:    "GCCPUFraction",
			MType: "gauge",
			Value: memStats.GCCPUFraction,
		},
		{
			ID:    "GCSys",
			MType: "gauge",
			Value: float64(memStats.GCSys),
		},
		{
			ID:    "HeapAlloc",
			MType: "gauge",
			Value: float64(memStats.HeapAlloc),
		},
		{
			ID:    "HeapIdle",
			MType: "gauge",
			Value: float64(memStats.HeapIdle),
		},
		{
			ID:    "HeapInuse",
			MType: "gauge",
			Value: float64(memStats.HeapInuse),
		},
		{
			ID:    "HeapObjects",
			MType: "gauge",
			Value: float64(memStats.HeapObjects),
		},
		{
			ID:    "HeapReleased",
			MType: "gauge",
			Value: float64(memStats.HeapReleased),
		},
		{
			ID:    "HeapSys",
			MType: "gauge",
			Value: float64(memStats.HeapSys),
		},
		{
			ID:    "LastGC",
			MType: "gauge",
			Value: float64(memStats.LastGC),
		},
		{
			ID:    "Lookups",
			MType: "gauge",
			Value: float64(memStats.Lookups),
		},
		{
			ID:    "MCacheInuse",
			MType: "gauge",
			Value: float64(memStats.MCacheInuse),
		},
		{
			ID:    "MCacheSys",
			MType: "gauge",
			Value: float64(memStats.MCacheSys),
		},
		{
			ID:    "MSpanInuse",
			MType: "gauge",
			Value: float64(memStats.MSpanInuse),
		},
		{
			ID:    "MSpanSys",
			MType: "gauge",
			Value: float64(memStats.MSpanSys),
		},
		{
			ID:    "Mallocs",
			MType: "gauge",
			Value: float64(memStats.Mallocs),
		},
		{
			ID:    "NextGfloat64C",
			MType: "gauge",
			Value: float64(memStats.NextGC),
		},
		{
			ID:    "NumForcedGC",
			MType: "gauge",
			Value: float64(memStats.NumForcedGC),
		},
		{
			ID:    "NumGC",
			MType: "gauge",
			Value: float64(uint64(memStats.NumGC)),
		},
		{
			ID:    "OtherSys",
			MType: "gauge",
			Value: float64(memStats.OtherSys),
		},
		{
			ID:    "PauseTotalNs",
			MType: "gauge",
			Value: float64(memStats.PauseTotalNs),
		},
		{
			ID:    "StackInuse",
			MType: "gauge",
			Value: float64(memStats.StackInuse),
		},
		{
			ID:    "StackSys",
			MType: "gauge",
			Value: float64(memStats.StackSys),
		},
		{
			ID:    "Sys",
			MType: "gauge",
			Value: float64(memStats.Sys),
		},
		{
			ID:    "TotalAlloc",
			MType: "gauge",
			Value: float64(memStats.TotalAlloc),
		},
		{
			ID:    "PollCount",
			MType: "counter",
			Delta: 1,
		},
		{
			ID:    "RandomValue",
			MType: "gauge",
			Value: rand.Float64(),
		},
	}
}
