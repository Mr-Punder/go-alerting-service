package telemetry

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"sync"

	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/sync/errgroup"

	"github.com/Mr-Punder/go-alerting-service/internal/logger"
	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
)

type Telemetry struct {
	log       logger.Logger
	address   string
	key       string
	rateLimit int
}

func NewTelemetry(adr string, key string, rateLimit int, logger logger.Logger) *Telemetry {
	return &Telemetry{
		log:       logger,
		address:   adr,
		key:       key,
		rateLimit: rateLimit,
	}

}

func (t *Telemetry) Run(pollInt, repInt time.Duration) error {

	// pollTicker := time.NewTicker(pollInt)
	// reportTicker := time.NewTicker(repInt)
	// defer pollTicker.Stop()
	// defer reportTicker.Stop()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metricschan1 := t.CollectMetrics(ctx, pollInt)
	metricschan2 := t.CollectPsutilMetrics(ctx, pollInt)
	metrChanel := t.fanIn(ctx, metricschan1, metricschan2)
	t.log.Infof("Chanel created")

	worker := func(ctx context.Context, metricChan <-chan []metrics.Metrics) error {

		for metrics := range metricChan {
			t.log.Info("I'm working")
			err := t.SendMetrics(metrics)
			if err != nil {
				t.log.Errorf("Error sending metrics")
				return err
			}
		}

		return nil
	}

	g := new(errgroup.Group)

	for w := 0; w < t.rateLimit; w++ {
		g.Go(func() error {
			err := worker(ctx, metrChanel)
			t.log.Info("Worker is done")
			if err != nil {
				cancel()
				return err
			}
			return nil
		})
		t.log.Info("worker created")
	}

	time.Sleep(time.Second * 10)
	if err := g.Wait(); err != nil {
		t.log.Errorf("Error sending metrics %s", err)
		return err
	}

	return nil
}

func (t *Telemetry) fanIn(ctx context.Context, metrChns ...chan []metrics.Metrics) chan []metrics.Metrics {
	metricCh := make(chan []metrics.Metrics, t.rateLimit)
	var wg sync.WaitGroup
	wg.Add(len(metrChns))
	for _, ch := range metrChns {
		metCh := ch
		go func() {
			defer wg.Done()

			for m := range metCh {
				select {
				case <-ctx.Done():
					return
				case metricCh <- m:
				}

			}
		}()

	}

	go func() {
		wg.Wait()
		close(metricCh)
	}()
	return metricCh
}

func (t *Telemetry) CollectPsutilMetrics(ctx context.Context, pollInt time.Duration) chan []metrics.Metrics {
	metricsch := make(chan []metrics.Metrics, t.rateLimit+1)

	pollTicker := time.NewTicker(pollInt)
	go func() {
		defer close(metricsch)
		defer pollTicker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-pollTicker.C:
				metr := []metrics.Metrics{}
				vm, err := mem.VirtualMemory()
				if err == nil {
					totalMemory := float64(vm.Total)
					freeMemory := float64(vm.Free)
					metr = append(metr, metrics.Metrics{
						ID:    "TotalMemory",
						MType: "gauge",
						Value: &totalMemory,
					})
					metr = append(metr, metrics.Metrics{
						ID:    "FreeMemory",
						MType: "gauge",
						Value: &freeMemory,
					})
				}

				cpuInfo, err := cpu.Info()
				if err == nil && len(cpuInfo) > 0 {
					cpuCount := float64(len(cpuInfo))

					metr = append(metr, metrics.Metrics{
						ID:    "CPUCount",
						MType: "gauge",
						Value: &cpuCount,
					})
				}
				metricsch <- metr

			}
		}
	}()
	return metricsch
}

// Collect returns slice of current runtime metrics
func (t *Telemetry) CollectMetrics(ctx context.Context, pollInt time.Duration) chan []metrics.Metrics {
	metricsch := make(chan []metrics.Metrics, t.rateLimit+1)

	pollTicker := time.NewTicker(pollInt)
	go func() {
		defer close(metricsch)
		defer pollTicker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-pollTicker.C:
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

				metr := []metrics.Metrics{
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
				metricsch <- metr

			}
		}
	}()
	return metricsch
}

func (t *Telemetry) SendMetrics(metr []metrics.Metrics) error {
	address := "http://" + t.address
	t.log.Info(fmt.Sprintf("sending metrics to %s", address))

	client := http.Client{}
	t.log.Info("client initialized")

	url := fmt.Sprintf("http://%s/updates/", t.address)
	body, err := json.Marshal(metr)
	if err != nil {
		return err
	}

	metricstr := string(body)
	t.log.Info(fmt.Sprintf("Metrics  encoded to %s", metricstr))
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, err = zb.Write(body)
	if err != nil {
		return err
	}
	err = zb.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return err
	}

	if t.key != "" {
		h := hmac.New(sha256.New, []byte(t.key))
		h.Write(body)

		hash := h.Sum(nil)

		hashStr := hex.EncodeToString(hash[:])
		t.log.Infof("Calculated sha256 hash: %s", hashStr)
		req.Header.Set("HashSHA256", hashStr)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Encoding", "gzip")
	resp, err := client.Do(req)
	t.log.Info(fmt.Sprintf("Send request, err : %s", err))
	if err == nil {
		defer resp.Body.Close() // statictest thinks that I have to put it exactly here
	}
	retries := 2
	for i := 1; i <= retries; i++ {

		if err != nil {

			time.Sleep(time.Duration(i*40) * time.Millisecond)
			req, err = http.NewRequest("POST", url, bytes.NewBufferString(metricstr))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Del("Accept-Encoding")
			resp, err = client.Do(req)
			t.log.Info(fmt.Sprintf("Repeated request, err: %s", err))
			if err == nil {
				defer resp.Body.Close() // statictest thinks that I have to put it exactly here
			}
			if err != nil {
				i++
			} else {
				break
			}
		}

	}

	if err != nil {
		t.log.Errorf("Sending error: %s", err)

	}
	if resp.StatusCode != http.StatusOK {
		t.log.Error(fmt.Sprintf("Unexpected code %d", resp.StatusCode))

		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	gzipEncoding := resp.Header.Get("Content-Encoding")
	var ans []byte
	if strings.Contains(gzipEncoding, "gzip") {

		zr, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
		ans, err = io.ReadAll(zr)
		if err != nil {
			return err
		}
	} else {
		ans, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}

	t.log.Info(fmt.Sprintf("recievd: %s", string(ans)))

	return nil
}
