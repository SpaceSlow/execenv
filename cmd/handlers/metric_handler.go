package handlers

import (
	"errors"
	"net/http"

	"github.com/SpaceSlow/execenv/cmd/metrics"
	"github.com/SpaceSlow/execenv/cmd/storages"
	"github.com/go-chi/chi/v5"
)

type MetricHandler struct {
	Storage storages.MetricStorage
}

func (h MetricHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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

	_ = h.Storage.Add(metric)

	res.WriteHeader(http.StatusOK)
}
