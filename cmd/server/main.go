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

	storage := &storage.MemStorage{
		GaugeStorage:   make(map[string]float64),
		CounterStorage: make(map[string]int64)}

	mux := http.NewServeMux()

	mux.HandleFunc(`/update/counter/`, handlers.CounterUpd(storage))
	mux.HandleFunc(`/update/gauge/`, handlers.GaugeUpd(storage))
	mux.HandleFunc(`/`, handlers.DefoultHandler)

	return http.ListenAndServe(`:8080`, mux)

}
