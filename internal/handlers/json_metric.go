package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/SpaceSlow/execenv/internal/metrics"
	"github.com/SpaceSlow/execenv/internal/storages"
)

// JSONMetricHandler хэндлер для обработки запросов в JSON-формате.
type JSONMetricHandler struct {
	MetricStorage storages.MetricStorage
}

func (h JSONMetricHandler) Post(res http.ResponseWriter, req *http.Request) {
	metric := &metrics.Metric{}
	if err := json.NewDecoder(req.Body).Decode(metric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if err := req.Body.Close(); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	var err error
	if metric, err = h.MetricStorage.Add(metric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	metricJSON, err := metric.MarshalJSON()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(metricJSON)
}

func (h JSONMetricHandler) BatchPost(res http.ResponseWriter, req *http.Request) {
	metricSlice := make([]metrics.Metric, 0)
	if err := json.NewDecoder(req.Body).Decode(&metricSlice); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if err := req.Body.Close(); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	var err error
	if err = h.MetricStorage.Batch(metricSlice); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func (h JSONMetricHandler) Get(res http.ResponseWriter, req *http.Request) {
	var jsonMetric *metrics.JSONMetric
	res.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(req.Body).Decode(&jsonMetric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Body.Close(); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	mType, err := metrics.ParseMetricType(jsonMetric.MType)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	metric, ok := h.MetricStorage.Get(mType, jsonMetric.ID)
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	metricJSON, err := metric.MarshalJSON()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(metricJSON)
}
