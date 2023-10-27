package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
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
	Info(mes string)
	Error(mes string)
}

type metricsAllGetter interface {
	GetAll() map[string]metrics.Metrics
}

type metricsGetter interface {
	Get(metric metrics.Metrics) (metrics.Metrics, bool)
}

type metricsSetter interface {
	Set(metric metrics.Metrics) error
}

type metricsDeleter interface {
	DeleteGouge(metric metrics.Metrics)
}

// Memstorer is a general metrics storage interface
type metricsStorer interface {
	metricsDeleter
	metricsGetter
	metricsSetter
	metricsAllGetter
}

// Handler type contains MemStorer and HttpLogger
type Handler struct {
	stor   metricsStorer
	logger httpLogger
}

func NewHandler(stor metricsStorer, logger httpLogger) *Handler {
	return &Handler{stor, logger}
}

func MetricRouter(storage metricsStorer, logger httpLogger) chi.Router {
	r := chi.NewRouter()

	handler := NewHandler(storage, logger)

	return r.Route("/", func(r chi.Router) {
		r.Get("/", handler.ShowAllHandler)
		r.Route("/update", func(r chi.Router) {
			r.Post("/", handler.JSONUpdHandler)
			r.Post("/{type}/{name}/{value}", handler.UpdHandler)
		})
		r.Route("/value", func(r chi.Router) {
			r.Post("/", handler.JSONValueHandler)
			r.Get("/{type}/{name}", handler.ValueHandler)
		})
		r.Get("/favicon.ico", handler.FaviconHandler)
		r.Get("/{}", handler.DefoultHandler)
		r.Post("/{}", handler.DefoultHandler)
	})
}

// JSONUpdHandler updates metric via json POST request
func (h *Handler) JSONUpdHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Entered JSONUpdHandler")
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
		h.logger.Error("wrong request method")
		http.Error(w, "Only POST requests are allowed for update!", http.StatusMethodNotAllowed)
		rData.size = 0
		rData.status = http.StatusBadRequest
		return
	}
	h.logger.Info("Method checked")

	metric := metrics.Metrics{}

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		h.logger.Error(fmt.Sprintf("json decoding error %e", err))
		http.Error(w, "wrong requests", http.StatusBadRequest)

		rData.size = 0
		rData.status = http.StatusBadRequest
		return
	}
	h.logger.Info(fmt.Sprintf("Decoded metric to stor %v", metric))
	if metric.MType != "gauge" && metric.MType != "counter" {
		h.logger.Error(fmt.Sprintf("wrong type %s", metric.MType))

		http.Error(w, "wrong type", http.StatusBadRequest)

		rData.size = 0
		rData.status = http.StatusBadRequest
		return
	}

	str := "Metrics on server: "
	for key := range h.stor.GetAll() {
		str += key + ", "
	}
	h.logger.Info(str)

	if err := h.stor.Set(metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error("Cann't find metric")
		rData.size = 0
		rData.status = http.StatusBadRequest
		return
	}

	h.logger.Info(fmt.Sprintf("Metric %s stored", metric.ID))

	respMetric, _ := h.stor.Get(metric)
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(respMetric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("json marhsaling error")
		rData.size = 0
		rData.status = http.StatusInternalServerError
		return
	}

	rData.size, _ = w.Write(body)
	rData.status = http.StatusOK

}

