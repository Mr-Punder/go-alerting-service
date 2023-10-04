package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	pollInterval   time.Duration
	reportInterval time.Duration
	serverAddress  string
)

func parseFlags() {
	var rawPollInterval, rawReportInterval int
	flag.StringVar(&serverAddress, "a", "localhost:8080", "addres and port to connect")
	flag.IntVar(&rawPollInterval, "r", 2, "poll interval")
	flag.IntVar(&rawReportInterval, "p", 10, "report interval")
	flag.Parse()

	if envAddrs, ok := os.LookupEnv("ADDRESS"); ok {
		serverAddress = envAddrs
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

	pollInterval = time.Duration(rawPollInterval) * time.Second
	reportInterval = time.Duration(rawReportInterval) * time.Second

}
