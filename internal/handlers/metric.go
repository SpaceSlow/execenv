package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/SpaceSlow/execenv/internal/metrics"
	"github.com/SpaceSlow/execenv/internal/storages"
)

// MetricHandler хэндлер для обработки запросов text/plain-формата.
type MetricHandler struct {
	MetricStorage storages.MetricStorage
}

func (h MetricHandler) Post(res http.ResponseWriter, req *http.Request) {
	mType, err := metrics.ParseMetricType(chi.URLParam(req, "type"))
	if errors.Is(err, metrics.ErrIncorrectMetricTypeOrValue) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	name := chi.URLParam(req, "name")
	value := chi.URLParam(req, "value")
	metric, err := metrics.NewMetric(mType, name, value)
	if errors.Is(err, metrics.ErrIncorrectMetricTypeOrValue) {
		res.WriteHeader(http.StatusBadRequest)
		return
	} else if errors.Is(err, metrics.ErrEmptyMetricName) {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := h.MetricStorage.Add(metric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (h MetricHandler) Get(res http.ResponseWriter, req *http.Request) {
	mType, err := metrics.ParseMetricType(chi.URLParam(req, "type"))
	if errors.Is(err, metrics.ErrIncorrectMetricTypeOrValue) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	metric, ok := h.MetricStorage.Get(mType, chi.URLParam(req, "name"))
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(metric.ValueAsString()))
}

func (h MetricHandler) List(res http.ResponseWriter, _ *http.Request) {
	result := strings.Builder{}

	for _, metric := range h.MetricStorage.List() {
		result.WriteString(metric.String())
		result.WriteString("\n")
	}

	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(result.String()))
}
