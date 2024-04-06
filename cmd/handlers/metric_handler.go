package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"

	"github.com/SpaceSlow/execenv/cmd/metrics"
	"github.com/SpaceSlow/execenv/cmd/storages"
)

type MetricHandler struct {
	MetricStorage storages.MetricStorage
}

func (h MetricHandler) Post(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	mType, err := metrics.ParseMetricType(chi.URLParam(req, "type"))
	if errors.Is(err, &metrics.IncorrectMetricTypeOrValueError{}) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	name := chi.URLParam(req, "name")
	value := chi.URLParam(req, "value")
	metric, err := metrics.NewMetric(mType, name, value)
	if errors.Is(err, &metrics.IncorrectMetricTypeOrValueError{}) {
		res.WriteHeader(http.StatusBadRequest)
		return
	} else if errors.Is(err, &metrics.EmptyMetricNameError{}) {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = h.MetricStorage.Add(metric)

	res.WriteHeader(http.StatusOK)
}

func (h MetricHandler) Get(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	mType, err := metrics.ParseMetricType(chi.URLParam(req, "type"))
	if errors.Is(err, &metrics.IncorrectMetricTypeOrValueError{}) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	metric, ok := h.MetricStorage.Get(mType, chi.URLParam(req, "name"))
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.Write([]byte(metric.ValueAsString()))
}

func (h MetricHandler) List(res http.ResponseWriter, _ *http.Request) {
	result := strings.Builder{}

	for _, metric := range h.MetricStorage.List() {
		result.WriteString(metric.String())
		result.WriteString("\n")
	}

	res.Write([]byte(result.String()))
}
