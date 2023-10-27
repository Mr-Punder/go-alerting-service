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

type simpleLogger interface {
	Info(mes string)
	Error(mes string)
}

func sendMetrics(metrics []metrics.Metrics, addres string, logger simpleLogger) error {
	init := fmt.Sprintf("%s/update", addres)
	client := resty.New()
	logger.Info("client initialized")

	for _, metric := range metrics {
		url := init
		body, err := json.Marshal(metric)
		if err != nil {
			panic(err)
		}
		if metric.MType == "gauge" {
			logger.Info(fmt.Sprintf("metric to encode %s %s %f", metric.ID, metric.MType, *metric.Value))
		} else {
			logger.Info(fmt.Sprintf("metric to encode %s %s %d", metric.ID, metric.MType, *metric.Delta))

		}
		logger.Info(fmt.Sprintf("Metric %s encoded to %s", metric.ID, body))

		resp, err := client.R().SetHeader("Content-Type", "application/json").SetBody(body).Post(url)
		logger.Info(fmt.Sprintf("Send request, err : %v", err))

		retries := 3
		for i := 0; i < retries; i++ {
			if err != nil {

				time.Sleep(40 * time.Millisecond)
				resp, err = client.R().SetHeader("Content-Type", "application/json").SetBody(body).Post(url)
				logger.Info(fmt.Sprintf("Repeated request, err: %v", err))
				if err != nil {
					i++
				} else {
					break
				}
			}
		}

		if err != nil {

			logger.Error(fmt.Sprintf("sending error %e", err))
			return err
		}
		if resp.StatusCode() != http.StatusOK {
			logger.Error(fmt.Sprintf("Unexpected code %d", resp.StatusCode()))

			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}
		ans := resp.Body()
		logger.Info(fmt.Sprintf("recievd: %s", string(ans)))

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
	//time.Sleep(11 * time.Second)
	zapLogger.Info("agent started")

	metric := Collect()
	zapLogger.Info("metrics collected")
	for {
		select {
		case <-pollTicker.C:
			metric = Collect()

		case <-reportTicker.C:
			address := "http://" + config.ServerAddress
			zapLogger.Info(fmt.Sprintf("sending metrics to %s", address))
			err := sendMetrics(metric, address, zapLogger)
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

	Alloc := float64(memStats.Alloc)
	BuckHashSys := float64(memStats.BuckHashSys)
	Frees := float64(memStats.Frees)
	GCCPUFraction := memStats.GCCPUFraction
	GCSys := float64(memStats.GCSys)
	HeapAlloc := float64(memStats.HeapAlloc)
	HeapIdle := float64(memStats.HeapIdle)
	HeapInuse := float64(memStats.HeapInuse)
	HeapObjects := float64(memStats.HeapObjects)
	HeapReleased := float64(memStats.HeapReleased)
	HeapSys := float64(memStats.HeapSys)
	LastGC := float64(memStats.LastGC)
	Lookups := float64(memStats.Lookups)
	MCacheInuse := float64(memStats.MCacheInuse)
	MCacheSys := float64(memStats.MCacheSys)
	MSpanInuse := float64(memStats.MSpanInuse)
	MSpanSys := float64(memStats.MSpanSys)
	Mallocs := float64(memStats.Mallocs)
	NextGC := float64(memStats.NextGC)
	NumForcedGC := float64(memStats.NumForcedGC)
	NumGC := float64(uint64(memStats.NumGC))
	OtherSys := float64(memStats.OtherSys)
	PauseTotalNs := float64(memStats.PauseTotalNs)
	StackInuse := float64(memStats.StackInuse)
	StackSys := float64(memStats.StackSys)
	Sys := float64(memStats.Sys)
	TotalAlloc := float64(memStats.TotalAlloc)
	var PollCount int64 = 1
	RandomValue := rand.Float64()

	return []metrics.Metrics{
		{
			ID:    "Alloc",
			MType: "gauge",
			Value: &Alloc,
		},
		{
			ID:    "BuckHashSys",
			MType: "gauge",
			Value: &BuckHashSys,
		},
		{
			ID:    "Frees",
			MType: "gauge",
			Value: &Frees,
		},
		{
			ID:    "GCCPUFraction",
			MType: "gauge",
			Value: &GCCPUFraction,
		},
		{
			ID:    "GCSys",
			MType: "gauge",
			Value: &GCSys,
		},
		{
			ID:    "HeapAlloc",
			MType: "gauge",
			Value: &HeapAlloc,
		},
		{
			ID:    "HeapIdle",
			MType: "gauge",
			Value: &HeapIdle,
		},
		{
			ID:    "HeapInuse",
			MType: "gauge",
			Value: &HeapInuse,
		},
		{
			ID:    "HeapObjects",
			MType: "gauge",
			Value: &HeapObjects,
		},
		{
			ID:    "HeapReleased",
			MType: "gauge",
			Value: &HeapReleased,
		},
		{
			ID:    "HeapSys",
			MType: "gauge",
			Value: &HeapSys,
		},
		{
			ID:    "LastGC",
			MType: "gauge",
			Value: &LastGC,
		},
		{
			ID:    "Lookups",
			MType: "gauge",
			Value: &Lookups,
		},
		{
			ID:    "MCacheInuse",
			MType: "gauge",
			Value: &MCacheInuse,
		},
		{
			ID:    "MCacheSys",
			MType: "gauge",
			Value: &MCacheSys,
		},
		{
			ID:    "MSpanInuse",
			MType: "gauge",
			Value: &MSpanInuse,
		},
		{
			ID:    "MSpanSys",
			MType: "gauge",
			Value: &MSpanSys,
		},
		{
			ID:    "Mallocs",
			MType: "gauge",
			Value: &Mallocs,
		},
		{
			ID:    "NextGC",
			MType: "gauge",
			Value: &NextGC,
		},
		{
			ID:    "NumForcedGC",
			MType: "gauge",
			Value: &NumForcedGC,
		},
		{
			ID:    "NumGC",
			MType: "gauge",
			Value: &NumGC,
		},
		{
			ID:    "OtherSys",
			MType: "gauge",
			Value: &OtherSys,
		},
		{
			ID:    "PauseTotalNs",
			MType: "gauge",
			Value: &PauseTotalNs,
		},
		{
			ID:    "StackInuse",
			MType: "gauge",
			Value: &StackInuse,
		},
		{
			ID:    "StackSys",
			MType: "gauge",
			Value: &StackSys,
		},
		{
			ID:    "Sys",
			MType: "gauge",
			Value: &Sys,
		},
		{
			ID:    "TotalAlloc",
			MType: "gauge",
			Value: &TotalAlloc,
		},
		{
			ID:    "PollCount",
			MType: "counter",
			Delta: &PollCount,
		},
		{
			ID:    "RandomValue",
			MType: "gauge",
			Value: &RandomValue,
		},
	}
}
