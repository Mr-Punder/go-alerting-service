package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	FlagRunAddr     string
	LogLevel        string
	LogOutputPath   string
	LogErrorPath    string
	StoreInterval   int64
	FileStoragePath string
	Restore         bool
	DBstring        string
}

// New from environment and consol parameters
func New() *Config {
	var (
		flagRunAddr, logLevel, logOutputPath, fileStoragePath, logErrortPath, dbString string
		storeInterval                                                                  int64
		restore                                                                        bool
	)

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "addres and port to run server")
	flag.StringVar(&logLevel, "l", "info", "level of logging")
	flag.StringVar(&logOutputPath, "lp", "stdout", "log output path")
	flag.StringVar(&logErrortPath, "le", "stderr", "log error output path")
	flag.Int64Var(&storeInterval, "i", 300, "metrics saving interval")
	flag.StringVar(&fileStoragePath, "f", "/tmp/metrics-db.json", "storage filename")
	flag.BoolVar(&restore, "r", true, "restore metrics from storage")
	flag.StringVar(&dbString, "d", "", "databese opening string")
	// host=localhost user=metrics password=metrics_password dbname=metrics

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
	if envStoreInterval, ok := os.LookupEnv("STORE_INTERVAL"); ok {

		storeInterval, _ = strconv.ParseInt(envStoreInterval, 10, 64)
	}
	if envFileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {

		fileStoragePath = envFileStoragePath
	}
	if envRestor, ok := os.LookupEnv("RESTORE"); ok {

		restore, _ = strconv.ParseBool(envRestor)
	}
	if envDBstring, ok := os.LookupEnv("DATABASE_DSN"); ok {
		dbString = envDBstring
	}

	return &Config{flagRunAddr, logLevel, logOutputPath, logErrortPath, storeInterval, fileStoragePath, restore, dbString}
}
