package main

import (
	cnf "github.com/Mr-Punder/go-alerting-service/internal/agent/config"
	"github.com/Mr-Punder/go-alerting-service/internal/agent/config/telemetry"
	zaplogger "github.com/Mr-Punder/go-alerting-service/internal/logger/zap"
)

func main() {
	config := cnf.New()

	if err := run(config); err != nil {
		panic(err)
	}
}

func run(config cnf.Config) error {
	//zapLogger, err := logger.NewLogZap("info", "stdout", "stderr")
	Log, err := zaplogger.New(config.LogLevel, config.LogOutputPath)
	if err != nil {
		return err
	}
	Log.Info("agent started")

	tel := telemetry.NewTelemetry(config.ServerAddress, nil, Log)
	tel.CollectMetrics()
	Log.Info("metrics collected")

	err = tel.Run(config.PollInterval, config.ReportInterval)
	return err

}
