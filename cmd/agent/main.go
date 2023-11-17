package main

import (
	cnf "github.com/Mr-Punder/go-alerting-service/internal/agent/config"
	"github.com/Mr-Punder/go-alerting-service/internal/logger"
	"github.com/Mr-Punder/go-alerting-service/internal/telemetry"
)

func main() {
	config := cnf.New()

	Log, err := logger.NewZapLogger(config.LogLevel, config.LogOutputPath)
	if err != nil {
		panic(err)
	}
	Log.Info("agent started")

	tel := telemetry.NewTelemetry(config.ServerAddress, nil, Log)
	tel.CollectMetrics()
	Log.Info("metrics collected")

	err = tel.Run(config.PollInterval, config.ReportInterval)
	if err != nil {
		panic(err)
	}
}
