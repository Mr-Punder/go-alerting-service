package main

import (
	"net/http"
	"strconv"
	"strings"
)

type MemStor interface {
	SetGouge(string, float64) error
	SetCounter(string, int64) error
}

type MemStorage struct {
	gaugeStorage   map[string]float64
	counterStorage map[string]int64
}

func (stor *MemStorage) SetGouge(key string, val float64) error {
	stor.gaugeStorage[key] = val
	return nil
}

func (stor *MemStorage) SetCounter(key string, val int64) error {
	stor.counterStorage[key] += val
	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	storage := &MemStorage{
		gaugeStorage:   make(map[string]float64),
		counterStorage: make(map[string]int64)}

	mux := http.NewServeMux()

	mux.HandleFunc(`/update/counter/`, counterUpd(storage))
	mux.HandleFunc(`/update/gauge/`, gaugeUpd(storage))
	mux.HandleFunc(`/`, http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusBadRequest)
	}))

	return http.ListenAndServe(`:8080`, mux)

}
func gaugeUpd(storage MemStor) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodPost {
			http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		path := req.URL.Path
		params := strings.Split(path, `/`)

		if len(params) < 5 {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		name := params[3]
		val := params[4]

		fval, err := strconv.ParseFloat(val, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := storage.SetGouge(name, fval); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.WriteHeader(http.StatusOK)

	}

}

func counterUpd(storage MemStor) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		path := req.URL.Path
		params := strings.Split(path, `/`)

		if len(params) < 5 {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		name := params[3]
		val := params[4]

		ival, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := storage.SetCounter(name, ival); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.WriteHeader(http.StatusOK)

	}
}
