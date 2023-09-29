package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Mr-Punder/go-alerting-service/internal/storage"
)

func GaugeUpd(stor storage.MemStor) http.HandlerFunc {
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

		if err := stor.SetGouge(name, fval); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.WriteHeader(http.StatusOK)

	}

}

func DefoultHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}

func CounterUpd(stor storage.MemStor) http.HandlerFunc {
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

		if err := stor.SetCounter(name, ival); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.WriteHeader(http.StatusOK)

	}
}
