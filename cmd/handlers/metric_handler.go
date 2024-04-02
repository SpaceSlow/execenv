package handlers

import (
	"errors"
	"net/http"

	"github.com/SpaceSlow/execenv/cmd/metrics"
	"github.com/SpaceSlow/execenv/cmd/storages"
)

type MetricHandler struct {
	MetricType metrics.MetricType
	Storage    storages.Storage
}

func (h MetricHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metric, err := metrics.ParseMetricFromPath(req.URL.Path, h.MetricType)
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
