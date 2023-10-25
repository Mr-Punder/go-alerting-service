package handlers

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type responseData struct {
	status int
	size   int
}

// httpLogger logs information about http requests and responses
type httpLogger interface {
	RequestLog(method string, path string, duration time.Duration)
	ResponseLog(status int, size int)
}

type memAllGetter interface {
	GetAllGauge() map[string]float64
	GetAllCounter() map[string]int64
}

type memGetter interface {
	GetGouge(string) (float64, bool)
	GetCounter(string) (int64, bool)
}

type memSetter interface {
	SetGouge(string, float64) error
	SetCounter(string, int64) error
}

type memDeleter interface {
	DeleteGouge(string)
	DeleteCouner(string)
}

// Memstorer is a general metrics storage
type memStorer interface {
	memDeleter
	memGetter
	memSetter
	memAllGetter
}

// Handler type contains MemStorer and HttpLogger
type Handler struct {
	stor   memStorer
	logger httpLogger
}

func NewHandler(stor memStorer, logger httpLogger) *Handler {
	return &Handler{stor, logger}
}

func MetricRouter(storage memStorer, logger httpLogger) chi.Router {
	r := chi.NewRouter()

	handler := NewHandler(storage, logger)

	return r.Route("/", func(r chi.Router) {
		r.Get("/", handler.ShowAllHandler)
		r.Route("/update", func(r chi.Router) {
			r.Post("/{type}/{name}/{value}", handler.UpdHandler)
		})
		r.Get("/value/{type}/{name}", handler.ValueHandler)
		r.Get("/favicon.ico", handler.FaviconHandler)
		r.Get("/{}", handler.DefoultHandler)
		r.Post("/{}", handler.DefoultHandler)
	})
}

// UpdHandler updates one metric or creates new one with name
func (h *Handler) UpdHandler(w http.ResponseWriter, r *http.Request) {
	path := r.RequestURI
	method := r.Method
	start := time.Now()
	rData := responseData{}
	defer func() {
		duration := time.Since(start)

		h.logger.RequestLog(method, path, duration)
		h.logger.ResponseLog(rData.status, rData.size)
	}()

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed for update!", http.StatusMethodNotAllowed)
		rData.size = 0
		rData.status = http.StatusBadRequest
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
			rData.size = 0
			rData.status = http.StatusBadRequest
			return
		}

		if err := h.stor.SetGouge(name, fval); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			rData.size = 0
			rData.status = http.StatusBadRequest
			return
		}
	case "counter":
		ival, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			http.Error(w, "wrong format value", http.StatusBadRequest)
			rData.size = 0
			rData.status = http.StatusBadRequest
			return
		}

		if err := h.stor.SetCounter(name, ival); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			rData.size = 0
			rData.status = http.StatusBadRequest
			return
		}
	default:
		http.Error(w, "wrong type", http.StatusBadRequest)
		rData.size = 0
		rData.status = http.StatusBadRequest
		return
	}

	w.WriteHeader(http.StatusOK)
	rData.size = 0
	rData.status = http.StatusOK

}

// DefoultHandler for incorrect requests
func (h *Handler) DefoultHandler(w http.ResponseWriter, r *http.Request) {
	path := r.RequestURI
	method := r.Method
	start := time.Now()
	rData := responseData{}
	defer func() {
		duration := time.Since(start)

		h.logger.RequestLog(method, path, duration)
		h.logger.ResponseLog(rData.status, rData.size)
	}()

	http.Error(w, "wrong requests", http.StatusBadRequest)
	rData.size = 0
	rData.status = http.StatusBadRequest
}

// ValueHandler returns value of metric by name if the metric exists
func (h *Handler) ValueHandler(w http.ResponseWriter, r *http.Request) {
	path := r.RequestURI
	method := r.Method
	start := time.Now()
	rData := responseData{}
	defer func() {
		duration := time.Since(start)

		h.logger.RequestLog(method, path, duration)
		h.logger.ResponseLog(rData.status, rData.size)
	}()

	if r.Method != http.MethodGet {
		http.Error(w, "Only Get requests are allowed for value!", http.StatusMethodNotAllowed)
		rData.size = 0
		rData.status = http.StatusMethodNotAllowed
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
			rData.size = 0
			rData.status = http.StatusNotFound
			return
		}
		rData.size, _ = w.Write([]byte(strconv.FormatFloat(val, 'f', -1, 64)))
		rData.status = http.StatusOK
	case "counter":
		val, ok := h.stor.GetCounter(name)
		if !ok {
			http.Error(w, fmt.Sprintf("%s not found", name), http.StatusNotFound)
			rData.size = 0
			rData.status = http.StatusNotFound
			return
		}
		rData.size, _ = w.Write([]byte(strconv.FormatInt(val, 10)))
		rData.status = http.StatusOK
	default:
		http.Error(w, "wrong type", http.StatusBadRequest)
		rData.size = 0
		rData.status = http.StatusBadRequest
	}

}

// ShowAllHandler returns html with all known metrics
func (h *Handler) ShowAllHandler(w http.ResponseWriter, r *http.Request) {
	path := r.RequestURI
	method := r.Method
	start := time.Now()
	rData := responseData{}
	defer func() {
		duration := time.Since(start)

		h.logger.RequestLog(method, path, duration)
		h.logger.ResponseLog(rData.status, rData.size)
	}()

	if r.Method != http.MethodGet {
		http.Error(w, "Only Get requests are allowed for value!", http.StatusMethodNotAllowed)
		rData.size = 0
		rData.status = http.StatusMethodNotAllowed
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

	rData.size, _ = w.Write([]byte(html))
	rData.status = http.StatusOK

}

// FaviconHandler returns Gopher!!!!
func (h *Handler) FaviconHandler(w http.ResponseWriter, r *http.Request) {
	path := r.RequestURI
	method := r.Method
	start := time.Now()
	rData := responseData{}
	defer func() {
		duration := time.Since(start)

		h.logger.RequestLog(method, path, duration)
		h.logger.ResponseLog(rData.status, rData.size)
	}()

	icon, err := os.ReadFile("../../images/gopher.png")
	if err != nil {
		http.Error(w, "Иконка не найдена", http.StatusNotFound)
		rData.size = 0
		rData.status = http.StatusNotFound
		return
	}

	w.Header().Set("Content-Type", "image/png")

	rData.size, _ = w.Write(icon)
	rData.status = http.StatusOK

}
