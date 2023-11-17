package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/Mr-Punder/go-alerting-service/internal/logger"
	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
	"github.com/Mr-Punder/go-alerting-service/internal/storage"
	"github.com/go-chi/chi/v5"
)

// Handler type contains MemStorer and HttpLogger
type Handler struct {
	stor    storage.MetricsStorer
	logger  logger.Logger
	timeout time.Duration
}

func NewHandler(stor storage.MetricsStorer, logger logger.Logger) *Handler {
	return &Handler{stor, logger, 3 * time.Second}
}

func NewMetricRouter(storage storage.MetricsStorer, logger logger.Logger) chi.Router {
	r := chi.NewRouter()

	handler := NewHandler(storage, logger)

	return r.Route("/", func(r chi.Router) {
		r.Get("/", handler.ShowAllHandler)
		r.Post("/updates/", handler.JSONUpdAllHandler)
		r.Route("/update", func(r chi.Router) {
			r.Post("/", handler.JSONUpdHandler)
			r.Post("/{type}/{name}/{value}", handler.UpdHandler)
		})
		r.Route("/value", func(r chi.Router) {
			r.Post("/", handler.JSONValueHandler)
			r.Get("/{type}/{name}", handler.ValueHandler)
		})
		r.Get("/ping", handler.PingHandler)
		r.Get("/favicon.ico", handler.FaviconHandler)
		r.Get("/{}", handler.DefoultHandler)
		r.Post("/{}", handler.DefoultHandler)
	})
}

func (h *Handler) JSONUpdAllHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Entered JSONUpdAllHandler")
	if r.Method != http.MethodPost {
		h.logger.Error("wrong request method")
		http.Error(w, "Only POST requests are allowed for update!", http.StatusMethodNotAllowed)

		return
	}
	h.logger.Info("Method checked")

	metrics := make([]metrics.Metrics, 0)
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		h.logger.Error(fmt.Sprintf("json decoding error %e", err))
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "wrong requests", http.StatusBadRequest)

		return
	}

	for _, m := range metrics {
		if m.MType != "gauge" && m.MType != "counter" {
			h.logger.Error(fmt.Sprintf("wrong type %s", m.MType))

			http.Error(w, "wrong type", http.StatusBadRequest)

			return
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()
	if err := h.stor.SetAll(ctx, metrics); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Errorf("Cann't store metrics", err)

		return
	}

	h.logger.Info("Metrics stored")
	w.WriteHeader(http.StatusOK)
	h.logger.Info("JSONUpdHandler exited")

}

func (h *Handler) PingHandler(w http.ResponseWriter, r *http.Request) {

	h.logger.Info("Entered PingHandler")
	if r.Method != http.MethodGet {
		h.logger.Error("wrong request method")
		http.Error(w, "Only GET requests are allowed for update!", http.StatusMethodNotAllowed)

		return
	}
	h.logger.Info("Method checked")

	err := h.stor.Ping()
	if err != nil {
		h.logger.Errorf("Database does not ping %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Database does not ping", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	h.logger.Info("PingHandler exited")

}

// JSONUpdHandler updates metric via json POST request
func (h *Handler) JSONUpdHandler(w http.ResponseWriter, r *http.Request) {

	h.logger.Info("Entered JSONUpdHandler")
	if r.Method != http.MethodPost {
		h.logger.Error("wrong request method")
		http.Error(w, "Only POST requests are allowed for update!", http.StatusMethodNotAllowed)

		return
	}
	h.logger.Info("Method checked")

	metric := metrics.Metrics{}

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		h.logger.Error(fmt.Sprintf("json decoding error %e", err))
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "wrong requests", http.StatusBadRequest)

		return
	}
	h.logger.Info(fmt.Sprintf("Decoded metric to stor %v", metric))
	if metric.MType != "gauge" && metric.MType != "counter" {
		h.logger.Error(fmt.Sprintf("wrong type %s", metric.MType))

		http.Error(w, "wrong type", http.StatusBadRequest)

		return
	}

	str := "Metrics on server: "

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()
	for key := range h.stor.GetAll(ctx) {
		str += key + ", "
	}
	h.logger.Info(str)

	ctx, cancel = context.WithTimeout(r.Context(), h.timeout)
	defer cancel()
	if err := h.stor.Set(ctx, metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Errorf("Cann't store metric", err)

		return
	}

	h.logger.Infof("Metric %s stored", metric.ID)

	ctx, cancel = context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	respMetric, _ := h.stor.Get(ctx, metric)
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(respMetric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("json marhsaling error")

		return
	}

	w.Write(body)
	h.logger.Info("JSONUpdHandler exited")

}

