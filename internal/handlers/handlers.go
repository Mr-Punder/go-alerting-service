package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/Mr-Punder/go-alerting-service/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	stor storage.MemStor
}

func NewHandler(stor storage.MemStor) *Handler {
	return &Handler{stor}
}

func MetricRouter(storage storage.MemStor) chi.Router {
	r := chi.NewRouter()

	handler := NewHandler(storage)

	return r.Route("/", func(r chi.Router) {
		r.Get("/", handler.ShowAllHandler)
		r.Route("/update", func(r chi.Router) {
			r.Post("/{type}/{name}/{value}", handler.UpdHandler)
		})
		r.Get("/value/{type}/{name}", handler.ValueHandler)
		r.Get("/{}", handler.DefoultHandler)
		r.Post("/{}", handler.DefoultHandler)
	})
}

func (h *Handler) UpdHandler(w http.ResponseWriter, r *http.Request) {

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

		if err := h.stor.SetGouge(name, fval); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case "counter":
		ival, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			http.Error(w, "wrong format value", http.StatusBadRequest)
			return
		}

		if err := h.stor.SetCounter(name, ival); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "wrong type", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (h *Handler) DefoultHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "wrong requests", http.StatusBadRequest)
}

func (h *Handler) ValueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get requests are allowed for value!", http.StatusMethodNotAllowed)
		return
	}

	tp := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	w.Header().Set("Content-Type", "text/plain")
	switch tp {
	case "gauge":
		val, ok := h.stor.GetGouge(name)
		if !ok {
			http.Error(w, fmt.Sprintf("%s not found", name), http.StatusNotFound)
			return
		}
		w.Write([]byte(strconv.FormatFloat(val, 'f', -1, 64)))
	case "counter":
		val, ok := h.stor.GetCounter(name)
		if !ok {
			http.Error(w, fmt.Sprintf("%s not found", name), http.StatusNotFound)
			return
		}
		w.Write([]byte(strconv.FormatInt(val, 10)))
	default:
		http.Error(w, "wrong type", http.StatusBadRequest)
	}

}

func (h *Handler) ShowAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get requests are allowed for value!", http.StatusMethodNotAllowed)
		return
	}

	html := "<html><body>"

	html += "<h2>Gauge:</h2>"
	metrics := []string{}
	for key, val := range h.stor.GetAllGauge() {
		metrics = append(metrics, fmt.Sprintf("<p>%s: %f</p>", key, val))
	}
	sort.Strings(metrics)
	for _, str := range metrics {
		html += str
	}
	html += "<h2>Counter:</h2>"
	for key, val := range h.stor.GetAllCounter() {
		html += fmt.Sprintf("<p>%s: %d</p>", key, val)
	}
	html += "</body></html>"
	w.Header().Set("Content-Type", "text/html")

	w.Write([]byte(html))

}
