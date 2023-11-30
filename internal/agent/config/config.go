package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
	ServerAddress  string
	LogLevel       string
	LogOutputPath  string
	LogErrorPath   string
	HashKey        string
	RateLimit      int
}

func New() Config {
	var (
		rawPollInterval, rawReportInterval, rawRateLimit                          int
		rawServerAddress, rawlogLevel, rawlogOutputPath, rawlogErrortPath, rawKey string
	)
	flag.StringVar(&rawServerAddress, "a", "localhost:8080", "address and port to connect")
	flag.IntVar(&rawPollInterval, "r", 2, "poll interval")
	flag.IntVar(&rawReportInterval, "p", 10, "report interval")
	flag.IntVar(&rawRateLimit, "l", 1, "Limit http requests")
	flag.StringVar(&rawlogLevel, "ll", "info", "level of logging")
	flag.StringVar(&rawlogOutputPath, "lp", "stdout", "log output path")
	flag.StringVar(&rawlogErrortPath, "le", "stderr", "log error output path")
	flag.StringVar(&rawKey, "k", "", "Key for hash summ")

	flag.Parse()

	if envAddrs, ok := os.LookupEnv("ADDRESS"); ok {
		rawServerAddress = envAddrs
	}

	if envPollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		rawPollInterval64, err := strconv.ParseInt(envPollInterval, 10, 64)
		if err != nil {
			log.Fatal("wrong report interval")
		}
		rawPollInterval = int(rawPollInterval64)

	}

	if envReportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		rawReportInterval64, err := strconv.ParseInt(envReportInterval, 10, 64)
		if err != nil {
			log.Fatal("wrong report interval")
		}
		rawReportInterval = int(rawReportInterval64)

	}

	if envLogLevel, ok := os.LookupEnv("LOGLEVEL"); ok {

		rawlogLevel = envLogLevel
	}
	if envLogOutputPath, ok := os.LookupEnv("LOGPATH"); ok {

		rawlogOutputPath = envLogOutputPath
	}
	if envLogErrorPath, ok := os.LookupEnv("LOGERROR"); ok {

		rawlogErrortPath = envLogErrorPath
	}
	if envHashKey, ok := os.LookupEnv("KEY"); ok {

		rawKey = envHashKey
	}
	if envRateLimit, ok := os.LookupEnv("RATE_LIMIT"); ok {
		rawRateLimit, err := strconv.ParseInt(envRateLimit, 10, 64)
		if err != nil {
			log.Fatal("wrong report interval")
		}
		rawReportInterval = int(rawRateLimit)
	}

	return Config{
		ServerAddress:  rawServerAddress,
		PollInterval:   time.Duration(rawPollInterval) * time.Second,
		ReportInterval: time.Duration(rawReportInterval) * time.Second,
		LogLevel:       rawlogLevel,
		LogOutputPath:  rawlogOutputPath,
		LogErrorPath:   rawlogErrortPath,
		HashKey:        rawKey,
		RateLimit:      rawRateLimit,
	}
}