// JSONValueHandler returns metric via json POST request
func (h *Handler) JSONValueHandler(w http.ResponseWriter, r *http.Request) {

	h.logger.Info("Entered JSONValueHandler")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed for update!", http.StatusMethodNotAllowed)
		h.logger.Error("wrong request method")

		return
	}

	h.logger.Info("Method checked")

	metric := metrics.Metrics{}

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		h.logger.Error("json decoding error")
		http.Error(w, "wrong requests", http.StatusBadRequest)

		return
	}

	h.logger.Info(fmt.Sprintf("Decoded metric %v", metric))

	if metric.MType != "gauge" && metric.MType != "counter" {
		h.logger.Error(fmt.Sprintf("wrong type %s", metric.MType))
		http.Error(w, "wrong type", http.StatusBadRequest)

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	respMetric, ok := h.stor.Get(ctx, metric)
	if !ok {
		ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
		defer cancel()
		str := "Metrics on server: "
		for key := range h.stor.GetAll(ctx) {
			str += key + ", "
		}
		h.logger.Info(str)
		h.logger.Error("Cann't find metric")
		http.Error(w, fmt.Sprintf("%s not found", metric.ID), http.StatusNotFound)

		return
	}
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(respMetric)
	if err != nil {
		h.logger.Error("json marhsaling error")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Write(body)

}

// UpdHandler updates one metric or creates new one with name
func (h *Handler) UpdHandler(w http.ResponseWriter, r *http.Request) {

	h.logger.Info("Entered UpdHandler")

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

		metric := metrics.Metrics{
			ID:    name,
			MType: tp,
			Value: &fval,
		}

		ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
		defer cancel()

		if err := h.stor.Set(ctx, metric); err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

	case "counter":
		ival, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			http.Error(w, "wrong format value", http.StatusBadRequest)

			return
		}

		metric := metrics.Metrics{
			ID:    name,
			MType: tp,
			Delta: &ival,
		}

		ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
		defer cancel()
		if err := h.stor.Set(ctx, metric); err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}
	default:
		http.Error(w, "wrong type", http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusOK)

}

// DefoultHandler for incorrect requests
func (h *Handler) DefoultHandler(w http.ResponseWriter, r *http.Request) {

	h.logger.Info("Entered DefoultHandler")

	http.Error(w, "wrong requests", http.StatusBadRequest)

}

// ValueHandler returns value of metric by name if the metric exists
func (h *Handler) ValueHandler(w http.ResponseWriter, r *http.Request) {

	h.logger.Info("Entered ValueHandler")
	headers := r.Header
	h.logger.Info(fmt.Sprintf("Headers:  %v", headers))

	if r.Method != http.MethodGet {
		http.Error(w, "Only Get requests are allowed for value!", http.StatusMethodNotAllowed)

		return
	}

	tp := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	metric := metrics.Metrics{
		ID:    name,
		MType: tp,
	}
	w.Header().Set("Content-Type", "text/plain")

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()
	switch tp {
	case "gauge":
		val, ok := h.stor.Get(ctx, metric)
		if !ok {
			http.Error(w, fmt.Sprintf("%s not found", name), http.StatusNotFound)

			return
		}
		w.Write([]byte(strconv.FormatFloat(*val.Value, 'f', -1, 64)))
	case "counter":
		val, ok := h.stor.Get(ctx, metric)
		if !ok {
			http.Error(w, fmt.Sprintf("%s not found", name), http.StatusNotFound)

			return
		}
		w.Write([]byte(strconv.FormatInt(*val.Delta, 10)))

	default:
		http.Error(w, "wrong type", http.StatusBadRequest)

	}

}

// ShowAllHandler returns html with all known metrics
func (h *Handler) ShowAllHandler(w http.ResponseWriter, r *http.Request) {

	h.logger.Info("Entered ShowAllHandler")

	if r.Method != http.MethodGet {
		http.Error(w, "Only Get requests are allowed for value!", http.StatusMethodNotAllowed)

		return
	}
	h.logger.Info("method checked")

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	gaugeMetrics := []string{}
	counterMetrics := []string{}
	for key, val := range h.stor.GetAll(ctx) {
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

	w.Write([]byte(html))

}

// FaviconHandler returns Gopher!!!!
func (h *Handler) FaviconHandler(w http.ResponseWriter, r *http.Request) {

	h.logger.Info("Entered FaviconHandler")

	icon, err := os.ReadFile("../../images/gopher.png")
	if err != nil {
		http.Error(w, "Иконка не найдена", http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "image/png")

	w.Write(icon)
}
