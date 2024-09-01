package handlers

import (
	"net/http"

	"github.com/SpaceSlow/execenv/internal/storages"
)

type CheckConnectionHandler struct {
	storage storages.ICheckConnection
}

func NewCheckConnectionHandler(storage storages.MetricStorage) *CheckConnectionHandler {
	if s, ok := storage.(storages.ICheckConnection); ok {
		return &CheckConnectionHandler{storage: s}
	}
	return &CheckConnectionHandler{storage: nil}
}

// Ping проверяет соединение с БД.
func (h CheckConnectionHandler) Ping(res http.ResponseWriter, _ *http.Request) {
	if h.storage.CheckConnection() {
		res.WriteHeader(http.StatusOK)
		return
	}
	res.WriteHeader(http.StatusInternalServerError)
}
