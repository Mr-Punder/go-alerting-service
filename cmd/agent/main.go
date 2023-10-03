package main

import (
	"fmt"
	"net/http"

	"time"

	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
	"github.com/go-resty/resty/v2"
)

func main() {
	parseFlags()
	if err := run(); err != nil {
		panic(err)
	}
}

func sendMetrics(metrics []metrics.Metric, serverAddress string) error {
	init := fmt.Sprintf("http://%s/update", serverAddress)
	client := resty.New()

	for _, metric := range metrics {
		url := fmt.Sprintf("%s/%s/%s/%s", init, metric.Type, metric.Name, metric.Val)

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
