package handlers

import (
	"net/http"

	"github.com/SpaceSlow/execenv/cmd/storages"
)

type DBHandler struct {
	MetricStorage storages.MetricStorage
}

// Ping проверяет соединение с БД.
func (h DBHandler) Ping(res http.ResponseWriter, _ *http.Request) {
	dbStorage, ok := h.MetricStorage.(*storages.DBStorage)
	if ok && dbStorage.CheckConnection() {
		res.WriteHeader(http.StatusOK)
		return
	}
	res.WriteHeader(http.StatusInternalServerError)
}
