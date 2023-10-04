package main

import (
	"flag"
	"os"
)

var flagRunAddr string

func parseFlags() {

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "addres and port to run server")
	flag.Parse()

	if envAddrs, ok := os.LookupEnv("ADDRESS"); ok {

		flagRunAddr = envAddrs
	}

}
