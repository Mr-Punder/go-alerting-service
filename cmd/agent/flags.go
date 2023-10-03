package main

import (
	"flag"
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
	pollInterval = time.Duration(rawPollInterval) * time.Second
	reportInterval = time.Duration(rawReportInterval) * time.Second
}
