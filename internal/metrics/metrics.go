package metrics

import (
	"math/rand"
	"runtime"
	"strconv"
)

// Metric is a type of Go runtime parameter
type Metric struct {
	Name string
	Type string
	Val  string
}

// Collect returns slice of current runtime metrics
func Collect() []Metric {

	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	return []Metric{
		{
			Name: "Alloc",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.Alloc, 10),
		},
		{
			Name: "BuckHashSys",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.BuckHashSys, 10),
		},
		{
			Name: "Frees",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.Frees, 10),
		},
		{
			Name: "GCCPUFraction",
			Type: "gauge",
			Val:  strconv.FormatFloat(memStats.GCCPUFraction, 'f', -1, 64),
		},
		{
			Name: "GCSys",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.GCSys, 10),
		},
		{
			Name: "HeapAlloc",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.HeapAlloc, 10),
		},
		{
			Name: "HeapIdle",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.HeapIdle, 10),
		},
		{
			Name: "HeapInuse",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.HeapInuse, 10),
		},
		{
			Name: "HeapObjects",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.HeapObjects, 10),
		},
		{
			Name: "HeapReleased",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.HeapReleased, 10),
		},
		{
			Name: "HeapSys",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.HeapSys, 10),
		},
		{
			Name: "LastGC",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.LastGC, 10),
		},
		{
			Name: "Lookups",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.Lookups, 10),
		},
		{
			Name: "MCacheInuse",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.MCacheInuse, 10),
		},
		{
			Name: "MCacheSys",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.MCacheSys, 10),
		},
		{
			Name: "MSpanInuse",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.MSpanInuse, 10),
		},
		{
			Name: "MSpanSys",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.MSpanSys, 10),
		},
		{
			Name: "Mallocs",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.Mallocs, 10),
		},
		{
			Name: "NextGC",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.NextGC, 10),
		},
		{
			Name: "NumForcedGC",
			Type: "gauge",
			Val:  strconv.FormatUint(uint64(memStats.NumForcedGC), 10),
		},
		{
			Name: "NumGC",
			Type: "gauge",
			Val:  strconv.FormatUint(uint64(memStats.NumGC), 10),
		},
		{
			Name: "OtherSys",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.OtherSys, 10),
		},
		{
			Name: "PauseTotalNs",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.PauseTotalNs, 10),
		},
		{
			Name: "StackInuse",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.StackInuse, 10),
		},
		{
			Name: "StackSys",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.StackSys, 10),
		},
		{
			Name: "Sys",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.Sys, 10),
		},
		{
			Name: "TotalAlloc",
			Type: "gauge",
			Val:  strconv.FormatUint(memStats.TotalAlloc, 10),
		},
		{
			Name: "PollCount",
			Type: "counter",
			Val:  "1",
		},
		{
			Name: "RandomValue",
			Type: "gauge",
			Val:  strconv.FormatFloat(rand.Float64(), 'f', -1, 64),
		},
	}
}
