package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mr-Punder/go-alerting-service/internal/handlers"
	logger "github.com/Mr-Punder/go-alerting-service/internal/logger/zap"
	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
	"github.com/Mr-Punder/go-alerting-service/internal/middleware"
	"github.com/Mr-Punder/go-alerting-service/internal/server/config"
	"github.com/Mr-Punder/go-alerting-service/internal/storage"
)

func main() {
	conf := config.GetConfig()
	if err := run(conf); err != nil {
		panic(err)
	}
}

func run(conf *config.Config) error {

	//zapLogger, err := logger.NewLogZap(conf.LogLevel, conf.LogOutputPath, conf.LogErrorPath)
	Log, err := logger.NewLogZap(conf.LogLevel, conf.LogOutputPath)
	if err != nil {
		return err
	}

	Log.Info("Initialized logger")
	met := make(map[string]metrics.Metrics, 0)

	if conf.Restore && conf.FileStoragePath != "" {
		data, err := os.ReadFile(conf.FileStoragePath)
		if err != nil {
			if os.IsNotExist(err) {
				Log.Errorf("File %s does not exist", conf.FileStoragePath)
			} else {
				Log.Errorf("Cann't open file %s", conf.FileStoragePath)
			}
		} else {
			if len(data) == 0 {
				Log.Errorf("file %s is empty", conf.FileStoragePath)
			} else {
				err := json.Unmarshal(data, &met)
				if err != nil {
					Log.Error("json decoding error")
					return err
				}
				Log.Infof("Metrics restored from file %s", conf.FileStoragePath)
			}
		}
	}

	syncSave := conf.StoreInterval <= 0 && conf.FileStoragePath != ""

	stor, err := storage.NewMemStorage(met, syncSave, conf.FileStoragePath, Log)
	if err != nil {
		Log.Error("Cann't create storage")
		return err
	}
	defer stor.Close()
	Log.Info("Storage created")

	if conf.StoreInterval > 0 && conf.FileStoragePath != "" {
		go func() {
			for range time.Tick(time.Duration(conf.StoreInterval) * time.Second) {
				err := stor.Save()
				if err != nil {
					log.Printf("Ошибка сохранения метрик на диск: %v", err)
				}
			}
		}()
		Log.Info("Started metric saving goroutine")
	}

	comp := middleware.NewGzipCompressor(Log)
	Log.Info("Initialized compressor")

	hLogger := middleware.NewHTTPLoger(Log)
	Log.Info("Initialized middleware functions")

	server := &http.Server{
		Addr:    conf.FlagRunAddr,
		Handler: hLogger.HTTPLogHandler(comp.CompressHandler(handlers.NewMetricRouter(stor, Log))),
	}

	go func() {
		Log.Info(fmt.Sprintf("Starting server on %s", conf.FlagRunAddr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Log.Errorf("starting server on %s", conf.FlagRunAddr)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	Log.Info("Initialized shutdown")
	if err := server.Shutdown(context.Background()); err != nil {
		Log.Errorf("Cann't stop server %s", err)
	}

	Log.Info("Server closed")

	err = stor.Save()
	if err != nil {
		Log.Error("Cann't Save metrics")
	}
	Log.Info("Metrics Saved")

	return nil

}
