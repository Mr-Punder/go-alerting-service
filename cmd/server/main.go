package main

import (
	"net/http"

	"github.com/Mr-Punder/go-alerting-service/internal/handlers"
	"github.com/Mr-Punder/go-alerting-service/internal/storage"
)

func main() {
	parseFlags()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	storage := new(storage.MemStorage)

	return http.ListenAndServe(flagRunAddr, handlers.MetricRouter(storage))

}
