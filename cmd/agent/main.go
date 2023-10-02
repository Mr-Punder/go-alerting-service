package main

import (
	"fmt"
	"net/http"

	"time"

	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
	"github.com/go-resty/resty/v2"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	serverAddress  = "http://localhost:8080"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func sendMetrics(metrics []metrics.Metric, serverAddress string) error {
	initUrl := fmt.Sprintf("%s/update", serverAddress)
	client := resty.New()

	for _, metric := range metrics {
		url := fmt.Sprintf("%s/%s/%s/%s", initUrl, metric.Type, metric.Name, metric.Val)

		resp, err := client.R().SetHeader("Content-Type", "text/plain").Post(url)
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
	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	metric := metrics.Collect()
	for {
		select {
		case <-pollTicker.C:
			metric = metrics.Collect()

		case <-reportTicker.C:

			err := sendMetrics(metric, serverAddress)
			if err != nil {
				return err
			}
		}
	}

}
