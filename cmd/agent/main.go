package main

import (
	"fmt"
	"net/http"

	"time"

	"github.com/Mr-Punder/go-alerting-service/internal/config"
	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
	"github.com/go-resty/resty/v2"
)

func main() {
	config.ParseFlags()
	if err := run(); err != nil {
		panic(err)
	}
}

func sendMetrics(metrics []metrics.Metric, addres string) error {
	init := fmt.Sprintf("%s/update", addres)
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
	pollTicker := time.NewTicker(config.PollInterval)
	reportTicker := time.NewTicker(config.ReportInterval)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	metric := metrics.Collect()
	for {
		select {
		case <-pollTicker.C:
			metric = metrics.Collect()

		case <-reportTicker.C:
			address := "http://" + config.ServerAddress
			err := sendMetrics(metric, address)
			if err != nil {
				return err
			}
		}
	}

}
