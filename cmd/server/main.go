package main

import (
	"net/http"

	"github.com/Mr-Punder/go-alerting-service/internal/handlers"
	"github.com/Mr-Punder/go-alerting-service/internal/storage"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	storage := new(storage.MemStorage)

	return http.ListenAndServe(`:8080`, handlers.MetricRouter(storage))

}
