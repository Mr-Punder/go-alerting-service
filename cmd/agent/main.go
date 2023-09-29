package main

import (
	"fmt"
	"net/http"

	"time"

	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	serverAddress  = "http://localhost:8080"
)

func sendMetrics(metrics []metrics.Metric, serverAddress string) error {
	url := fmt.Sprintf("%s/update", serverAddress)
	client := &http.Client{}

	for _, metric := range metrics {
		url = fmt.Sprintf("%s/%s/%s/%s", url, metric.Type, metric.Name, metric.Val)

		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "text/plain")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
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
