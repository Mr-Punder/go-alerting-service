package config

import (
	"flag"
	"os"
)

type Config struct {
	FlagRunAddr   string
	LogLevel      string
	LogOutputPath string
	LogErrorPath  string
}

// GetConfig from environment and consol parameters
func GetConfig() *Config {
	var flagRunAddr string
	var logLevel string
	var logOutputPath string
	var logErrortPath string

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "addres and port to run server")
	flag.StringVar(&logLevel, "l", "info", "level of logging")
	flag.StringVar(&logOutputPath, "lp", "stdout", "log output path")
	flag.StringVar(&logErrortPath, "le", "stderr", "log error output path")

	flag.Parse()

	if envAddrs, ok := os.LookupEnv("ADDRESS"); ok {

		flagRunAddr = envAddrs
	}
	if envLogLevel, ok := os.LookupEnv("LOGLEVEL"); ok {

		logLevel = envLogLevel
	}
	if envLogOutputPath, ok := os.LookupEnv("LOGPATH"); ok {

		logOutputPath = envLogOutputPath
	}
	if envLogErrorPath, ok := os.LookupEnv("LOGERROR"); ok {

		logErrortPath = envLogErrorPath
	}

	return &Config{flagRunAddr, logLevel, logOutputPath, logErrortPath}
}
