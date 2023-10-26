package main

import (
	"fmt"
	"net/http"

	"github.com/Mr-Punder/go-alerting-service/internal/handlers"
	"github.com/Mr-Punder/go-alerting-service/internal/logger"
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

	zapLogger, err := logger.NewLogZap(conf.LogLevel, conf.LogOutputPath, conf.LogErrorPath)
	if err != nil {
		return err
	}

	storage := new(storage.MemStorage)

	zapLogger.Info(fmt.Sprintf("Started server on %s", conf.FlagRunAddr))
	return http.ListenAndServe(conf.FlagRunAddr, handlers.MetricRouter(storage, zapLogger))

}
