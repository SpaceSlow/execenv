package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/SpaceSlow/execenv/cmd/metrics"
	"github.com/SpaceSlow/execenv/cmd/storages"
)

type JSONMetricHandler struct {
	MetricStorage storages.MetricStorage
}

func (h JSONMetricHandler) Post(res http.ResponseWriter, req *http.Request) {
	var metric *metrics.Metric
	if err := json.NewDecoder(req.Body).Decode(&metric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	_ = h.MetricStorage.Add(metric)

	updMetric, _ := h.MetricStorage.Get(metric.Type, metric.Name)
	metricJSON, err := updMetric.MarshalJSON()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(metricJSON)
}

func (h JSONMetricHandler) Get(res http.ResponseWriter, req *http.Request) {
	var jsonMetric *metrics.JSONMetric
	if err := json.NewDecoder(req.Body).Decode(&jsonMetric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
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

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(metricJSON)
}
