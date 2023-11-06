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
	"github.com/Mr-Punder/go-alerting-service/internal/postgre"
	"github.com/Mr-Punder/go-alerting-service/internal/server/config"
	"github.com/Mr-Punder/go-alerting-service/internal/storage"
)

func main() {
	conf := config.New()
	log, err := logger.New(conf.LogLevel, conf.LogOutputPath)
	if err != nil {
		panic(err)
	}
	log.Info("Initialized logger")

	met := make(map[string]metrics.Metrics, 0)
	if conf.Restore && conf.FileStoragePath != "" {
		if err = metricserver.RestoreMetric(conf.FileStoragePath, &met, log); err != nil {
			log.Errorf("failed to restor metrics %s", err)
		}

	}

	syncSave := conf.StoreInterval <= 0 && conf.FileStoragePath != ""
	stor, err := storage.NewMemStorage(met, syncSave, conf.FileStoragePath, log)
	if err != nil {
		log.Error("Cann't create storage")
		panic(err)
	}
	defer stor.Close()
	log.Info("Storage created")

	if conf.StoreInterval > 0 && conf.FileStoragePath != "" {
		go metricserver.SaveMetrics(stor, conf.StoreInterval, log)
		log.Info("Started metric saving goroutine")
	}

	db, err := postgre.NewPostgreDB(conf.DBstring, log)
	if err != nil {
		log.Errorf("Error opening database", err)
		panic(err)
	}
	log.Infof("Database is opened with dsn %s", conf.DBstring)
	env := os.Environ()
	log.Infof("env values: %s", env)

	log.Infof("config parametrs: %s", *conf)

	router := handlers.NewMetricRouter(stor, db, log)
	mserver := metricserver.NewMetricServer(conf.FlagRunAddr, router, log)

	comp := middleware.NewGzipCompressor(log)
	log.Info("Initialized compressor")

	hLogger := middleware.NewHTTPLoger(log)
	log.Info("Initialized middleware functions")

	mserver.AddMidleware(comp.CompressHandler, hLogger.HTTPLogHandler)

	go mserver.RunServer()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Info("Initialized shutdown")
	if err := mserver.Shutdown(context.Background()); err != nil {
		log.Errorf("Cann't stop server %s", err)
	}

}
