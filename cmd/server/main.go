package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mr-Punder/go-alerting-service/internal/handlers"
	"github.com/Mr-Punder/go-alerting-service/internal/logger"
	"github.com/Mr-Punder/go-alerting-service/internal/metricserver"
	"github.com/Mr-Punder/go-alerting-service/internal/middleware"
	"github.com/Mr-Punder/go-alerting-service/internal/server/config"
	"github.com/Mr-Punder/go-alerting-service/internal/storage"
)

func main() {
	conf := config.New()
	log, err := logger.NewZapLogger(conf.LogLevel, conf.LogOutputPath)
	if err != nil {
		panic(err)
	}
	log.Info("Initialized logger")

	env := os.Environ()
	log.Infof("Env values: %s", env)
	log.Infof("Config parametrs: %s", *conf)

	stor, closeFunc, err := storage.NewStorage(conf, log)
	if err != nil {
		log.Errorf("ERrro creating storage %s", err)
		panic(err)
	}
	defer func() {
		closeFunc()
		log.Info("Storage is closed")
	}()

	router := handlers.NewMetricRouter(stor, log)
	mserver := metricserver.NewMetricServer(conf.FlagRunAddr, router, log)

	comp := middleware.NewGzipCompressor(log)
	log.Info("Initialized compressor")

	hashHandler := middleware.NewHashSum(conf.HashKey, log)
	log.Info("Initialized SHA256 Handler")
	mserver.AddMidleware()

	hLogger := middleware.NewHTTPLoger(log)
	log.Info("Initialized middleware functions")

	mserver.AddMidleware(hashHandler.HashSummHandler, comp.CompressHandler, hLogger.HTTPLogHandler)

	go mserver.RunServer()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Info("Initialized shutdown")
	if err := mserver.Shutdown(context.Background()); err != nil {
		log.Errorf("Cann't stop server %s", err)
	}

}
