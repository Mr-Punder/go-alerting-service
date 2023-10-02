package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mr-Punder/go-alerting-service/internal/storage"
	"github.com/go-chi/chi/v5"
)

func MetricRouter(storage storage.MemStor) chi.Router {
	r := chi.NewRouter()

	return r.Route("/", func(r chi.Router) {
		r.Get("/", ShowAllHandler(storage))
		r.Route("/update", func(r chi.Router) {
			r.Post("/{type}/{name}/{value}", UpdHandler(storage))
		})
		r.Get("/value/{type}/{name}", ValueHandler(storage))
		r.Get("/{}", DefoultHandler)
		r.Post("/{}", DefoultHandler)
	})
}

func UpdHandler(stor storage.MemStor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed for update!", http.StatusMethodNotAllowed)
			return
		}

		tp := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")
		val := chi.URLParam(r, "value")

		switch tp {
		case "gauge":
			fval, err := strconv.ParseFloat(val, 64)
			if err != nil {
				http.Error(w, "wrong format value", http.StatusBadRequest)
				return
			}

			if err := stor.SetGouge(name, fval); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		case "counter":
			ival, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				http.Error(w, "wrong format value", http.StatusBadRequest)
				return
			}

			if err := stor.SetCounter(name, ival); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "wrong type", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

	}

}

func DefoultHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "wrong requests", http.StatusBadRequest)
}

func ValueHandler(stor storage.MemStor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only Get requests are allowed for value!", http.StatusMethodNotAllowed)
			return
		}

		tp := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")
		w.Header().Set("Content-Type", "text/plain")
		switch tp {
		case "gauge":
			val, ok := stor.GetGouge(name)
			if !ok {
				http.Error(w, fmt.Sprintf("%s not found", name), http.StatusNotFound)
				return
			}
			w.Write([]byte(strconv.FormatFloat(val, 'f', -1, 64)))
		case "counter":
			val, ok := stor.GetCounter(name)
			if !ok {
				http.Error(w, fmt.Sprintf("%s not found", name), http.StatusNotFound)
				return
			}
			w.Write([]byte(strconv.FormatInt(val, 10)))
		default:
			http.Error(w, "wrong type", http.StatusBadRequest)
		}

	}
}

func ShowAllHandler(stor storage.MemStor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only Get requests are allowed for value!", http.StatusMethodNotAllowed)
			return
		}

		html := "<html><body>"

		html += "<h2>Gauge:</h2>"
		for key, val := range stor.GetAllGauge() {
			html += fmt.Sprintf("<p>%s: %f</p>", key, val)
		}
		html += "<h2>Counter:</h2>"
		for key, val := range stor.GetAllCounter() {
			html += fmt.Sprintf("<p>%s: %d</p>", key, val)
		}
		html += "</body></html>"
		w.Header().Set("Content-Type", "text/html")

		w.Write([]byte(html))

	}
}