// JSONValueHandler returns metric via json POST request
func (h *Handler) JSONValueHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Entered JSONValueHandler")

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
		h.logger.Error("wrong request method")
		rData.size = 0
		rData.status = http.StatusBadRequest
		return
	}

	h.logger.Info("Method checked")

	metric := metrics.Metrics{}

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		h.logger.Error("json decoding error")
		http.Error(w, "wrong requests", http.StatusBadRequest)
		rData.size = 0
		rData.status = http.StatusBadRequest
		return
	}

	h.logger.Info(fmt.Sprintf("Decoded metric %v", metric))

	if metric.MType != "gauge" && metric.MType != "counter" {
		h.logger.Error(fmt.Sprintf("wrong type %s", metric.MType))
		http.Error(w, "wrong type", http.StatusBadRequest)
		rData.size = 0
		rData.status = http.StatusBadRequest
		return
	}

	respMetric, ok := h.stor.Get(metric)
	if !ok {
		str := "Metrics on server: "
		for key := range h.stor.GetAll() {
			str += key + ", "
		}
		h.logger.Info(str)
		h.logger.Error("Cann't find metric")
		http.Error(w, fmt.Sprintf("%s not found", metric.ID), http.StatusNotFound)
		rData.size = 0
		rData.status = http.StatusNotFound
		return
		// respMetric = metrics.Metrics{
		// 	ID:    metric.ID,
		// 	MType: metric.MType,
		// 	Delta: func() *int64 { var d int64 = 1; return &d }(),
		// 	Value: func() *float64 { var v = 1.5; return &v }(),
		// }
		// h.stor.Set(respMetric)
	}
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(respMetric)
	if err != nil {
		h.logger.Error("json marhsaling error")
		w.WriteHeader(http.StatusInternalServerError)
		rData.size = 0
		rData.status = http.StatusInternalServerError
		return
	}

	rData.size, _ = w.Write(body)
	rData.status = http.StatusOK
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

		metric := metrics.Metrics{
			ID:    name,
			MType: tp,
			Value: &fval,
		}

		if err := h.stor.Set(metric); err != nil {
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

		metric := metrics.Metrics{
			ID:    name,
			MType: tp,
			Delta: &ival,
		}

		if err := h.stor.Set(metric); err != nil {
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
	metric := metrics.Metrics{
		ID:    name,
		MType: tp,
	}
	w.Header().Set("Content-Type", "text/plain")
	switch tp {
	case "gauge":
		val, ok := h.stor.Get(metric)
		if !ok {
			http.Error(w, fmt.Sprintf("%s not found", name), http.StatusNotFound)
			rData.size = 0
			rData.status = http.StatusNotFound
			return
		}
		rData.size, _ = w.Write([]byte(strconv.FormatFloat(*val.Value, 'f', -1, 64)))
		rData.status = http.StatusOK
	case "counter":
		val, ok := h.stor.Get(metric)
		if !ok {
			http.Error(w, fmt.Sprintf("%s not found", name), http.StatusNotFound)
			rData.size = 0
			rData.status = http.StatusNotFound
			return
		}
		rData.size, _ = w.Write([]byte(strconv.FormatInt(*val.Delta, 10)))
		rData.status = http.StatusOK
	default:
		http.Error(w, "wrong type", http.StatusBadRequest)
		rData.size = 0
		rData.status = http.StatusBadRequest
	}

}

// ShowAllHandler returns html with all known metrics
func (h *Handler) ShowAllHandler(w http.ResponseWriter, r *http.Request) {

	h.logger.Info("entered ShowAllHandler")
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
	h.logger.Info("method checked")

	gaugeMetrics := []string{}
	counterMetrics := []string{}
	for key, val := range h.stor.GetAll() {
		if val.MType == "gauge" {
			var value = 0.0
			if val.Value != nil {
				value = *val.Value
			}
			h.logger.Info(fmt.Sprintf("trying to add gauge %v", val))
			gaugeMetrics = append(gaugeMetrics, fmt.Sprintf("<p>%s: %f</p>", key, value))
		} else if val.MType == "counter" {
			h.logger.Info(fmt.Sprintf("trying to add counter %v", val))
			var value int64 = 0
			if val.Delta != nil {
				value = *val.Delta
			}
			counterMetrics = append(counterMetrics, fmt.Sprintf("<p>%s: %d</p>", key, value))
		}
	}
	h.logger.Info("Metrics collected")

	html := "<html><body>"

	html += "<h2>Gauge:</h2>"
	sort.Strings(gaugeMetrics)
	for _, str := range gaugeMetrics {
		html += str
	}
	html += "<h2>Counter:</h2>"
	sort.Strings(counterMetrics)
	for _, str := range counterMetrics {
		html += str
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
