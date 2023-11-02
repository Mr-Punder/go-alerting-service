package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mr-Punder/go-alerting-service/internal/handlers"
	logger "github.com/Mr-Punder/go-alerting-service/internal/logger/zap"
	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
	"github.com/Mr-Punder/go-alerting-service/internal/metricserver"
	"github.com/Mr-Punder/go-alerting-service/internal/middleware"
	"github.com/Mr-Punder/go-alerting-service/internal/server/config"
	"github.com/Mr-Punder/go-alerting-service/internal/storage"
)

func main() {
	conf := config.New()
	Log, err := logger.New(conf.LogLevel, conf.LogOutputPath)
	if err != nil {
		panic(err)
	}
	Log.Info("Initialized logger")

	met := make(map[string]metrics.Metrics, 0)
	if conf.Restore && conf.FileStoragePath != "" {
		if err = metricserver.RestoreMetric(conf.FileStoragePath, &met, Log); err != nil {
			Log.Errorf("failed to restor metrics %s", err)
		}

	}

	syncSave := conf.StoreInterval <= 0 && conf.FileStoragePath != ""
	stor, err := storage.NewMemStorage(met, syncSave, conf.FileStoragePath, Log)
	if err != nil {
		Log.Error("Cann't create storage")
		panic(err)
	}
	defer stor.Close()
	Log.Info("Storage created")

	if conf.StoreInterval > 0 && conf.FileStoragePath != "" {
		go metricserver.SaveMetrics(stor, conf.StoreInterval, Log)
		Log.Info("Started metric saving goroutine")
	}

	router := handlers.NewMetricRouter(stor, Log)
	mserver := metricserver.NewMetricServer(conf.FlagRunAddr, router, Log)

	comp := middleware.NewGzipCompressor(Log)
	Log.Info("Initialized compressor")

	hLogger := middleware.NewHTTPLoger(Log)
	Log.Info("Initialized middleware functions")

	mserver.AddMidleware(comp.CompressHandler, hLogger.HTTPLogHandler)

	go mserver.RunServer()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	Log.Info("Initialized shutdown")
	if err := mserver.Shutdown(context.Background()); err != nil {
		Log.Errorf("Cann't stop server %s", err)
	}

}
