package handlers

import (
	"net/http"

	"github.com/SpaceSlow/execenv/cmd/storages"
)

type DBHandler struct {
	MetricStorage storages.MetricStorage
}

func (h DBHandler) Ping(res http.ResponseWriter, _ *http.Request) {
	switch h.MetricStorage.(type) {
	case storages.DBStorage:
		storage := h.MetricStorage.(storages.DBStorage)
		if storage.CheckConnection() {
			res.WriteHeader(http.StatusOK)
			return
		}
	}
	res.WriteHeader(http.StatusInternalServerError)
}
