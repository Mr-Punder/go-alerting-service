package main

import (
	"fmt"
	"net/http"

	"github.com/Mr-Punder/go-alerting-service/internal/handlers"
	"github.com/Mr-Punder/go-alerting-service/internal/logger"
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
	ruslog, err := logger.NewLogLogrus("info", "stdout")

	if err != nil {
		return err
	}

	storage := new(storage.MemStorage)

	ruslog.Info("Initialized logger and storage")
	// gzipw := gzipcomp.NewEmptyGzipCompressWriter()
	// gzipr := gzipcomp.NewEmptyGzipCompressReader()
	// comp := middleware.NewCompressor(gzipw, gzipr, zapLogger)

	ruslog.Info("Initialized compressor")

	hLogger := middleware.NewHTTPLoger(ruslog)
	comp := middleware.NewGzipCompressor(ruslog)

	ruslog.Info("Initialized middleware functions")

	ruslog.Info(fmt.Sprintf("Starting server on %s", conf.FlagRunAddr))
	return http.ListenAndServe(conf.FlagRunAddr, hLogger.HTTPLogHandler(comp.CompressHandler(handlers.NewMetricRouter(storage, ruslog))))

}
